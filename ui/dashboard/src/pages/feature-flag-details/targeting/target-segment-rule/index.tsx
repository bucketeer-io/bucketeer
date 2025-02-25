import {
  forwardRef,
  Ref,
  useCallback,
  useImperativeHandle,
  useMemo,
  useRef
} from 'react';
import { Form, FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import { Fragment } from 'react/jsx-runtime';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import * as yup from 'yup';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { SubmitRef } from '..';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import { initialSegmentCondition } from '../constants';
import { TargetSegmentItem } from '../types';
import Condition from './condition';
import SegmentVariation from './variation';

interface Props {
  targetSegmentRules: TargetSegmentItem[];
  setTargetSegmentRules: (value: TargetSegmentItem[]) => void;
}

const formSchema = yup.object().shape({
  situation: yup
    .string()
    .oneOf(['compare', 'user-segment', 'date', 'feature-flag'])
    .required(),
  conditioner: yup.string().required(),
  firstValue: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && situation === 'compare')
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });

    return true;
  }),
  secondValue: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && situation === 'compare')
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });

    return true;
  }),
  value: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && ['user-segment', 'date'].includes(situation))
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });

    return true;
  }),
  date: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && situation === 'date')
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });
    return true;
  }),
  flagId: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && situation === 'feature-flag')
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });

    return true;
  }),
  variation: yup.string().test('required', (value, context) => {
    const situation = context.from && context.from[0].value.situation;
    if (!value && situation === 'feature-flag')
      return context.createError({
        message: `This field is required.`,
        path: context.path
      });

    return true;
  })
});

const TargetSegmentRule = forwardRef(
  (
    { targetSegmentRules, setTargetSegmentRules }: Props,
    ref: Ref<SubmitRef>
  ) => {
    const { t } = useTranslation(['table', 'form']);
    const submitBtnRef = useRef<HTMLButtonElement>(null);

    const form = useForm({
      resolver: yupResolver(formSchema),
      defaultValues: initialSegmentCondition
    });

    useImperativeHandle(
      ref,
      () => ({
        isFormValid: false,
        submit: () => submitBtnRef?.current?.click()
      }),
      []
    );

    const cloneTargetSegmentRules = useMemo(
      () => cloneDeep(targetSegmentRules),
      [targetSegmentRules]
    );

    const onAddCondition = useCallback(
      (segmentIndex: number, ruleIndex: number) => {
        cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions.push(
          initialSegmentCondition
        );
        setTargetSegmentRules(cloneTargetSegmentRules);
      },
      [targetSegmentRules, cloneTargetSegmentRules]
    );

    const onDeleteCondition = useCallback(
      (segmentIndex: number, ruleIndex: number, conditionIndex: number) => {
        const cloneTargetSegmentRules = cloneDeep(targetSegmentRules);
        cloneTargetSegmentRules[segmentIndex].rules[
          ruleIndex
        ].conditions.splice(conditionIndex, 1);
        setTargetSegmentRules(cloneTargetSegmentRules);
      },
      [targetSegmentRules, cloneTargetSegmentRules]
    );

    const onChangeFormField = useCallback(
      (
        segmentIndex: number,
        ruleIndex: number,
        field: string,
        value: string | number | boolean,
        conditionIndex?: number
      ) => {
        if (typeof conditionIndex === 'number') {
          cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions[
            conditionIndex
          ] = {
            ...cloneTargetSegmentRules[segmentIndex].rules[ruleIndex]
              .conditions[conditionIndex],
            [field]: value
          };
          return setTargetSegmentRules(cloneTargetSegmentRules);
        }
        cloneTargetSegmentRules[segmentIndex].rules[ruleIndex] = {
          ...cloneTargetSegmentRules[segmentIndex].rules[ruleIndex],
          [field]: value
        };
        setTargetSegmentRules(cloneTargetSegmentRules);
      },
      [cloneTargetSegmentRules]
    );

    const onSubmit = async () => {
      // console.log(submitRef.current);
      // submitRef.current?.submit();
      // console.log(submitRef.current?.isFormValid)
    };

    return (
      targetSegmentRules.length > 0 && (
        <FormProvider {...form}>
          <Form onSubmit={onSubmit} className="w-full">
            <Button ref={submitBtnRef} type="submit" className="hidden">
              Submit
            </Button>
            {targetSegmentRules.map((segment, index) => (
              <div
                className="flex flex-col w-full gap-y-6"
                key={`segment-${index}`}
              >
                <Card>
                  <div>
                    <div className="flex items-center gap-x-2">
                      <p className="typo-para-medium leading-4 text-gray-700">
                        {t('feature-flags.rules')}
                      </p>
                      <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                    </div>
                  </div>
                  <Card className="shadow-none border border-gray-400">
                    <div className="flex items-center justify-between w-full">
                      <p className="typo-para-medium leading-5 text-gray-700">
                        <Trans
                          i18nKey={'table:feature-flags.rule-index'}
                          values={{
                            index: index + 1
                          }}
                        />
                      </p>
                      {targetSegmentRules.length > 1 && (
                        <div className="flex items-center gap-x-1">
                          {index !== targetSegmentRules.length - 1 && (
                            <Icon
                              icon={IconArrowDownwardFilled}
                              color="gray-500"
                              size={'sm'}
                            />
                          )}
                          {index !== 0 && (
                            <Icon
                              icon={IconArrowUpwardFilled}
                              color="gray-500"
                              size={'sm'}
                            />
                          )}
                        </div>
                      )}
                    </div>
                    {segment.rules.map((rule, ruleIndex) => (
                      <Fragment key={`rule-${ruleIndex}`}>
                        {rule.conditions.map((condition, conditionIndex) => (
                          <Condition
                            key={`condition-${conditionIndex}`}
                            isDisabledDelete={rule.conditions.length <= 1}
                            type={conditionIndex === 0 ? 'if' : 'and'}
                            condition={condition}
                            onDeleteCondition={() =>
                              onDeleteCondition(
                                index,
                                ruleIndex,
                                conditionIndex
                              )
                            }
                            onChangeFormField={(field, value) =>
                              onChangeFormField(
                                index,
                                ruleIndex,
                                field,
                                value,
                                conditionIndex
                              )
                            }
                          />
                        ))}
                        <Button
                          type="button"
                          variant={'text'}
                          className="w-fit gap-x-2 h-6 !p-0"
                          onClick={() => onAddCondition(index, ruleIndex)}
                        >
                          <Icon
                            icon={IconPlus}
                            color="primary-500"
                            className="flex-center"
                            size={'sm'}
                          />{' '}
                          {t('form:feature-flags.add-condition')}
                        </Button>
                        <SegmentVariation variation={rule.variation} />
                      </Fragment>
                    ))}
                    <AddRuleButton />
                  </Card>
                </Card>
              </div>
            ))}
          </Form>
        </FormProvider>
      )
    );
  }
);

export default TargetSegmentRule;
