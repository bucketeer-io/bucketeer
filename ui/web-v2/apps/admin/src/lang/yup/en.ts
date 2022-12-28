export const localEn = {
  mixed: {
    default: 'Input Error',
    required: 'Required field',
    oneOf: ({ values }) => `Must be one of the following values: ${values}`,
    notOneOf: ({ values }) =>
      `Must not be one of the following values: ${values}`,
    isNumber: 'Different format',
  },
  string: {
    required: 'Required field',
    length: ({ length }) => `Please enter ${length} characters`,
    min: ({ min }) => `Please enter at least ${min} characters`,
    max: ({ max }) => `Please enter within ${max} characters`,
    matches: 'Different format',
    email: 'Different format',
    url: 'Different format',
    trim: 'Do not put spaces before or after',
    lowercase: 'Must be lowercase',
    uppercase: 'Must be uppercase',
  },
  number: {
    min: ({ min }) => `Please enter a value greater than or equal to ${min}`,
    max: ({ max }) => `Please enter a value less than ${max}`,
    lessThan: ({ less }) => `Please enter a value less than ${less}`,
    moreThan: ({ more }) => `Please enter a value greater than ${more}`,
    notEqual: ({ notEqual }) =>
      `Please enter a value different from ${notEqual}`,
    positive: 'Please enter a positive number',
    negative: 'Please enter a negative number',
    integer: 'Please enter an integer',
  },
  date: {
    default: 'Different format',
    min: ({ min }) => `Please enter a date greater than or equal to ${min}`,
    max: ({ max }) => `Please enter a date less than or equal to ${max}`,
  },
  object: {
    noUnknown: 'Please enter data with a valid key',
  },
  array: {
    min: ({ min }) => `Please enter at least ${min} values`,
    max: ({ max }) => `Please enter no more than ${max} values`,
  },
};
