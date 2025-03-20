import { useCallback } from 'react';
import { useFormContext } from 'react-hook-form';
import { Fragment } from 'react/jsx-runtime';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { Feature } from '@types';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Card from '../../elements/card';
import { initialPrerequisite } from '../constants';
import { PrerequisiteSchema } from '../types';
import ConditionForm from './condition';
import PrerequisiteBanner from './prequisite-banner';

interface Props {
  feature: Feature;
  features: Feature[];
  prerequisites: PrerequisiteSchema[];
  onChangePrerequisites: (value: PrerequisiteSchema[]) => void;
}

const PrerequisiteRule = ({
  feature,
  features,
  prerequisites,
  onChangePrerequisites
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

  const methods = useFormContext();

  const { control, setValue } = methods;

  const onAddCondition = useCallback(() => {
    prerequisites.push(cloneDeep(initialPrerequisite));
    onChangePrerequisites(prerequisites);
  }, [prerequisites]);

  const onDeleteCondition = useCallback(
    (currentIndex: number) => {
      prerequisites.splice(currentIndex, 1);
      onChangePrerequisites(prerequisites);
    },
    [prerequisites]
  );

  const onChangeFormField = useCallback(
    (
      prerequisiteIndex: number,
      field: string,
      value: string | number | boolean
    ) => {
      prerequisites[prerequisiteIndex] = {
        ...prerequisites[prerequisiteIndex],
        [field]: value
      };
      setValue('prerequisites', [...prerequisites]);
      return onChangePrerequisites([...prerequisites]);
    },
    [prerequisites]
  );

  return (
    prerequisites.length > 0 && (
      <div className="flex flex-col gap-y-6 w-full">
        {feature?.prerequisites?.length > 0 && (
          <PrerequisiteBanner
            features={features}
            prerequisite={feature.prerequisites}
          />
        )}

        <div className="flex flex-col w-full gap-y-6">
          <Card>
            <div>
              <div className="flex items-center gap-x-2">
                <p className="typo-para-medium leading-4 text-gray-700">
                  {t('form:feature-flags.prerequisites')}
                </p>
                <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
              </div>
            </div>
            {prerequisites.map((prerequisite, prerequisiteIndex) => (
              <Form.Field
                key={`prerequisite-rule-${prerequisiteIndex}`}
                control={control}
                name={`prerequisites.${prerequisiteIndex}`}
                render={({ field }) => (
                  <Fragment>
                    <ConditionForm
                      prerequisites={prerequisites}
                      features={features}
                      prerequisite={prerequisite}
                      prerequisiteIndex={prerequisiteIndex}
                      type={prerequisiteIndex === 0 ? 'if' : 'and'}
                      isDisabledDelete={prerequisites.length <= 1}
                      onChangeFormField={(field, value) =>
                        onChangeFormField(prerequisiteIndex, field, value)
                      }
                      onDeleteCondition={() =>
                        onDeleteCondition(prerequisiteIndex)
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
              onClick={() => onAddCondition()}
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
      </div>
    )
  );
};

export default PrerequisiteRule;
