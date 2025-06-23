import { useMemo } from 'react';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import {
  IconPlus,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import { RuleCategory } from '../types';

const AddRule = ({
  isDisableAddPrerequisite,
  isDisableAddIndividualRules,
  onAddRule
}: {
  isDisableAddPrerequisite: boolean;
  isDisableAddIndividualRules: boolean;
  onAddRule: (rule: RuleCategory) => void;
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

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
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
        disabled={!editable}
        showArrow={false}
        className="w-full [&>div]:flex-center border-dashed !shadow-none"
      />
      <DropdownMenuContent>
        {options.map((item, index) => (
          <Tooltip
            side="right"
            sideOffset={10}
            align="start"
            className="w-[180px] p-3 bg-white typo-para-small text-gray-600 shadow-card"
            key={index}
            content={item.tooltip}
            showArrow={false}
            trigger={
              <DropdownMenuItem
                icon={item.icon}
                label={item.label}
                value={item.value}
                disabled={item?.disabled}
                onSelectOption={value => onAddRule(value as RuleCategory)}
              />
            }
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default AddRule;
