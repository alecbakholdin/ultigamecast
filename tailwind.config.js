/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.{templ, html, js}"],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui')],
}

