import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import { Fragment } from 'react/jsx-runtime';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import { initialSegmentCondition } from '../constants';
import { TargetSegmentItem } from '../types';
import Condition from './condition';
import SegmentVariation from './variation';

interface Props {
  targetSegmentRules: TargetSegmentItem[];
  onChangeTargetSegmentRules: (value: TargetSegmentItem[]) => void;
  onAddRule: () => void;
}

const TargetSegmentRule = ({
  targetSegmentRules,
  onChangeTargetSegmentRules,
  onAddRule
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const cloneTargetSegmentRules = useMemo(
    () => cloneDeep(targetSegmentRules),
    [targetSegmentRules]
  );
  const methods = useFormContext();

  const { control } = methods;

  const onAddCondition = useCallback(
    (segmentIndex: number, ruleIndex: number) => {
      cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions.push(
        initialSegmentCondition
      );
      onChangeTargetSegmentRules(cloneTargetSegmentRules);
    },
    [targetSegmentRules, cloneTargetSegmentRules]
  );

  const onDeleteCondition = useCallback(
    (segmentIndex: number, ruleIndex: number, conditionIndex: number) => {
      const cloneTargetSegmentRules = cloneDeep(targetSegmentRules);
      cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions.splice(
        conditionIndex,
        1
      );
      onChangeTargetSegmentRules(cloneTargetSegmentRules);
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
          ...cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions[
            conditionIndex
          ],
          [field]: value
        };
        return onChangeTargetSegmentRules(cloneTargetSegmentRules);
      }
      cloneTargetSegmentRules[segmentIndex].rules[ruleIndex] = {
        ...cloneTargetSegmentRules[segmentIndex].rules[ruleIndex],
        [field]: value
      };
      onChangeTargetSegmentRules(cloneTargetSegmentRules);
    },
    [cloneTargetSegmentRules]
  );

  return (
    targetSegmentRules.length > 0 && (
      <div className="w-full">
        {targetSegmentRules.map((segment, segmentIndex) => (
          <div key={segmentIndex} className="flex flex-col w-full gap-y-6">
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
                        index: segmentIndex + 1
                      }}
                    />
                  </p>
                  {targetSegmentRules.length > 1 && (
                    <div className="flex items-center gap-x-1">
                      {segmentIndex !== targetSegmentRules.length - 1 && (
                        <Icon
                          icon={IconArrowDownwardFilled}
                          color="gray-500"
                          size={'sm'}
                        />
                      )}
                      {segmentIndex !== 0 && (
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
                  <Form.Field
                    key={`rule-${ruleIndex}`}
                    control={control}
                    name={`targetSegmentRules.${segmentIndex}.rules.${ruleIndex}`}
                    render={({ field }) => (
                      <Fragment>
                        {rule.conditions.map((condition, conditionIndex) => (
                          <Condition
                            key={`condition-${conditionIndex}`}
                            isDisabledDelete={rule.conditions.length <= 1}
                            segmentIndex={segmentIndex}
                            ruleIndex={ruleIndex}
                            conditionIndex={conditionIndex}
                            type={conditionIndex === 0 ? 'if' : 'and'}
                            condition={condition}
                            onDeleteCondition={() =>
                              onDeleteCondition(
                                segmentIndex,
                                ruleIndex,
                                conditionIndex
                              )
                            }
                            onChangeFormField={(field, value) =>
                              onChangeFormField(
                                segmentIndex,
                                ruleIndex,
                                field,
                                value,
                                conditionIndex
                              )
                            }
                            {...field}
                          />
                        ))}
                        <Button
                          type="button"
                          variant={'text'}
                          className="w-fit gap-x-2 h-6 !p-0"
                          onClick={() =>
                            onAddCondition(segmentIndex, ruleIndex)
                          }
                        >
                          <Icon
                            icon={IconPlus}
                            color="primary-500"
                            className="flex-center"
                            size={'sm'}
                          />{' '}
                          {t('form:feature-flags.add-condition')}
                        </Button>
                        <SegmentVariation
                          segmentIndex={segmentIndex}
                          ruleIndex={ruleIndex}
                        />
                      </Fragment>
                    )}
                  />
                ))}
                <AddRuleButton onAddRule={onAddRule} />
              </Card>
            </Card>
          </div>
        ))}
      </div>
    )
  );
};

export default TargetSegmentRule;
