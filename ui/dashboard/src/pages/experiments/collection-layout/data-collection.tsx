import {
  IconArchiveOutlined,
  IconEditOutlined,
  IconMoreHorizOutlined
} from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import type { ColumnDef } from '@tanstack/react-table';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment, ExperimentStatus } from '@types';
import { truncateTextCenter } from 'utils/converts';
import { formatLongDateTime } from 'utils/date-time';
import { copyToClipBoard } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { cn } from 'utils/style';
import {
  IconCopy,
  IconExperiment,
  IconStartExperiment,
  IconStopExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import Icon from 'components/icon';
import { Popover } from 'components/popover';
import DateTooltip from 'elements/date-tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';
import { ExperimentActionsType } from '../types';

export const ExperimentStatuses = ({
  status
}: {
  status: ExperimentStatus;
}) => {
  const isRunning = status === 'RUNNING',
    isWaiting = status === 'WAITING',
    isStopped = ['STOPPED', 'FORCE_STOPPED'].includes(status);

  return (
    <div
      className={cn(
        'flex-center w-fit gap-x-1 px-2 py-1.5 typo-para-small rounded whitespace-nowrap capitalize',
        {
          'bg-primary-50 text-primary-500': isRunning,
          'bg-accent-orange-50 text-accent-orange-500': isWaiting,
          'bg-accent-red-50 text-accent-red-500': isStopped
        }
      )}
    >
      <Icon
        icon={
          isRunning
            ? IconExperiment
            : isWaiting
              ? IconWaitingExperiment
              : IconStoppedExperiment
        }
        size={'xxs'}
        className="flex-center [&>svg]:size-4"
      />
      {status.replace('_', ' ').toLowerCase()}
    </div>
  );
};

export const useColumns = ({
  onActions
}: {
  onActions: (item: Experiment, type: ExperimentActionsType) => void;
}): ColumnDef<Experiment>[] => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const { searchOptions } = useSearchParams();

  const { consoleAccount } = useAuth();
  const currenEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const { notify } = useToast();

  const handleCopyId = (id: string) => {
    copyToClipBoard(id);
    notify({
      toastType: 'toast',
      message: (
        <span>
          <b>ID</b> {` has been successfully copied!`}
        </span>
      )
    });
  };

  return [
    {
      id: 'name',
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 400,
      cell: ({ row }) => {
        const experiment = row.original;
        const { id, name } = experiment;
        return (
          <div className="flex flex-col gap-0.5">
            <NameWithTooltip
              id={id}
              content={<NameWithTooltip.Content content={name} id={id} />}
              trigger={
                <Link
                  to={`/${currenEnvironment.urlCode}${PAGE_PATH_EXPERIMENTS}/${id}/results`}
                >
                  <NameWithTooltip.Trigger
                    id={id}
                    name={name}
                    className="min-w-[300px]"
                  />
                </Link>
              }
            />
            <div className="flex items-center h-5 gap-x-2 typo-para-tiny text-gray-500 group select-none">
              {truncateTextCenter(id)}
              <div onClick={() => handleCopyId(id)}>
                <Icon
                  icon={IconCopy}
                  size={'sm'}
                  className="opacity-0 group-hover:opacity-100 cursor-pointer"
                />
              </div>
            </div>
          </div>
        );
      }
    },
    {
      id: 'goals_count',
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
            onClick={() => onActions(experiment, 'GOALS-CONNECTION')}
          >
            {experiment?.goalIds?.length || 0}
          </div>
        );
      }
    },
    {
      id: 'start_at',
      accessorKey: 'startAt',
      header: `${t('form:start-at')}`,
      size: 200,
      cell: ({ row }) => {
        const experiment = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium min-w-[150px]">
                {formatLongDateTime({
                  value: experiment.startAt,
                  overrideOptions: {
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false
                  },
                  locale: 'ja-JP'
                })}
              </div>
            }
            date={experiment.startAt}
          />
        );
      }
    },
    {
      id: 'stop_at',
      accessorKey: 'stopAt',
      header: `${t('form:end-at')}`,
      size: 200,
      cell: ({ row }) => {
        const experiment = row.original;
        return (
          <DateTooltip
            trigger={
              <div className="text-gray-700 typo-para-medium min-w-[150px]">
                {formatLongDateTime({
                  value: experiment.stopAt,
                  overrideOptions: {
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: false
                  },
                  locale: 'ja-JP'
                })}
              </div>
            }
            date={experiment.stopAt}
          />
        );
      }
    },
    {
      id: 'status',
      accessorKey: 'statuses',
      header: `${t('status')}`,
      size: 200,
      cell: ({ row }) => {
        const experiment = row.original;

        return <ExperimentStatuses status={experiment.status} />;
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
        const { status } = experiment;
        return (
          editable && (
            <Popover
              options={[
                {
                  label: `${t('table:popover.edit-experiment')}`,
                  icon: IconEditOutlined,
                  value: 'EDIT'
                },
                ...(['WAITING', 'RUNNING'].includes(status)
                  ? [
                      status === 'WAITING'
                        ? {
                            label: `${t('table:popover.start-experiment')}`,
                            icon: IconStartExperiment,
                            value: 'START'
                          }
                        : {
                            label: `${t('table:popover.stop-experiment')}`,
                            icon: IconStopExperiment,
                            value: 'STOP'
                          }
                    ]
                  : []),
                searchOptions.status === 'ARCHIVED'
                  ? {
                      label: `${t('table:popover.unarchive-experiment')}`,
                      icon: IconArchiveOutlined,
                      value: 'UNARCHIVE'
                    }
                  : {
                      label: `${t('table:popover.archive-experiment')}`,
                      icon: IconArchiveOutlined,
                      value: 'ARCHIVE',
                      disabled: ['RUNNING', 'WAITING'].includes(
                        row.original.status
                      )
                    }
              ]}
              icon={IconMoreHorizOutlined}
              onClick={value =>
                onActions(experiment, value as ExperimentActionsType)
              }
              align="end"
            />
          )
        );
      }
    }
  ];
};
