import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { APIKey } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import { APIKeyActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: APIKey, type: APIKeyActionsType) => void;
}): ColumnDef<APIKey>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const apiKey = row.original;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <button
              onClick={() => onActions(apiKey, 'EDIT')}
              className="underline text-primary-500 break-all typo-para-medium text-left"
            >
              {apiKey.name}
            </button>
            <div className="typo-para-tiny text-gray-500">
              {truncateTextCenter(apiKey.id)}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'type',
      header: `${t('type')}`,
      size: 150,
      cell: () => {
        return (
          <div className="typo-para-small text-accent-blue-500 bg-accent-blue-50 px-2 py-[3px] w-fit rounded">
            {`API Key`}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 150,
      cell: ({ row }) => {
        const apiKey = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(apiKey.createdAt)}
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
        const apiKey = row.original;

        return (
          <Popover
            options={compact([
              {
                label: `${t('table:popover.edit-api-key')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              }
            ])}
            icon={IconMoreHorizOutlined}
            onClick={value => onActions(apiKey, value as APIKeyActionsType)}
            align="end"
          />
        );
      }
    }
  ];
};
