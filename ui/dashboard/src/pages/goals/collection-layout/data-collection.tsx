import {
  IconArchiveOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
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
import { cn } from 'utils/style';
import { IconCopy } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import { GoalActions } from '../types';

const Tag = ({
  tag,
  type,
  onClick
}: {
  tag: string;
  type: ConnectionType;
  onClick?: () => void;
}) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] text-center rounded-[3px] capitalize cursor-pointer',
        {
          'px-[19.5px] text-gray-600 bg-gray-100 cursor-default':
            type === 'UNKNOWN',
          'text-primary-500 bg-primary-50': type === 'EXPERIMENT',
          'text-accent-pink-500 bg-accent-pink-50': type === 'OPERATION'
        }
      )}
      onClick={onClick}
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

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <Link
              to={`/${currentEnvironment.urlCode}${PAGE_PATH_GOALS}/${goal.id}`}
              className="underline text-primary-500 break-all typo-para-medium text-left truncate"
            >
              {goal.name}
            </Link>
            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {truncateTextCenter(goal.id)}
              <div onClick={() => handleCopyId(goal.id)}>
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
      accessorKey: 'connections',
      header: `${t('table:goals.connections')}`,
      size: 150,
      cell: ({ row }) => {
        const goal = row.original;
        const experimentLength = goal.experiments?.length;
        const { connectionType } = goal;

        if (!goal.isInUseStatus || !experimentLength)
          return <Tag tag={'not in use'} type="UNKNOWN" />;
        return (
          <Tag
            tag={`${experimentLength} ${connectionType === 'EXPERIMENT' ? 'Experiment' : 'Operation'}${experimentLength > 1 ? 's' : ''}`}
            type={connectionType}
            onClick={() => experimentLength && onActions(goal, 'CONNECTION')}
          />
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

        return (
          <Popover
            options={compact([
              {
                label: `${t('archive-goal')}`,
                icon: IconArchiveOutlined,
                value: 'ARCHIVE'
              }
            ])}
            icon={IconMoreHorizOutlined}
            onClick={value => onActions(goal, value as GoalActions)}
            align="end"
          />
        );
      }
    }
  ];
};
