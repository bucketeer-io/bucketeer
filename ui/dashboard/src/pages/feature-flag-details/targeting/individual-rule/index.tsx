import { useCallback, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Fragment } from 'react/jsx-runtime';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { IconInfo } from '@icons';
import { CreatableSelect } from 'components/creatable-select';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import { IndividualRuleItem } from '../types';

interface Props {
  individualRules: IndividualRuleItem[];
  onChangeIndividualRules: (value: IndividualRuleItem[]) => void;
}

const IndividualRule = ({
  individualRules,
  onChangeIndividualRules
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const cloneIndividualRules = useMemo(
    () => cloneDeep(individualRules),
    [individualRules]
  );
  const methods = useFormContext();

  const { control } = methods;

  const onChangeFormField = useCallback(
    (individualIndex: number, field: string, value: string[]) => {
      cloneIndividualRules[individualIndex] = {
        ...cloneIndividualRules[individualIndex],
        [field]: value
      };
      return onChangeIndividualRules(cloneIndividualRules);
    },
    [cloneIndividualRules]
  );

  return (
    individualRules.length > 0 && (
      <div className="flex flex-col gap-y-6 w-full">
        {individualRules.map((_, individualIndex) => (
          <Card key={individualIndex}>
            <div>
              <div className="flex items-center gap-x-2">
                <p className="typo-para-medium leading-4 text-gray-700">
                  {t('form:feature-flags.prerequisites')}
                </p>
                <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
              </div>
            </div>

            <Form.Field
              key={`individual-rule-${individualIndex}`}
              control={control}
              name={`targetIndividualRules.${individualIndex}.on`}
              render={({ field }) => (
                <Fragment>
                  <CreatableSelect
                    defaultValues={field.value?.map((item: string) => ({
                      label: item,
                      value: item
                    }))}
                    onChange={value => {
                      field.onChange(value);
                      onChangeFormField(
                        individualIndex,
                        'on',
                        value?.map(item => item.value)
                      );
                    }}
                  />
                </Fragment>
              )}
            />
          </Card>
        ))}
      </div>
    )
  );
};

export default IndividualRule;
