import {
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Environment } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { checkEnvironmentEmptyId } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import DateTooltip from 'elements/date-tooltip';
import DisabledPopoverTooltip from 'elements/disabled-popover-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { EnvironmentActionsType } from '../types';

export const useColumns = ({
  onActions
}: {
  onActions: (item: Environment, type: EnvironmentActionsType) => void;
}): ColumnDef<Environment>[] => {
  const { searchOptions } = useSearchParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { t } = useTranslation(['common', 'table']);
  const formatDateTime = useFormatDateTime();

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 350,
      cell: ({ row }) => {
        const environment = row.original;
        const { id, name } = environment;

        return (
          <NameWithTooltip
            id={id}
            content={<NameWithTooltip.Content content={name} id={id} />}
            trigger={
              <NameWithTooltip.Trigger
                id={id}
                name={name}
                maxLines={1}
                className="min-w-[300px]"
                onClick={() => onActions(environment, 'EDIT')}
              />
            }
            maxLines={1}
          />
        );
      }
    },
    {
      accessorKey: 'featureFlagCount',
      header: `${t('table:flags')}`,
      size: 250,
      cell: ({ row }) => {
        const environment = row.original;
        return (
          <div className="text-gray-700 typo-para-medium">
            {environment.featureFlagCount}
          </div>
        );
      }
    },
    {
      accessorKey: 'createdAt',
      header: `${t('table:created-at')}`,
      size: 160,
      cell: ({ row }) => {
        const environment = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatDateTime(environment.createdAt)}
              </div>
            }
            date={environment.createdAt}
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
        const environment = row.original;
        const isDisabled =
          checkEnvironmentEmptyId(currentEnvironment.id) === environment.id;
        return (
          <DisabledPopoverTooltip
            isNeedAdminAccess
            content={
              isDisabled ? t('table:disabled-archive-current-env') : undefined
            }
            options={[
              {
                label: `${t('table:popover.edit-env')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-env')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE',
                    disabled: isDisabled
                  }
            ]}
            onClick={value =>
              onActions(environment, value as EnvironmentActionsType)
            }
          />
        );
      }
    }
  ];
};
