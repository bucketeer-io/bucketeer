import { i18n } from 'i18n';
import * as yup from 'yup';
import { StrategyType } from '@types';
import { RuleClauseType } from './types';

export type RuleClauseSchema = yup.InferType<typeof ruleClauseSchema>;
const translation = i18n.t;
const requiredMessage = translation('message:required-field');

const ruleClauseSchema = yup.object().shape({
  id: yup.string().required(),
  type: yup
    .string()
    .oneOf([
      RuleClauseType.COMPARE,
      RuleClauseType.SEGMENT,
      RuleClauseType.DATE,
      RuleClauseType.FEATURE_FLAG
    ])
    .required(requiredMessage),
  attribute: yup.string().when('type', {
    is: (type: string) => type === RuleClauseType.SEGMENT,
    then: schema => schema,
    otherwise: schema => schema.required(requiredMessage)
  }),
  operator: yup.string().required(requiredMessage),
  values: yup
    .array()
    .of(yup.string())
    .min(1, requiredMessage)
    .required(requiredMessage)
});

const strategySchema = yup.object().shape({
  currentOption: yup.string(),
  type: yup
    .string()
    .oneOf([StrategyType.FIXED, StrategyType.ROLLOUT, StrategyType.MANUAL])
    .required(requiredMessage),
  fixedStrategy: yup
    .object()
    .shape({
      variation: yup.string()
    })
    .when('type', {
      is: (type: string) => type !== StrategyType.FIXED,
      then: schema => schema,
      otherwise: schema => schema.required(requiredMessage)
    }),
  rolloutStrategy: yup
    .object()
    .shape({
      audience: yup.object().shape({
        percentage: yup
          .number()
          .transform(value => (isNaN(value) ? undefined : value))
          .min(0)
          .max(100, translation('message:validation.percentage-less-than-100')),
        defaultVariation: yup.string().when('percentage', {
          is: (percentage: number) => percentage > 0 && percentage !== 100,
          then: schema => schema.required(requiredMessage),
          otherwise: schema => schema
        })
      }),
      variations: yup.array().of(
        yup.object().shape({
          variation: yup.string().required(requiredMessage),
          weight: yup.number().required(requiredMessage)
        })
      )
    })
    .test('sum', function (value) {
      const { type } = this.parent || {};
      if (type !== StrategyType.ROLLOUT) {
        return true;
      }
      if (value.variations) {
        const total = value.variations
          .map(v => Number(v.weight))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        if (total !== 100)
          return this.createError({
            message: translation('message:validation.should-be-percent'),
            path: `${this.path}.variations`
          });
      }
      return true;
    })
});
export type StrategySchema = yup.InferType<typeof strategySchema>;

export const rulesSchema = yup.object().shape({
  id: yup.string().required(),
  clauses: yup.array().of(ruleClauseSchema).required(),
  strategy: strategySchema
});

export type RuleSchema = yup.InferType<typeof rulesSchema>;

const defaultAudienceRuleSchema = yup.object().shape({
  rule: yup.string().required(requiredMessage),
  inExperiment: yup.number().required(requiredMessage),
  notInExperiment: yup.number().required(requiredMessage),
  served: yup.boolean().required(requiredMessage),
  variationReassignment: yup.boolean().required(requiredMessage)
});

export type DefaultAudienceRuleSchema = yup.InferType<
  typeof defaultAudienceRuleSchema
>;

export const defaultRuleSchema = yup.object().shape({
  audienceRules: yup.array().of(defaultAudienceRuleSchema),
  currentOption: yup.string(),
  type: yup
    .string()
    .oneOf([StrategyType.FIXED, StrategyType.ROLLOUT, StrategyType.MANUAL])
    .required(requiredMessage),
  fixedStrategy: yup
    .object()
    .shape({
      variation: yup.string()
    })
    .when('type', {
      is: (type: StrategyType) => type !== StrategyType.FIXED,
      then: schema => schema,
      otherwise: schema => schema.required(requiredMessage)
    }),
  rolloutStrategy: yup
    .object()
    .shape({
      audience: yup.object().shape({
        percentage: yup
          .number()
          .transform(value => (isNaN(value) ? undefined : value))
          .min(0)
          .max(100, translation('message:validation.percentage-less-than-100')),
        defaultVariation: yup.string().when('percentage', {
          is: (percentage: number) => percentage > 0 && percentage !== 100,
          then: schema => schema.required(requiredMessage),
          otherwise: schema => schema
        })
      }),
      variations: yup.array().of(
        yup.object().shape({
          variation: yup.string().required(requiredMessage),
          weight: yup.number().required(requiredMessage)
        })
      )
    })
    .test('sum', function (value) {
      const { type } = this.parent || {};
      if (type !== StrategyType.MANUAL) {
        return true;
      }
      if (value.variations) {
        const total = value.variations
          .map(v => Number(v.weight))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        if (total !== 100)
          return this.createError({
            message: translation('message:validation.should-be-percent'),
            path: `${this.path}.variations`
          });
      }
      return true;
    })
});

export type DefaultRuleSchema = yup.InferType<typeof defaultRuleSchema>;

export const formSchema = yup.object().shape({
  prerequisites: yup.array().of(
    yup.object().shape({
      featureId: yup.string().required(requiredMessage),
      variationId: yup.string().required(requiredMessage)
    })
  ),
  individualRules: yup.array().of(
    yup.object().shape({
      variationId: yup.string().required(requiredMessage),
      name: yup.string(),
      users: yup.array().required(requiredMessage)
    })
  ),
  segmentRules: yup.array().of(rulesSchema),
  defaultRule: defaultRuleSchema,
  enabled: yup.boolean().required(requiredMessage),
  requireComment: yup.boolean(),
  resetSampling: yup.boolean(),
  comment: yup.string().when('requireComment', {
    is: (requireComment: boolean) => requireComment,
    then: schema => schema.required(requiredMessage),
    otherwise: schema => schema
  }),
  scheduleType: yup.string().oneOf(['ENABLE', 'DISABLE', 'SCHEDULE']),
  scheduleAt: yup.string().test('test', function (value, context) {
    const scheduleType = context.from && context.from[0].value.scheduleType;
    if (scheduleType === 'SCHEDULE') {
      if (!value)
        return context.createError({
          message: requiredMessage,
          path: context.path
        });
      if (+value * 1000 < new Date().getTime())
        return context.createError({
          message: translation(
            'message:validation.operation.later-than-current-time'
          ),
          path: context.path
        });
    }
    return true;
  }),
  offVariation: yup.string()
});

export type TargetingSchema = yup.InferType<typeof formSchema>;
