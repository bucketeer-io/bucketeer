import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { ExperimentStatus, FeatureCountByStatus } from '@types';
import { cn } from 'utils/style';
import { IconActiveFlags, IconInactiveFlags, IconTotalFlags } from '@icons';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { SummaryType } from '../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof FeatureCountByStatus;
  color: OverviewIconColor;
  icon: FunctionComponent;
  summaryFilterValue: SummaryType;
  filterValues: ExperimentStatus[];
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'feature-flags.total-flags',
    countKey: 'total',
    color: 'brand',
    icon: IconTotalFlags,
    summaryFilterValue: 'TOTAL',
    filterValues: []
  },
  {
    titleKey: 'feature-flags.active-flags',
    countKey: 'active',
    color: 'green',
    icon: IconActiveFlags,
    summaryFilterValue: 'ACTIVE',
    filterValues: []
  },
  {
    titleKey: 'feature-flags.inactive-flags',
    countKey: 'inactive',
    color: 'yellow',
    icon: IconInactiveFlags,
    summaryFilterValue: 'INACTIVE',
    filterValues: []
  }
];

const Overview = ({
  summary,
  filterBySummary
}: {
  summary?: FeatureCountByStatus;
  filterBySummary?: SummaryType;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <div className="w-full px-6">
      <div className="flex flex-wrap items-center w-full gap-6 pb-8">
        {overviewOptions.map((item, index) => (
          <OverviewCard
            key={index}
            title={t(item.titleKey)}
            count={
              summary && item.countKey ? Number(summary[item.countKey]) : 0
            }
            color={item.color}
            icon={item.icon}
            className={cn('border border-transparent', {
              'border-gray-300':
                filterBySummary && item.summaryFilterValue === filterBySummary
            })}
          />
        ))}
      </div>
    </div>
  );
};

export default Overview;
