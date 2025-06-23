import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Push, Tag } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { IconTrash } from '@icons';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
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

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
      cell: ({ row }) => {
        const push = row.original;
        const { id, name } = push;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit min-w-[300px]">
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
            <div className="typo-para-tiny text-gray-500 break-all line-clamp-1">
              {truncateTextCenter(id)}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'tags',
      header: `${t('tags')}`,
      enableSorting: false,
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
          <Switch
            checked={!push.disabled}
            onCheckedChange={value =>
              onActions(push, value ? 'ENABLE' : 'DISABLE')
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
          <Popover
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
            icon={IconMoreHorizOutlined}
            onClick={value => onActions(push, value as PushActionsType)}
            align="end"
          />
        );
      }
    }
  ];
};
