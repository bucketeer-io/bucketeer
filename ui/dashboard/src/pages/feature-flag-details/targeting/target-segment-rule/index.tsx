import { useCallback } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import { Fragment } from 'react/jsx-runtime';
import { useTranslation } from 'i18n';
import { v4 as uuid } from 'uuid';
import {
  Feature,
  FeatureRuleClauseOperator,
  RuleStrategyVariation
} from '@types';
import { IconClose, IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import { RuleClauseSchema, RuleClauseType, RuleSchema } from '../form-schema';
import Condition from './condition';
import SegmentVariation from './variation';

interface Props {
  feature: Feature;
  features: Feature[];
  defaultRolloutStrategy: RuleStrategyVariation[];
  targetSegmentRules: RuleSchema[];
  onChangeTargetSegmentRules: (value: RuleSchema[]) => void;
  onAddRule: () => void;
  onDeleteRule: (index: number) => void;
}

const TargetSegmentRule = ({
  feature,
  features,
  defaultRolloutStrategy,
  targetSegmentRules,
  onChangeTargetSegmentRules,
  onAddRule,
  onDeleteRule
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const methods = useFormContext();

  const { control } = methods;

  const onAddCondition = useCallback(
    (segmentIndex: number) => {
      targetSegmentRules[segmentIndex].clauses.push({
        id: uuid(),
        type: RuleClauseType.COMPARE,
        attribute: '',
        operator: FeatureRuleClauseOperator.EQUALS,
        values: []
      });
      onChangeTargetSegmentRules([...targetSegmentRules]);
    },
    [targetSegmentRules]
  );

  const onDeleteCondition = useCallback(
    (ruleIndex: number, conditionIndex: number) => {
      targetSegmentRules[ruleIndex].clauses.splice(conditionIndex, 1);
      onChangeTargetSegmentRules(targetSegmentRules);
    },
    [targetSegmentRules]
  );

  const onChangeFormField = useCallback(
    (
      ruleIndex: number,
      field: keyof RuleClauseSchema,
      value: string | string[],
      clauseIndex: number
    ) => {
      targetSegmentRules[ruleIndex].clauses[clauseIndex] = {
        ...targetSegmentRules[ruleIndex].clauses[clauseIndex],
        [field]: value
      };
      return onChangeTargetSegmentRules([...targetSegmentRules]);
    },

    [targetSegmentRules]
  );

  const handleChangeIndexRule = useCallback(
    (type: 'increase' | 'decrease', currentIndex: number) => {
      const isIncrease = type === 'increase';
      const _rules = targetSegmentRules.map((rule, index) => {
        const targetIndex = isIncrease ? currentIndex + 1 : currentIndex - 1;
        if ([currentIndex, targetIndex].includes(index))
          return targetSegmentRules[
            index === currentIndex ? targetIndex : currentIndex
          ];
        return rule;
      });
      return onChangeTargetSegmentRules(_rules);
    },
    [targetSegmentRules]
  );

  return (
    targetSegmentRules.length > 0 && (
      <div className="flex flex-col w-full gap-y-6">
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
                  <div className="flex items-center gap-x-2">
                    <div
                      className="flex-center cursor-pointer group"
                      onClick={() => onDeleteRule(segmentIndex)}
                    >
                      <Icon
                        icon={IconClose}
                        size={'sm'}
                        className="flex-center text-gray-500 group-hover:text-gray-700"
                      />
                    </div>
                    {targetSegmentRules.length > 1 && (
                      <div className="flex items-center gap-x-1">
                        {segmentIndex !== targetSegmentRules.length - 1 && (
                          <div
                            className="flex-center group cursor-pointer"
                            onClick={() =>
                              handleChangeIndexRule('increase', segmentIndex)
                            }
                          >
                            <Icon
                              icon={IconArrowDownwardFilled}
                              size={'sm'}
                              className="text-gray-500 group-hover:text-gray-700"
                            />
                          </div>
                        )}
                        {segmentIndex !== 0 && (
                          <div
                            className="flex-center group cursor-pointer"
                            onClick={() =>
                              handleChangeIndexRule('decrease', segmentIndex)
                            }
                          >
                            <Icon
                              icon={IconArrowUpwardFilled}
                              size={'sm'}
                              className="text-gray-500 group-hover:text-gray-700"
                            />
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </div>
                <Fragment>
                  {segment.clauses.map((clause, clauseIndex) => (
                    <Form.Field
                      key={`clause-${clauseIndex}`}
                      control={control}
                      name={`rules.${segmentIndex}.clauses.${clauseIndex}`}
                      render={({ field }) => (
                        <Condition
                          features={features}
                          isDisabledDelete={segment.clauses.length <= 1}
                          segmentIndex={segmentIndex}
                          clauseIndex={clauseIndex}
                          type={clauseIndex === 0 ? 'if' : 'and'}
                          clause={clause}
                          onDeleteCondition={() =>
                            onDeleteCondition(segmentIndex, clauseIndex)
                          }
                          onChangeFormField={(field, value) =>
                            onChangeFormField(
                              segmentIndex,
                              field,
                              value,
                              clauseIndex
                            )
                          }
                          {...field}
                        />
                      )}
                    />
                  ))}
                  <Button
                    type="button"
                    variant={'text'}
                    className="w-fit gap-x-2 h-6 !p-0"
                    onClick={() => onAddCondition(segmentIndex)}
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
                    feature={feature}
                    defaultRolloutStrategy={defaultRolloutStrategy}
                    segmentIndex={segmentIndex}
                    targetSegmentRules={targetSegmentRules}
                    onChangeTargetSegmentRules={onChangeTargetSegmentRules}
                  />
                </Fragment>
              </Card>
              <AddRuleButton onAddRule={onAddRule} />
            </Card>
          </div>
        ))}
      </div>
    )
  );
};

export default TargetSegmentRule;
