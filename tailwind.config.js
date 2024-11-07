const { title } = require('process');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["**/*/*.templ", ".././*.html", "'./*.templ'"],
  theme: {
    extend: {
      fontFamily: {
        title: ['Vazirmatn', 'sans-serif'],

      }
    },
  },
  plugins: [],
}

