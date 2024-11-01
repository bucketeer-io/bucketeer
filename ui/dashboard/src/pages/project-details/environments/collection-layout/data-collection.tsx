import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { Popover } from 'components/popover';

export const useColumns = ({
  onActionHandler
}: {
  onActionHandler: (type: string, v: Environment) => void;
}): ColumnDef<Environment>[] => {
  const { searchOptions } = useSearchParams();
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const environment = row.original;
        return (
          <div className="underline text-primary-500 typo-para-medium">
            {environment.name}
          </div>
        );
      }
    },
    {
      accessorKey: 'featureFlagCount',
      header: `${t('table:flags')}`,
      size: 250,
      cell: ({ row }) => {
        const environment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {environment.featureFlagCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 160,
      cell: ({ row }) => {
        const environment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(environment.createdAt)}
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
        const environment = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-env')}`,
                icon: IconEditOutlined,
                value: 'EDIT_ENVIRONMENT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE_ENVIRONMENT'
                  }
                : {
                    label: `${t('table:popover.archive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVED_ENVIRONMENT'
                  }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={value => onActionHandler(value as string, environment)}
            align="end"
          />
        );
      }
    }
  ];
};
