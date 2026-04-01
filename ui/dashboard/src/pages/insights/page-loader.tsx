import { useMemo, useState } from 'react';
import { useQueryEnvironments } from '@queries/environments';
import { useQueryInsightsMonthlySummary } from '@queries/insights';
import { useQueryProjects } from '@queries/projects';
import { DateTime } from 'luxon';
import { InsightSourceId, InsightApiId } from '@types';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

export type TimeRangePreset = '1h' | '6h' | '24h' | '7d' | '30d' | 'this_month';

export interface InsightsFilters {
  projectId: string;
  environmentId: string;
  sourceId: InsightSourceId | '';
  apiId: InsightApiId | '';
  timeRange: TimeRangePreset;
}

const presetMap: Record<TimeRangePreset, (now: DateTime) => DateTime> = {
  '1h': now => now.minus({ hours: 1 }),
  '6h': now => now.minus({ hours: 6 }),
  '24h': now => now.minus({ hours: 24 }),
  '7d': now => now.minus({ days: 7 }),
  '30d': now => now.minus({ days: 30 }),
  this_month: now => now.startOf('month')
};

const computeTimeRange = (
  preset: TimeRangePreset
): { startAt: string; endAt: string } => {
  const now = DateTime.now();
  const startAt = presetMap[preset](now);
  return {
    startAt: String(Math.floor(startAt.toSeconds())),
    endAt: String(Math.floor(now.toSeconds()))
  };
};

const PageLoader = () => {
  const [filters, setFilters] = useState<InsightsFilters>({
    projectId: '',
    environmentId: '',
    sourceId: '',
    apiId: '',
    timeRange: '24h'
  });

  const { data: projectsData, isLoading: projectsLoading } = useQueryProjects({
    params: { cursor: '' }
  });

  const { data: environmentsData, isLoading: environmentsLoading } =
    useQueryEnvironments({
      params: { cursor: '' }
    });

  const monthlySummaryParams = useMemo(() => {
    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : undefined,
      sourceIds: filters.sourceId ? [filters.sourceId] : undefined
    };
  }, [filters.environmentId, filters.sourceId]);

  const { data: monthlySummary, isLoading: monthlySummaryLoading } =
    useQueryInsightsMonthlySummary({ params: monthlySummaryParams });

  const timeRangeParams = useMemo(() => {
    const { startAt, endAt } = computeTimeRange(filters.timeRange);
    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : undefined,
      sourceIds: filters.sourceId ? [filters.sourceId] : undefined,
      apiIds: filters.apiId ? [filters.apiId] : undefined,
      startAt,
      endAt
    };
  }, [filters]);

  if (projectsLoading || environmentsLoading) {
    return <PageLayout.LoadingState />;
  }

  return (
    <PageContent
      projects={projectsData?.projects ?? []}
      environments={environmentsData?.environments ?? []}
      monthlySummary={monthlySummary}
      monthlySummaryLoading={monthlySummaryLoading}
      timeRangeParams={timeRangeParams}
      filters={filters}
      onFiltersChange={setFilters}
    />
  );
};

export default PageLoader;
