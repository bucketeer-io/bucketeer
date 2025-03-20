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
import { v4 as uuid } from 'uuid';
import { Feature, FeatureRuleClauseOperator } from '@types';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import { RuleClauseType, RuleSchema } from '../form-schema';
import Condition from './condition';

interface Props {
  features: Feature[];
  targetSegmentRules: RuleSchema[];
  onChangeTargetSegmentRules: (value: RuleSchema[]) => void;
  onAddRule: () => void;
}

const TargetSegmentRule = ({
  features,
  targetSegmentRules,
  onChangeTargetSegmentRules,
  onAddRule
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const cloneTargetSegmentRules: RuleSchema[] = useMemo(
    () => cloneDeep(targetSegmentRules),
    [targetSegmentRules]
  );
  const methods = useFormContext();

  const { control } = methods;

  const onAddCondition = useCallback(
    (ruleIndex: number) => {
      cloneTargetSegmentRules[ruleIndex].clauses.push({
        id: uuid(),
        type: RuleClauseType.COMPARE,
        attribute: '',
        operator: FeatureRuleClauseOperator.EQUALS,
        values: []
      });
      onChangeTargetSegmentRules(cloneTargetSegmentRules);
    },
    [targetSegmentRules, cloneTargetSegmentRules]
  );

  const onDeleteCondition = useCallback(
    (ruleIndex: number, conditionIndex: number) => {
      const cloneTargetSegmentRules = cloneDeep(targetSegmentRules);
      cloneTargetSegmentRules[ruleIndex].clauses.splice(conditionIndex, 1);
      onChangeTargetSegmentRules(cloneTargetSegmentRules);
    },
    [targetSegmentRules, cloneTargetSegmentRules]
  );

  const onChangeFormField = useCallback(
    (
      ruleIndex: number,
      field: string,
      value: string | string[],
      conditionIndex: number
    ) => {
      cloneTargetSegmentRules[ruleIndex].clauses[conditionIndex] = {
        ...cloneTargetSegmentRules[ruleIndex].clauses[conditionIndex],
        [field]: value
      };
      return onChangeTargetSegmentRules(cloneTargetSegmentRules);
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
                  {/* <SegmentVariation
                    segmentIndex={segmentIndex}
                    ruleIndex={ruleIndex}
                  /> */}
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
