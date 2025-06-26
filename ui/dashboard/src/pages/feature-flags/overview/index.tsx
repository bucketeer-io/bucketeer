import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { FeatureCountByStatus } from '@types';
import { cn } from 'utils/style';
import { IconActiveFlags, IconInactiveFlags, IconTotalFlags } from '@icons';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { StatusFilterType } from '../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof FeatureCountByStatus;
  color: OverviewIconColor;
  icon: FunctionComponent;
  filterValue: StatusFilterType | undefined;
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'feature-flags.total-flags',
    countKey: 'total',
    color: 'brand',
    icon: IconTotalFlags,
    filterValue: undefined
  },
  {
    titleKey: 'feature-flags.active-flags',
    countKey: 'active',
    color: 'green',
    icon: IconActiveFlags,
    filterValue: StatusFilterType.ACTIVE
  },
  {
    titleKey: 'feature-flags.inactive-flags',
    countKey: 'inactive',
    color: 'yellow',
    icon: IconInactiveFlags,
    filterValue: StatusFilterType.NO_ACTIVITY
  }
];

const Overview = ({
  summary,
  statusFilter,
  onChangeFilters
}: {
  summary?: FeatureCountByStatus;
  statusFilter?: StatusFilterType;
  onChangeFilters: (filterValue: StatusFilterType | undefined) => void;
}) => {
  const { t } = useTranslation(['table']);

  return (
    <div className="w-full px-6 lg:pr-2">
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
              'border-gray-300': item.filterValue === statusFilter
            })}
            onClick={() => onChangeFilters(item.filterValue)}
          />
        ))}
      </div>
    </div>
  );
};

export default Overview;
