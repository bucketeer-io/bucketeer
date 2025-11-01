import { IconUndoOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { IconInfo, IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import Card from '../../elements/card';
import { DiscardChangesType, PrerequisiteSchema, RuleCategory } from '../types';
import ConditionForm from './condition';
import PrerequisiteBanner from './prerequisite-banner';

interface Props {
  feature: Feature;
  features: Feature[];
  prerequisites: PrerequisiteSchema[];
  hasPrerequisiteFlags: Feature[];
  isDisableAddPrerequisite: boolean;
  onRemovePrerequisite: (index: number) => void;
  onAddPrerequisite: () => void;
  handleCheckEdit?: (type: RuleCategory) => boolean;
  handleDiscardChanges: (type: DiscardChangesType) => void;
}

const PrerequisiteRule = ({
  feature,
  features,
  prerequisites,
  hasPrerequisiteFlags,
  isDisableAddPrerequisite,
  onRemovePrerequisite,
  onAddPrerequisite,
  handleCheckEdit,
  handleDiscardChanges
}: Props) => {
  const { t } = useTranslation(['table', 'form']);
  const editPrerequisite = handleCheckEdit?.(RuleCategory.PREREQUISITE);
  return (
    <div className="flex flex-col gap-y-6 w-full">
      {hasPrerequisiteFlags?.length > 0 && (
        <PrerequisiteBanner hasPrerequisiteFlags={hasPrerequisiteFlags} />
      )}
      {prerequisites.length > 0 && (
        <div className="flex flex-col w-full gap-y-6">
          <Card>
            <div className="w-full h-8 flex items-center justify-between">
              <div className="flex items-center gap-x-2">
                <p className="typo-para-medium leading-4 text-gray-700">
                  {t('form:feature-flags.prerequisites')}
                </p>
                <Tooltip
                  align="start"
                  alignOffset={-105}
                  content={t('form:targeting.tooltip.prerequisites')}
                  trigger={
                    <div className="flex-center size-fit">
                      <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                    </div>
                  }
                  className="max-w-[400px]"
                />
              </div>
              {editPrerequisite && (
                <div
                  className="flex-center h-8 w-8 px-2 rounded-md cursor-pointer group border border-gray-300 hover:border-gray-800"
                  onClick={() =>
                    handleDiscardChanges(DiscardChangesType.PREREQUISITE)
                  }
                >
                  <Icon
                    icon={IconUndoOutlined}
                    size={'sm'}
                    className="flex-center text-gray-500 group-hover:text-gray-700"
                  />
                </div>
              )}
            </div>
            {prerequisites.map((item, prerequisiteIndex) => (
              <ConditionForm
                key={`${item.featureId}.${prerequisiteIndex}`}
                features={features}
                featureId={feature.id}
                prerequisiteIndex={prerequisiteIndex}
                type={prerequisiteIndex === 0 ? 'if' : 'and'}
                onDeleteCondition={() =>
                  onRemovePrerequisite(prerequisiteIndex)
                }
              />
            ))}

            <Button
              type="button"
              variant={'text'}
              className="w-fit gap-x-2 h-6 !p-0"
              disabled={isDisableAddPrerequisite}
              onClick={() => onAddPrerequisite()}
            >
              <Icon icon={IconPlus} className="flex-center" size={'sm'} />
              {t('form:feature-flags.add-prerequisites')}
            </Button>
          </Card>
        </div>
      )}
    </div>
  );
};

export default PrerequisiteRule;
