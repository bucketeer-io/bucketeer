import { useState } from 'react';
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
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import TruncationWithTooltip from '../../../elements/truncation-with-tooltip';
import { PushActionsType } from '../types';

export const Tag = ({ value }: { value: string }) => (
  <div className="flex-center px-2 py-1.5 bg-primary-100/70 text-primary-500 typo-para-small leading-[14px] rounded whitespace-nowrap">
    {value}
  </div>
);

export const renderTag = ({
  tags,
  isExpanded,
  onExpand
}: {
  tags: string[];
  isExpanded: boolean;
  onExpand: () => void;
}) => {
  return (
    <div
      className={cn(
        'flex items-center w-full gap-x-2 transition-all duration-300',
        {
          'items-start': isExpanded
        }
      )}
    >
      <div className="flex items-center flex-wrap gap-2 max-w-fit transition-all duration-300">
        {(isExpanded ? tags : tags.slice(0, 3))?.map((tag, index) => (
          <Tag value={tag} key={index} />
        ))}
        {tags.length > 3 && <Tag value={`+${tags.length - 3}`} />}
      </div>
      {tags.length > 3 && (
        <div
          className={cn('flex-center cursor-pointer hover:bg-gray-200 rounded')}
          onClick={onExpand}
        >
          <Icon
            icon={IconChevronDown}
            size={'sm'}
            className={cn('flex-center rotate-0', {
              'rotate-180': isExpanded
            })}
          />
        </div>
      )}
    </div>
  );
};

export const useColumns = ({
  onActions
}: {
  onActions: (item: Push, type: PushActionsType) => void;
}): ColumnDef<Push>[] => {
  const { t } = useTranslation(['common', 'table']);
  const [expandedTags, setExpandedTags] = useState<string[]>([]);
  const formatDateTime = useFormatDateTime();

  const handleExpandTag = (pushId: string) => {
    setExpandedTags(
      expandedTags.includes(pushId)
        ? expandedTags.filter(item => item !== pushId)
        : [...expandedTags, pushId]
    );
  };

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
      cell: ({ row }) => {
        const push = row.original;

        return (
          <div className="flex flex-col gap-0.5 max-w-fit">
            <button
              onClick={() => onActions(push, 'EDIT')}
              className="underline text-primary-500 break-all line-clamp-2 typo-para-medium text-left"
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
      size: 350,
      cell: ({ row }) => {
        const push = row.original;

        return renderTag({
          tags: push.tags,
          isExpanded: expandedTags.includes(push.id),
          onExpand: () => handleExpandTag(push.id)
        });
      }
    },
    {
      accessorKey: 'environment',
      header: `${t('environment')}`,
      size: 250,
      maxSize: 250,
      cell: ({ row }) => {
        const push = row.original;
        return (
          <TruncationWithTooltip
            elementId={`env-${push.id}`}
            maxSize={250}
            content={push.environmentName}
            trigger={
              <div
                id={`env-${push.id}`}
                className={cn('text-gray-700 typo-para-medium w-fit')}
              >
                {push.environmentName}
              </div>
            }
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
