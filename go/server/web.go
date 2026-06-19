/*
Copyright 2024 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

//go:embed templates/web
var webTemplatesFS embed.FS

//go:embed static/web
var webStaticFS embed.FS

// webRenderer implements gin's render.HTMLRender. Each public page is parsed
// into its own template set (base layout + shared partials + the page body),
// keyed by the page's file name. Rendering always executes the "base" template,
// which pulls in the page's "content" block.
type webRenderer struct {
	templates map[string]*template.Template
}

func (r webRenderer) Instance(name string, data any) render.Render {
	// A name of the form "page#block" renders just that defined block from the
	// page's template set (used for HTMX partial swaps); a bare "page" renders
	// the full "base" layout.
	page, block := name, "base"
	if i := strings.IndexByte(name, '#'); i != -1 {
		page, block = name[:i], name[i+1:]
	}
	return render.HTML{
		Template: r.templates[page],
		Name:     block,
		Data:     data,
	}
}

var webFuncMap = template.FuncMap{
	"hasPrefix":  strings.HasPrefix,
	"badgeClass": badgeClass,
}

type webNavItem struct {
	To    string
	Title string
}

// webNavItems is the shared top navigation, mirroring the React Navbar.
var webNavItems = []webNavItem{
	{To: "/status", Title: "Status"},
	{To: "/daily", Title: "Daily"},
	{To: "/compare", Title: "Compare"},
	{To: "/fk", Title: "Foreign Keys"},
	{To: "/pr", Title: "PR"},
	{To: "/history", Title: "History"},
}

// webPage builds the data common to every server-rendered page (title, current
// path for active-link styling, and the nav items) and merges in page-specific
// values.
func webPage(c *gin.Context, title string, extra gin.H) gin.H {
	data := gin.H{
		"Title":    title,
		"Path":     c.Request.URL.Path,
		"NavItems": webNavItems,
	}
	for k, v := range extra {
		data[k] = v
	}
	return data
}

// renderWeb renders the page within the full base layout, or — when the request
// was issued by HTMX — just the named block, so a form/link with hx-get swaps a
// single region in place instead of reloading the whole page.
func renderWeb(c *gin.Context, page, block string, data gin.H) {
	name := page
	if c.GetHeader("HX-Request") == "true" {
		name = page + "#" + block
	}
	c.HTML(http.StatusOK, name, data)
}

// loadWebTemplates parses every page under templates/web/pages together with the
// base layout and all shared partials. It panics on a parse error because a
// broken template is a programming error that must fail loudly at startup.
func loadWebTemplates() webRenderer {
	const (
		base        = "templates/web/base.html"
		partialGlob = "templates/web/partials/*.html"
		pageGlob    = "templates/web/pages/*.html"
	)

	partials, err := fs.Glob(webTemplatesFS, partialGlob)
	if err != nil {
		panic(err)
	}
	pages, err := fs.Glob(webTemplatesFS, pageGlob)
	if err != nil {
		panic(err)
	}

	r := webRenderer{templates: make(map[string]*template.Template, len(pages))}
	for _, page := range pages {
		files := make([]string, 0, len(partials)+2)
		files = append(files, base)
		files = append(files, partials...)
		files = append(files, page)

		t := template.Must(template.New(path.Base(base)).Funcs(webFuncMap).ParseFS(webTemplatesFS, files...))
		name := strings.TrimSuffix(path.Base(page), ".html")
		r.templates[name] = t
	}
	return r
}

// registerWebRoutes wires the server-side rendered public website onto the Gin
// router, alongside the existing JSON API. It must be called after the router
// is created in Run().
func (s *Server) registerWebRoutes() {
	s.router.HTMLRender = loadWebTemplates()

	staticSub, err := fs.Sub(webStaticFS, "static/web")
	if err != nil {
		panic(err)
	}
	s.router.StaticFS("/static", http.FS(staticSub))

	s.router.GET("/", s.webHome)
	s.router.GET("/home", s.webHome)
	s.router.GET("/status", s.webStatus)
	s.router.GET("/daily", s.webDaily)
	s.router.GET("/compare", s.webCompare)
	s.router.GET("/macrobench/queries/compare", s.webMacroQueriesCompare)
	s.router.GET("/pr", s.webPRs)
	s.router.GET("/pr/:nb", s.webPR)
	s.router.GET("/fk", s.webForeignKeys)
	s.router.GET("/history", s.webHistory)
}

// webDailyCard holds the data for a single daily-summary card on the home page.
type webDailyCard struct {
	Name string
	// QPSJSON is a JSON array of total-QPS values (e.g. "[1,2,3]") read from a
	// data attribute by app.js to draw the card's sparkline. It contains only
	// digits/brackets/commas, so it is safe in an HTML attribute.
	QPSJSON string
	HasData bool
	Error   bool
}

func (s *Server) webHome(c *gin.Context) {
	workloads := []string{"OLTP", "TPCC"}
	results, err := macrobench.SearchForLast30DaysQPSOnly(s.dbClient, workloads, macrobench.Gen4Planner)

	cards := make([]webDailyCard, 0, len(workloads))
	for _, name := range workloads {
		card := webDailyCard{Name: name}
		if err != nil {
			card.Error = true
		} else {
			qps := make([]float64, 0, len(results[name]))
			for _, r := range results[name] {
				qps = append(qps, r.TotalQPS.Center)
			}
			card.HasData = len(qps) > 0
			if b, mErr := json.Marshal(qps); mErr == nil {
				card.QPSJSON = string(b)
			} else {
				card.QPSJSON = "[]"
			}
		}
		cards = append(cards, card)
	}
	if err != nil {
		slog.Error(err)
	}

	c.HTML(http.StatusOK, "home", webPage(c, "Vitess | Benchmark", gin.H{
		"Cards": cards,
	}))
}
