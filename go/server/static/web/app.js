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

// Top navigation shared with the Go template (see web.go webNavItems).
var NAV_ITEMS = [
  { to: "/status", title: "Status" },
  { to: "/daily", title: "Daily" },
  { to: "/compare", title: "Compare" },
  { to: "/fk", title: "Foreign Keys" },
  { to: "/pr", title: "PR" },
  { to: "/history", title: "History" },
  { to: "/", title: "Home" },
];

// setTheme applies and persists the light/dark/system theme. "system" clears the
// stored preference and follows the OS setting.
function setTheme(mode) {
  var root = document.documentElement;
  if (mode === "system") {
    localStorage.removeItem("theme");
    var prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    root.classList.toggle("dark", prefersDark);
  } else {
    localStorage.setItem("theme", mode);
    root.classList.toggle("dark", mode === "dark");
  }
}

// commandPalette is the Alpine.js data factory for the Cmd/Ctrl+K palette.
function commandPalette() {
  return {
    open: false,
    query: "",
    items: NAV_ITEMS,
    show() {
      this.open = true;
      this.query = "";
      this.$nextTick(() => {
        if (this.$refs.search) this.$refs.search.focus();
      });
    },
    filtered() {
      var q = this.query.trim().toLowerCase();
      if (!q) return this.items;
      return this.items.filter(function (i) {
        return i.title.toLowerCase().includes(q) || i.to.toLowerCase().includes(q);
      });
    },
    enter() {
      var f = this.filtered();
      if (f.length) window.location.href = f[0].to;
    },
  };
}

// cssHSL resolves a Tailwind design-token variable (e.g. "--primary", stored as
// "24.6 95% 53.1%") into a CSS hsl() color usable by Chart.js.
function cssHSL(varName, fallback) {
  var v = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
  if (!v) return fallback;
  return "hsl(" + v.replace(/,/g, " ") + ")";
}

// initSparklines draws a minimal QPS line chart into every uninitialized
// [data-sparkline] canvas, reading its series from the data-qps attribute.
function initSparklines(root) {
  if (typeof Chart === "undefined") return;
  var scope = root || document;
  var canvases = scope.querySelectorAll("canvas[data-sparkline]:not([data-rendered])");
  canvases.forEach(function (canvas) {
    var values;
    try {
      values = JSON.parse(canvas.dataset.qps || "[]");
    } catch (e) {
      values = [];
    }
    canvas.setAttribute("data-rendered", "true");
    var color = cssHSL("--primary", "hsl(24.6 95% 53.1%)");
    new Chart(canvas, {
      type: "line",
      data: {
        labels: values.map(function (_, i) {
          return i;
        }),
        datasets: [
          {
            data: values,
            borderColor: color,
            backgroundColor: color,
            borderWidth: 1,
            pointRadius: 0,
            pointHoverRadius: 4,
            tension: 0.3,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false }, tooltip: { enabled: true } },
        scales: { x: { display: false }, y: { display: false } },
      },
    });
  });
}

// initBarCharts draws a small bar chart into every uninitialized [data-barchart]
// canvas, reading its series from the data-values attribute (e.g. the Status
// page's "benchmarks over the last 7 days").
function initBarCharts(root) {
  if (typeof Chart === "undefined") return;
  var scope = root || document;
  var canvases = scope.querySelectorAll("canvas[data-barchart]:not([data-rendered])");
  canvases.forEach(function (canvas) {
    var values;
    try {
      values = JSON.parse(canvas.dataset.values || "[]");
    } catch (e) {
      values = [];
    }
    canvas.setAttribute("data-rendered", "true");
    var color = cssHSL("--primary", "hsl(24.6 95% 53.1%)");
    new Chart(canvas, {
      type: "bar",
      data: {
        labels: values.map(function (_, i) {
          return i;
        }),
        datasets: [
          {
            data: values,
            backgroundColor: color,
            borderColor: color,
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false }, tooltip: { enabled: true } },
        scales: { x: { display: false }, y: { display: false } },
      },
    });
  });
}

