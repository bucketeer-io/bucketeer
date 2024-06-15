module.exports = {
  mode: 'jit',
  purge: ['./src/**/*.html', './src/**/*.ts', './src/**/*.tsx'],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        primary: '#5d3597',
      },
    },
  },
  variants: {
    extend: {
      backgroundColor: ['checked'],
      borderColor: ['checked'],
      inset: ['checked'],
      zIndex: ['hover', 'active'],
      opacity: ['disabled'],
    },
  },
  plugins: [require('@tailwindcss/forms'), require('@tailwindcss/line-clamp')],
};
