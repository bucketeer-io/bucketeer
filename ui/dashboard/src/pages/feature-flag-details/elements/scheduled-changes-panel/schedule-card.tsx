import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import {
  ScheduledChangeCategories,
  ScheduledFlagChange,
  ScheduledFlagChangeStatuses
} from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import {
  IconCalendar,
  IconChecked,
  IconClose,
  IconToastWarning,
  IconToastWarningDynamic,
  IconWatch
} from '@icons';
import Icon from 'components/icon';
import { Popover, PopoverOption, PopoverValue } from 'components/popover';
import { Tooltip } from 'components/tooltip';
import DateTooltip from 'elements/date-tooltip';

export type ScheduleCardAction =
  | 'APPLY_NOW'
  | 'EDIT_SCHEDULE'
  | 'CANCEL_SCHEDULE';

const formatScheduledDate = (timestamp: string | number): string => {
  const date = new Date(Number(timestamp) * 1000);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}/${month}/${day} ${hours}:${minutes}`;
};

const STATUS_STYLE: Record<
  string,
  { bg: string; text: string; labelKey: string }
> = {
  [ScheduledFlagChangeStatuses.PENDING]: {
    bg: 'bg-accent-blue-50',
    text: 'text-accent-blue-500',
    labelKey: 'status-pending'
  },
  [ScheduledFlagChangeStatuses.CONFLICT]: {
    bg: 'bg-accent-yellow-50',
    text: 'text-accent-yellow-500',
    labelKey: 'status-conflict'
  },
  [ScheduledFlagChangeStatuses.EXECUTED]: {
    bg: 'bg-accent-green-50',
    text: 'text-accent-green-500',
    labelKey: 'status-executed'
  },
  [ScheduledFlagChangeStatuses.FAILED]: {
    bg: 'bg-accent-red-50',
    text: 'text-accent-red-500',
    labelKey: 'status-failed'
  },
  [ScheduledFlagChangeStatuses.CANCELLED]: {
    bg: 'bg-gray-200',
    text: 'text-gray-600',
    labelKey: 'status-cancelled'
  }
};

const CATEGORY_LABEL_KEY: Record<string, string> = {
  [ScheduledChangeCategories.TARGETING]: 'category-targeting',
  [ScheduledChangeCategories.VARIATIONS]: 'category-variations',
  [ScheduledChangeCategories.SETTINGS]: 'category-settings',
  [ScheduledChangeCategories.MIXED]: 'category-mixed'
};

interface ScheduleCardProps {
  schedule: ScheduledFlagChange;
  onAction: (action: ScheduleCardAction, schedule: ScheduledFlagChange) => void;
}

const ScheduleCard = ({ schedule, onAction }: ScheduleCardProps) => {
  const { t } = useTranslation(['form', 'common', 'scheduled-changes']);
  const formatDateTime = useFormatDateTime();

  const statusStyle =
    STATUS_STYLE[schedule.status] ||
    STATUS_STYLE[ScheduledFlagChangeStatuses.PENDING];
  const categoryKey = CATEGORY_LABEL_KEY[schedule.category] || 'category-mixed';
  const isPending = schedule.status === ScheduledFlagChangeStatuses.PENDING;
  const isConflict = schedule.status === ScheduledFlagChangeStatuses.CONFLICT;
  const isFailed = schedule.status === ScheduledFlagChangeStatuses.FAILED;
  const isActionable = isPending || isConflict;

  const menuOptions: PopoverOption<PopoverValue>[] = [
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
    <div className="flex flex-col w-full rounded shadow-card bg-white overflow-hidden">
      <div className="flex items-center justify-between px-4 py-3">
        <div className="flex items-center gap-x-3">
          <div className="flex items-center gap-x-1.5">
            <Icon icon={IconCalendar} size="xs" color="primary-500" />
            <span className="typo-para-medium font-medium text-gray-800">
              {formatScheduledDate(schedule.scheduledAt)}
            </span>
          </div>
          <span
            className={cn(
              'px-2 py-[3px] rounded typo-para-small whitespace-nowrap',
              statusStyle.bg,
              statusStyle.text
            )}
          >
            {t(`feature-flags.${statusStyle.labelKey}`)}
          </span>
          <span className="px-2 py-[3px] rounded bg-gray-200 text-gray-600 typo-para-small whitespace-nowrap">
            {t(`feature-flags.${categoryKey}`)}
          </span>
        </div>
        {isActionable && (
          <Popover
            align="end"
            options={menuOptions}
            trigger={
              <span className="flex-center text-gray-500 cursor-pointer">
                <IconMoreHorizOutlined style={{ fontSize: 20 }} />
              </span>
            }
            onClick={value => onAction(value as ScheduleCardAction, schedule)}
          />
        )}
      </div>

      <div className="flex flex-col px-4 py-3 gap-y-2 border-t border-gray-100">
        {schedule.changeSummaries?.length > 0 ? (
          <ul className="flex flex-col gap-y-1.5">
            {schedule.changeSummaries.map((summary, index) => (
              <li
                key={index}
                className="flex items-start gap-x-2 typo-para-small text-gray-700"
              >
                <span className="text-gray-400 mt-0.5">•</span>
                <span>
                  {t(
                    `ScheduledChange.${summary.messageKey.replace('ScheduledChange.', '')}`,
                    { ...summary.values, ns: 'scheduled-changes' }
                  )}
                </span>
              </li>
            ))}
          </ul>
        ) : (
          <p className="typo-para-small text-gray-400 italic">
            {t('common:no-data')}
          </p>
        )}

        {schedule.comment && (
          <div className="flex items-start gap-x-2 pt-2 border-t border-gray-100">
            <Tooltip
              align="start"
              content={schedule.comment}
              trigger={
                <p className="typo-para-small text-gray-500 truncate max-w-full cursor-default">
                  {t('form:comment-for-update')}: {schedule.comment}
                </p>
              }
              className="max-w-[400px]"
            />
          </div>
        )}

        {isConflict && (
          <div className="flex items-start gap-x-2 px-4 py-3 rounded bg-accent-yellow-50 border-l-4 border-accent-yellow-500 mt-1">
            <Icon
              icon={IconToastWarning}
              size="xxs"
              color="accent-yellow-500"
              className="mt-0.5 flex-shrink-0"
            />
            <div className="flex flex-col gap-y-1">
              <p className="typo-para-medium text-accent-yellow-500">
                {t('form:feature-flags.conflict-warning')}
              </p>
              {schedule.conflicts?.map((conflict, index) => (
                <p
                  key={index}
                  className="typo-para-small text-accent-yellow-500"
                >
                  {conflict.description}
                </p>
              ))}
            </div>
          </div>
        )}

        {isFailed && schedule.failureReason && (
          <div className="flex items-start gap-x-2 px-4 py-3 rounded bg-accent-red-50 border-l-4 border-accent-red-500 mt-1">
            <Icon
              icon={IconToastWarningDynamic}
              size="xxs"
              color="accent-red-500"
              className="mt-0.5 flex-shrink-0"
            />
            <div className="flex flex-col gap-y-1">
              <p className="typo-para-medium text-accent-red-500">
                {t('form:feature-flags.failed-warning')}
              </p>
              <p className="typo-para-small text-accent-red-500">
                {schedule.failureReason}
              </p>
            </div>
          </div>
        )}
      </div>

      <div className="flex items-center justify-between px-4 py-2 border-t border-gray-100">
        <p className="typo-para-small text-gray-500 truncate mr-2">
          {schedule.createdBy}
        </p>
        <DateTooltip
          side="top"
          align="end"
          trigger={
            <div className="flex items-center gap-x-1.5 flex-shrink-0">
              <Icon icon={IconWatch} size="xxs" color="gray-500" />
              <p className="typo-para-small text-gray-500 whitespace-nowrap">
                {formatDateTime(schedule.createdAt)}
              </p>
            </div>
          }
          date={schedule.createdAt}
        />
      </div>
    </div>
  );
};

export default ScheduleCard;