// initLineCharts draws a multi-series line chart into every uninitialized
// [data-linechart] canvas. The attribute holds a JSON {labels, series} payload
// (see web_pages.go buildDailyCharts), where each series is {label, color, data}.
function initLineCharts(root) {
  if (typeof Chart === "undefined") return;
  var scope = root || document;
  var canvases = scope.querySelectorAll("canvas[data-linechart]:not([data-rendered])");
  canvases.forEach(function (canvas) {
    var cfg;
    try {
      cfg = JSON.parse(canvas.dataset.linechart || "{}");
    } catch (e) {
      return;
    }
    canvas.setAttribute("data-rendered", "true");
    var axisColor = cssHSL("--muted-foreground", "#888");
    var gridColor = cssHSL("--border", "#ddd");
    var datasets = (cfg.series || []).map(function (s) {
      return {
        label: s.label,
        data: s.data,
        borderColor: s.color,
        backgroundColor: s.color,
        borderWidth: 2,
        pointRadius: 3,
        pointHoverRadius: 6,
        tension: 0.3,
      };
    });
    new Chart(canvas, {
      type: "line",
      data: { labels: cfg.labels || [], datasets: datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: "index", intersect: false },
        plugins: {
          legend: { display: true, position: "bottom", labels: { color: axisColor, usePointStyle: true } },
          tooltip: { enabled: true },
        },
        scales: {
          x: { ticks: { color: axisColor }, grid: { color: gridColor } },
          y: { ticks: { color: axisColor }, grid: { color: gridColor } },
        },
      },
    });
  });
}

function initCharts(root) {
  initSparklines(root);
  initBarCharts(root);
  initLineCharts(root);
}

// copyCompareMarkdown copies the server-rendered markdown summary (held in the
// hidden #compareMarkdown textarea) to the clipboard and flips the button label.
function copyCompareMarkdown(btn) {
  var ta = document.getElementById("compareMarkdown");
  if (!ta) return;
  navigator.clipboard.writeText(ta.value).then(function () {
    var label = btn.querySelector("[data-copy-label]");
    var icon = btn.querySelector("i");
    if (label) label.textContent = "Copied!";
    if (icon) {
      icon.classList.remove("fa-copy");
      icon.classList.add("fa-check");
    }
    setTimeout(function () {
      if (label) label.textContent = "Copy as markdown";
      if (icon) {
        icon.classList.remove("fa-check");
        icon.classList.add("fa-copy");
      }
    }, 2000);
  });
}

// queryPlansTable is the Alpine.js data factory for the Compare Query Plans
// table (see templates/web/pages/macro_queries_compare.html). It reads the row
// data from the host element's data-rows attribute and handles the text/operator
// filtering, exec-time sorting, pagination, and per-row plan dialog client-side,
// mirroring the React TanStack table + dialog.
function queryPlansTable() {
  return {
    rows: [],
    query: "",
    operators: [],
    sortDir: "desc",
    page: 0,
    pageSize: 10,
    modalOpen: false,
    current: {},
    showCol: { query: true, execTime: true },
    init() {
      try {
        this.rows = JSON.parse(this.$el.dataset.rows || "[]");
      } catch (e) {
        this.rows = [];
      }
    },
    toggleOperator(op) {
      var idx = this.operators.indexOf(op);
      if (idx === -1) {
        this.operators.push(op);
      } else {
        this.operators.splice(idx, 1);
      }
      this.page = 0;
    },
    reset() {
      this.query = "";
      this.operators = [];
      this.page = 0;
    },
    toggleSort() {
      this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      this.page = 0;
    },
    filteredAll() {
      var q = this.query.trim().toLowerCase();
      var ops = this.operators;
      var res = this.rows.filter(function (r) {
        var k = (r.key || "").toLowerCase();
        if (q && k.indexOf(q) === -1) return false;
        if (ops.length && !ops.some(function (o) { return k.indexOf(o) !== -1; })) return false;
        return true;
      });
      var dir = this.sortDir === "asc" ? 1 : -1;
      res.sort(function (a, b) {
        return (a.execTimeDiff - b.execTimeDiff) * dir;
      });
      return res;
    },
    pageCount() {
      return Math.ceil(this.filteredAll().length / this.pageSize);
    },
    paged() {
      var start = this.page * this.pageSize;
      return this.filteredAll().slice(start, start + this.pageSize);
    },
    prevPage() {
      if (this.page > 0) this.page--;
    },
    nextPage() {
      if (this.page + 1 < this.pageCount()) this.page++;
    },
    openModal(r) {
      this.current = r;
      this.modalOpen = true;
    },
    badgeClasses(diff) {
      var base =
        "inline-flex items-center rounded-full border border-transparent px-2.5 py-0.5 text-xs font-semibold ";
      if (diff > 0) return base + "bg-success text-success-foreground";
      if (diff === 0) return base + "bg-warning text-warning-foreground";
      return base + "bg-destructive text-destructive-foreground";
    },
  };
}

