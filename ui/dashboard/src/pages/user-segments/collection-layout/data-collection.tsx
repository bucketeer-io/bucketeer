import {
  IconDeleteOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { Popover } from 'components/popover';
import { UserSegmentsActionsType } from '../types';

export const useColumns = ({
  onActionHandler
}: {
  onActionHandler: (value: UserSegment, type: UserSegmentsActionsType) => void;
}): ColumnDef<UserSegment>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 250,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <Link
            to={`/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}/${segment.id}`}
            className="underline text-primary-500 typo-para-medium"
          >
            {segment.name}
          </Link>
        );
      }
    },
    {
      accessorKey: 'includedUserCount',
      header: `${t('users')}`,
      size: 350,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {segment.includedUserCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'connections',
      header: `${t('connections')}`,
      size: 120,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div
            className="flex-center w-fit px-2 py-1.5 rounded bg-primary-50 text-primary-500 typo-para-medium cursor-pointer"
            onClick={() =>
              segment?.features?.length && onActionHandler(segment, 'FLAG')
            }
          >
            {segment?.features?.length}
            {` ${segment?.features?.length === 1 ? 'Flag' : 'Flags'}`}
          </div>
        );
      }
    },
    {
      accessorKey: 'updatedAt',
      header: t('table:updated-at'),
      size: 100,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {Number(segment.updatedAt) === 0
              ? t('never')
              : formatDateTime(segment.updatedAt)}
          </div>
        );
      }
    },
    {
      accessorKey: 'status',
      header: `${t('status')}`,
      size: 150,
      cell: ({ row }) => {
        const segment = row.original;
        return (
          <div
            className={cn(
              'typo-para-small text-accent-green-500 bg-accent-green-50 px-2 py-[3px] w-fit rounded',
              {
                'bg-gray-200 text-gray-600': !segment.isInUseStatus
              }
            )}
          >
            {segment.isInUseStatus ? 'In Use' : 'Not In Use'}
          </div>
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
        const segment = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-segment')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.delete-segment')}`,
                icon: IconDeleteOutlined,
                value: 'DELETE'
              }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={value =>
              onActionHandler(segment, value as UserSegmentsActionsType)
            }
            align="end"
          />
        );
      }
    }
  ];
};
