import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { FeatureCountByStatus } from '@types';
import { cn } from 'utils/style';
import { IconActiveFlags, IconFlagNoTraffic, IconTotalFlags } from '@icons';
import { Tooltip } from 'components/tooltip';
import OverviewCard, { OverviewIconColor } from 'elements/overview-card';
import { StatusFilterType } from '../types';

interface OverviewOption {
  titleKey: string;
  countKey?: keyof FeatureCountByStatus;
  color: OverviewIconColor;
  icon: FunctionComponent;
  filterValue: StatusFilterType | undefined;
  tooltipKey: string;
}

const overviewOptions: OverviewOption[] = [
  {
    titleKey: 'feature-flags.total-flags',
    countKey: 'total',
    color: 'brand',
    icon: IconTotalFlags,
    filterValue: undefined,
    tooltipKey: 'feature-flags.total-flags-description'
  },
  {
    titleKey: 'feature-flags.receiving-traffic-flags',
    countKey: 'active',
    color: 'green',
    icon: IconActiveFlags,
    filterValue: StatusFilterType.RECEIVING_TRAFFIC,
    tooltipKey: 'feature-flags.receiving-traffic-description'
  },
  {
    titleKey: 'feature-flags.no-recent-traffic-flags',
    countKey: 'inactive',
    color: 'yellow',
    icon: IconFlagNoTraffic,
    filterValue: StatusFilterType.NO_RECENT_TRAFFIC,
    tooltipKey: 'feature-flags.no-recent-traffic-description'
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
    <div className="w-full px-3 sm:px-6">
      <div className="flex flex-nowrap sm:flex-wrap overflow-x-scroll sm:overflow-visible px-2 sm:px-0 hidden-scroll items-center w-full gap-6 pb-8">
        {overviewOptions.map((item, index) => (
          <Tooltip
            key={index}
            content={t(item.tooltipKey)}
            trigger={
              <div className="flex flex-1 w-full min-w-[268px]">
                <OverviewCard
                  title={t(item.titleKey)}
                  count={
                    summary && item.countKey
                      ? Number(summary[item.countKey])
                      : 0
                  }
                  color={item.color}
                  icon={item.icon}
                  className={cn('border border-transparent', {
                    'border-gray-300': item.filterValue === statusFilter
                  })}
                  onClick={() => onChangeFilters(item.filterValue)}
                />
              </div>
            }
            className="max-w-[300px]"
          />
        ))}
      </div>
    </div>
  );
};

export default Overview;
