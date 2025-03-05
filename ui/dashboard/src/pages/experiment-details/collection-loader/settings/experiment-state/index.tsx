import { Trans } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { experimentUpdater, ExperimentUpdaterParams } from '@api/experiment';
import { invalidateExperimentDetails } from '@queries/experiment-details';
import { invalidateExperiments } from '@queries/experiments';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import {
  IconExperiment,
  IconNotStartedExperiment,
  IconStartExperiment,
  IconStopExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import ConfirmModal from 'elements/confirm-modal';

const ExperimentState = ({ experiment }: { experiment: Experiment }) => {
  const { t } = useTranslation(['table', 'form']);
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const params = useParams();
  const { notify } = useToast();
  const isRunning = experiment.status === 'RUNNING',
    isWaiting = experiment.status === 'WAITING',
    isStopped = ['STOPPED', 'FORCE_STOPPED'].includes(experiment.status);

  const [
    openToggleExperimentModal,
    onOpenToggleExperimentModal,
    onCloseToggleExperimentModal
  ] = useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (params: ExperimentUpdaterParams) => {
      return experimentUpdater(params);
    },
    onSuccess: () => {
      onCloseToggleExperimentModal();

      invalidateExperiments(queryClient);
      invalidateExperimentDetails(queryClient, {
        environmentId: currentEnvironment.id,
        id: params?.experimentId ?? ''
      });
      mutation.reset();
    },
    onError: error => {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: error?.message || 'Something went wrong.'
      });
    }
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
        <div className="pl-3 typo-para-medium text-gray-700 whitespace-nowrap">
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
      </div>
      <Button
        variant={'text'}
        className={cn('typo-sm h-10 whitespace-nowrap', {
          'text-accent-red-500 hover:text-accent-red-600':
            isRunning || isWaiting
        })}
        onClick={onOpenToggleExperimentModal}
      >
        <Icon
          icon={
            isRunning || isWaiting ? IconStopExperiment : IconStartExperiment
          }
          size={'sm'}
        />
        {t(
          isRunning || isWaiting
            ? `popover.stop-experiment`
            : `popover.start-experiment`
        )}
      </Button>
      {openToggleExperimentModal && (
        <ConfirmModal
          isOpen={openToggleExperimentModal}
          onClose={onCloseToggleExperimentModal}
          onSubmit={onToggleExperiment}
          title={
            isRunning || isWaiting
              ? t(`table:popover.stop-experiment`)
              : t(`table:popover.start-experiment`)
          }
          description={
            <Trans
              i18nKey={
                isRunning || isWaiting
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
