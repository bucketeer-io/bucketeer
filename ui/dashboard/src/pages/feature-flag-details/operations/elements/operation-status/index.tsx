import { FunctionComponent, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { cva } from 'class-variance-authority';
import { IconDisable, IconOperationClone, IconOperationDetails } from '@icons';
import { Popover, PopoverOption, PopoverValue } from 'components/popover';
import { OpsTypeMap } from '../../types';

const statusVariants = cva(
  'flex-center px-2 py-1.5 rounded-[3px] typo-para-small',
  {
    variants: {
      status: {
        [OpsTypeMap.SCHEDULE]: 'bg-accent-green-50 text-accent-green-800',
        [OpsTypeMap.EVENT_RATE]: 'bg-primary-50 text-primary-700',
        [OpsTypeMap.ROLLOUT]: 'bg-accent-orange-50 text-orange-600'
      }
    }
  }
);

const Status = ({ status }: { status: OpsTypeMap }) => {
  const { t } = useTranslation(['form']);

  const getStatusText = useCallback(
    (status: OpsTypeMap) => {
      switch (status) {
        case OpsTypeMap.SCHEDULE:
          return 'schedule';
        case OpsTypeMap.EVENT_RATE:
          return 'event-rate';
        case OpsTypeMap.ROLLOUT:
          return 'progressive-rollout';
        default:
          break;
      }
    },
    [status]
  );

  return (
    <div className={statusVariants({ status })}>
      {t(`feature-flags.${getStatusText(status)}`)}
    </div>
  );
};

const OperationStatus = ({
  isCompleted,
  isKillSwitch,
  title,
  status,
  onActions
}: {
  title: string;
  status: OpsTypeMap;
  isCompleted?: boolean;
  isKillSwitch?: boolean;
  onActions: () => void;
}) => {
  const { t } = useTranslation(['form']);

  const completedOptions: PopoverOption<PopoverValue>[] = useMemo(
    () => [
      {
        label: t('feature-flags.see-details'),
        icon: IconOperationDetails,
        value: 'DETAILS'
      },
      {
        label: t('feature-flags.clone-operation'),
        icon: IconOperationClone,
        value: 'CLONE'
      }
    ],
    []
  );

  const operationOptions: PopoverOption<PopoverValue>[] = useMemo(() => {
    const translationKey = isKillSwitch ? 'kill-switch' : 'operation';
    return [
      {
        label: t(`feature-flags.edit-${translationKey}`),
        icon: IconEditOutlined as FunctionComponent,
        value: 'EDIT'
      },
      {
        label: t(`feature-flags.stop-${translationKey}`),
        icon: IconDisable,
        value: 'STOP'
      },
      {
        label: (
          <p className="text-accent-red-500">
            {t(`feature-flags.delete-${translationKey}`)}
          </p>
        ),
        icon: IconDisable,
        value: 'DELETE',
        color: 'accent-red-500'
      }
    ];
  }, [isKillSwitch]);

  const popoverOptions = useMemo(
    () => (isCompleted ? completedOptions : operationOptions),
    [isCompleted, isKillSwitch, operationOptions, completedOptions]
  );

  return (
    <div className="flex items-center w-full justify-between pb-4 border-b border-gray-200">
      <p className="typo-head-bold-big text-gray-700">{title}</p>
      <div className="flex items-center gap-x-4">
        <Status status={status} />
      </div>
      <Popover
        options={popoverOptions}
        icon={IconMoreHorizOutlined}
        onClick={onActions}
        align="end"
      />
    </div>
  );
};

export default OperationStatus;
