import { yupLocale } from '../../lang/yup';
import {
  isArraySorted,
  isTimestampArraySorted
} from '../../utils/isArraySorted';
import { areIntervalsApart } from '../../utils/areIntervalsApart';
import * as yup from 'yup';

import { AUTOOPS_MAX_MIN_COUNT } from '../../constants/autoops';
import {
  FEATURE_DESCRIPTION_MAX_LENGTH,
  FEATURE_ID_MAX_LENGTH,
  FEATURE_NAME_MAX_LENGTH,
  FEATURE_TAG_MIN_SIZE
} from '../../constants/feature';
import {
  VARIATION_DESCRIPTION_MAX_LENGTH,
  VARIATION_NAME_MAX_LENGTH,
  VARIATION_VALUE_MAX_LENGTH,
  VARIATION_NUMBER_VALUE_MAX_LENGTH
} from '../../constants/variation';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Feature } from '../../proto/feature/feature_pb';
import { Strategy } from '../../proto/feature/strategy_pb';
import { isJsonString } from '../../utils/validate';
import {
  ActionTypeMap,
  ProgressiveRolloutTemplateScheduleClause
} from '../../proto/autoops/clause_pb';
import { OpsType, OpsTypeMap } from '../../proto/autoops/auto_ops_rule_pb';
import { ProgressiveRollout } from '../../proto/autoops/progressive_rollout_pb';

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
          (value, context) => {
            if (
              context.from[1].value.variationType ==
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
          (value, context) => {
            if (
              context.from[1].value.variationType ==
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
            max: `${VARIATION_VALUE_MAX_LENGTH}`
          }),
          (value, context) => {
            const type = context.from[1].value.variationType;
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
            max: `${VARIATION_NUMBER_VALUE_MAX_LENGTH}`
          }),
          (value, context) => {
            if (
              context.from[1].value.variationType ==
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
          function (value, context) {
            const unChangedList = context.from[1].value.variations.filter(
              (val) => val.id != context.from[0].value.id
            );
            return !unChangedList.find((val) => val.value === value);
          }
        ),
      name: yup.string().required().max(VARIATION_NAME_MAX_LENGTH),
      description: yup.string().max(VARIATION_DESCRIPTION_MAX_LENGTH)
    })
    .required()
);

const schedulesListSchema = yup.array().of(
  yup.object().shape({
    weight: yup
      .number()
      .transform((value) => (isNaN(value) ? undefined : value))
      .required()
      .min(1)
      .max(100)
      .test('isAscending', '', (_, context) => {
        if (
          context.from[3].value.progressiveRolloutType ===
          ProgressiveRollout.Type.MANUAL_SCHEDULE
        ) {
          return isArraySorted(
            context.from[3].value.progressiveRollout.manual.schedulesList.map(
              (d) => Number(d.weight)
            )
          );
        }
        return true;
      }),
    executeAt: yup.object().shape({
      time: yup
        .date()
        .test(
          'isLaterThanNow',
          intl.formatMessage(messages.input.error.notLaterThanCurrentTime),
          (value, context) => {
            if (
              context.from[4].value.progressiveRolloutType ===
              ProgressiveRollout.Type.MANUAL_SCHEDULE
            ) {
              return value.getTime() > new Date().getTime();
            }
            return true;
          }
        )
        .test('isAscending', '', (_, context) => {
          if (
            context.from[4].value.progressiveRolloutType ===
            ProgressiveRollout.Type.MANUAL_SCHEDULE
          ) {
            return isTimestampArraySorted(
              context.from[4].value.progressiveRollout.manual.schedulesList.map(
                (d) => d.executeAt.time.getTime()
              )
            );
          }
          return true;
        })
        .test('timeIntervals', '', (_, context) => {
          if (
            context.from[4].value.progressiveRolloutType ===
            ProgressiveRollout.Type.MANUAL_SCHEDULE
          ) {
            return areIntervalsApart(
              context.from[4].value.progressiveRollout.manual.schedulesList.map(
                (d) => d.executeAt.time.getTime()
              ),
              5
            );
          }
          return true;
        })
    })
  })
);

