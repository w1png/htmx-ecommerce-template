/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{templ,html,js}"],
  theme: {
    extend: {
      transitionProperty: {
        height: "height",
        width: "width",
      },
      colors: {
        "admin-accent": "#000000",
        "admin-accent-dark": "#e2e8f0",
        "admin-text": "#000000",
        "admin-text-dark": "#e2e8f0",
        "admin-background": "#ffffff",
        "admin-background-dark": "#1F1F1F",
        "admin-background-secondary": "#e2e8f0",
        "admin-background-secondary-dark": "#393939",
        "admin-submit": "#0ea5e9",
        "admin-submit-dark": "#7dd3fc",
        "admin-delete": "#dc2626",
        "admin-delete-dark": "#dc2626",
      },
    },
  },
  plugins: [],
};
