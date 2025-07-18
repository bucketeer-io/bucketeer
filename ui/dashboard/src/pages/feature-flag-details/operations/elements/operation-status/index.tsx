import { FunctionComponent, useCallback, useMemo } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { IconEditOutlined } from 'react-icons-material-design';
import { hasEditable, useAuth } from 'auth';
import { cva } from 'class-variance-authority';
import { Color, RolloutStoppedBy } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import {
  IconCalendar,
  IconDisable,
  IconFlagOperation,
  IconOperationDetails,
  IconStoppedByUser,
  IconTrash,
  IconWatch
} from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import { PopoverOption, PopoverValue } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import { OperationModalState } from '../..';
import {
  OperationActionType,
  OperationCombinedType,
  OpsTypeMap
} from '../../types';

const statusVariants = cva(
  'flex-center px-2 py-1.5 rounded-[3px] typo-para-small',
  {
    variants: {
      status: {
        [OpsTypeMap.SCHEDULE]: 'bg-accent-green-50 text-accent-green-800',
        [OpsTypeMap.EVENT_RATE]: 'bg-primary-50 text-primary-700',
        [OpsTypeMap.ROLLOUT]: 'bg-accent-orange-50 text-accent-orange-600'
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

interface StoppedByData {
  icon: FunctionComponent;
  textKey: string;
  iconColor: Color;
}

const getStoppedByData = (stoppedBy: RolloutStoppedBy): StoppedByData => {
  switch (stoppedBy) {
    case 'USER':
      return {
        icon: IconStoppedByUser,
        textKey: 'user',
        iconColor: 'accent-blue-500'
      };

    case 'OPS_SCHEDULE':
      return {
        icon: IconCalendar,
        textKey: 'schedule',
        iconColor: 'primary-500'
      };
    case 'OPS_KILL_SWITCH':
      return {
        icon: IconFlagOperation,
        textKey: 'event-rate',
        iconColor: 'accent-pink-500'
      };
    case 'UNKNOWN':
    default:
      return {} as StoppedByData;
  }
};

const StoppedBy = ({ stoppedBy }: { stoppedBy: RolloutStoppedBy }) => {
  const { t } = useTranslation(['form']);
  const isUnknown = useMemo(() => stoppedBy === 'UNKNOWN', [stoppedBy]);
  if (isUnknown) return <></>;
  const { icon, iconColor, textKey } = getStoppedByData(stoppedBy);
  return (
    <>
      <Icon icon={icon} size={'xxs'} color={iconColor} />
      <p className="typo-para-small text-gray-700">
        {t(`feature-flags.${textKey}`)}
      </p>
    </>
  );
};

const OperationStatus = ({
  isFinished,
  operation,
  onActions
}: {
  isFinished: boolean;
  operation: OperationCombinedType;
  onActions: (data: OperationModalState) => void;
}) => {
  const { t } = useTranslation(['form']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);

  const isRollout = useMemo(
    () => ['MANUAL_SCHEDULE', 'TEMPLATE_SCHEDULE'].includes(operation.type),
    [operation.type]
  );

  const operationType = useMemo(
    () => (isRollout ? 'ROLLOUT' : operation.opsType) as OpsTypeMap,
    [operation, isRollout]
  );

  const isKillSwitch = useMemo(
    () => operation.opsType === OpsTypeMap.EVENT_RATE,
    [operation]
  );

  const isStopped = useMemo(() => operation.status === 'STOPPED', [operation]);
  const stoppedAt = useMemo(() => {
    const isNever = Number(operation.stoppedAt) === 0;
    return isNever ? null : operation.stoppedAt;
  }, [operation]);

  const completedOptions: PopoverOption<PopoverValue>[] = useMemo(() => {
    return [
      {
        label: t('feature-flags.operation-details'),
        icon: IconOperationDetails,
        value: 'DETAILS'
      },
      {
        label: (
          <p className={editable ? 'text-accent-red-500' : ''}>
            {t(
              `feature-flags.delete-${isKillSwitch ? 'kill-switch' : isRollout ? 'rollout' : 'schedule'}`
            )}
          </p>
        ),
        icon: IconTrash,
        value: 'DELETE',
        color: editable ? 'accent-red-500' : undefined
      }
    ];
  }, [isRollout, isKillSwitch, editable]);

  const operationOptions: PopoverOption<PopoverValue>[] = useMemo(() => {
    const translationKey = isRollout
      ? 'rollout'
      : isKillSwitch
        ? 'kill-switch'
        : 'schedule';
    return [
      ...(isRollout
        ? []
        : [
            {
              label: t(`feature-flags.edit-${translationKey}`),
              icon: IconEditOutlined as FunctionComponent,
              value: 'UPDATE',
              type: isKillSwitch
            }
          ]),
      {
        label: t(`feature-flags.stop-${translationKey}`),
        icon: IconDisable,
        value: 'STOP'
      },
      {
        label: (
          <p className={editable ? 'text-accent-red-500' : ''}>
            {t(`feature-flags.delete-${translationKey}`)}
          </p>
        ),
        icon: IconTrash,
        value: 'DELETE',
        color: editable ? 'accent-red-500' : undefined
      }
    ];
  }, [isKillSwitch, isRollout, editable]);

  const popoverOptions = useMemo(
    () => (isFinished ? completedOptions : operationOptions),
    [isFinished, isKillSwitch, operationOptions, completedOptions, editable]
  );

  const titleKey = useMemo(() => {
    if (isRollout) return 'enable-operation';
    if (isKillSwitch) return 'kill-switch-operation';
    return 'schedule-operation';
  }, [isRollout, isKillSwitch]);

  return (
    <div className="flex flex-col w-full gap-y-4">
      <div className="flex items-center w-full justify-between gap-x-4">
        <p className="typo-head-bold-big text-gray-700">
          {t(`feature-flags.${titleKey}`)}
        </p>
        <div className="flex items-center gap-x-4">
          <Status status={operationType} />
          <DisabledPopoverTooltip
            options={popoverOptions}
            onClick={value =>
              onActions({
                actionType: value as OperationActionType,
                operationType,
                selectedData: operation
              })
            }
          />
        </div>
      </div>
      <Divider />
      <div className="flex items-center w-full justify-between gap-x-4">
        <p className="flex items-center typo-head-bold-medium text-gray-700 h-9">
          {t('feature-flags.progress-information')}
        </p>
        {isStopped && (
          <div className="flex items-center gap-x-1.5">
            <DateTooltip
              trigger={
                <div className="flex items-center typo-para-small text-gray-700 gap-x-1.5">
                  <Icon icon={IconWatch} size={'xxs'} />
                  <Trans
                    i18nKey={'form:feature-flags.stopped-at'}
                    values={{
                      stoppedAt: formatDateTime(operation.stoppedAt)
                    }}
                  />
                </div>
              }
              date={stoppedAt}
            />
            <StoppedBy stoppedBy={operation.stoppedBy} />
          </div>
        )}
      </div>
    </div>
  );
};

export default OperationStatus;
