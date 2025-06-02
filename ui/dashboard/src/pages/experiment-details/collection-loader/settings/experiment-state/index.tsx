import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { experimentUpdater, ExperimentUpdaterParams } from '@api/experiment';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment, ExperimentResult } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import {
  IconExperiment,
  IconMember,
  IconNotStartedExperiment,
  IconStartExperiment,
  IconStopExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import ConfirmModal from 'elements/confirm-modal';

const ExperimentState = ({
  experimentResult,
  experiment
}: {
  experiment: Experiment;
  experimentResult?: ExperimentResult;
}) => {
  const { t } = useTranslation(['table', 'form', 'common', 'message']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { notify, errorNotify } = useToast();
  const isRunning = experiment.status === 'RUNNING',
    isWaiting = experiment.status === 'WAITING',
    isStopped = ['STOPPED', 'FORCE_STOPPED'].includes(experiment.status);

  const totalUsers = useMemo(() => {
    const formatNumber = (num: number): string => {
      if (!num) return '';
      if (num >= 1e9) return (num / 1e9).toFixed(1) + 'B';
      if (num >= 1e6) return (num / 1e6).toFixed(1) + 'M';
      if (num >= 1e3) return (num / 1e3).toFixed(1) + 'K';
      return num.toString();
    };

    const total = formatNumber(
      experimentResult ? +experimentResult?.totalEvaluationUserCount : 0
    );
    return total || 0;
  }, [experimentResult]);

  const [
    openToggleExperimentModal,
    onOpenToggleExperimentModal,
    onCloseToggleExperimentModal
  ] = useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (params: ExperimentUpdaterParams) => {
      return experimentUpdater(params);
    },
    onSuccess: ({ experiment }, params) => {
      onCloseToggleExperimentModal();

      invalidateExperiments(queryClient);
      invalidateExperimentDetails(queryClient, {
        environmentId: currentEnvironment.id,
        id: experiment?.id ?? ''
      });
      mutation.reset();
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:source-type.experiment'),
          action: t(
            params?.status?.status === 'FORCE_STOPPED'
              ? 'common:stopped'
              : 'common:started'
          )
        })
      });
    },
    onError: error => errorNotify(error)
  });

  const onToggleExperiment = () => {
    if (experiment?.id) {
      mutation.mutate({
        id: experiment?.id,
        environmentId: currentEnvironment.id,
        status: {
          status: isRunning ? 'FORCE_STOPPED' : 'RUNNING'
        }
      });
    }
  };

  return (
    <div className="flex items-center justify-between w-full min-w-fit px-4 py-2 gap-x-4 bg-gray-100 rounded-lg">
      <div className="flex items-center">
        <div className="flex items-center gap-x-2 pr-3 border-r border-gray-400">
          <Icon
            icon={
              isRunning
                ? IconExperiment
                : isWaiting
                  ? IconWaitingExperiment
                  : isStopped
                    ? IconStoppedExperiment
                    : IconNotStartedExperiment
            }
            size={'md'}
            className="[&>svg]:size-6"
          />
          <p className="typo-head-bold-small text-gray-700 whitespace-nowrap">
            {t(
              isRunning
                ? `experiment.running-experiment`
                : isWaiting
                  ? `experiment.scheduled-experiment`
                  : isStopped
                    ? `experiment.stopped-experiment`
                    : `experiment.not-started-experiment`
            )}
          </p>
        </div>
        <div className="px-3 typo-para-small text-gray-700 whitespace-nowrap">
          <Trans
            i18nKey={
              isRunning
                ? 'table:experiment.running-experiment-desc'
                : isWaiting
                  ? 'table:experiment.scheduled-experiment-desc'
                  : isStopped
                    ? 'table:experiment.stopped-experiment-desc'
                    : 'table:experiment.not-started-experiment-desc'
            }
            values={{
              date: formatLongDateTime({
                value: experiment.stopAt,
                overrideOptions: {
                  month: '2-digit',
                  day: '2-digit',
                  hour: '2-digit',
                  minute: '2-digit',
                  hour12: false
                },
                locale: 'ja-JP'
              })?.replace(' ', ' - ')
            }}
          />
        </div>
        <div className="flex items-center gap-x-3 pl-3 typo-para-small text-gray-700 whitespace-nowrap border-l border-gray-400">
          <Icon icon={IconMember} size="sm" />
          <p>
            <Trans
              i18nKey={`table:results.total-users-use`}
              components={{
                bold: <strong className="text-gray-700" />
              }}
              values={{
                value: totalUsers
              }}
            />
          </p>
        </div>
      </div>
      <Button
        disabled={isStopped}
        variant={'text'}
        className={cn('!typo-para-small h-10 whitespace-nowrap', {
          'text-accent-red-500 hover:text-accent-red-600': isRunning
        })}
        onClick={onOpenToggleExperimentModal}
      >
        <Icon
          icon={isRunning ? IconStopExperiment : IconStartExperiment}
          size={'sm'}
        />
        {t(isRunning ? `popover.stop-experiment` : `popover.start-experiment`)}
      </Button>
      {openToggleExperimentModal && (
        <ConfirmModal
          isOpen={openToggleExperimentModal}
          onClose={onCloseToggleExperimentModal}
          onSubmit={onToggleExperiment}
          title={
            isRunning
              ? t(`table:popover.stop-experiment`)
              : t(`table:popover.start-experiment`)
          }
          description={
            <Trans
              i18nKey={
                isRunning
                  ? 'table:experiment.confirm-stop-desc'
                  : 'table:experiment.confirm-start-desc'
              }
              values={{ name: experiment?.name }}
              components={{ bold: <strong /> }}
            />
          }
          loading={mutation.isPending}
        />
      )}
    </div>
  );
};

export default ExperimentState;
