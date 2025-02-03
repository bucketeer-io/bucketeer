import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useTranslation } from 'i18n';
import { Experiment, ExperimentStatus } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import { ExperimentActionsType } from '../types';

export const ExperimentStatuses = ({
  status
}: {
  status: ExperimentStatus;
}) => (
  <div
    className={cn(
      'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] rounded whitespace-nowrap capitalize',
      {
        'bg-accent-green-50 text-accent-green-500': status === 'RUNNING',
        'bg-accent-orange-50 text-accent-orange-500': status === 'WAITING',
        'bg-accent-red-50 text-accent-red-500': [
          'STOPPED',
          'FORCE_STOPPED'
        ].includes(status)
      }
    )}
  >
    {status.toLowerCase()}
  </div>
);

export const useColumns = ({
  onActions
}: {
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
}): ColumnDef<Experiment>[] => {
  const { t } = useTranslation(['common', 'table']);
  const { searchOptions } = useSearchParams();

  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const experiment = row.original;

        return (
          <div className="flex flex-col gap-0.5">
            <Link
              to={`/${currenEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${experiment.id}/results`}
              className="underline text-primary-500 typo-para-medium text-left w-fit max-w-full line-clamp-1 break-all"
            >
              {experiment.name}
            </Link>
            <div className="h-5 gap-x-2 typo-para-tiny text-gray-500 line-clamp-1">
              {experiment.id}
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'goalIds',
      header: `${t('navigation.goals')}`,
      size: 150,
      cell: ({ row }) => {
        const experiment = row.original;
        return (
          <div
            className={cn(
              'underline text-primary-500 break-all typo-para-medium text-left',
              {
                'cursor-pointer': experiment.goalIds?.length
              }
            )}
          >
            {experiment?.goalIds?.length || 0}
          </div>
        );
      }
    },
    {
      accessorKey: 'startAt',
      header: `${t('form:start-at')}`,
      size: 150,
      cell: ({ row }) => {
        const experiment = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatLongDateTime({
                  value: experiment.startAt,
                  overrideOptions: { month: 'numeric' },
                  locale: 'en-CA'
                })}
              </div>
            }
            date={experiment.startAt}
          />
        );
      }
    },
    {
      accessorKey: 'stopAt',
      header: `${t('form:end-at')}`,
      size: 150,
      cell: ({ row }) => {
        const experiment = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium">
                {formatLongDateTime({
                  value: experiment.stopAt,
                  overrideOptions: { month: 'numeric' },
                  locale: 'en-CA'
                })}
              </div>
            }
            date={experiment.stopAt}
          />
        );
      }
    },
    {
      accessorKey: 'statuses',
      header: `${t('status')}`,
      size: 120,
      cell: ({ row }) => {
        const experiment = row.original;

        return (
          <div className="flex items-center gap-x-2">
            <ExperimentStatuses status={experiment.status} />
            <Icon
              icon={IconChevronRight}
              className="rotate-90"
              color="gray-500"
              size={'sm'}
            />
          </div>
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
        const experiment = row.original;

        return (
          <Popover
            options={[
              {
                label: `${t('table:popover.edit-experiment')}`,
                icon: IconEditOutlined,
                value: 'EDIT'
              },
              searchOptions.status === 'ARCHIVED'
                ? {
                    label: `${t('table:popover.unarchive-experiment')}`,
                    icon: IconArchiveOutlined,
                    value: 'UNARCHIVE'
                  }
                : {
                    label: `${t('table:popover.archive-experiment')}`,
                    icon: IconArchiveOutlined,
                    value: 'ARCHIVE'
                  }
            ]}
            icon={IconMoreHorizOutlined}
            onClick={value =>
              onActions(experiment, value as ExperimentActionsType)
            }
            align="end"
          />
        );
      }
    }
  ];
};
