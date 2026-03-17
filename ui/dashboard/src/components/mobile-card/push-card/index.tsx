import {
  IconEditOutlined,
  IconMoreVertOutlined
} from 'react-icons-material-design';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { compact } from 'lodash';
import { Push, Tag } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { IconTrash } from '@icons';
import { PushActionsType } from 'pages/pushes/types';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';
import { Card } from '../elements';

interface PushCardProps {
  data: Push;
  tags: Tag[];
  onActions: (item: Push, type: PushActionsType) => void;
}

export const PushCard: React.FC<PushCardProps> = ({
  data,
  tags,
  onActions
}) => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);
  const envId = `env-${data.id}`;

  const formattedTags = data.tags?.map(
    item => tags.find(tag => tag.id === item)?.name || item
  );
  const { name, id, environmentName, disabled, createdAt } = data || {};

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
          <div className="flex items-center gap-0.5 max-w-fit min-w-[230px]">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={id}
                  name={name}
                  onClick={() => onActions(data, 'EDIT')}
                  maxLines={1}
                />
              }
              maxLines={1}
            />
          </div>
        }
      >
        <Card.Action>
          <DisabledPopoverTooltip
            icon={IconMoreVertOutlined}
            options={compact([
              {
                label: `${t('table:popover.edit-push')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              {
                label: `${t('table:popover.delete-push')}`,
                icon: IconTrash,
                value: 'DELETE'
              }
            ])}
            onClick={value => onActions(data, value as PushActionsType)}
          />
        </Card.Action>
      </Card.Header>

      <Card.Meta>
        <div className="flex flex-col gap-y-2 items-center justify-between pb-5">
          <ExpandableTag
            tags={formattedTags}
            rowId={data.id}
            className="!max-w-[250px] truncate"
          />
        </div>
        <div className="flex flex-wrap h-full w-full items-stretch justify-between gap-3 typo-para-medium">
          <div className="flex-1 p-3 rounded-xl bg-gray-100 text-nowrap">
            <div className="flex-1">
              <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
                <span>{t('environment')}</span>
              </p>
              <div className="mt-2 flex items-center gap-2">
                <NameWithTooltip
                  id={id}
                  align="center"
                  content={
                    <NameWithTooltip.Content
                      content={environmentName}
                      id={envId}
                    />
                  }
                  trigger={
                    <NameWithTooltip.Trigger
                      id={envId}
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
          <div className="flex-1 p-3 rounded-xl bg-gray-100">
            <p className="flex items-center gap-2 uppercase typo-para-tiny text-gray-500">
              <span>{t('table:created-at')}</span>
            </p>
            <div className="mt-2 flex items-center gap-2">
              <DateTooltip
                trigger={
                  <div className="text-gray-700 typo-para-medium">
                    {formatDateTime(createdAt)}
                  </div>
                }
                date={createdAt}
              />
            </div>
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};
