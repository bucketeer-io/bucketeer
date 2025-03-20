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
    otherwise: schema => schema.required()
  }),
  operator: yup.string().required('This field is required.'),
  values: yup
    .array()
    .of(yup.string())
    .min(1)
    .required('This field is required.')
});

export type RuleSchema = yup.InferType<typeof rulesSchema>;

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
    .test('sum', (value, context) => {
      if (context.parent.option.value != StrategyType.ROLLOUT) {
        return true;
      }
      const total = value
        .map(v => Number(v.percentage))
        .reduce((total, current) => {
          return total + (current || 0);
        }, 0);
      if (total == 100)
        return context.createError({
          message: `Total should be 100%.`,
          path: context.path
        });

      return true;
    })
});

export const rulesSchema = yup.object().shape({
  id: yup.string(),
  clauses: yup.array().of(ruleClauseSchema).required(),
  strategy: strategySchema
});

export const formSchema = yup.object().shape({
  prerequisitesRules: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        index: yup.number().required(),
        rules: yup
          .array()
          .required()
          .of(
            yup.object().shape({
              featureFlag: yup.string().required('This field is required.'),
              variation: yup.string().required('This field is required.')
            })
          )
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
  rules: yup.array().of(rulesSchema)
});
