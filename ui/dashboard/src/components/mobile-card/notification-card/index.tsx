import {
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { Notification } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { IconTrash } from '@icons';
import { NotificationActionsType } from 'pages/notifications/types';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface NotificationCardProps {
  data: Notification;
  onActions: (item: Notification, type: NotificationActionsType) => void;
}

export const NotificationCard: React.FC<NotificationCardProps> = ({
  data,
  onActions
}) => {
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { name, id, disabled, environmentName, createdAt } = data;
  const isNever = Number(createdAt) === 0;

  return (
    <Card>
      <Card.Header
        icon={
          <DisabledButtonTooltip
            align="center"
            hidden={editable}
            trigger={
              <div className="w-fit">
                <Switch
                  disabled={!editable}
                  checked={!disabled}
                  onCheckedChange={value =>
                    onActions(data, value ? 'ENABLE' : 'DISABLE')
                  }
                />
              </div>
            }
          />
        }
        triger={
          <NameWithTooltip
            id={id}
            content={<NameWithTooltip.Content content={name} id={id} />}
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={name}
                onClick={() => {
                  onActions(data, 'EDIT');
                }}
                maxLines={1}
                className="min-w-[230px]"
              />
            }
            maxLines={1}
          />
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            options={compact([
              {
                label: `${t('table:popover.edit-notification')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.delete-notification')}`,
                icon: IconTrash,
                value: 'DELETE'
              }
            ])}
            onClick={value => onActions(data, value as NotificationActionsType)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex flex-wrap h-full w-full items-stretch justify-between gap-3 typo-para-medium">
          <div className="flex-1 p-3 rounded-xl bg-gray-100 text-nowrap">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
                <span>{t('environment')}</span>
              </p>
              <div className="mt-2 flex items-center gap-2">
                <div className="text-gray-700 typo-para-medium">
                  {environmentName}
                </div>
              </div>
            </div>
          </div>
          <div className="flex-1 p-3 rounded-xl bg-gray-100">
            <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
              <span>{t('table:created-at')}</span>
            </p>
            <div className="mt-2 flex items-center gap-2">
              <DateTooltip
                trigger={
                  <div className="text-gray-700 typo-para-medium text-nowrap">
                    {isNever ? t('never') : formatDateTime(createdAt)}
                  </div>
                }
                date={isNever ? null : createdAt}
              />
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
