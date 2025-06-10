import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Notification } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
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
        const { id, name } = notification;
        return (
          <NameWithTooltip
            id={id}
            content={<NameWithTooltip.Content content={name} id={id} />}
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={name}
                onClick={() => onActions(notification, 'EDIT')}
                maxLines={1}
              />
            }
            maxLines={1}
          />
        );
      }
    },
    {
      accessorKey: 'environment',
      header: `${t('environment')}`,
      size: 228,
      cell: ({ row }) => {
        const notification = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {notification.environmentName}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 140,
      cell: ({ row }) => {
        const notification = row.original;
        const isNever = Number(notification.createdAt) === 0;

        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {isNever ? t('never') : formatDateTime(notification.createdAt)}
              </div>
            }
            date={isNever ? null : notification.createdAt}
          />
        );
      }
    },
    {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 76,
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
      size: 20,
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
