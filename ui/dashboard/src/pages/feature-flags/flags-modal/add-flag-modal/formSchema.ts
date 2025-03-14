import {
  FEATURE_DESCRIPTION_MAX_LENGTH,
  FEATURE_NAME_MAX_LENGTH,
  VARIATION_DESCRIPTION_MAX_LENGTH,
  VARIATION_NAME_MAX_LENGTH,
  VARIATION_NUMBER_VALUE_MAX_LENGTH,
  VARIATION_VALUE_MAX_LENGTH
} from 'constants/feature-flag';
import * as yup from 'yup';
import { FeatureVariation, FeatureVariationType } from '@types';
import { isNumber } from 'utils/chart';
import { isJsonString } from 'utils/converts';

const nameSchema = yup.string().max(FEATURE_NAME_MAX_LENGTH).required();
const descriptionSchema = yup.string().max(FEATURE_DESCRIPTION_MAX_LENGTH);
const variationsSchema = yup
  .array()
  .required()
  .of(
    yup
      .object()
      .shape({
        id: yup.string().required('This field is required'),
        value: yup
          .string()
          .required('This field is required')
          .test('isNumber', (value, context) => {
            if (
              context?.from &&
              context.from[1].value.variationType === 'NUMBER' &&
              !isNumber(+value)
            ) {
              return context.createError({
                message: 'This must be a number.',
                path: context.path
              });
            }
            return true;
          })
          .test('isJson', (value, context) => {
            if (
              context?.from &&
              context.from[1].value.variationType === 'JSON' &&
              !isJsonString(value)
            ) {
              return context.createError({
                message: 'This must be a JSON.',
                path: context.path
              });
            }
            return true;
          })
          .test('maxLength', (value, context) => {
            const type = context.from && context.from[1].value.variationType;
            if (
              ['string', 'json'].includes(type) &&
              value.length >= VARIATION_VALUE_MAX_LENGTH
            ) {
              return context.createError({
                message: `The maximum length for this field is ${VARIATION_VALUE_MAX_LENGTH} characters.`,
                path: context.path
              });
            }
            return true;
          })
          .test('maxLengthNumber', (value, context) => {
            if (
              context.from &&
              context.from[1].value.variationType === 'NUMBER' &&
              value.length >= VARIATION_NUMBER_VALUE_MAX_LENGTH
            ) {
              return context.createError({
                message: `The maximum length for this field is ${VARIATION_NUMBER_VALUE_MAX_LENGTH} numbers.`,
                path: context.path
              });
            }
            return true;
          })
          .test('isUnique', function (_, context) {
            const variations: FeatureVariation[] =
              context.from && context.from[1].value.variations;
            const currentVariation: FeatureVariation =
              context.from && context.from[0].value;
            if (
              variations?.filter(item => item.value === currentVariation?.value)
                .length > 1
            ) {
              return context.createError({
                message: `This must be unique.`,
                path: context.path
              });
            }
            return true;
          }),
        name: yup
          .string()
          .required('This field is required')
          .max(VARIATION_NAME_MAX_LENGTH),
        description: yup.string().max(VARIATION_DESCRIPTION_MAX_LENGTH)
      })
      .required()
  );

export const formSchema = yup.object().shape({
  name: nameSchema,
  flagId: yup
    .string()
    .required()
    .matches(
      /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
      "urlCode must start with a letter or number and only contain letters, numbers, or '-'"
    ),
  description: descriptionSchema,
  tags: yup.array().min(1).required(),
  variationType: yup.mixed<FeatureVariationType>().required(),
  variations: variationsSchema,
  defaultOnVariation: yup.string().required(),
  defaultOffVariation: yup.string().required()
});
