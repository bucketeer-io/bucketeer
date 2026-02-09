import { useCallback } from 'react';
import { IconEditOutlined } from 'react-icons-material-design';
import { useAuthAccess } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { APIKey, APIKeyRole } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { IconCopy } from '@icons';
import { APIKeyActionsType } from 'pages/api-keys/types';
import Divider from 'components/divider';
import Icon from 'components/icon';
import Switch from 'components/switch';
import { Tooltip } from 'components/tooltip';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface ApiCardProps {
  data: APIKey;
  onActions: (item: APIKey, type: APIKeyActionsType) => void;
}

export const ApiCard: React.FC<ApiCardProps> = ({ data, onActions }) => {
  const { t } = useTranslation(['common', 'table']);
  const { notify } = useToast();
  const formatDateTime = useFormatDateTime();
  const {
    name,
    id,
    apiKey,
    role,
    disabled,
    environmentName,
    lastUsedAt,
    createdAt
  } = data;
  const getAPIKeyRole = (role: APIKeyRole) => {
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
  };
  const { role: roleTooltip, roleTooltipContent } = getAPIKeyRole(role);
  const idEnv = `env-${id}`;
  const isNever = !lastUsedAt || lastUsedAt === '0';

  const { envEditable, isOrganizationAdmin } = useAuthAccess();
  const handleCopyId = useCallback((id: string) => {
    copyToClipBoard(id);
    notify({
      message: t('message:copied')
    });
  }, []);
  return (
    <Card>
      <Card.Header
        icon={
          <div className="flex items-center typo-para-tiny font-bold text-primary-500 gap-1">
            <DisabledButtonTooltip
              align="center"
              type={!isOrganizationAdmin ? 'admin' : 'editor'}
              hidden={envEditable && isOrganizationAdmin}
              trigger={
                <div className="w-fit">
                  <Switch
                    disabled={!envEditable || !isOrganizationAdmin}
                    checked={!disabled}
                    onCheckedChange={value =>
                      onActions(data, value ? 'ENABLE' : 'DISABLE')
                    }
                  />
                </div>
              }
            />
          </div>
        }
        triger={
          <div className="flex flex-col gap-0.5 max-w-fit">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={id}
                  name={name}
                  maxLines={1}
                  className="min-w-[150px]"
                  onClick={() => onActions(data, 'EDIT')}
                />
              }
              maxLines={1}
            />

            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {truncateTextCenter(apiKey)}
              <div onClick={() => handleCopyId(apiKey)}>
                <Icon
                  icon={IconCopy}
                  size={'sm'}
                  className="opacity-0 group-hover:opacity-100 cursor-pointer"
                />
              </div>
            </div>
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            isNeedAdminAccess
            options={compact([
              {
                label: `${t('table:popover.edit-api-key')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              }
            ])}
            onClick={value => onActions(data, value as APIKeyActionsType)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex flex-wrap h-full w-full items-stretch justify-between gap-3 pb-5 typo-para-medium">
          <div className="flex-1 p-3 rounded-xl bg-gray-100">
            <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
              <span>{t('role')}</span>
            </p>
            <div className="mt-2 flex items-center gap-2">
              <Tooltip
                content={roleTooltipContent}
                trigger={
                  <div className="typo-para-small text-accent-blue-500 py-[3px] w-fit whitespace-nowrap">
                    {roleTooltip}
                  </div>
                }
                className="max-w-[300px]"
              />
            </div>
          </div>

          <div className="flex-1 p-3 rounded-xl bg-gray-100 text-nowrap">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
                <span>{t('environment')}</span>
              </p>
              <div className="mt-2 flex items-center gap-2">
                <NameWithTooltip
                  id={idEnv}
                  align="center"
                  content={
                    <NameWithTooltip.Content
                      content={environmentName}
                      id={id}
                      className="!max-w-[300px]"
                    />
                  }
                  trigger={
                    <NameWithTooltip.Trigger
                      id={id}
                      name={environmentName}
                      maxLines={1}
                      haveAction={false}
                    />
                  }
                  maxLines={1}
                />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <div className="flex flex-wrap h-full w-full pt-5 items-stretch justify-between gap-3 typo-para-medium">
          <div className="flex-1 p-3 rounded-xl">
            <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
              <span>{t('table:last-used-at')}</span>
            </p>
            <div className="mt-2 text-nowrap">
              <DateTooltip
                trigger={
                  <div className="text-gray-500 typo-para-medium">
                    {isNever ? t('never') : formatDateTime(lastUsedAt)}
                  </div>
                }
                date={isNever ? null : lastUsedAt}
              />
            </div>
          </div>

          <div className="flex-1 p-3 rounded-xl">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
                <span>{t('table:created-at')}</span>
              </p>
              <div className="mt-2 text-nowrap">
                {' '}
                <DateTooltip
                  trigger={
                    <div className="text-gray-500 typo-para-medium">
                      {formatDateTime(createdAt)}
                    </div>
                  }
                  date={createdAt}
                />
              </div>
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
