/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html,js}"],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
    },
  },
  plugins: [],
};
