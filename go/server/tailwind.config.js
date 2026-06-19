// Tailwind config for the server-rendered website, compiled at image-build time
// by the standalone Tailwind CLI (no Node) into static/web/app.css, which is then
// embedded into the binary. Mirrors the design tokens previously configured inline
// via the Play CDN in base.html.
//
// `safelist` pins the badge classes built by the Go badgeClass() helper
// (web_pages.go): Tailwind only scans HTML/JS, so a variant passed dynamically
// (e.g. `badgeClass .StatusVariant`) never appears as a literal anywhere and would
// otherwise be purged, leaving status/profile badges colorless.
module.exports = {
  darkMode: "class",
  content: [
    "./go/server/templates/web/**/*.html",
    "./go/server/static/web/*.js",
  ],
  safelist: [
    "bg-primary", "text-primary-foreground",
    "bg-success", "text-success-foreground",
    "bg-destructive", "text-destructive-foreground",
    "bg-warning", "text-warning-foreground",
    "bg-progress", "text-progress-foreground",
    // Per-row grouping classes for the macro compare table, built by
    // compareRows() in web_pages.go and never present literally in a template.
    "bg-muted/80", "bg-background", "border-foreground", "dark:border-none",
    "transition-colors", "hover:bg-muted/50",
  ],
  theme: {
    container: { center: true, padding: "2rem", screens: { "2xl": "1400px" } },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        front: "hsl(var(--front))",
        back: "hsl(var(--back))",
        theme: "hsl(var(--theme))",
        primary: { DEFAULT: "hsl(var(--primary))", foreground: "hsl(var(--primary-foreground))" },
        secondary: { DEFAULT: "hsl(var(--secondary))", foreground: "hsl(var(--secondary-foreground))" },
        destructive: { DEFAULT: "hsl(var(--destructive))", foreground: "hsl(var(--destructive-foreground))" },
        muted: { DEFAULT: "hsl(var(--muted))", foreground: "hsl(var(--muted-foreground))" },
        accent: { DEFAULT: "hsl(var(--accent))", foreground: "hsl(var(--accent-foreground))" },
        warning: { DEFAULT: "hsl(var(--warning))", foreground: "hsl(var(--warning-foreground))" },
        success: { DEFAULT: "hsl(var(--success))", foreground: "hsl(var(--success-foreground))" },
        progress: { DEFAULT: "hsl(var(--progress))", foreground: "hsl(var(--progress-foreground))" },
        popover: { DEFAULT: "hsl(var(--popover))", foreground: "hsl(var(--popover-foreground))" },
        card: { DEFAULT: "hsl(var(--card))", foreground: "hsl(var(--card-foreground))" },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      fontFamily: { opensans: '"Open Sans", sans-serif' },
      screens: {
        mobile: { max: "780px" },
        midscreen: { max: "1536px" },
        widescreen: { min: "1536px" },
      },
    },
  },
};
