/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    darkMode: 'class',
    theme: {
        extend: {
            colors: {
                primary: {
                    DEFAULT: '#1f3b73',
                    light: '#2c5aa0',
                    dark: '#162a52',
                },
                secondary: '#2c5aa0',
                critical: '#ef4444', // Red-500
                high: '#f97316',     // Orange-500
                medium: '#eab308',   // Yellow-500
                low: '#3b82f6',      // Blue-500
                info: '#22c55e',     // Green-500
                background: '#0f172a', // Slate-900
                surface: '#1e293b',    // Slate-800
                border: '#334155',     // Slate-700
            },
            fontFamily: {
                sans: ['Inter', 'sans-serif'],
                mono: ['JetBrains Mono', 'monospace'],
            }
        },
    },
    plugins: [],
}
