import {
  IconArchiveOutlined,
  IconMoreHorizOutlined,
  IconSaveAsFilled
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { useFormatDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import {
  IconFlagSwitch,
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString
} from '@icons';
import { Popover } from 'components/popover';
import Switch from 'components/switch';
import DateTooltip from 'elements/date-tooltip';
import { FlagsTemp, FlagDataType, FlagActionType } from '../types';
import {
  FlagNameElement,
  FlagOperationsElement,
  FlagTagsElement,
  FlagVariationsElement
} from './elements';

export const getDataTypeIcon = (type: FlagDataType) => {
  if (type === 'boolean') return IconFlagSwitch;
  if (type === 'string') return IconFlagString;
  if (type === 'number') return IconFlagNumber;
  return IconFlagJSON;
};

export const useColumns = ({
  onActions
}: {
  onActions: (item: FlagsTemp, type: FlagActionType) => void;
}): ColumnDef<FlagsTemp>[] => {
  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const { id, name, type, status } = row.original;

        return (
          <FlagNameElement
            link={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${id}/targeting`}
            name={name}
            id={id}
            icon={getDataTypeIcon(type)}
            status={status}
            viewType="LIST_VIEW"
          />
        );
      }
    },
    {
      accessorKey: 'tags',
      header: `${t('tags')}`,
      size: 150,
      cell: ({ row }) => {
        const { tags } = row.original;

        return <FlagTagsElement tags={tags} />;
      }
    },
    {
      accessorKey: 'variations',
      header: `${t('variations')}`,
      size: 190,
      cell: () => {
        return <FlagVariationsElement />;
      }
    },
    {
      accessorKey: 'updatedAt',
      header: `${t('table:updated-at')}`,
      size: 150,
      cell: ({ row }) => {
        const goal = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {Number(goal.updatedAt) === 0
                  ? t('never')
                  : formatDateTime(goal.updatedAt)}
              </div>
            }
            date={goal.updatedAt}
          />
        );
      }
    },
    {
      accessorKey: 'operations',
      header: `${t('operations')}`,
      size: 150,
      cell: () => {
        return <FlagOperationsElement />;
      }
    },
    {
      accessorKey: 'disabled',
      header: `${t('state')}`,
      size: 150,
      cell: ({ row }) => {
        const { disabled } = row.original;
        return <Switch checked={!disabled} />;
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
        const flag = row.original;

        return (
          <Popover
            options={compact([
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-flag')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-flag')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE'
                  },
              {
                label: `${t('table:popover.clone-flag')}`,
                icon: IconSaveAsFilled,
                value: 'CLONE'
              }
            ])}
            icon={IconMoreHorizOutlined}
            onClick={value => onActions(flag, value as FlagActionType)}
            align="end"
          />
        );
      }
    }
  ];
};
