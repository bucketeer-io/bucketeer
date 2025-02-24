const flagOptions = [
  {
    label: 'feature-flag-test-boolean-name',
    value: 'feature-flag-test-boolean'
  },
  {
    label: 'feature-flag-test-string-name',
    value: 'feature-flag-test-string'
  },
  {
    label: 'feature-flag-test-number-name',
    value: 'feature-flag-test-number'
  },
  {
    label: 'feature-flag-test-json-name',
    value: 'feature-flag-test-json'
  }
];

const stringVariations = [
  { label: 'A', value: 'a' },
  { label: 'B', value: 'b' },
  { label: 'C', value: 'c' }
];
const numberVariations = [
  { label: '1', value: 1 },
  { label: '2', value: 2 },
  { label: '3', value: 3 }
];

const booleanVariations = [
  { label: 'True', value: 'true' },
  { label: 'False', value: 'false' }
];

const jsonVariations = [
  {
    label: JSON.stringify({ name: 'yuichi-alessandro' }),
    value: JSON.stringify({ name: 'yuichi-alessandro' })
  },
  {
    label: JSON.stringify({ name: 'kenta - kozuka' }),
    value: JSON.stringify({ name: 'kenta-kozuka' })
  }
];

export {
  flagOptions,
  stringVariations,
  numberVariations,
  booleanVariations,
  jsonVariations
};