// historyTable is the Alpine.js data factory for the History table (see
// templates/web/pages/history.html). Like queryPlansTable it reads its rows from
// the host element's data-rows attribute and handles the text filter, source
// faceted filter, column visibility, and pagination client-side. The data-initial
// attribute seeds the text filter from the ?gitRef= query param so a "Benchmarks
// History" row action deep-links to a pre-filtered table.
function historyTable() {
  return {
    rows: [],
    query: "",
    sources: [],
    page: 0,
    pageSize: 10,
    showCol: { sha: true, source: true, workloads: true, started: true },
    init() {
      try {
        this.rows = JSON.parse(this.$el.dataset.rows || "[]");
      } catch (e) {
        this.rows = [];
      }
      this.query = this.$el.dataset.initial || "";
    },
    toggleSource(src) {
      var idx = this.sources.indexOf(src);
      if (idx === -1) {
        this.sources.push(src);
      } else {
        this.sources.splice(idx, 1);
      }
      this.page = 0;
    },
    reset() {
      this.query = "";
      this.sources = [];
      this.page = 0;
    },
    filteredAll() {
      var q = this.query.trim().toLowerCase();
      var srcs = this.sources;
      return this.rows.filter(function (r) {
        if (q && (r.sha || "").toLowerCase().indexOf(q) === -1) return false;
        if (srcs.length && srcs.indexOf(r.source) === -1) return false;
        return true;
      });
    },
    pageCount() {
      return Math.ceil(this.filteredAll().length / this.pageSize);
    },
    paged() {
      var start = this.page * this.pageSize;
      return this.filteredAll().slice(start, start + this.pageSize);
    },
    prevPage() {
      if (this.page > 0) this.page--;
    },
    nextPage() {
      if (this.page + 1 < this.pageCount()) this.page++;
    },
  };
}

// tooltip is a small Alpine.js data factory for a hover tooltip, modeled on the
// React shadcn/Radix Tooltip used on the Status page timestamps. Following Radix
// (@radix-ui/react-tooltip + react-popper):
//   - opens after a 200ms delay on pointer enter, closes on pointer leave;
//   - default side "top", centered, with a 4px offset (Radix sideOffset);
//   - flips to "bottom" when there isn't room above (collision detection);
//   - closes on scroll (Radix dismisses on scroll), which also prevents a
//     position:fixed bubble from being stranded after the page scrolls. We
//     dismiss on wheel/touchmove too, not just scroll, because browsers defer
//     scroll events until scrolling settles — wheel/touchmove fire at the very
//     start of the gesture, so the bubble vanishes immediately.
// The bubble is position:fixed (our equivalent of Radix's portal) with coords
// from the trigger's bounding rect, so it escapes the executions table's
// overflow:auto clipping. The content stays mounted (opacity-toggled) so it can
// be measured for the flip decision and never gets stuck in a display:none state.
var TOOLTIP_OFFSET = 4;

