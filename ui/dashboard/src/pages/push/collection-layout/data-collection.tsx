import {
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { Push } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { useFormatDateTime } from 'utils/date-time';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import { PushActionsType } from '../types';

export const tagOptions = [
  {
    label: 'Android',
    value: 'android'
  },
  {
    label: 'IOS',
    value: 'ios'
  },
  {
    label: 'Web',
    value: 'web'
  },
  {
    label: 'JSON',
    value: 'json'
  },
  {
    label: 'JSON New',
    value: 'json-new'
  },
  {
    label: 'Number',
    value: 'number'
  },
  {
    label: 'String',
    value: 'string'
  }
];

export const Tag = ({ tag }: { tag: string }) => (
  <div className="flex-center px-2 py-1.5 bg-primary-100/70 text-primary-500 typo-para-small leading-[14px] rounded whitespace-nowrap">
    {tagOptions.find(item => item.value === tag)?.label || tag}
  </div>
);

export const renderTag = (tags: string[]) => {
  return (
    <div className="flex items-center gap-2 max-w-fit">
      {tags.slice(0, 3)?.map((tag, index) => <Tag tag={tag} key={index} />)}
      {tags.length > 3 && <Tag tag={`+${tags.length - 3}`} />}
    </div>
  );
};

export const useColumns = ({
  onActions
}: {
  onActions: (item: Push, type: PushActionsType) => void;
}): ColumnDef<Push>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const push = row.original;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <button
              onClick={() => onActions(push, 'EDIT')}
              className="underline text-primary-500 break-all typo-para-medium text-left"
            >
              {push.name}
            </button>
            <div className="typo-para-tiny text-gray-500">
              {truncateTextCenter(push.name)}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'tags',
      header: `${t('tags')}`,
      size: 500,
      cell: ({ row }) => {
        const push = row.original;

        return renderTag(push.tags);
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 150,
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
            checked={push.deleted}
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
