import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import Card from '../../elements/card';
import { PrerequisiteSchema } from '../types';
import ConditionForm from './condition';
import PrerequisiteBanner from './prequisite-banner';

interface Props {
  feature: Feature;
  features: Feature[];
  prerequisites: PrerequisiteSchema[];
  onRemovePrerequisite: (index: number) => void;
  onAddPrerequisite: () => void;
}

const PrerequisiteRule = ({
  feature,
  features,
  prerequisites,
  onRemovePrerequisite,
  onAddPrerequisite
}: Props) => {
  const { t } = useTranslation(['table', 'form']);

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
              <ConditionForm
                key={prerequisiteIndex}
                features={features}
                featureId={feature.id}
                prerequisite={prerequisite}
                prerequisiteIndex={prerequisiteIndex}
                type={prerequisiteIndex === 0 ? 'if' : 'and'}
                isDisabledDelete={prerequisites.length <= 1}
                onDeleteCondition={() =>
                  onRemovePrerequisite(prerequisiteIndex)
                }
              />
            ))}

            <Button
              type="button"
              variant={'text'}
              className="w-fit gap-x-2 h-6 !p-0"
              onClick={() => onAddPrerequisite()}
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