function tooltip() {
  return {
    open: false,
    side: "top",
    x: 0,
    y: 0,
    _t: null,
    _dismiss: null,
    show(e) {
      var trigger = e.currentTarget;
      clearTimeout(this._t);
      this._t = setTimeout(
        function () {
          var t = trigger.getBoundingClientRect();
          var c = this.$refs.content.getBoundingClientRect();
          this.x = t.left + t.width / 2;
          if (t.top - c.height - TOOLTIP_OFFSET < 0) {
            this.side = "bottom";
            this.y = t.bottom + TOOLTIP_OFFSET;
          } else {
            this.side = "top";
            this.y = t.top - TOOLTIP_OFFSET;
          }
          this.open = true;
          this._dismiss = this.hide.bind(this);
          window.addEventListener("scroll", this._dismiss, true);
          window.addEventListener("wheel", this._dismiss, { capture: true, passive: true });
          window.addEventListener("touchmove", this._dismiss, { capture: true, passive: true });
        }.bind(this),
        200
      );
    },
    hide() {
      clearTimeout(this._t);
      this.open = false;
      if (this._dismiss) {
        window.removeEventListener("scroll", this._dismiss, true);
        window.removeEventListener("wheel", this._dismiss, true);
        window.removeEventListener("touchmove", this._dismiss, true);
        this._dismiss = null;
      }
    },
  };
}

// compareWidget is the Alpine.js data factory for the draggable header of the
// global "Compare Versions" widget (see partials/compare_widget.html). It tracks
// a translate offset updated on mouse drag; the widget's visibility and the
// old/new refs themselves live in the global $store.compare.
function compareWidget() {
  return {
    dragging: false,
    x: 0,
    y: 0,
    ox: 0,
    oy: 0,
    sx: 0,
    sy: 0,
    startDrag(e) {
      this.dragging = true;
      this.sx = e.clientX;
      this.sy = e.clientY;
      this.ox = this.x;
      this.oy = this.y;
    },
    onDrag(e) {
      if (!this.dragging) return;
      this.x = this.ox + (e.clientX - this.sx);
      this.y = this.oy + (e.clientY - this.sy);
    },
    endDrag() {
      this.dragging = false;
    },
    style() {
      return "transform: translate(" + this.x + "px," + this.y + "px)";
    },
  };
}

// The global compare store backs the cross-page "Compare Versions" widget,
// mirroring the React CompareContext. Row actions on the History (and, later,
// Status) tables call addOld/addNew to stage commits; go() navigates to the
// /compare page. Vitess ref names are fetched lazily from the JSON API to
// populate the input datalist, failing silently if unavailable.
document.addEventListener("alpine:init", function () {
  Alpine.store("compare", {
    old: "",
    new: "",
    visible: false,
    refs: [],
    refsLoaded: false,
    addOld(ref) {
      this.old = ref;
      this.show();
    },
    addNew(ref) {
      this.new = ref;
      this.show();
    },
    show() {
      this.visible = true;
      this.loadRefs();
    },
    close() {
      this.old = "";
      this.new = "";
      this.visible = false;
    },
    go() {
      if (this.old && this.new) {
        window.location.href =
          "/compare?old=" + encodeURIComponent(this.old) + "&new=" + encodeURIComponent(this.new);
      }
    },
    loadRefs() {
      if (this.refsLoaded) return;
      this.refsLoaded = true;
      var self = this;
      fetch("/api/vitess/refs")
        .then(function (r) {
          return r.ok ? r.json() : null;
        })
        .then(function (d) {
          if (!d) return;
          var out = [];
          (d.branches || []).forEach(function (b) {
            if (b && b.name) out.push(b.name);
          });
          (d.tags || []).forEach(function (t) {
            if (t && t.name) out.push(t.name);
          });
          self.refs = out;
        })
        .catch(function () {});
    },
  });
});

document.addEventListener("DOMContentLoaded", function () {
  initCharts(document);
});

// Re-initialize charts inside HTMX-swapped fragments (events bubble to document).
document.addEventListener("htmx:afterSwap", function (e) {
  initCharts(e.target);
});
