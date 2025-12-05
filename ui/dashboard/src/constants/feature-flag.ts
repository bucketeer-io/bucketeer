export const FEATURE_ID_MAX_LENGTH = 100;
export const FEATURE_NAME_MAX_LENGTH = 100;
export const FEATURE_TAG_MIN_SIZE = 1;
export const FEATURE_VARIATION_MIN_SIZE = 100;
export const FEATURE_DESCRIPTION_MAX_LENGTH = 100;
export const FEATURE_UPDATE_COMMENT_MAX_LENGTH = 100;

export const VARIATION_VALUE_MAX_LENGTH = 10000;
export const VARIATION_NUMBER_VALUE_MAX_LENGTH = 309;
export const VARIATION_NAME_MAX_LENGTH = 100;
export const VARIATION_DESCRIPTION_MAX_LENGTH = 100;

export const getDefaultYamlValue = (index: number) =>
  `
key:
  variation_${index + 1}: value_1
  variation_${index + 2}: value_2
`.trim() + '\n';
