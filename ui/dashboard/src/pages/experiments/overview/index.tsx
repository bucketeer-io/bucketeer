import { useTranslation } from 'i18n';
import {
  IconExperiment,
  IconStoppedExperiment,
  IconWaitingExperiment
} from '@icons';
import OverviewCard from 'elements/overview-card';

const Overview = () => {
  const { t } = useTranslation(['table']);
  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-8">
      <OverviewCard
        title={t('experiment.waiting-experiments')}
        count={10}
        color="orange"
        icon={IconWaitingExperiment}
      />
      <OverviewCard
        title={t('experiment.running-experiments')}
        count={10}
        color="brand"
        icon={IconExperiment}
      />
      <OverviewCard
        title={t('experiment.stopped-experiments')}
        count={5}
        color="red"
        icon={IconStoppedExperiment}
      />
    </div>
  );
};

export default Overview;
