import * as yup from 'yup';
import { StrategyType } from '@types';

export enum RuleClauseType {
  COMPARE = 'compare',
  SEGMENT = 'segment',
  DATE = 'date',
  FEATURE_FLAG = 'feature-flag'
}

export type RuleClauseSchema = yup.InferType<typeof ruleClauseSchema>;

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
    .required('This field is required.'),
  attribute: yup.string().when('type', {
    is: (type: string) => type === RuleClauseType.SEGMENT,
    then: schema => schema,
    otherwise: schema => schema.required('This field is required.')
  }),
  operator: yup.string().required('This field is required.'),
  values: yup
    .array()
    .of(yup.string())
    .min(1, 'This field is required.')
    .required('This field is required.')
});

const strategySchema = yup.object().shape({
  currentOption: yup.string(),
  type: yup
    .string()
    .oneOf([StrategyType.FIXED, StrategyType.ROLLOUT])
    .required('This field is required.'),
  fixedStrategy: yup
    .object()
    .shape({
      variation: yup.string()
    })
    .when('type', {
      is: (type: string) => type === StrategyType.ROLLOUT,
      then: schema => schema,
      otherwise: schema => schema.required('This field is required.')
    }),
  rolloutStrategy: yup
    .array()
    .of(
      yup.object().shape({
        variation: yup.string().required('This field is required.'),
        weight: yup.number().required('This field is required.')
      })
    )
    .when('type', {
      is: (type: string) => type === StrategyType.FIXED,
      then: schema => schema,
      otherwise: schema => schema.required('This field is required.')
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
            message: `Total should be 100%.`,
            path: context.path
          });
      }
      return true;
    })
});
export type StrategySchema = yup.InferType<typeof strategySchema>;

export const rulesSchema = yup.object().shape({
  id: yup.string(),
  clauses: yup.array().of(ruleClauseSchema).required(),
  strategy: strategySchema
});

export type RuleSchema = yup.InferType<typeof rulesSchema>;

export const formSchema = yup.object().shape({
  prerequisites: yup.array().of(
    yup.object().shape({
      featureId: yup.string().required('This field is required.'),
      variationId: yup.string().required('This field is required.')
    })
  ),
  targetIndividualRules: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        variationId: yup.string().required('This field is required.'),
        name: yup.string(),
        users: yup
          .array()
          .required('This field is required.')
          .test('required', (value, context) => {
            if ((Array.isArray(value) && !value.length) || !value) {
              return context.createError({
                message: `This field is required.`,
                path: context.path
              });
            }

            return true;
          })
      })
    ),
  rules: yup.array().of(rulesSchema),
  defaultStrategy: strategySchema
});
