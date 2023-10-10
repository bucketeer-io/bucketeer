import { yupLocale } from '@/lang/yup';
import * as yup from 'yup';

import { ClauseType } from '../../components/FeatureAutoOpsRulesForm';
import { AUTOOPS_MAX_MIN_COUNT } from '../../constants/autoops';
import {
  FEATURE_DESCRIPTION_MAX_LENGTH,
  FEATURE_ID_MAX_LENGTH,
  FEATURE_NAME_MAX_LENGTH,
  FEATURE_TAG_MIN_SIZE,
} from '../../constants/feature';
import {
  VARIATION_DESCRIPTION_MAX_LENGTH,
  VARIATION_NAME_MAX_LENGTH,
  VARIATION_VALUE_MAX_LENGTH,
  VARIATION_NUMBER_VALUE_MAX_LENGTH,
} from '../../constants/variation';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { isJsonString } from '../../utils/validate';

yup.setLocale(yupLocale);

const regex = new RegExp('^[a-zA-Z0-9-]+$');
const idSchema = yup
  .string()
  .required()
  .matches(regex, intl.formatMessage(messages.input.error.invalidId))
  .max(FEATURE_ID_MAX_LENGTH);

const nameSchema = yup.string().max(FEATURE_NAME_MAX_LENGTH).required();
const descriptionSchema = yup.string().max(FEATURE_DESCRIPTION_MAX_LENGTH);
const commentSchema = yup.string().required();
const variationTypeSchema = yup.string();

const variationsSchema = yup.array().of(
  yup
    .object()
    .shape({
      id: yup.string(),
      value: yup
        .string()
        .required()
        .test(
          'isNumber',
          intl.formatMessage(messages.input.error.notNumber),
          function (value) {
            const { from } = this as any;
            if (
              from[1].value.variationType ==
              Feature.VariationType.NUMBER.toString()
            ) {
              return !isNaN(Number(value));
            }
            return true;
          }
        )
        .test(
          'isJson',
          intl.formatMessage(messages.input.error.notJson),
          function (value) {
            const { from } = this as any;
            if (
              from[1].value.variationType ==
              Feature.VariationType.JSON.toString()
            ) {
              return isJsonString(value);
            }
            return true;
          }
        )
        .test(
          'maxLength',
          intl.formatMessage(messages.input.error.maxLength, {
            max: `${VARIATION_VALUE_MAX_LENGTH}`,
          }),
          function (value) {
            const { from } = this as any;
            const type = from[1].value.variationType;
            if (
              type == Feature.VariationType.JSON.toString() ||
              type == Feature.VariationType.STRING.toString()
            ) {
              return value.length <= VARIATION_VALUE_MAX_LENGTH;
            }
            return true;
          }
        )
        .test(
          'maxLengthNumber',
          intl.formatMessage(messages.input.error.maxLength, {
            max: `${VARIATION_NUMBER_VALUE_MAX_LENGTH}`,
          }),
          function (value) {
            const { from } = this as any;
            if (
              from[1].value.variationType ==
              Feature.VariationType.NUMBER.toString()
            ) {
              return value.length <= VARIATION_NUMBER_VALUE_MAX_LENGTH;
            }
            return true;
          }
        )
        .test(
          'isUnique',
          intl.formatMessage(messages.input.error.mustBeUnique),
          function (value) {
            const { from } = this as any;
            const unChangedList = from[1].value.variations.filter(
              (val) => val.id != from[0].value.id
            );
            return !unChangedList.find((val) => val.value === value);
          }
        ),
      name: yup.string().required().max(VARIATION_NAME_MAX_LENGTH),
      description: yup.string().max(VARIATION_DESCRIPTION_MAX_LENGTH),
    })
    .required()
);

