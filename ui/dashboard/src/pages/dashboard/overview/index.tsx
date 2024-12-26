import { IconFeatureSwitch, IconUserCheckmark, IconUsersGroup } from '@icons';
import OverviewCard from 'elements/overview-card';
import { useTranslation } from 'i18n';

const Overview = () => {
  const {t} = useTranslation(['common'])
  return (
    <div className="flex flex-wrap items-center w-full gap-6">
      <OverviewCard
        icon={IconUsersGroup}
        color="blue"
        title={t("active-members")}
        count={5}
        description="out of 144 total members"
        showArrow={true}
      />
      <OverviewCard
        icon={IconUserCheckmark}
        color="pink"
        title={t("flags-ready-remove")}
        count={56}
        description="out of 531 temporary flags"
      />
      <OverviewCard
        icon={IconFeatureSwitch}
        color="brand"
        title={t("flags-ready-archive")}
        count={4}
        description="out of 531 temporary flags"
      />
    </div>
  );
};

export default Overview;
