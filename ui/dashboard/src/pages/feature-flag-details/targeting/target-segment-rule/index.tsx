import { useCallback, useMemo } from 'react';
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
import Icon from 'components/icon';
import Card from '../../elements/card';
import AddRuleButton from '../add-rule-button';
import { initialSegmentCondition } from '../constants';
import { TargetSegmentItem } from '../types';
import Condition from './condition';
import SegmentVariation from './variation';

const TargetSegmentRule = ({
  targetSegmentRules,
  setTargetSegmentRules
}: {
  targetSegmentRules: TargetSegmentItem[];
  setTargetSegmentRules: (value: TargetSegmentItem[]) => void;
}) => {
  const { t } = useTranslation(['table', 'form']);
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
      cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions.splice(
        conditionIndex,
        1
      );
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
          ...cloneTargetSegmentRules[segmentIndex].rules[ruleIndex].conditions[
            conditionIndex
          ],
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

  return (
    targetSegmentRules.length > 0 &&
    targetSegmentRules.map((segment, index) => (
      <div className="flex flex-col w-full gap-y-6" key={`segment-${index}`}>
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
                      onDeleteCondition(index, ruleIndex, conditionIndex)
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
    ))
  );
};

export default TargetSegmentRule;
