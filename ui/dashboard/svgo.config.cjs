module.exports = {
  js2svg: {
    indent: 2,
    pretty: true
  },
  plugins: [
    {
      name: 'preset-default',
      params: {
        removeViewBox: false
      }
    },
    {
      name: 'removeDimensions',
      active: true
    },
    {
      name: 'removeAttrs',
      params: {
        attrs: ['fill-rule', 'fill']
      }
    },
    {
      name: 'addAttributesToSVGElement',
      params: {
        attribute: {
          width: '1em',
          height: '1em',
          fill: 'currentColor'
        }
      }
    }
  ]
};
