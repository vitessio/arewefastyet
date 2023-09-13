/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "rgb(var(--color-primary) / <alpha-value>)",
        secondary: "rgb(var(--color-secondary) / <alpha-value>)",
        accent: "rgb(var(--color-accent) / <alpha-value>)",
        foreground: "rgb(var(--color-foreground) / <alpha-value>)",
        background: "rgb(var(--color-background) / <alpha-value>)",
        front: "rgb(var(--color-front) / <alpha-value>)",
        back: "rgb(var(--color-back) / <alpha-value>)",
      },
      screens: {
        mobile: { max: "780px" },
        widescreen: { min: "780px" },
      },
      fontFamily: {
        opensans: '"Open Sans", sans-serif',
      },
      borderRadius: {
        inherit: "inherit",
      },
      transitionDuration: {
        inherit: "inherit",
      },
      zIndex: {
        1: 1,
      },
    },
  },
  plugins: [],
};
