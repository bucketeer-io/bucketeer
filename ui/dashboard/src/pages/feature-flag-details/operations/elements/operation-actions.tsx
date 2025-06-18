import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import {
  IconArrowDown,
  IconCalendar,
  IconFlagOperation,
  IconOperationArrow,
  IconPlus
} from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import { OperationModalState } from '..';
import { OpsTypeMap } from '../types';

const OperationActions = ({
  disabled,
  onOperationActions
}: {
  disabled: boolean;
  onOperationActions: (data: OperationModalState) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const options = useMemo(
    () => [
      {
        label: t('form:feature-flags.schedule'),
        value: OpsTypeMap.SCHEDULE,
        icon: IconCalendar
      },
      {
        label: t('form:feature-flags.event-rate'),
        value: OpsTypeMap.EVENT_RATE,
        icon: IconFlagOperation
      },
      {
        label: t('form:feature-flags.progressive-rollout'),
        value: OpsTypeMap.ROLLOUT,
        icon: IconOperationArrow
      }
    ],
    []
  );

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        trigger={
          <DisabledButtonTooltip
            hidden={!disabled}
            trigger={
              <div
                className={cn(
                  'inline-flex animate-fade gap-2 items-center justify-center duration-300 ease-out whitespace-nowrap w-[215px] h-12',
                  'bg-primary-500 hover:bg-primary-700 text-gray-50',
                  'rounded-lg px-6 py-2',
                  {
                    'pointer-events-none bg-primary-200 text-primary-50':
                      disabled
                  }
                )}
              >
                <Icon icon={IconPlus} size={'sm'} color="gray-100" />
                <p>{t('new-operation')}</p>
                <Icon icon={IconArrowDown} size={'sm'} color="gray-100" />
              </div>
            }
          />
        }
        className="!shadow-none !border-none [&_p]:!text-white p-0"
        showArrow={false}
        disabled={disabled}
      />
      <DropdownMenuContent sideOffset={0} className="w-[215px]">
        {options.map((item, index) => (
          <DropdownMenuItem
            key={index}
            icon={item.icon}
            value={item.value}
            label={item.label}
            onSelectOption={value =>
              onOperationActions({
                operationType: value as OpsTypeMap,
                actionType: 'NEW'
              })
            }
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default OperationActions;
