import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import {
  IconPlus,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

const AddRule = () => {
  const { t } = useTranslation(['form', 'table']);

  const options = useMemo(
    () => [
      {
        label: t('feature-flags.prerequisites'),
        value: 'prerequisites',
        tooltip: t('targeting.prerequisite-tooltip'),
        icon: IconPrerequisite
      },
      {
        label: t('targeting.individual-targeting'),
        value: 'individual',
        tooltip: t('targeting.individual-tooltip'),
        icon: IconUserOutlined
      },
      {
        label: t('targeting.custom-rule'),
        value: 'custom',
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
          <Button variant="text" className="h-6 p-0">
            <Icon icon={IconPlus} size={'md'} />
            {t('table:feature-flags.add-rule')}
          </Button>
        }
        showArrow={false}
        className="w-full [&>div]:flex-center border-dashed !shadow-none"
      />
      <DropdownMenuContent>
        {options.map((item, index) => (
          <Tooltip
            side="right"
            align="start"
            className="max-w-[172px] bg-white text-gray-600 shadow-card"
            key={index}
            content={item.tooltip}
            trigger={
              <DropdownMenuItem
                icon={item.icon}
                label={item.label}
                value={item.value}
              />
            }
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default AddRule;
