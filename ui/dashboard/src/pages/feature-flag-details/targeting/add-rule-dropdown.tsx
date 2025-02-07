import { useTranslation } from 'i18n';
import {
  IconArrowDown,
  IconDocument,
  IconPlus,
  IconTargetSegments,
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

const AddRuleDropdown = () => {
  const { t } = useTranslation(['table']);

  const options = [
    {
      icon: IconTargetSegments,
      label: t('feature-flags.target-segments'),
      value: 'target-segments'
    },
    {
      icon: IconUserOutlined,
      label: t('feature-flags.target-individuals'),
      value: 'target-individuals'
    },
    {
      icon: IconDocument,
      label: t('feature-flags.set-prerequisites'),
      value: 'set-prerequisites'
    }
  ];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        trigger={
          <Button>
            <Icon
              icon={IconPlus}
              size={'sm'}
              className="flex-center text-white"
            />
            {t('feature-flags.add-rule')}
            <Icon
              icon={IconArrowDown}
              size={'sm'}
              className="flex-center text-white"
            />
          </Button>
        }
        showArrow={false}
        className="!border-none !shadow-none"
      />
      <DropdownMenuContent align="start" sideOffset={2}>
        {options.map((item, index) => (
          <DropdownMenuItem
            key={index}
            icon={item.icon}
            label={item.label}
            value={item.value}
            className="!text-gray-600"
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default AddRuleDropdown;
