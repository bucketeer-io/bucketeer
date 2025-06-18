import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { hasEditable, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { APIKey, APIKeyRole } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { APIKeyActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: APIKey, type: APIKeyActionsType) => void;
}): ColumnDef<APIKey>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();

  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);

  const getAPIkeyRole = (role: APIKeyRole) => {
    switch (role) {
      case 'SDK_CLIENT':
        return t('client');

      case 'SDK_SERVER':
        return t('server');

      default:
        return t('public-api');
    }
  };

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
    });
  };

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const apiKey = row.original;
        const { id, name } = apiKey;
        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={id}
                  name={name}
                  maxLines={1}
                  className="min-w-[300px]"
                  onClick={() => onActions(apiKey, 'EDIT')}
                />
              }
              maxLines={1}
            />

            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {truncateTextCenter(id)}
              <div onClick={() => handleCopyId(id)}>
                <Icon
                  icon={IconCopy}
                  size={'sm'}
                  className="opacity-0 group-hover:opacity-100 cursor-pointer"
                />
              </div>
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'role',
      header: `${t('role')}`,
      size: 150,
      cell: ({ row }) => {
        const apiKey = row.original;
        return (
          <div className="typo-para-small text-accent-blue-500 bg-accent-blue-50 px-2 py-[3px] w-fit rounded">
            {getAPIkeyRole(apiKey.role)}
          </div>
        );
      }
    },
    {
      accessorKey: 'environment',
      header: `${t('environment')}`,
      size: 250,
      maxSize: 250,
      cell: ({ row }) => {
        const apiKey = row.original;
        const id = `env-${apiKey.id}`;
        return (
          <NameWithTooltip
            id={id}
            align="center"
            content={
              <NameWithTooltip.Content
                content={apiKey.environmentName}
                id={id}
              />
            }
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={apiKey.environmentName}
                maxLines={1}
                haveAction={false}
              />
            }
            maxLines={1}
          />
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
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatDateTime(apiKey.createdAt)}
              </div>
            }
            date={apiKey.createdAt}
          />
        );
      }
    },

    {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 120,
      cell: ({ row }) => {
        const apiKey = row.original;

        return (
          <Switch
            disabled={!editable}
            checked={!apiKey.disabled}
            onCheckedChange={value =>
              onActions(apiKey, value ? 'ENABLE' : 'DISABLE')
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
