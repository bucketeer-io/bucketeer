import { useCallback } from 'react';
import { IconEditOutlined } from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useAuthAccess } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { APIKey, APIKeyRole } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import Icon from 'components/icon';
import Switch from 'components/switch';
import { Tooltip } from 'components/tooltip';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { APIKeyActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: APIKey, type: APIKeyActionsType) => void;
}): ColumnDef<APIKey>[] => {
  const { t } = useTranslation(['common', 'table', 'message']);
  const formatDateTime = useFormatDateTime();
  const { notify } = useToast();

  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const getAPIKeyRole = useCallback((role: APIKeyRole) => {
    let roleKey = '';
    let roleTooltipKey = '';
    switch (role) {
      case 'SDK_CLIENT':
        roleKey = 'client-sdk';
        roleTooltipKey = 'client-sdk-desc';
        break;
      case 'SDK_SERVER':
        roleKey = 'server-sdk';
        roleTooltipKey = 'server-sdk-desc';
        break;
      case 'PUBLIC_API_READ_ONLY':
        roleKey = 'public-api';
        roleTooltipKey = 'public-read-only-desc';
        break;
      case 'PUBLIC_API_WRITE':
        roleKey = 'public-api';
        roleTooltipKey = 'public-read-write-desc';
        break;
      case 'PUBLIC_API_ADMIN':
        roleKey = 'public-api';
        roleTooltipKey = 'public-admin-desc';
        break;
      case 'UNKNOWN':
      default:
        roleKey = 'unknown';
        roleTooltipKey = 'unknown';
        break;
    }
    return {
      role: t(
        role === 'UNKNOWN' ? 'form:unknown' : `table:api-keys.${roleKey}`
      ),
      roleTooltipContent: t(
        role === 'UNKNOWN' ? 'form:unknown' : `table:api-keys.${roleTooltipKey}`
      )
    };
  }, []);

  const handleCopyId = useCallback((id: string) => {
    copyToClipBoard(id);
    notify({
      message: t('message:copied')
    });
  }, []);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
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
        const { role, roleTooltipContent } = getAPIKeyRole(apiKey.role);
        return (
          <Tooltip
            content={roleTooltipContent}
            trigger={
              <div className="typo-para-small text-accent-blue-500 bg-accent-blue-50 px-2 py-[3px] w-fit rounded whitespace-nowrap">
                {role}
              </div>
            }
            className="max-w-[300px]"
          />
        );
      }
    },
    {
      accessorKey: 'environment',
      header: `${t('environment')}`,
      size: 350,
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
                className="!max-w-[300px]"
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
          <DisabledButtonTooltip
            align="center"
            type={!isOrganizationAdmin ? 'admin' : 'editor'}
            hidden={envEditable && isOrganizationAdmin}
            trigger={
              <div className="w-fit">
                <Switch
                  disabled={!envEditable || !isOrganizationAdmin}
                  checked={!apiKey.disabled}
                  onCheckedChange={value =>
                    onActions(apiKey, value ? 'ENABLE' : 'DISABLE')
                  }
                />
              </div>
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
          <DisabledPopoverTooltip
            isNeedAdminAccess
            options={compact([
              {
                label: `${t('table:popover.edit-api-key')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              }
            ])}
            onClick={value => onActions(apiKey, value as APIKeyActionsType)}
          />
        );
      }
    }
  ];
};
