import {
  FEATURE_DESCRIPTION_MAX_LENGTH,
  FEATURE_NAME_MAX_LENGTH,
  VARIATION_DESCRIPTION_MAX_LENGTH,
  VARIATION_NAME_MAX_LENGTH,
  VARIATION_NUMBER_VALUE_MAX_LENGTH,
  VARIATION_VALUE_MAX_LENGTH
} from 'constants/feature-flag';
import { FormSchemaProps } from 'hooks/use-form-schema';
import * as yup from 'yup';
import { FeatureVariation, FeatureVariationType } from '@types';
import { isNumber } from 'utils/chart';
import { isJsonString } from 'utils/converts';
import { FlagSwitchVariationType } from './types';

const nameSchema = ({ requiredMessage }: { requiredMessage: string }) =>
  yup.string().max(FEATURE_NAME_MAX_LENGTH).required(requiredMessage);
const descriptionSchema = yup.string().max(FEATURE_DESCRIPTION_MAX_LENGTH);

export interface VariationSchema {
  id: string;
  value: string;
  name: string;
  description?: string;
}

export const createVariationsSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) =>
  yup
    .array()
    .required()
    .of(
      yup
        .object()
        .shape({
          id: yup.string().required(requiredMessage),
          value: yup
            .string()
            .required(requiredMessage)
            .test('isNumber', (value, context) => {
              if (
                context?.from &&
                context.from[1].value.variationType === 'NUMBER' &&
                (!isNumber(+value) || value.startsWith(' '))
              ) {
                return context.createError({
                  message: translation('message:validation.must-be-number'),
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
                  message: translation('message:validation.must-be-json'),
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
                  message: translation('message:validation.max-length-string', {
                    count: VARIATION_VALUE_MAX_LENGTH
                  }),
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
                  message: translation('message:validation.max-length-number', {
                    count: VARIATION_NUMBER_VALUE_MAX_LENGTH
                  }),
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
                variations?.filter(
                  item => item.value === currentVariation?.value
                ).length > 1
              ) {
                return context.createError({
                  message: translation('message:validation.must-be-unique'),
                  path: context.path
                });
              }
              return true;
            }),
          name: yup
            .string()
            .required(requiredMessage)
            .max(VARIATION_NAME_MAX_LENGTH),
          description: yup.string().max(VARIATION_DESCRIPTION_MAX_LENGTH)
        })
        .required()
    );
export interface FlagFormSchema {
  name: string;
  flagId: string;
  description?: string;
  tags: string[];
  switchVariationType?: FlagSwitchVariationType;
  variationType: FeatureVariationType;
  variations: VariationSchema[];
  defaultOnVariation: string;
  defaultOffVariation: string;
}

export const createFlagFormSchema = ({
  requiredMessage,
  translation
}: FormSchemaProps) =>
  yup.object().shape({
    name: nameSchema({ requiredMessage }),
    flagId: yup
      .string()
      .required(requiredMessage)
      .matches(
        /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
        translation('message:validation.id-rule', {
          name: translation('common:source-type.feature-flag')
        })
      ),
    description: descriptionSchema,
    tags: yup.array().min(1).required(requiredMessage),
    switchVariationType: yup.mixed<FlagSwitchVariationType>(),
    variationType: yup.mixed<FeatureVariationType>().required(requiredMessage),
    variations: createVariationsSchema({ requiredMessage, translation }),
    defaultOnVariation: yup.string().required(requiredMessage),
    defaultOffVariation: yup.string().required(requiredMessage)
  });
