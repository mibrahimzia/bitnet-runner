/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#1a1a1a",
        surface: "#252525",
        primary: "#10b981", // Emerald green
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'), // <-- ADD THIS
  ],
}