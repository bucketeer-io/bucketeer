import { IconMoreHorizOutlined } from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import Icon from 'components/icon';

export const useColumns = (): ColumnDef<Account>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 300,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="underline text-primary-500 typo-para-medium">
            {account.name}
          </div>
        );
      }
    },
    {
      accessorKey: 'organizationRole',
      header: `${t('role')}`,
      size: 350,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {account.environmentRoles[0].role?.split('_')[1]}
          </div>
        );
      }
    },
    {
      accessorKey: 'environmentCount',
      header: `${t('environments')}`,
      size: 120,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {account.environmentRoles.length}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 180,
      cell: ({ row }) => {
        const account = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(account.createdAt)}
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
      cell: () => {
        return (
          <button className="flex-center text-gray-600">
            <Icon icon={IconMoreHorizOutlined} size="sm" />
          </button>
        );
      }
    }
  ];
};
