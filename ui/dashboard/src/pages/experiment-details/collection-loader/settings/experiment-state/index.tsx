import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Experiment } from '@types';
import { formatLongDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import {
  IconExperiment,
  IconStartExperiment,
  IconStopExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

const ExperimentState = ({ experiment }: { experiment: Experiment }) => {
  const { t } = useTranslation(['table', 'form']);

  const isRunning = experiment.status === 'RUNNING',
    isWaiting = experiment.status === 'WAITING';

  return (
    <div className="flex items-center justify-between w-full px-4 py-2 gap-x-4 bg-gray-100 rounded-lg">
      <div className="flex items-center">
        <div className="flex items-center gap-x-2 pr-3 border-r border-gray-400">
          <Icon
            icon={
              isRunning
                ? IconExperiment
                : isWaiting
                  ? IconWaitingExperiment
                  : IconStoppedExperiment
            }
            size={'md'}
            className="[&>svg]:size-6"
          />
          <p className="typo-head-bold-small text-gray-700">
            {t(
              isRunning
                ? `experiment.running-experiments`
                : isWaiting
                  ? `experiment.waiting-experiments`
                  : `experiment.stopped-experiments`
            ).slice(0, -1)}
          </p>
        </div>
        <div className="pl-3 typo-para-medium text-gray-700">
          <Trans
            i18nKey={
              isRunning
                ? 'table:experiment.running-experiment-desc'
                : isWaiting
                  ? 'table:experiment.waiting-experiment-desc'
                  : 'table:experiment.stopped-experiment-desc'
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
        className={cn('typo-sm h-10', {
          'text-accent-red-500 hover:text-accent-red-600': isRunning
        })}
      >
        <Icon
          icon={isRunning ? IconStopExperiment : IconStartExperiment}
          size={'sm'}
        />
        {t(
          isRunning
            ? `experiment.stopped-experiments`
            : `experiment.running-experiments`
        ).slice(0, -1)}
      </Button>
    </div>
  );
};

export default ExperimentState;
