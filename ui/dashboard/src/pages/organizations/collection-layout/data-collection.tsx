import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { Organization } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import Icon from 'components/icon';

export const useColumns = (): ColumnDef<Organization>[] => {
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: 'Name',
      size: 400,
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <Link
            to={`/organizations/${organization.id}`}
            className="underline text-primary-500 typo-para-medium"
          >
            {organization.name}
          </Link>
        );
      }
    },
    {
      accessorKey: 'projectCount',
      header: 'Projects',
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
      header: 'Environments',
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
      header: 'Users',
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
      header: 'Created at',
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
