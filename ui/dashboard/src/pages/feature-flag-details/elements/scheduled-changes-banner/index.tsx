import { useState } from 'react';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import {
  useGetScheduledFlagChangeSummary,
  useQueryScheduledFlagChanges
} from '@queries/scheduled-flag-changes';
import { useTranslation } from 'i18n';
import {
  ScheduledChangeCategories,
  ScheduledFlagChange,
  ScheduledFlagChangeStatuses
} from '@types';
import { cn } from 'utils/style';
import { IconCalendar, IconChecked, IconChevronDown, IconClose, IconInfoFilled } from '@icons';
import Icon from 'components/icon';
import { Popover, PopoverOption, PopoverValue } from 'components/popover';
import ApplyNowDialog from './apply-now-dialog';
import CancelScheduleDialog from './cancel-schedule-dialog';
import EditScheduleDialog from './edit-schedule-dialog';

const formatScheduledDate = (timestamp: string | number): string => {
  const date = new Date(Number(timestamp) * 1000);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}/${month}/${day} ${hours}:${minutes}`;
};

const BANNER_STATUSES = [
  ScheduledFlagChangeStatuses.PENDING,
  ScheduledFlagChangeStatuses.CONFLICT
];

type ScheduleAction = 'APPLY_NOW' | 'EDIT_SCHEDULE' | 'CANCEL_SCHEDULE';

interface ScheduledChangesBannerProps {
  featureId: string;
  environmentId: string;
}

const CATEGORY_I18N_KEY: Record<string, string> = {
  [ScheduledChangeCategories.TARGETING]: 'targeting-changes',
  [ScheduledChangeCategories.VARIATIONS]: 'variation-changes',
  [ScheduledChangeCategories.SETTINGS]: 'settings-changes',
  [ScheduledChangeCategories.MIXED]: 'mixed-changes'
};

const ScheduleActionsMenu = ({
  schedule,
  onAction
}: {
  schedule: ScheduledFlagChange;
  onAction: (action: ScheduleAction, schedule: ScheduledFlagChange) => void;
}) => {
  const { t } = useTranslation(['form']);

  const options: PopoverOption<PopoverValue>[] = [
    {
      label: t('feature-flags.apply-now'),
      icon: IconChecked,
      value: 'APPLY_NOW'
    },
    {
      label: t('feature-flags.edit-schedule'),
      icon: IconCalendar,
      value: 'EDIT_SCHEDULE'
    },
    {
      label: t('feature-flags.cancel-schedule'),
      icon: IconClose,
      value: 'CANCEL_SCHEDULE',
      color: 'accent-red-500'
    }
  ];

  return (
    <Popover
      align="end"
      options={options}
      trigger={
        <span className="flex-center text-gray-500 cursor-pointer">
          <IconMoreHorizOutlined style={{ fontSize: 20 }} />
        </span>
      }
      onClick={value =>
        onAction(value as ScheduleAction, schedule)
      }
    />
  );
};

const ScheduleItemSummary = ({
  schedule,
  onAction
}: {
  schedule: ScheduledFlagChange;
  onAction: (action: ScheduleAction, schedule: ScheduledFlagChange) => void;
}) => {
  const { t } = useTranslation(['form']);

  const changeCount = schedule.changeSummaries?.length ?? 0;
  const categoryKey =
    CATEGORY_I18N_KEY[schedule.category] || 'mixed-changes';

  return (
    <li className="flex items-center justify-between w-full typo-para-small text-gray-700 gap-x-2">
      <div className="flex items-center gap-x-2 min-w-0">
        <span className="whitespace-nowrap text-accent-blue-500 font-medium">
          {formatScheduledDate(schedule.scheduledAt)}
        </span>
        <span className="text-gray-400">|</span>
        <span className="truncate">
          {t(`feature-flags.scheduled-banner-${categoryKey}`, {
            count: changeCount
          })}
        </span>
      </div>
      <div className="flex items-center gap-x-2 flex-shrink-0">
        <ScheduleActionsMenu schedule={schedule} onAction={onAction} />
      </div>
    </li>
  );
};

const ScheduledChangesBanner = ({
  featureId,
  environmentId
}: ScheduledChangesBannerProps) => {
  const { t } = useTranslation(['common', 'form']);
  const [isExpanded, setIsExpanded] = useState(false);
  const [activeDialog, setActiveDialog] = useState<{
    type: ScheduleAction;
    schedule: ScheduledFlagChange;
  } | null>(null);

  const { data: summaryData } = useGetScheduledFlagChangeSummary({
    params: { environmentId, featureId },
    enabled: !!featureId && !!environmentId
  });

  const pendingCount = summaryData?.summary?.pendingCount ?? 0;
  const hasMultiple = pendingCount > 1;

  const { data: listData } = useQueryScheduledFlagChanges({
    params: {
      environmentId,
      featureId,
      statuses: [...BANNER_STATUSES],
      orderBy: 'SCHEDULED_AT',
      orderDirection: 'ASC',
      pageSize: 10
    },
    enabled: !!featureId && !!environmentId && pendingCount > 0
  });

  if (pendingCount === 0) return null;

  const schedules = listData?.scheduledFlagChanges ?? [];
  const nextScheduledAt = summaryData?.summary?.nextScheduledAt;
  const firstSchedule = schedules[0];

  const handleAction = (
    action: ScheduleAction,
    schedule: ScheduledFlagChange
  ) => {
    setActiveDialog({ type: action, schedule });
  };

  const handleCloseDialog = () => setActiveDialog(null);

  return (
    <>
      <div className="flex flex-col w-full rounded border-l-4 p-4 border-accent-blue-500 bg-accent-blue-50">
        {hasMultiple ? (
          <button
            type="button"
            onClick={() => setIsExpanded(prev => !prev)}
            className="flex items-center justify-between w-full gap-x-4 cursor-pointer"
          >
            <div className="flex items-center gap-x-2 min-w-0">
              <Icon
                icon={IconInfoFilled}
                size="xxs"
                color="accent-blue-500"
              />
              <p className="typo-para-small leading-[14px] text-accent-blue-500">
                {t('form:feature-flags.scheduled-changes-count', {
                  count: pendingCount
                })}
              </p>
            </div>
            <Icon
              icon={IconChevronDown}
              size="xxs"
              color="accent-blue-500"
              className={cn(
                'transition-transform duration-200 flex-shrink-0',
                isExpanded && 'rotate-180'
              )}
            />
          </button>
        ) : (
          <div className="flex items-center justify-between gap-x-4">
            <div className="flex items-center gap-x-2 min-w-0">
              <Icon
                icon={IconInfoFilled}
                size="xxs"
                color="accent-blue-500"
              />
              <p className="typo-para-small leading-[14px] text-accent-blue-500">
                {t('form:feature-flags.scheduled-changes-single', {
                  datetime: nextScheduledAt
                    ? formatScheduledDate(nextScheduledAt)
                    : ''
                })}
              </p>
            </div>
            {firstSchedule && (
              <div className="flex items-center gap-x-2 flex-shrink-0">
                <ScheduleActionsMenu
                  schedule={firstSchedule}
                  onAction={handleAction}
                />
              </div>
            )}
          </div>
        )}

        {isExpanded && schedules.length > 0 && (
          <ul className="flex flex-col gap-y-2 pl-6 pt-3 w-full max-w-full">
            {schedules.map(schedule => (
              <ScheduleItemSummary
                key={schedule.id}
                schedule={schedule}
                onAction={handleAction}
              />
            ))}
          </ul>
        )}
      </div>

      {activeDialog?.type === 'EDIT_SCHEDULE' && (
        <EditScheduleDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
      {activeDialog?.type === 'CANCEL_SCHEDULE' && (
        <CancelScheduleDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
      {activeDialog?.type === 'APPLY_NOW' && (
        <ApplyNowDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
    </>
  );
};

export default ScheduledChangesBanner;
