import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Fragment } from 'react/jsx-runtime';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import { initialPrerequisitesRule } from '../constants';
import { TargetPrerequisiteItem } from '../types';
import ConditionForm from './condition';

interface Props {
  prerequisitesRules: TargetPrerequisiteItem[];
  onChangePrerequisitesRules: (value: TargetPrerequisiteItem[]) => void;
}

const PrerequisiteRule = ({
  prerequisitesRules,
  onChangePrerequisitesRules
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const clonePrerequisitesRules = useMemo(
    () => cloneDeep(prerequisitesRules),
    [prerequisitesRules]
  );
  const methods = useFormContext();

  const { control } = methods;

  const onAddCondition = useCallback(
    (prerequisiteIndex: number) => {
      clonePrerequisitesRules[prerequisiteIndex].rules.push(
        initialPrerequisitesRule
      );
      onChangePrerequisitesRules(clonePrerequisitesRules);
    },
    [prerequisitesRules, clonePrerequisitesRules]
  );

  const onDeleteCondition = useCallback(
    (segmentIndex: number, ruleIndex: number) => {
      clonePrerequisitesRules[segmentIndex].rules.splice(ruleIndex, 1);
      onChangePrerequisitesRules(clonePrerequisitesRules);
    },
    [prerequisitesRules, clonePrerequisitesRules]
  );

  const onChangeFormField = useCallback(
    (
      segmentIndex: number,
      ruleIndex: number,
      field: string,
      value: string | number | boolean
    ) => {
      clonePrerequisitesRules[segmentIndex].rules[ruleIndex] = {
        ...clonePrerequisitesRules[segmentIndex].rules[ruleIndex],
        [field]: value
      };
      return onChangePrerequisitesRules(clonePrerequisitesRules);
    },
    [clonePrerequisitesRules]
  );

  return (
    prerequisitesRules.length > 0 && (
      <div className="flex flex-col gap-y-6 w-full">
        {prerequisitesRules.map((prerequisite, prerequisiteIndex) => (
          <div
            key={`prerequisite-${prerequisiteIndex}`}
            className="flex flex-col w-full gap-y-6"
          >
            <Card>
              <div>
                <div className="flex items-center gap-x-2">
                  <p className="typo-para-medium leading-4 text-gray-700">
                    {t('form:feature-flags.prerequisites')}
                  </p>
                  <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                </div>
              </div>
              {prerequisite.rules.map((rule, ruleIndex) => (
                <Form.Field
                  key={`prerequisite-rule-${ruleIndex}`}
                  control={control}
                  name={`prerequisitesRules.${prerequisiteIndex}.rules.${ruleIndex}`}
                  render={({ field }) => (
                    <Fragment>
                      <ConditionForm
                        condition={rule}
                        ruleIndex={ruleIndex}
                        prerequisiteIndex={prerequisiteIndex}
                        type={ruleIndex === 0 ? 'if' : 'and'}
                        isDisabledDelete={prerequisite.rules.length <= 1}
                        onChangeFormField={(field, value) =>
                          onChangeFormField(
                            prerequisiteIndex,
                            ruleIndex,
                            field,
                            value
                          )
                        }
                        onDeleteCondition={() =>
                          onDeleteCondition(prerequisiteIndex, ruleIndex)
                        }
                        {...field}
                      />
                    </Fragment>
                  )}
                />
              ))}
              <Button
                type="button"
                variant={'text'}
                className="w-fit gap-x-2 h-6 !p-0"
                onClick={() => onAddCondition(prerequisiteIndex)}
              >
                <Icon
                  icon={IconPlus}
                  color="primary-500"
                  className="flex-center"
                  size={'sm'}
                />{' '}
                {t('form:feature-flags.add-prerequisites')}
              </Button>
            </Card>
          </div>
        ))}
      </div>
    )
  );
};

export default PrerequisiteRule;
