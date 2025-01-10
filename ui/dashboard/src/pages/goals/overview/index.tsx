import { IconExperiment, IconNotInUse, IconOperation } from '@icons';
import OverviewCard from './overview-card';

const Overview = () => {
  return (
    <div className="flex flex-wrap items-center w-full gap-6 pb-8">
      <OverviewCard
        title="Experiments Connected"
        count={10}
        highlightText="+12%"
        description="from last month"
        color="brand"
        icon={IconExperiment}
      />
      <OverviewCard
        title="Operations Connected"
        count={10}
        highlightText="+12%"
        description="from last month"
        color="pink"
        icon={IconOperation}
      />
      <OverviewCard
        title="Not In Use"
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
