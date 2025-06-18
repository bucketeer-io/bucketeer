import { IconArchiveOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_GOALS } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { ConnectionType, Goal } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconArrowDown, IconCopy, IconTrash } from '@icons';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { GoalActions } from '../types';

const Tag = ({ tag, type }: { tag: string; type: ConnectionType }) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] text-center rounded-[3px] capitalize whitespace-nowrap',
        {
          'px-4 text-gray-600 bg-gray-100': type === 'UNKNOWN',
          'text-primary-500 bg-primary-50': type === 'EXPERIMENT',
          'text-accent-pink-500 bg-accent-pink-50': type === 'OPERATION'
        }
      )}
    >
      {tag}
    </div>
  );
};

export const useColumns = ({
  onActions
}: {
  onActions: (item: Goal, type: GoalActions) => void;
}): ColumnDef<Goal>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
    });
  };

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const goal = row.original;
        const { id, name } = goal;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}/${id}`}
                >
                  <NameWithTooltip.Trigger id={id} name={name} />
                </Link>
              }
            />
            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {truncateTextCenter(id)}
              <div onClick={() => handleCopyId(id)}>
                <Icon
                  icon={IconCopy}
                  size={'sm'}
                  className="opacity-0 group-hover:opacity-100 cursor-pointer"
                />
              </div>
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'connectionType',
      header: `${t('table:goals.connections')}`,
      size: 150,
      cell: ({ row }) => {
        const goal = row.original;

        const connectionCount =
          goal.connectionType === 'EXPERIMENT'
            ? goal.experiments?.length
            : goal?.autoOpsRules?.length;

        if (!goal.isInUseStatus && goal.connectionType === 'UNKNOWN')
          return <Tag tag={'not in use'} type="UNKNOWN" />;
        return (
          <button
            disabled={!connectionCount}
            onClick={() => connectionCount && onActions(goal, 'CONNECTION')}
            className="flex items-center gap-1"
          >
            <Tag
              tag={
                goal.connectionType === 'EXPERIMENT'
                  ? t('form:experiment', { count: connectionCount })
                  : t('form:operation', { count: connectionCount })
              }
              type={goal.connectionType}
            />
            {connectionCount > 0 && <Icon icon={IconArrowDown} />}
          </button>
        );
      }
    },
    {
      accessorKey: 'updatedAt',
      header: `${t('table:updated-at')}`,
      size: 150,
      cell: ({ row }) => {
        const goal = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {Number(goal.updatedAt) === 0
                  ? t('never')
                  : formatDateTime(goal.updatedAt)}
              </div>
            }
            date={goal.updatedAt}
          />
        );
      }
    },
    {
      accessorKey: 'action',
      size: 60,
      header: '',
      meta: {
        align: 'center',
        style: { textAlign: 'center', fitContent: true }
      },
      enableSorting: false,
      cell: ({ row }) => {
        const goal = row.original;
        const { isInUseStatus } = goal;
        return (
          <DisabledPopoverTooltip
            options={compact([
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-goal')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-goal')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE',
                    disabled: isInUseStatus,
                    tooltip: isInUseStatus
                      ? t('form:goal-details.archive-warning-desc')
                      : ''
                  },
              {
                label: `${t('table:popover.delete-goal')}`,
                icon: IconTrash,
                value: 'DELETE',
                disabled: isInUseStatus,
                tooltip: isInUseStatus
                  ? t('form:goal-details.delete-warning-desc')
                  : ''
              }
            ])}
            onClick={value => onActions(goal, value as GoalActions)}
          />
        );
      }
    }
  ];
};
