import { useEffect, useMemo, useState } from 'react';
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

  /**
   * Filter environments by selected project
   */
  const filteredEnvironments = useMemo(() => {
    const envs = environmentsData?.environments ?? [];

    return filters.projectId
      ? envs.filter(env => env.projectId === filters.projectId)
      : envs;
  }, [environmentsData?.environments, filters.projectId]);

  /**
   * Auto-select first environment when project changes
   */
  useEffect(() => {
    const firstEnvId = filters.projectId
      ? (filteredEnvironments[0]?.id ?? '')
      : '';

    setFilters(prev => {
      if (prev.environmentId === firstEnvId) return prev;

      return {
        ...prev,
        environmentId: firstEnvId
      };
    });
  }, [filters.projectId, filteredEnvironments]);

  /**
   * Params for monthly summary
   */
  const monthlySummaryParams = useMemo(() => {
    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : (['UNKNOWN'] as string[]),
      sourceIds: filters.sourceId
        ? [filters.sourceId]
        : (['UNKNOWN'] as InsightSourceId[])
    };
  }, [filters.environmentId, filters.sourceId]);

  const { data: monthlySummary, isLoading: monthlySummaryLoading } =
    useQueryInsightsMonthlySummary({
      params: monthlySummaryParams
    });

  /**
   * Params for time range based queries
   */
  const timeRangeParams = useMemo(() => {
    const { startAt, endAt } = computeTimeRange(filters.timeRange);

    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : (['UNKNOWN'] as string[]),
      sourceIds: filters.sourceId
        ? [filters.sourceId]
        : (['UNKNOWN'] as InsightSourceId[]),
      apiIds: filters.apiId
        ? [filters.apiId]
        : (['UNKNOWN_API'] as InsightApiId[]),
      startAt,
      endAt
    };
  }, [
    filters.environmentId,
    filters.sourceId,
    filters.apiId,
    filters.timeRange
  ]);

  if (projectsLoading || environmentsLoading) {
    return <PageLayout.LoadingState />;
  }

  return (
    <PageContent
      projects={projectsData?.projects ?? []}
      environments={filteredEnvironments}
      monthlySummary={monthlySummary}
      monthlySummaryLoading={monthlySummaryLoading}
      timeRangeParams={timeRangeParams}
      filters={filters}
      onFiltersChange={setFilters}
    />
  );
};

export default PageLoader;
