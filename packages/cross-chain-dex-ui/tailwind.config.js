/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Core brand (surge.wtf light theme)
        'surge-primary': '#172342',        // Deep navy — primary dark / CTA base
        'surge-secondary': '#7bacff',      // Periwinkle blue — secondary accent
        'surge-accent': '#be9eff',         // Lavender — tertiary accent

        // Surface tokens (names preserved; values flipped to light)
        'surge-dark': '#fcfcfc',           // Page background (off-white)
        'surge-card': '#ffffff',           // Card surface (white)
        'surge-card-hover': '#f5f7fb',     // Hovered card surface
        'surge-border': '#e5e7eb',         // Subtle light border

        // Text
        'surge-text': '#172342',           // Primary text (navy)
        'surge-muted': '#5e6575',          // Muted / secondary text

        // Pastel accent set from surge.wtf
        'surge-mint': '#8ce8ab',           // Success / positive
        'surge-aqua': '#90eee4',           // Info / neutral highlight
        'surge-peach': '#fabeab',          // Warning / warm accent
        'surge-amber': '#f1aa47',          // Strong warn / CTA contrast
        'surge-lavender': '#be9eff',       // Brand accent (same as accent)
      },
      backgroundImage: {
        'surge-gradient': 'linear-gradient(135deg, #ffffff 0%, #f5f7fb 50%, #ffffff 100%)',
        'surge-glow': 'radial-gradient(ellipse at top, rgba(123, 172, 255, 0.18) 0%, transparent 55%)',
      },
    },
  },
  plugins: [],
}
