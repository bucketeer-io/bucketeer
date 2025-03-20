import { useTranslation } from 'i18n';
import {
  IconArrowDown,
  IconDocument,
  IconPlus,
  IconTargetSegments,
  IconUserOutlined
} from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import { RuleCategory } from './types';

interface Props {
  onAddRule: (type: RuleCategory) => void;
}

const AddRuleDropdown = ({ onAddRule }: Props) => {
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
          <div className="inline-flex animate-fade gap-2 items-center justify-center duration-300 ease-out whitespace-nowrap bg-primary-500 text-gray-50 rounded-lg px-6 py-2 hover:bg-primary-700 disabled:bg-primary-200 disabled:text-primary-50 h-12">
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
          </div>
        }
        showArrow={false}
        className="!border-none !shadow-none p-0"
      />
      <DropdownMenuContent align="start" sideOffset={2}>
        {options.map((item, index) => (
          <DropdownMenuItem
            key={index}
            icon={item.icon}
            label={item.label}
            value={item.value}
            className="!text-gray-600"
            onSelectOption={value => onAddRule(value as RuleCategory)}
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default AddRuleDropdown;
