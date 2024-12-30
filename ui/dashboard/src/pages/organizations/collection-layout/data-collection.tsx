import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import { OrganizationActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: Organization, type: OrganizationActionsType) => void;
}): ColumnDef<Organization>[] => {
  const { searchOptions } = useSearchParams();
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
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatDateTime(organization.createdAt)}
              </div>
            }
            date={organization.createdAt}
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
        const organization = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-org')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-org')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-org')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE'
                  }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={value =>
              onActions(organization, value as OrganizationActionsType)
            }
            align="end"
          />
        );
      }
    }
  ];
};
