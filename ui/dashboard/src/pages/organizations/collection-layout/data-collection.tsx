import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import Icon from 'components/icon';

export const useColumns = (): ColumnDef<Organization>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <Link
            to={`${organization.id}`}
            className="underline text-primary-500 typo-para-medium"
          >
            {organization.name}
          </Link>
        );
      }
    },
    {
      accessorKey: 'projectCount',
      header: `${t('projects')}`,
      size: 170,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {organization.projectCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'environmentCount',
      header: `${t('environments')}`,
      size: 170,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {organization.environmentCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'userCount',
      header: `${t('users')}`,
      size: 160,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {organization.userCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 180,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(organization.createdAt)}
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
