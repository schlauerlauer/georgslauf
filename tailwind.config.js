/** @type {import('tailwindcss').Config} */

const { addDynamicIconSelectors } = require('@iconify/tailwind');

module.exports = {
	content: [
		"internal/handler/templates/*.templ",
	],
	theme: {
		extend: {
			fontFamily: {
				sans: ['-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'Roboto', 'Oxygen', 'Ubuntu', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', 'sans-serif']
			},
		},
	},
	plugins: [
		addDynamicIconSelectors(),
		require('@tailwindcss/forms')({
			strategy: "class",
		}),
		require("daisyui"),
	],
	daisyui: {
		themes: ["pastel", "cupcake", "dracula", "sunset"], // false: only light + dark | true: all themes | array: specific themes like this ["light", "dark", "cupcake"]
		darkTheme: "dracula", // name of one of the included themes for dark mode
		base: true, // applies background color and foreground color for root element by default
		styled: true, // include daisyUI colors and design decisions for all components
		utils: true, // adds responsive and modifier utility classes
		prefix: "", // prefix for daisyUI classnames (components, modifiers and responsive class names. Not colors)
		logs: true, // Shows info about daisyUI version and used config in the console when building your CSS
		themeRoot: ":root", // The element that receives theme color CSS variables
	},
}
