import { IconEditOutlined } from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { hasEditable, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Push, Tag } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { IconTrash } from '@icons';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import ExpandableTag from 'elements/expandable-tag';
import NameWithTooltip from 'elements/name-with-tooltip';
import { PushActionsType } from '../types';

export const useColumns = ({
  onActions,
  tags
}: {
  onActions: (item: Push, type: PushActionsType) => void;
  tags: Tag[];
}): ColumnDef<Push>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { consoleAccount } = useAuth();
  const editable = hasEditable(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
      cell: ({ row }) => {
        const push = row.original;
        const { id, name } = push;

        return (
          <div className="flex items-center gap-0.5 max-w-fit min-w-[300px]">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <NameWithTooltip.Trigger
                  id={id}
                  name={name}
                  onClick={() => onActions(push, 'EDIT')}
                  maxLines={1}
                />
              }
              maxLines={1}
            />
          </div>
        );
      }
    },
    {
      accessorKey: 'tags',
      enableSorting: false,
      header: `${t('tags')}`,
      size: 350,
      cell: ({ row }) => {
        const push = row.original;
        const formattedTags = push.tags?.map(
          item => tags.find(tag => tag.id === item)?.name || item
        );

        return (
          <ExpandableTag
            tags={formattedTags}
            rowId={push.id}
            className="!max-w-[250px] truncate"
          />
        );
      }
    },
    {
      accessorKey: 'environment',
      header: `${t('environment')}`,
      size: 250,
      maxSize: 250,
      cell: ({ row }) => {
        const push = row.original;
        const id = `env-${push.id}`;
        return (
          <NameWithTooltip
            id={id}
            align="center"
            content={
              <NameWithTooltip.Content content={push.environmentName} id={id} />
            }
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={push.environmentName}
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
      size: 200,
      cell: ({ row }) => {
        const push = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatDateTime(push.createdAt)}
              </div>
            }
            date={push.createdAt}
          />
        );
      }
    },
    {
      accessorKey: 'state',
      header: `${t('state')}`,
      size: 120,
      cell: ({ row }) => {
        const push = row.original;

        return (
          <DisabledButtonTooltip
            align="center"
            hidden={editable}
            trigger={
              <div className="w-fit">
                <Switch
                  disabled={!editable}
                  checked={!push.disabled}
                  onCheckedChange={value =>
                    onActions(push, value ? 'ENABLE' : 'DISABLE')
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
        const push = row.original;

        return (
          <DisabledPopoverTooltip
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
            onClick={value => onActions(push, value as PushActionsType)}
          />
        );
      }
    }
  ];
};
