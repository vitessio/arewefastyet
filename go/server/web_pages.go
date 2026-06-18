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
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/vitessio/arewefastyet/go/exec"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
)

// ---- formatting helpers (mirroring website/src/utils/Utils.tsx) ----

// webFixed mirrors the React `fixed` helper: zero/NaN render as "0", otherwise
// the value is rendered with a fixed number of fraction digits.
func webFixed(v float64, digits int) string {
	if v == 0 || math.IsNaN(v) {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', digits, 64)
}

func webSecondToMicrosecond(v float64) string {
	return webFixed(v*1000000, 2) + "μs"
}

// webFormatByte mirrors the npm `bytes` formatter used by the React table:
// base-2 units, up to 2 decimals with trailing zeros stripped, no separator.
func webFormatByte(v float64) string {
	if v == 0 || math.IsNaN(v) {
		return "0"
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	neg := v < 0
	if neg {
		v = -v
	}
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if neg {
		s = "-" + s
	}
	return s + units[i]
}

func webRange(r macrobench.Range) string {
	if r.Infinite {
		return "∞"
	}
	if r.Unknown {
		return "?"
	}
	return "±" + webFixed(r.Value, 1) + "%"
}

// ---- Status page ----

type webStatCard struct {
	Title string
	Value string
}

type webQueueRow struct {
	SHAShort  string
	Source    string
	Workload  string
	PRDisplay string
	Profile   string
}

type webExecutionRow struct {
	UUIDShort     string
	SHAShort      string
	Source        string
	Workload      string
	StartedRel    string
	StartedTitle  string
	FinishedRel   string
	FinishedTitle string
	PRDisplay     string
	Golang        string
	Status        string
	StatusVariant string
	Profile       string
	GitRef        string
}

// badgeClass returns the Tailwind classes for a shadcn/ui Badge variant, so the
// server-rendered badges match the React component.
func badgeClass(variant string) string {
	const base = "inline-flex items-center rounded-full border border-transparent px-2.5 py-0.5 text-xs font-semibold"
	switch variant {
	case "success":
		return base + " bg-success text-success-foreground"
	case "destructive":
		return base + " bg-destructive text-destructive-foreground"
	case "warning":
		return base + " bg-warning text-warning-foreground"
	case "progress":
		return base + " bg-progress text-progress-foreground"
	default:
		return base + " bg-primary text-primary-foreground"
	}
}

func statusVariant(status string) string {
	switch status {
	case "failed":
		return "destructive"
	case "canceled":
		return "warning"
	case "started":
		return "progress"
	default:
		return "success"
	}
}

func shortStr(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

func prDisplay(nb int) string {
	if nb == 0 {
		return "N/A"
	}
	return strconv.Itoa(nb)
}

func profileDisplay(binary, mode string) string {
	if binary != "" && mode != "" {
		return binary + "|" + mode
	}
	return ""
}

const absTimeLayout = "Jan 2, 2006, 3:04 PM MST"

func (s *Server) webStatus(c *gin.Context) {
	extra := gin.H{}

	stats, err := exec.GetBenchmarkStats(s.dbClient)
	if err != nil {
		slog.Error(err)
		extra["StatsError"] = true
	} else {
		extra["Cards"] = []webStatCard{
			{Title: "Benchmark in total", Value: humanize.Comma(int64(stats.Total))},
			{Title: "Benchmark this month", Value: humanize.Comma(int64(stats.Last30Days))},
			{Title: "Commits benchmarked", Value: humanize.Comma(int64(stats.Commits))},
			{Title: "Execution duration (avg)", Value: webFixed(stats.AvgDuration, 2) + " min"},
		}
		if b, mErr := json.Marshal(stats.Last7Days); mErr == nil {
			extra["Last7DaysJSON"] = string(b)
		} else {
			extra["Last7DaysJSON"] = "[]"
		}
	}

	// Execution queue (in-memory, no DB needed).
	queueRows := make([]webQueueRow, 0, len(queue))
	for _, e := range queue {
		if e.Executing {
			continue
		}
		queueRows = append(queueRows, webQueueRow{
			SHAShort:  shortStr(e.identifier.GitRef, 8),
			Source:    e.identifier.Source,
			Workload:  e.identifier.Workload,
			PRDisplay: prDisplay(e.identifier.PullNb),
			Profile:   profileProto(e.identifier),
		})
	}
	sort.Slice(queueRows, func(i, j int) bool {
		return queueRows[i].SHAShort > queueRows[j].SHAShort && queueRows[i].Source > queueRows[j].Source
	})
	extra["QueueRows"] = queueRows

	// Previous (recent) executions.
	execs, err := exec.GetRecentExecutions(s.dbClient)
	if err != nil {
		slog.Error(err)
		extra["ExecutionsError"] = true
	} else {
		rows := make([]webExecutionRow, 0, len(execs))
		for _, e := range execs {
			row := webExecutionRow{
				UUIDShort:     shortStr(e.RawUUID, 8),
				SHAShort:      shortStr(e.GitRef, 8),
				Source:        e.Source,
				Workload:      e.Workload,
				PRDisplay:     prDisplay(e.PullNB),
				Golang:        e.GolangVersion,
				Status:        e.Status,
				StatusVariant: statusVariant(e.Status),
				Profile:       profileDisplay(e.ProfileInformation.Binary, e.ProfileInformation.Mode),
				GitRef:        e.GitRef,
			}
			if e.StartedAt != nil {
				row.StartedRel = humanize.Time(*e.StartedAt)
				row.StartedTitle = e.StartedAt.Format(absTimeLayout)
			}
			if e.FinishedAt != nil {
				row.FinishedRel = humanize.Time(*e.FinishedAt)
				row.FinishedTitle = e.FinishedAt.Format(absTimeLayout)
			} else {
				row.FinishedRel = "N/A"
			}
			rows = append(rows, row)
		}
		extra["ExecutionRows"] = rows
	}

	c.HTML(http.StatusOK, "status", webPage(c, "Status | Vitess Benchmark", extra))
}

// profileProto formats the profile badge for a queued execution.
func profileProto(id executionIdentifier) string {
	if id.Profile != nil {
		return profileDisplay(id.Profile.Binary, id.Profile.Mode)
	}
	return ""
}

// ---- Pull Requests page ----

type webPRRow struct {
	ID          int
	Title       string
	Author      string
	OpenedRel   string
	OpenedTitle string
}

func (s *Server) webPRs(c *gin.Context) {
	extra := gin.H{}

	prNumbers, err := exec.GetPullRequestList(s.dbClient)
	if err != nil {
		slog.Error(err)
		extra["Error"] = true
		c.HTML(http.StatusOK, "prs", webPage(c, "Pull Requests | Vitess Benchmark", extra))
		return
	}

	rows := make([]webPRRow, 0, len(prNumbers))
	for _, nb := range prNumbers {
		info, err := s.ghApp.GetPullRequestInfo(nb)
		if err != nil {
			slog.Error(err)
			extra["Error"] = true
			c.HTML(http.StatusOK, "prs", webPage(c, "Pull Requests | Vitess Benchmark", extra))
			return
		}
		row := webPRRow{ID: info.ID, Title: info.Title, Author: info.Author}
		if info.CreatedAt != nil {
			row.OpenedRel = humanize.Time(*info.CreatedAt)
			row.OpenedTitle = info.CreatedAt.Format(absTimeLayout)
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID > rows[j].ID })
	extra["Rows"] = rows

	c.HTML(http.StatusOK, "prs", webPage(c, "Pull Requests | Vitess Benchmark", extra))
}

// ---- Pull Request detail page ----

type webPRDetail struct {
	Number       int
	Title        string
	Author       string
	CreatedRel   string
	CreatedTitle string
	Base         string
	Head         string
	CanCompare   bool
}

func (s *Server) webPR(c *gin.Context) {
	nbStr := c.Param("nb")
	extra := gin.H{"Number": nbStr}

	nb, err := strconv.Atoi(nbStr)
	if err != nil {
		extra["NotFound"] = true
		c.HTML(http.StatusOK, "pr_detail", webPage(c, "Pull Request | Vitess Benchmark", extra))
		return
	}
	extra["Number"] = nb

	info, err := s.ghApp.GetPullRequestInfo(nb)
	if err != nil {
		// GitHub returns 404 for an unknown PR; everything else is a real failure.
		if strings.Contains(err.Error(), "404") {
			extra["NotFound"] = true
		} else {
			slog.Error(err)
			extra["Error"] = true
		}
		c.HTML(http.StatusOK, "pr_detail", webPage(c, "Pull Request | Vitess Benchmark", extra))
		return
	}

	gitPRInfo, err := exec.GetPullRequestInfo(s.dbClient, nb, info)
	if err != nil {
		slog.Error(err)
		extra["Error"] = true
		c.HTML(http.StatusOK, "pr_detail", webPage(c, "Pull Request | Vitess Benchmark", extra))
		return
	}

	detail := webPRDetail{
		Number:     nb,
		Title:      info.Title,
		Author:     info.Author,
		Base:       gitPRInfo.Base,
		Head:       gitPRInfo.Head,
		CanCompare: gitPRInfo.Base != "" && gitPRInfo.Head != "",
	}
	if info.CreatedAt != nil {
		detail.CreatedRel = humanize.Time(*info.CreatedAt)
		detail.CreatedTitle = info.CreatedAt.Format(absTimeLayout)
	}
	extra["PR"] = detail

	c.HTML(http.StatusOK, "pr_detail", webPage(c, "Pull Request | Vitess Benchmark", extra))
}

// ---- Foreign Keys page (macro comparison of two workloads on one commit) ----

// webCompareRow is one row of the macro-benchmark comparison table.
type webCompareRow struct {
	Title          string
	Old            string
	New            string
	P              string
	PBg            string
	PFg            string
	PText          string
	Delta          string
	ShowDeltaBadge bool
	DeltaVariant   string
}

type compareMetric struct {
	title         string
	res           macrobench.StatisticalResult
	format        func(float64) string
	lowerIsBetter bool
}

func significance(p float64) (bg, fg, text string) {
	switch {
	case p <= 0.01:
		return "#2E7D32", "#FFFFFF", "Statistically Significant"
	case p <= 0.05:
		return "#388E3C", "#FFFFFF", "Moderate Significance"
	case p <= 0.10:
		return "#6A9A1F", "#FFFFFF", "Marginal Significance"
	default:
		return "#9E9E9E", "#000000", "Not Statistically Significant"
	}
}

func deltaVariant(delta float64, lowerIsBetter bool) string {
	if delta == 0 {
		return "warning"
	}
	if lowerIsBetter {
		if delta < 0 {
			return "success"
		}
		return "destructive"
	}
	if delta > 0 {
		return "success"
	}
	return "destructive"
}

func plainFmt(v float64) string { return webFixed(v, 2) }
func msFmt(v float64) string    { return webFixed(v, 2) + "ms" }

// compareRows turns the statistical comparison into the ordered, formatted rows
// the template renders, mirroring the React MacroBenchmarkTable.
func compareRows(r macrobench.StatisticalCompareResults) []webCompareRow {
	metrics := []compareMetric{
		{"QPS Total", r.TotalQPS, plainFmt, false},
		{"Reads", r.ReadsQPS, plainFmt, false},
		{"Writes", r.WritesQPS, plainFmt, false},
		{"Other", r.OtherQPS, plainFmt, false},
		{"TPS", r.TPS, plainFmt, false},
		{"P95 Latency", r.Latency, msFmt, true},
		{"Errors / Second", r.Errors, plainFmt, true},
		{"Total CPU / Query", r.TotalComponentsCPUTime, webSecondToMicrosecond, true},
		{"vtgate", r.ComponentsCPUTime["vtgate"], webSecondToMicrosecond, true},
		{"vttablet", r.ComponentsCPUTime["vttablet"], webSecondToMicrosecond, true},
		{"Total Allocated / Query", r.TotalComponentsMemStatsAllocBytes, webFormatByte, true},
		{"vtgate", r.ComponentsMemStatsAllocBytes["vtgate"], webFormatByte, true},
		{"vttablet", r.ComponentsMemStatsAllocBytes["vttablet"], webFormatByte, true},
	}

	rows := make([]webCompareRow, 0, len(metrics))
	for _, m := range metrics {
		bg, fg, text := significance(m.res.P)
		row := webCompareRow{
			Title:          m.title,
			Old:            m.format(m.res.Old.Center) + " (" + webRange(m.res.Old.Range) + ")",
			New:            m.format(m.res.New.Center) + " (" + webRange(m.res.New.Range) + ")",
			P:              webFixed(m.res.P, 3),
			PBg:            bg,
			PFg:            fg,
			PText:          text,
			Delta:          webFixed(m.res.Delta, 3) + "%",
			ShowDeltaBadge: m.res.P <= 0.10,
			DeltaVariant:   deltaVariant(m.res.Delta, m.lowerIsBetter),
		}
		rows = append(rows, row)
	}
	return rows
}

func (s *Server) webForeignKeys(c *gin.Context) {
	sha := c.Query("sha")
	oldWorkload := c.Query("oldWorkload")
	newWorkload := c.Query("newWorkload")

	// Only TPCC workloads support the managed/unmanaged foreign-key comparison.
	tpccWorkloads := make([]string, 0)
	for _, w := range s.workloads {
		if strings.Contains(w, "TPCC") {
			tpccWorkloads = append(tpccWorkloads, w)
		}
	}
	sort.Strings(tpccWorkloads)

	extra := gin.H{
		"Workloads":   tpccWorkloads,
		"SHA":         sha,
		"OldWorkload": oldWorkload,
		"NewWorkload": newWorkload,
		"HasQuery":    sha != "" && oldWorkload != "" && newWorkload != "",
	}

	if extra["HasQuery"].(bool) {
		results, err := macrobench.CompareFKs(s.dbClient, oldWorkload, newWorkload, sha, macrobench.Gen4Planner)
		if err != nil {
			slog.Error(err)
			extra["Error"] = true
		} else if results.MissingResults {
			extra["MissingResults"] = true
		} else {
			extra["Rows"] = compareRows(results)
			extra["OldLabel"] = oldWorkload
			extra["NewLabel"] = newWorkload
		}
	}

	renderWeb(c, "fk", "fk_results", webPage(c, "Foreign Keys | Vitess Benchmark", extra))
}

// ---- Daily page (30-day summary sparklines + per-workload trend charts) ----

// webDailySummaryCard is one clickable sparkline card (one per workload) in the
// daily summary row. Selected highlights the card for the workload currently
// being charted below.
type webDailySummaryCard struct {
	Name     string
	QPSJSON  string
	HasData  bool
	Selected bool
}

// webChartSeries is one line of a multi-series chart. Color is a CSS color
// string consumed directly by Chart.js (see app.js initLineCharts).
type webChartSeries struct {
	Label string    `json:"label"`
	Color string    `json:"color"`
	Data  []float64 `json:"data"`
}

// webDailyChart is one collapsible chart card. DataJSON is the marshaled
// {labels, series} payload read from the canvas's data-linechart attribute.
type webDailyChart struct {
	Title    string
	DataJSON string
}

func round2(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return math.Round(v*100) / 100
}

// buildDailyCharts shapes the 30-day series into the five charts the Daily page
// renders, mirroring website DailyCharts.tsx (metrics, labels, μs conversion).
func buildDailyCharts(data []macrobench.StatisticalSingleResult) []webDailyChart {
	labels := make([]string, 0, len(data))
	var qpsReads, qpsWrites, qpsOther, qpsTotal, tps, latency []float64
	var cpuTotal, cpuVtgate, cpuVttablet, memTotal, memVtgate, memVttablet []float64
	for _, r := range data {
		labels = append(labels, shortStr(r.GitRef, 8))
		qpsReads = append(qpsReads, r.ReadsQPS.Center)
		qpsWrites = append(qpsWrites, r.WritesQPS.Center)
		qpsOther = append(qpsOther, r.OtherQPS.Center)
		qpsTotal = append(qpsTotal, r.TotalQPS.Center)
		tps = append(tps, r.TPS.Center)
		latency = append(latency, r.Latency.Center)
		cpuTotal = append(cpuTotal, round2(r.TotalComponentsCPUTime.Center*1000000))
		cpuVtgate = append(cpuVtgate, round2(r.ComponentsCPUTime["vtgate"].Center*1000000))
		cpuVttablet = append(cpuVttablet, round2(r.ComponentsCPUTime["vttablet"].Center*1000000))
		memTotal = append(memTotal, r.TotalComponentsMemStatsAllocBytes.Center)
		memVtgate = append(memVtgate, r.ComponentsMemStatsAllocBytes["vtgate"].Center)
		memVttablet = append(memVttablet, r.ComponentsMemStatsAllocBytes["vttablet"].Center)
	}

	// Colors ported from website/src/assets/styles/tailwind.css (--chart-*).
	const (
		colReads    = "hsl(220 100% 50%)"
		colWrites   = "hsl(0 59% 41%)"
		colOther    = "hsl(300 100% 25%)"
		colTotal    = "hsl(39 100% 50%)"
		colVtgate   = "hsl(220 100% 50%)"
		colVttablet = "hsl(0 59% 41%)"
	)

	specs := []struct {
		title  string
		series []webChartSeries
	}{
		{"QPS (Queries per second)", []webChartSeries{
			{"Reads", colReads, qpsReads},
			{"Writes", colWrites, qpsWrites},
			{"Other", colOther, qpsOther},
			{"Total", colTotal, qpsTotal},
		}},
		{"TPS (Transactions per second)", []webChartSeries{
			{"TPS", colTotal, tps},
		}},
		{"Latency (ms)", []webChartSeries{
			{"Latency", colTotal, latency},
		}},
		{"CPU / query (μs)", []webChartSeries{
			{"Vtgate", colVtgate, cpuVtgate},
			{"Vttablet", colVttablet, cpuVttablet},
			{"Total", colTotal, cpuTotal},
		}},
		{"Allocated / query (bytes)", []webChartSeries{
			{"Vtgate", colVtgate, memVtgate},
			{"Vttablet", colVttablet, memVttablet},
			{"Total", colTotal, memTotal},
		}},
	}

	charts := make([]webDailyChart, 0, len(specs))
	for _, sp := range specs {
		payload := struct {
			Labels []string         `json:"labels"`
			Series []webChartSeries `json:"series"`
		}{labels, sp.series}
		chart := webDailyChart{Title: sp.title}
		if b, err := json.Marshal(payload); err == nil {
			chart.DataJSON = string(b)
		} else {
			chart.DataJSON = "{}"
		}
		charts = append(charts, chart)
	}
	return charts
}

func (s *Server) webDaily(c *gin.Context) {
	workload := c.Query("workload")
	if workload == "" {
		workload = "OLTP"
	}
	extra := gin.H{"Workload": workload}

	// Summary sparkline cards: total-QPS over the last 30 days for every workload.
	summary, err := macrobench.SearchForLast30DaysQPSOnly(s.dbClient, s.workloads, macrobench.Gen4Planner)
	if err != nil {
		slog.Error(err)
		extra["SummaryError"] = true
	} else {
		cards := make([]webDailySummaryCard, 0, len(s.workloads))
		for _, name := range s.workloads {
			qps := make([]float64, 0, len(summary[name]))
			for _, r := range summary[name] {
				qps = append(qps, r.TotalQPS.Center)
			}
			card := webDailySummaryCard{Name: name, HasData: len(qps) > 0, Selected: name == workload}
			if b, mErr := json.Marshal(qps); mErr == nil {
				card.QPSJSON = string(b)
			} else {
				card.QPSJSON = "[]"
			}
			cards = append(cards, card)
		}
		sort.Slice(cards, func(i, j int) bool { return cards[i].Name < cards[j].Name })
		extra["Cards"] = cards
	}

	// Charts for the selected workload.
	data, err := macrobench.SearchForLast30Days(s.dbClient, workload, macrobench.Gen4Planner)
	if err != nil {
		slog.Error(err)
		extra["ChartsError"] = true
	} else if len(data) == 0 {
		extra["ChartsError"] = true
	} else {
		extra["Charts"] = buildDailyCharts(data)
	}

	renderWeb(c, "daily", "daily_body", webPage(c, "Daily | Vitess Benchmark", extra))
}

// ---- Compare page (macro comparison of two git refs across all workloads) ----

// webRefOption is one selectable vitess ref (a branch/release name and its
// commit hash) used to populate the compare form and to resolve names ↔ SHAs.
type webRefOption struct {
	Name string
	Hash string
}

// webCompareCard is one per-workload comparison card on the Compare page.
type webCompareCard struct {
	Workload       string
	MissingResults bool
	OldLabel       string
	NewLabel       string
	OldHref        string
	NewHref        string
	Rows           []webCompareRow
	QueryPlanHref  string
}

// vitessRefOptions builds the list of named vitess refs (main + release
// branches + comparable release tags), mirroring getLatestVitessGitRef. It
// degrades gracefully: any source that errors (e.g. the vitess clone is not
// ready) is simply skipped so raw-SHA comparison still works.
func (s *Server) vitessRefOptions() []webRefOption {
	var opts []webRefOption
	if sha, err := exec.GetLatestDailyJobForMacrobenchmarks(s.dbClient); err == nil && sha != "" {
		opts = append(opts, webRefOption{Name: "main", Hash: sha})
	}
	if branches, err := git.GetLatestVitessReleaseBranchCommitHash(s.getVitessPath()); err == nil {
		for _, b := range branches {
			opts = append(opts, webRefOption{Name: b.Name, Hash: b.CommitHash})
		}
	}
	if tags, err := git.GetAllComparableVitessReleases(s.getVitessPath()); err == nil {
		for _, t := range tags {
			opts = append(opts, webRefOption{Name: t.Name, Hash: t.CommitHash})
		}
	}
	return opts
}

// resolveRef turns a ref name (e.g. "main", "v19.0.0") into its commit hash, or
// returns the input unchanged when it is already a SHA / an unknown ref.
func resolveRef(ref string, opts []webRefOption) string {
	if ref == "" {
		return ""
	}
	for _, o := range opts {
		if o.Name == ref {
			return o.Hash
		}
	}
	return ref
}

// refLabel returns the friendly ref name for a commit hash, or the hash itself
// when it does not correspond to a named ref (mirrors Utils.getRefName).
func refLabel(sha string, opts []webRefOption) string {
	if sha == "" {
		return ""
	}
	for _, o := range opts {
		if o.Hash != "" && strings.Contains(o.Hash, sha) {
			return o.Name
		}
	}
	return sha
}

// compareMarkdown renders the per-workload comparison as a markdown document,
// mirroring ComparePage.generateCompareMarkdown for the "Copy as markdown" button.
func compareMarkdown(oldLabel, oldSHA, newLabel, newSHA string, cards []webCompareCard) string {
	formatRef := func(name, hash string) string {
		if name != hash {
			return name + " (`" + hash + "`)"
		}
		return "`" + hash + "`"
	}
	var b strings.Builder
	b.WriteString("# arewefastyet - Vitess benchmark comparison\n\n")
	b.WriteString("**Old:** " + formatRef(oldLabel, oldSHA) + "\n")
	b.WriteString("**New:** " + formatRef(newLabel, newSHA) + "\n\n")
	b.WriteString("---\n\n")
	for _, card := range cards {
		if card.MissingResults {
			continue
		}
		b.WriteString("## " + card.Workload + "\n\n")
		b.WriteString("| Metric | " + oldLabel + " | " + newLabel + " | P | Delta |\n")
		b.WriteString("|---|---|---|---|---|\n")
		for _, row := range card.Rows {
			b.WriteString("| " + row.Title + " | " + row.Old + " | " + row.New + " | " + row.P + " | " + row.Delta + " |\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (s *Server) webCompare(c *gin.Context) {
	oldRef := c.Query("old")
	newRef := c.Query("new")

	opts := s.vitessRefOptions()
	oldSHA := resolveRef(oldRef, opts)
	newSHA := resolveRef(newRef, opts)
	oldLabel := refLabel(oldSHA, opts)
	newLabel := refLabel(newSHA, opts)

	extra := gin.H{
		"OldRef":   oldRef,
		"NewRef":   newRef,
		"Refs":     opts,
		"HasQuery": oldSHA != "" && newSHA != "",
	}

	if oldSHA != "" && newSHA != "" {
		results, err := macrobench.Compare(s.dbClient, oldSHA, newSHA, s.workloads, macrobench.Gen4Planner)
		if err != nil {
			slog.Error(err)
			extra["Error"] = true
		} else {
			workloads := make([]string, 0, len(results))
			for w := range results {
				workloads = append(workloads, w)
			}
			sort.Strings(workloads)

			cards := make([]webCompareCard, 0, len(workloads))
			for _, w := range workloads {
				res := results[w]
				card := webCompareCard{
					Workload:       w,
					MissingResults: res.MissingResults,
					OldLabel:       oldLabel,
					NewLabel:       newLabel,
					OldHref:        "https://github.com/vitessio/vitess/commit/" + oldSHA,
					NewHref:        "https://github.com/vitessio/vitess/commit/" + newSHA,
					QueryPlanHref:  "/macrobench/queries/compare?old=" + url.QueryEscape(oldRef) + "&new=" + url.QueryEscape(newRef) + "&workload=" + url.QueryEscape(w),
				}
				if !res.MissingResults {
					card.Rows = compareRows(res)
				}
				cards = append(cards, card)
			}
			extra["Cards"] = cards
			extra["OldLabel"] = oldLabel
			extra["NewLabel"] = newLabel
			extra["Markdown"] = compareMarkdown(oldLabel, oldSHA, newLabel, newSHA, cards)
		}
	}

	renderWeb(c, "compare", "compare_results", webPage(c, "Compare | Vitess Benchmark", extra))
}

// ---- Macro query-plan compare page (/macrobench/queries/compare) ----

// webQueryPlanRow is one comparison row for a single normalized query. The
// table renders Key + ExecTimeDiff; the rest populate the per-row dialog. It is
// marshaled to JSON and consumed client-side by app.js queryPlansTable().
type webQueryPlanRow struct {
	Key          string `json:"key"`
	ExecTimeDiff int    `json:"execTimeDiff"`
	RowsDiff     int    `json:"rowsDiff"`
	ErrorsDiff   int    `json:"errorsDiff"`
	OldExecTime  string `json:"oldExecTime"`
	NewExecTime  string `json:"newExecTime"`
	OldRows      string `json:"oldRows"`
	NewRows      string `json:"newRows"`
	OldErrors    string `json:"oldErrors"`
	NewErrors    string `json:"newErrors"`
	OldPlan      string `json:"oldPlan"`
	NewPlan      string `json:"newPlan"`
	HasOld       bool   `json:"hasOld"`
	HasNew       bool   `json:"hasNew"`
}

// prettyJSON renders a query plan's Instructions (a JSON string, or already a
// decoded value) as indented JSON for display in the plan dialog.
func prettyJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		var tmp interface{}
		if err := json.Unmarshal([]byte(s), &tmp); err == nil {
			if b, err := json.MarshalIndent(tmp, "", "  "); err == nil {
				return string(b)
			}
		}
		return s
	}
	if b, err := json.MarshalIndent(v, "", "  "); err == nil {
		return string(b)
	}
	return ""
}

// buildQueryPlanRows turns the VTGate plan comparison into the rows the page
// renders, mirroring website MacroQueriesComparePage Columns/Dialog.
func buildQueryPlanRows(cmp []macrobench.VTGateQueryPlanComparer) []webQueryPlanRow {
	rows := make([]webQueryPlanRow, 0, len(cmp))
	for _, c := range cmp {
		row := webQueryPlanRow{
			Key:          c.Key,
			ExecTimeDiff: c.ExecTimeDiff,
			RowsDiff:     c.RowsReturnedDiff,
			ErrorsDiff:   c.ErrorsDiff,
			OldExecTime:  "N/A",
			NewExecTime:  "N/A",
			OldRows:      "N/A",
			NewRows:      "N/A",
			OldErrors:    "N/A",
			NewErrors:    "N/A",
		}
		if c.Left != nil {
			row.HasOld = true
			row.OldExecTime = strconv.Itoa(c.Left.Value.ExecTime)
			row.OldRows = strconv.Itoa(c.Left.Value.RowsReturned)
			row.OldErrors = strconv.Itoa(c.Left.Value.Errors)
			row.OldPlan = prettyJSON(c.Left.Value.Instructions)
		}
		if c.Right != nil {
			row.HasNew = true
			row.NewExecTime = strconv.Itoa(c.Right.Value.ExecTime)
			row.NewRows = strconv.Itoa(c.Right.Value.RowsReturned)
			row.NewErrors = strconv.Itoa(c.Right.Value.Errors)
			row.NewPlan = prettyJSON(c.Right.Value.Instructions)
		}
		rows = append(rows, row)
	}
	return rows
}

func (s *Server) webMacroQueriesCompare(c *gin.Context) {
	oldRef := c.Query("old")
	newRef := c.Query("new")
	workload := c.Query("workload")

	opts := s.vitessRefOptions()
	oldSHA := resolveRef(oldRef, opts)
	newSHA := resolveRef(newRef, opts)

	extra := gin.H{
		"OldRef":   oldRef,
		"NewRef":   newRef,
		"Workload": workload,
		"OldSHA":   oldSHA,
		"NewSHA":   newSHA,
		"OldLabel": refLabel(oldSHA, opts),
		"NewLabel": refLabel(newSHA, opts),
		"HasQuery": oldSHA != "" && newSHA != "" && workload != "",
	}

	if extra["HasQuery"].(bool) {
		oldPlans, err := macrobench.GetVTGateSelectQueryPlansWithFilter(oldSHA, macrobench.Workload(workload), macrobench.Gen4Planner, s.dbClient)
		if err != nil {
			slog.Error(err)
			extra["Error"] = true
		} else {
			newPlans, err := macrobench.GetVTGateSelectQueryPlansWithFilter(newSHA, macrobench.Workload(workload), macrobench.Gen4Planner, s.dbClient)
			if err != nil {
				slog.Error(err)
				extra["Error"] = true
			} else {
				rows := buildQueryPlanRows(macrobench.CompareVTGateQueryPlans(oldPlans, newPlans))
				extra["HasRows"] = len(rows) > 0
				if b, mErr := json.Marshal(rows); mErr == nil {
					extra["RowsJSON"] = string(b)
				} else {
					extra["RowsJSON"] = "[]"
				}
			}
		}
	}

	c.HTML(http.StatusOK, "macro_queries_compare", webPage(c, "Compare Query Plans | Vitess Benchmark", extra))
}

// ---- History page (/history) ----

// webHistoryRow is one benchmarked-commit row (a git_ref+source group). It is
// marshaled to JSON and rendered client-side by app.js historyTable().
type webHistoryRow struct {
	SHA          string `json:"sha"`
	SHAShort     string `json:"shaShort"`
	Source       string `json:"source"`
	Workloads    int    `json:"workloads"`
	StartedRel   string `json:"startedRel"`
	StartedTitle string `json:"startedTitle"`
}

func (s *Server) webHistory(c *gin.Context) {
	gitRef := c.Query("gitRef")
	extra := gin.H{"InitialGitRef": gitRef}

	history, err := exec.GetHistory(s.dbClient)
	if err != nil {
		slog.Error(err)
		extra["Error"] = true
		c.HTML(http.StatusOK, "history", webPage(c, "History | Vitess Benchmark", extra))
		return
	}

	rows := make([]webHistoryRow, 0, len(history))
	sourceSet := make(map[string]struct{})
	for _, h := range history {
		row := webHistoryRow{
			SHA:       h.SHA,
			SHAShort:  shortStr(h.SHA, 8),
			Source:    h.Source,
			Workloads: h.WorkloadsBenchmarked,
		}
		if h.StartedAt != nil {
			row.StartedRel = humanize.Time(*h.StartedAt)
			row.StartedTitle = h.StartedAt.Format(absTimeLayout)
		}
		rows = append(rows, row)
		sourceSet[h.Source] = struct{}{}
	}

	sources := make([]string, 0, len(sourceSet))
	for src := range sourceSet {
		sources = append(sources, src)
	}
	sort.Strings(sources)

	extra["Sources"] = sources
	extra["HasRows"] = len(rows) > 0
	if b, mErr := json.Marshal(rows); mErr == nil {
		extra["RowsJSON"] = string(b)
	} else {
		extra["RowsJSON"] = "[]"
	}

	c.HTML(http.StatusOK, "history", webPage(c, "History | Vitess Benchmark", extra))
}
