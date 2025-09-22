import { useMemo } from 'react';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import {
  IconPlus,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import { RuleCategory } from '../types';

const AddRule = ({
  isDisableAddPrerequisite,
  isDisableAddIndividualRules,
  isInsertSegmentRule,
  indexInsertSegmentRule,
  onAddRule
}: {
  isDisableAddPrerequisite: boolean;
  isDisableAddIndividualRules: boolean;
  isInsertSegmentRule?: boolean;
  indexInsertSegmentRule?: number;
  onAddRule: (rule: RuleCategory, index?: number) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);

  const options = useMemo(
    () => [
      {
        label: t('feature-flags.prerequisites'),
        value: RuleCategory.PREREQUISITE,
        tooltip: t('targeting.prerequisite-tooltip'),
        icon: IconPrerequisite,
        disabled: isDisableAddPrerequisite
      },
      {
        label: t('targeting.individual-targeting'),
        value: RuleCategory.INDIVIDUAL,
        tooltip: t('targeting.individual-tooltip'),
        icon: IconUserOutlined,
        disabled: isDisableAddIndividualRules
      },
      {
        label: t('targeting.custom-rule'),
        value: RuleCategory.CUSTOM,
        tooltip: t('targeting.custom-rule-tooltip'),
        icon: IconSetting
      }
    ],
    [isDisableAddIndividualRules, isDisableAddPrerequisite]
  );

  const getRuleCategoryCall = (value: RuleCategory) => {
    if (value === RuleCategory.CUSTOM) {
      if (isInsertSegmentRule) {
        return onAddRule(RuleCategory.CUSTOM, indexInsertSegmentRule);
      }
    }
    return onAddRule(value);
  };

  return (
    <Dropdown
      trigger={
        <DisabledButtonTooltip
          align="center"
          hidden={editable}
          trigger={
            <div className="flex items-center gap-x-2 h-6 p-0 typo-para-medium !text-primary-500">
              <Icon icon={IconPlus} size={'md'} />
              {t('table:feature-flags.add-rule')}
            </div>
          }
        />
      }
      isTooltip={true}
      options={options}
      disabled={!editable}
      showArrow={false}
      onChange={value => onAddRule(value as RuleCategory)}
      alignContent="center"
      className="w-full [&>div]:flex-center border-dashed !shadow-none"
    />
  );
};

export default AddRule;
