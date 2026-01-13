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
import Dropdown from 'components/dropdown';
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
    <Dropdown
      trigger={
        <DisabledButtonTooltip
          hidden={!disabled}
          trigger={
            <div
              className={cn(
                'inline-flex animate-fade gap-2 items-center justify-center duration-300 ease-out whitespace-nowrap w-[215px] h-12',
                'bg-primary-500 hover:bg-primary-700 text-gray-50',
                'rounded-lg px-3 sm:px-6 py-2',
                {
                  'cursor-not-allowed bg-primary-200 hover:bg-primary-200 text-primary-50':
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
      wrapTriggerStyle="w-fit"
      className="!shadow-none !border-none [&_p]:!text-white p-0"
      showArrow={false}
      disabled={disabled}
      options={options}
      onChange={value =>
        onOperationActions({
          operationType: value as OpsTypeMap,
          actionType: 'NEW'
        })
      }
    />
  );
};

export default OperationActions;
