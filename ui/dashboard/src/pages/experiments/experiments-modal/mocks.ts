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
  { label: 'A', value: 'f0876bba-bf65-42a1-b9d8-64e7935534e1' },
  { label: 'B', value: 'cd80d7d1-01e5-475f-a66a-37b1f2703bb4' },
  { label: 'C', value: '60ad078a-ad77-49ea-8557-611b1ca11dff' }
];
const numberVariations = [
  { label: '1', value: '7b2ffcf7-4567-407c-90bb-eeeb2d17d3af' },
  { label: '2', value: 'b3440386-2480-418b-a163-ba54567091ee' },
  { label: '3', value: '422f0ac3-3632-4432-a937-63e7041ce0a4' }
];

const booleanVariations = [
  { label: 'True', value: '8ca38e63-048a-43e1-aa24-f383772f94b2' },
  { label: 'False', value: 'd27224a6-638a-43bb-b9a0-b0ecc49c6345' }
];

const jsonVariations = [
  {
    label: JSON.stringify({ name: 'yuichi-alessandro' }),
    value: '8d47d048-73a2-4b8b-bd95-4eba16592046'
  },
  {
    label: JSON.stringify({ name: 'kenta - kozuka' }),
    value: 'b32b8560-1987-4dda-b3fa-6a460cf64b03'
  }
];

export {
  flagOptions,
  stringVariations,
  numberVariations,
  booleanVariations,
  jsonVariations
};
