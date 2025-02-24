import { useTranslation } from 'i18n';
import { ExperimentCollection, ExperimentStatus } from '@types';
import {
  IconExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import OverviewCard from 'elements/overview-card';

const Overview = ({
  summary,
  onChangeFilters
}: {
  summary?: ExperimentCollection['summary'];
  onChangeFilters: (statuses: ExperimentStatus[]) => void;
}) => {
  const { t } = useTranslation(['table']);
  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-8">
      <OverviewCard
        title={t('experiment.waiting-experiments')}
        count={Number(summary?.totalWaitingCount || 0)}
        color="orange"
        icon={IconWaitingExperiment}
        onClick={() => onChangeFilters(['WAITING'])}
      />
      <OverviewCard
        title={t('experiment.running-experiments')}
        count={Number(summary?.totalRunningCount || 0)}
        color="brand"
        icon={IconExperiment}
        onClick={() => onChangeFilters(['RUNNING'])}
      />
      <OverviewCard
        title={t('experiment.stopped-experiments')}
        count={Number(summary?.totalStoppedCount || 0)}
        color="red"
        icon={IconStoppedExperiment}
        onClick={() => onChangeFilters(['STOPPED', 'FORCE_STOPPED'])}
      />
    </div>
  );
};

export default Overview;
