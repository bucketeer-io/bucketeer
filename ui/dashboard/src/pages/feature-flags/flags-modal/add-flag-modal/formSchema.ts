import {
  FEATURE_DESCRIPTION_MAX_LENGTH,
  FEATURE_NAME_MAX_LENGTH,
  VARIATION_DESCRIPTION_MAX_LENGTH,
  VARIATION_NAME_MAX_LENGTH,
  VARIATION_NUMBER_VALUE_MAX_LENGTH,
  VARIATION_VALUE_MAX_LENGTH
} from 'constants/feature-flag';
import * as yup from 'yup';
import { isJsonString } from 'utils/converts';
import { FlagDataType } from 'pages/feature-flags/types';
import { VariationType } from './variations';

const nameSchema = yup.string().max(FEATURE_NAME_MAX_LENGTH).required();
const descriptionSchema = yup.string().max(FEATURE_DESCRIPTION_MAX_LENGTH);
const variationsSchema = yup.array().of(
  yup
    .object()
    .shape({
      id: yup.string().required(),
      value: yup
        .string()
        .required()
        .test('isNumber', (value, context) => {
          if (
            context?.from &&
            context.from[1].value.flagType === 'number' &&
            isNaN(Number(value))
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
            context.from[1].value.flagType === 'json' &&
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
          const type = context.from && context.from[1].value.flagType;
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
            context.from[1].value.flagType === 'number' &&
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
          const variations: VariationType[] =
            context.from && context.from[1].value.variations;
          const currentVariation: VariationType =
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
      name: yup.string().max(VARIATION_NAME_MAX_LENGTH),
      description: yup.string().max(VARIATION_DESCRIPTION_MAX_LENGTH)
    })
    .required()
);

export const formSchema = yup.object().shape({
  name: nameSchema,
  flagId: yup.string().required(),
  description: descriptionSchema,
  tags: yup.array().min(1).required(),
  flagType: yup.mixed<FlagDataType>().required(),
  variations: variationsSchema,
  serveOn: yup.object().shape({
    id: yup.string(),
    value: yup.string()
  }),
  serveOff: yup.object().shape({
    id: yup.string().required(),
    value: yup.string().required()
  })
});
