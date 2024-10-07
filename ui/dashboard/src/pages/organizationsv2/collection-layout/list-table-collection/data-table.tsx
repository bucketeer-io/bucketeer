import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { Organization } from '@types';
import { useFormatDateTime } from 'utils/date-time';

export const useColumns = (): ColumnDef<Organization>[] => {
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: 'Name',
      cell: ({ row }) => {
        const organization = row.original;
        return (
          <Link
            to={`/organization/${organization.id}`}
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
      header: '',
      meta: {
        align: 'center',
        style: { textAlign: 'center', fitContent: true }
      },
      enableSorting: false,
      cell: () => {
        return (
          <div className="px-4">
            <IconMoreHorizOutlined />
          </div>
        );
      }
    }
  ];
};
