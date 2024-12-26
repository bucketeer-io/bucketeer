import { IconFeatureSwitch, IconUserCheckmark, IconUsersGroup } from '@icons';
import OverviewCard from 'elements/overview-card';

const Overview = () => {
  return (
    <div className="flex flex-wrap items-center w-full gap-6">
      <OverviewCard
        icon={IconUsersGroup}
        color="blue"
        title="Active members"
        count={5}
        description="out of 144 total members"
        showArrow={true}
      />
      <OverviewCard
        icon={IconUserCheckmark}
        color="pink"
        title="Flags ready for code removal"
        count={56}
        description="out of 531 temporary flags"
      />
      <OverviewCard
        icon={IconFeatureSwitch}
        color="brand"
        title="Flags ready to archive"
        count={4}
        description="out of 531 temporary flags"
      />
    </div>
  );
};

export default Overview;
