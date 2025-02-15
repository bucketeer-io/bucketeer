import { useTranslation } from 'i18n';
import { IconExperiment, IconNotInUse, IconOperation } from '@icons';
import OverviewCard from 'elements/overview-card';

const Overview = () => {
  const { t } = useTranslation(['common']);
  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-8">
      <OverviewCard
        title={t('experiments-connected')}
        count={10}
        highlightText="+12%"
        description="from last month"
        color="brand"
        icon={IconExperiment}
      />
      <OverviewCard
        title={t('operations-connected')}
        count={10}
        highlightText="+12%"
        description="from last month"
        color="pink"
        icon={IconOperation}
      />
      <OverviewCard
        title={t('not-in-use')}
        count={5}
        highlightText="+12%"
        description="from last month"
        color="brand"
        icon={IconNotInUse}
      />
    </div>
  );
};

export default Overview;
