import { useCallback, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import Form from 'components/form';
import AddRuleDropdown from './add-rule-dropdown';
import {
  initialIndividualRule,
  initialPrerequisitesRule,
  initialSegmentCondition
} from './constants';
import DefaultRule from './default-rule';
import TargetSegmentRule from './target-segment-rule';
import TargetingState from './targeting-state';
import {
  RuleCategory,
  TargetIndividualItem,
  TargetPrerequisiteItem,
  TargetSegmentItem
} from './types';

export interface SubmitRef {
  isFormValid: boolean;
  submit: () => void;
}

const formSchema = yup.object().shape({
  targetSegmentRules: yup
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
              conditions: yup
                .array()
                .required()
                .of(
                  yup.object().shape({
                    situation: yup
                      .string()
                      .oneOf([
                        'compare',
                        'user-segment',
                        'date',
                        'feature-flag'
                      ])
                      .required(),
                    conditioner: yup.string().required(),
                    firstValue: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'compare')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      }),
                    secondValue: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'compare')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      }),
                    value: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (
                        !value &&
                        ['user-segment', 'date'].includes(situation)
                      )
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });

                      return true;
                    }),
                    date: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (!value && situation === 'date')
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });
                      return true;
                    }),
                    flagId: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (!value && situation === 'feature-flag')
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });

                      return true;
                    }),
                    variation: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'feature-flag')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      })
                  })
                ),
              variation: yup.string().required()
            })
          )
      })
    )
});

const Targeting = () => {
  const [targetSegmentRules, setTargetSegmentRules] = useState<
    TargetSegmentItem[]
  >([]);
  const [targetIndividualRules, setTargetIndividualRules] = useState<
    TargetIndividualItem[]
  >([]);
  const [prerequisitesRules, setPrerequisitesRules] = useState<
    TargetPrerequisiteItem[]
  >([]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      targetSegmentRules: [
        {
          index: 1,
          rules: [
            {
              variation: '',
              conditions: [initialSegmentCondition]
            }
          ]
        }
      ]
    }
  });

  const onAddRule = useCallback(
    (type: RuleCategory) => {
      if (type === 'target-segments') {
        return setTargetSegmentRules([
          ...targetSegmentRules,
          {
            index: targetSegmentRules.length + 1,
            rules: [
              {
                variation: true,
                conditions: [initialSegmentCondition]
              }
            ]
          }
        ]);
      }
      if (type === 'target-individuals') {
        return setTargetIndividualRules([
          ...targetIndividualRules,
          {
            index: targetIndividualRules.length + 1,
            rules: [initialIndividualRule]
          }
        ]);
      }
      setPrerequisitesRules([
        ...prerequisitesRules,
        {
          index: prerequisitesRules.length + 1,
          rules: [initialPrerequisitesRule]
        }
      ]);
    },
    [targetSegmentRules, targetIndividualRules, prerequisitesRules]
  );

  const onSubmit = async values => {
    console.log(values);
  };

  return (
    <>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="flex flex-col size-full gap-y-6 overflow-visible">
            <TargetingState />
            <AddRuleDropdown onAddRule={onAddRule} />
            <Form.Field
              control={form.control}
              {...form.register('targetSegmentRules')}
              render={() => (
                <Form.Item>
                  <Form.Control>
                    <TargetSegmentRule
                      targetSegmentRules={targetSegmentRules}
                      setTargetSegmentRules={setTargetSegmentRules}
                    />
                  </Form.Control>
                </Form.Item>
              )}
            />
            <DefaultRule />
          </div>
        </Form>
      </FormProvider>
    </>
  );
};

export default Targeting;
