/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        // Core brand (surge.wtf light theme)
        "surge-primary": "#172342",
        "surge-secondary": "#7bacff",
        "surge-accent": "#be9eff",

        // Surface tokens
        "surge-dark": "#fcfcfc",
        "surge-card": "#ffffff",
        "surge-card-hover": "#f5f7fb",
        "surge-border": "#e5e7eb",

        // Text
        "surge-text": "#172342",
        "surge-muted": "#5e6575",

        // Pastel accent set
        "surge-mint": "#8ce8ab",
        "surge-aqua": "#90eee4",
        "surge-peach": "#fabeab",
        "surge-amber": "#f1aa47",
        "surge-lavender": "#be9eff",
      },
      backgroundImage: {
        "surge-gradient":
          "linear-gradient(135deg, #ffffff 0%, #f5f7fb 50%, #ffffff 100%)",
        "surge-glow":
          "radial-gradient(ellipse at top, rgba(123, 172, 255, 0.18) 0%, transparent 55%)",
      },
    },
  },
  plugins: [],
};
