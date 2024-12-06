import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Notification } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import { NotificationActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: Notification, type: NotificationActionsType) => void;
}): ColumnDef<Notification>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const notification = row.original;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <button
              onClick={() => onActions(notification, 'EDIT')}
              className="underline text-primary-500 break-all typo-para-medium text-left"
            >
              {notification.name}
            </button>
            <div className="typo-para-tiny text-gray-500">
              {truncateTextCenter(notification.name)}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 150,
      cell: ({ row }) => {
        const notification = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {formatDateTime(notification.createdAt)}
          </div>
        );
      }
    },
    {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 120,
      cell: ({ row }) => {
        const notification = row.original;

        return (
          <Switch
            checked={!notification.disabled}
            onCheckedChange={value =>
              onActions(notification, value ? 'ENABLE' : 'DISABLE')
            }
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
        const notification = row.original;

        return (
          <Popover
            options={compact([
              {
                label: `${t('table:popover.edit-notification')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              }
            ])}
            icon={IconMoreHorizOutlined}
            onClick={value =>
              onActions(notification, value as NotificationActionsType)
            }
            align="end"
          />
        );
      }
    }
  ];
};