export const operationFormSchema = yup.object().shape({
  opsType: yup.string().required(),
  clauseType: yup.string().required(),
  datetime: yup.object().shape({
    time: yup
      .date()
      .test(
        'isLaterThanNow',
        intl.formatMessage(messages.input.error.notLaterThanCurrentTime),
        function (value) {
          const { from } = this as any;
          if (from[1].value.clauseType === ClauseType.DATETIME) {
            return value.getTime() > new Date().getTime();
          }
          return true;
        }
      ),
  }),
  eventRate: yup.object().shape({
    variation: yup.string(),
    goal: yup
      .string()
      .nullable()
      .test(
        'required',
        intl.formatMessage(messages.input.error.required),
        function (value) {
          const { from } = this as any;
          if (from[1].value.clauseType == ClauseType.EVENT_RATE) {
            return value != null;
          }
          return true;
        }
      ),
    minCount: yup
      .number()
      .transform((value) => (isNaN(value) ? undefined : value))
      .required()
      .min(1)
      .max(AUTOOPS_MAX_MIN_COUNT),
    threadsholdRate: yup
      .number()
      .transform((value) => (isNaN(value) ? undefined : value))
      .required()
      .moreThan(0)
      .max(100),
    operator: yup.string(),
  }),
});

const tagsSchema = yup.array().min(FEATURE_TAG_MIN_SIZE).of(yup.string());

export const switchEnabledFormSchema = yup.object().shape({
  featureId: idSchema,
  enabled: yup.boolean().required(),
  comment: commentSchema,
});

export const archiveFormSchema = yup.object().shape({
  featureId: idSchema,
  comment: commentSchema,
});

export const cloneSchema = yup.object().shape({
  // Since some old environments have empty id, so we don't require it
  // destinationEnvironmentId: yup.string().required(),
});

export const onVariationSchema = yup.object().shape({
  id: yup.string(),
  value: yup.string(),
  label: yup.string(),
});

export const offVariationSchema = yup.object().shape({
  id: yup.string(),
  value: yup.string(),
  label: yup.string(),
});

export const addFormSchema = yup.object().shape({
  id: idSchema,
  name: nameSchema,
  description: descriptionSchema,
  tags: tagsSchema,
  variationType: variationTypeSchema,
  variations: variationsSchema,
  onVariation: onVariationSchema,
  offVariation: offVariationSchema,
});

export const variationsFormSchema = yup.object().shape({
  onVariation: onVariationSchema,
  variations: variationsSchema,
  resetSampling: yup.bool(),
  comment: commentSchema,
});

export const settingsFormSchema = yup.object().shape({
  name: nameSchema,
  description: descriptionSchema,
  tags: tagsSchema,
  comment: commentSchema,
});

export const strategySchema = yup.object().shape({
  option: yup.object().shape({
    value: yup.string(),
    label: yup.string(),
  }),
  rolloutStrategy: yup
    .array()
    .of(
      yup.object().shape({
        id: yup.string(),
        percentage: yup.number(),
      })
    )
    .required()
    .test(
      'sum',
      intl.formatMessage(messages.input.error.not100Percentage),
      (variations, context) => {
        if (context.parent.option.value != Strategy.Type.ROLLOUT) {
          return true;
        }
        const total = variations
          .map((v: any) => Number(v.percentage))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        return total == 100;
      }
    ),
});

export const targetingFormSchema = yup.object().shape({
  prerequisites: yup.array().of(
    yup.object().shape({
      featureId: yup.string().required(),
      variationId: yup.string().required(),
    })
  ),
  enabled: yup.bool(),
  targets: yup.array().of(
    yup.object().shape({
      variationId: yup.string().required(),
      users: yup.array().of(yup.string()),
    })
  ),
  rules: yup.array().of(
    yup.object().shape({
      id: yup.string(),
      clauses: yup.array().of(
        yup.object().shape({
          id: yup.string(),
          type: yup.string(),
          attribute: yup
            .string()
            .test(
              'required',
              intl.formatMessage(messages.input.error.required),
              (value, context) => {
                if (context.parent.type === 'segment') {
                  return true;
                }
                return !!value;
              }
            ),
          operator: yup.string(),
          values: yup.array().of(yup.string()).min(1),
        })
      ),
      strategy: strategySchema,
    })
  ),
  defaultStrategy: strategySchema,
  offVariation: yup.object().shape({
    id: yup.string(),
    value: yup.string(),
    label: yup.string(),
  }),
  resetSampling: yup.bool(),
  comment: yup.string().required(),
});
