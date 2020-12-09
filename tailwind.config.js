module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        primary: {
          200: 'rgba(calc(var(--color-primary-r) + 150), calc(var(--color-primary-g) + 150), calc(var(--color-primary-b) + 150))',
          300: 'rgba(calc(var(--color-primary-r) + 96), calc(var(--color-primary-g) + 96), calc(var(--color-primary-b) + 96))',
          400: 'rgba(calc(var(--color-primary-r) + 64), calc(var(--color-primary-g) + 64), calc(var(--color-primary-b) + 64))',
          500: 'rgba(var(--color-primary-r), var(--color-primary-g), var(--color-primary-b))',
          600: 'rgba(calc(var(--color-primary-r) * 0.85), calc(var(--color-primary-g) * 0.85), calc(var(--color-primary-b) * 0.85))',
          700: 'rgba(calc(var(--color-primary-r) * 0.65), calc(var(--color-primary-g) * 0.65), calc(var(--color-primary-b) * 0.65))',
          800: 'rgba(calc(var(--color-primary-r) * 0.40), calc(var(--color-primary-g) * 0.40), calc(var(--color-primary-b) * 0.40))',
        }
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
