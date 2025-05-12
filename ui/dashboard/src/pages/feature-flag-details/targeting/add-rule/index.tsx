import { useMemo } from 'react';
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
import { RuleCategory } from '../types';

const AddRule = ({
  onAddRule
}: {
  onAddRule: (rule: RuleCategory) => void;
}) => {
  const { t } = useTranslation(['form', 'table']);

  const options = useMemo(
    () => [
      {
        label: t('feature-flags.prerequisites'),
        value: RuleCategory.PREREQUISITE,
        tooltip: t('targeting.prerequisite-tooltip'),
        icon: IconPrerequisite
      },
      {
        label: t('targeting.individual-targeting'),
        value: RuleCategory.INDIVIDUAL,
        tooltip: t('targeting.individual-tooltip'),
        icon: IconUserOutlined
      },
      {
        label: t('targeting.custom-rule'),
        value: RuleCategory.CUSTOM,
        tooltip: t('targeting.custom-rule-tooltip'),
        icon: IconSetting
      }
    ],
    []
  );

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        trigger={
          <div className="flex items-center gap-x-2 h-6 p-0 typo-para-medium text-primary-500">
            <Icon icon={IconPlus} size={'md'} />
            {t('table:feature-flags.add-rule')}
          </div>
        }
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
