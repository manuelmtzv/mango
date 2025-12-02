/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,ts,tsx}"],
  theme: {
    fontFamily: {
      sans: ["Nunito Sans", "sans-serif"],
    },
    extend: {
      colors: {
        background: {
          default: "#F9F5F6",
        },
        header: {
          default: "#FFF2CC",
        },
        entry: {
          default: "#FFFFFF",
        },
        navigation: {
          default: "#FFD966",
          hover: "#ffe492",
        },
        delete: {
          default: "#eb7d7d",
          hover: "#ee9090",
        },
        edit: {
          default: "#A4D0A4",
          hover: "#b6d9b6",
        },
        tech: {
          vue: {
            default: "#337357",
          },
          pinia: {
            default: "#FFC73B",
          },
          typescript: {
            default: "#007ACC",
          },
        },
      },
    },
  },
  plugins: [],
};
