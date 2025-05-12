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
    .array()
    .of(
      yup.object().shape({
        variation: yup.string().required(requiredMessage),
        weight: yup.number().required(requiredMessage)
      })
    )
    .when('type', {
      is: (type: string) => type !== StrategyType.ROLLOUT,
      then: schema => schema,
      otherwise: schema => schema.required(requiredMessage)
    })
    .test('sum', (value, context) => {
      if (context.parent?.type !== StrategyType.ROLLOUT) {
        return true;
      }
      if (value) {
        const total = value
          .map(v => Number(v.weight))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        if (total !== 100)
          return context.createError({
            message: translation('message:validation.should-be-percent'),
            path: context.path
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

export const defaultRuleSchema = yup.object().shape({
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
  manualStrategy: yup
    .array()
    .of(
      yup.object().shape({
        variation: yup.string().required(requiredMessage),
        weight: yup.number().required(requiredMessage)
      })
    )
    .when('type', {
      is: (type: string) => type !== StrategyType.MANUAL,
      then: schema => schema,
      otherwise: schema => schema.required(requiredMessage)
    })
    .test('sum', (value, context) => {
      if (context.parent?.type !== StrategyType.MANUAL) {
        return true;
      }
      if (value) {
        const total = value
          .map(v => Number(v.weight))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        if (total !== 100)
          return context.createError({
            message: translation('message:validation.should-be-percent'),
            path: context.path
          });
      }
      return true;
    }),
  rolloutStrategy: yup
    .array()
    .of(
      yup.object().shape({
        variation: yup.string().required(requiredMessage),
        weight: yup.number().required(requiredMessage)
      })
    )
    .when('type', {
      is: (type: string) => type !== StrategyType.ROLLOUT,
      then: schema => schema,
      otherwise: schema => schema.required(requiredMessage)
    })
    .test('sum', (value, context) => {
      if (context.parent?.type !== StrategyType.ROLLOUT) {
        return true;
      }
      if (value) {
        const total = value
          .map(v => Number(v.weight))
          .reduce((total, current) => {
            return total + (current || 0);
          }, 0);
        if (total !== 100)
          return context.createError({
            message: translation('message:validation.should-be-percent'),
            path: context.path
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
      users: yup
        .array()
        .required(requiredMessage)
        .test('required', (value, context) => {
          if ((Array.isArray(value) && !value.length) || !value) {
            return context.createError({
              message: requiredMessage,
              path: context.path
            });
          }

          return true;
        })
    })
  ),
  segmentRules: yup.array().of(rulesSchema),
  defaultRule: defaultRuleSchema,
  enabled: yup.boolean().required(requiredMessage),
  isShowRules: yup.boolean().required(requiredMessage)
});

export type TargetingSchema = yup.InferType<typeof formSchema>;