export const operationFormSchema = yup.object().shape({
  opsType: yup.mixed<OpsTypeMap[keyof OpsTypeMap]>().required(),
  datetimeClausesList: yup.array().of(
    yup.object().shape({
      id: yup.string(),
      actionType: yup.mixed<ActionTypeMap[keyof ActionTypeMap]>().required(),
      time: yup.date()
    })
  ),
  eventRate: yup.object().shape({
    variation: yup.string(),
    goal: yup
      .string()
      .nullable()
      .test(
        'required',
        intl.formatMessage(messages.input.error.required),
        (value, context) => {
          if (context.from[1].value.opsType === OpsType.EVENT_RATE) {
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
    operator: yup.string()
  }),
  progressiveRolloutType: yup
    .mixed<ProgressiveRollout.TypeMap[keyof ProgressiveRollout.TypeMap]>()
    .required(),
  progressiveRollout: yup.object().shape({
    template: yup.object().shape({
      variationId: yup.string().required(),
      increments: yup
        .number()
        .transform((value) => (isNaN(value) ? undefined : value))
        .required()
        .min(1)
        .max(100),
      datetime: yup.object().shape({
        time: yup
          .date()
          .test(
            'isLaterThanNow',
            intl.formatMessage(messages.input.error.notLaterThanCurrentTime),
            (value, context) => {
              if (
                context.from[3].value.progressiveRolloutType ===
                ProgressiveRollout.Type.TEMPLATE_SCHEDULE
              ) {
                return value.getTime() > new Date().getTime();
              }
              return true;
            }
          )
      }),
      schedulesList: schedulesListSchema,
      interval: yup
        .mixed<
          ProgressiveRolloutTemplateScheduleClause.IntervalMap[keyof ProgressiveRolloutTemplateScheduleClause.IntervalMap]
        >()
        .required()
    }),
    manual: yup.object().shape({
      variationId: yup.string().required(),
      schedulesList: schedulesListSchema
    })
  })
});
export type OperationForm = yup.InferType<typeof operationFormSchema>;

const tagsSchema = yup.array().min(FEATURE_TAG_MIN_SIZE).of(yup.string());
const settingsTagsSchema = yup.array().of(yup.string());

export const switchEnabledFormSchema = (requireComment: boolean) => {
  return yup.object().shape({
    featureId: idSchema,
    enabled: yup.boolean().required(),
    comment: requireComment ? commentSchema : yup.string()
  });
};

export const archiveFormSchema = (requireComment: boolean) =>
  yup.object().shape({
    featureId: idSchema,
    comment: requireComment ? commentSchema : yup.string()
  });

export const cloneSchema = yup.object().shape({
  // Since some old environments have empty id, so we don't require it
  // destinationEnvironmentId: yup.string().required(),
});

export const onVariationSchema = yup.object().shape({
  id: yup.string(),
  value: yup.string(),
  label: yup.string()
});

export const offVariationSchema = yup.object().shape({
  id: yup.string(),
  value: yup.string(),
  label: yup.string()
});

export const addFormSchema = yup.object().shape({
  id: idSchema,
  name: nameSchema,
  description: descriptionSchema,
  tags: tagsSchema,
  variationType: variationTypeSchema,
  variations: variationsSchema,
  onVariation: onVariationSchema,
  offVariation: offVariationSchema
});
export type AddForm = yup.InferType<typeof addFormSchema>;

export const variationsFormSchema = yup.object().shape({
  variationType: variationTypeSchema,
  onVariation: onVariationSchema,
  offVariation: onVariationSchema,
  variations: variationsSchema,
  resetSampling: yup.bool(),
  requireComment: yup.bool(),
  comment: yup.string().when('requireComment', {
    is: (requireComment: boolean) => requireComment,
    then: (schema) => schema.required()
  })
});
export type VariationForm = yup.InferType<typeof variationsFormSchema>;

export const settingsFormSchema = (requireComment: boolean) =>
  yup.object().shape({
    name: nameSchema,
    description: descriptionSchema,
    tags: settingsTagsSchema,
    comment: requireComment ? commentSchema : yup.string()
  });

const strategySchema = yup.object().shape({
  option: yup.object().shape({
    value: yup.string(),
    label: yup.string()
  }),
  rolloutStrategy: yup
    .array()
    .of(
      yup.object().shape({
        id: yup.string(),
        percentage: yup.number()
      })
    )
    .required()
    .test(
      'sum',
      intl.formatMessage(messages.input.error.not100Percentage),
      (value, context) => {
        if (context.parent.option.value != Strategy.Type.ROLLOUT) {
          return true;
        }
        const total = value
          .map((v) => Number(v.percentage))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        return total == 100;
      }
    )
});
export type StrategySchema = yup.InferType<typeof strategySchema>;

export const ruleClauseType = {
  COMPARE: 'compare',
  SEGMENT: 'segment',
  DATE: 'date',
  FEATURE_FLAG: 'feature_flag'
} as const;
export type RuleClauseType =
  (typeof ruleClauseType)[keyof typeof ruleClauseType];

const ruleClauseSchema = yup.object().shape({
  id: yup.string(),
  type: yup.string(),
  attribute: yup.string().when('type', {
    is: (type: string) => type === ruleClauseType.SEGMENT,
    then: (schema) => schema,
    otherwise: (schema) => schema.required()
  }),
  operator: yup.string(),
  values: yup.array().of(yup.string()).min(1)
});
export type RuleClauseSchema = yup.InferType<typeof ruleClauseSchema>;

export const rulesSchema = yup.object().shape({
  id: yup.string(),
  clauses: yup.array().of(ruleClauseSchema),
  strategy: strategySchema
});
export type RuleSchema = yup.InferType<typeof rulesSchema>;

export const targetingFormSchema = yup.object().shape({
  prerequisites: yup.array().of(
    yup.object().shape({
      featureId: yup.string().required(),
      variationId: yup.string().required()
    })
  ),
  enabled: yup.bool(),
  targets: yup.array().of(
    yup.object().shape({
      variationId: yup.string().required(),
      users: yup.array().of(yup.string())
    })
  ),
  rules: yup.array().of(rulesSchema),
  defaultStrategy: strategySchema,
  offVariation: offVariationSchema,
  resetSampling: yup.bool(),
  requireComment: yup.bool(),
  comment: yup.string().when('requireComment', {
    is: (requireComment: boolean) => requireComment,
    then: (schema) => schema.required()
  })
});
export type TargetingForm = yup.InferType<typeof targetingFormSchema>;

export const triggerFormSchema = yup.object().shape({
  triggerType: yup.string().nullable().required(),
  action: yup.string().nullable().required(),
  description: yup.string()
});

export const shortcutFormSchema = yup.object().shape({
  name: yup.string().required()
});
