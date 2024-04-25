/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.{templ, html, js}"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["winter", "dark"]
  },
  plugins: [require('daisyui')],
}

