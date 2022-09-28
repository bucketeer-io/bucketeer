module.exports = {
  mode: 'jit',
  purge: [
    './apps/admin/src/**/*.html',
    './apps/admin/src/**/*.ts',
    './apps/admin/src/**/*.tsx',
  ],
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
  plugins: [require('@tailwindcss/forms')],
};
