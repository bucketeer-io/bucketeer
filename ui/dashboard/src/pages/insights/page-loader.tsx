import { useEffect, useRef, useMemo, useState } from 'react';
import { useQueryInsightsMonthlySummary } from '@queries/insights';
import { useQueryProjects } from '@queries/projects';
import {
  getCurrentEnvironment,
  getCurrentProject,
  getUniqueProjects,
  useAuth
} from 'auth';
import { ALL_API_IDS, ALL_SOURCE_IDS } from 'constants/insight';
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
  customStartAt?: string;
  customEndAt?: string;
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
  preset: TimeRangePreset,
  customStartAt?: string,
  customEndAt?: string
): { startAt: string; endAt: string } => {
  if (customStartAt && customEndAt) {
    return {
      startAt: customStartAt,
      endAt: customEndAt
    };
  }
  const now = DateTime.now();
  const startAt = presetMap[preset](now);

  return {
    startAt: String(Math.floor(startAt.toSeconds())),
    endAt: String(Math.floor(now.toSeconds()))
  };
};

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const currentProject = getCurrentProject(
    consoleAccount!.environmentRoles,
    currentEnvironment.id
  );

  const [filters, setFilters] = useState<InsightsFilters>({
    projectId: currentProject?.id ?? '',
    environmentId: '',
    sourceId: '',
    apiId: '',
    timeRange: '24h'
  });

  const userProjects = getUniqueProjects(consoleAccount!.environmentRoles);
  const userEnvironments = consoleAccount!.environmentRoles.map(
    r => r.environment
  );

  const { data: projectsData, isLoading: projectsLoading } = useQueryProjects({
    params: { cursor: '', organizationId: currentEnvironment.organizationId }
  });

  const filteredEnvironments = useMemo(() => {
    return filters.projectId
      ? userEnvironments.filter(env => env.projectId === filters.projectId)
      : userEnvironments;
  }, [userEnvironments, filters.projectId]);

  const prevProjectIdRef = useRef(filters.projectId);

  useEffect(() => {
    if (prevProjectIdRef.current === filters.projectId) return;
    prevProjectIdRef.current = filters.projectId;

    setFilters(prev => ({ ...prev, environmentId: '' }));
  }, [filters.projectId]);

  const hasEnvironments = filters.environmentId
    ? true
    : filteredEnvironments.length > 0;

  /**
   * Params for monthly summary
   */
  const monthlySummaryParams = useMemo(() => {
    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : filteredEnvironments.map(e => e.id),
      sourceIds: filters.sourceId ? [filters.sourceId] : ALL_SOURCE_IDS
    };
  }, [filters.environmentId, filters.sourceId, filteredEnvironments]);

  const { data: monthlySummary, isLoading: monthlySummaryLoading } =
    useQueryInsightsMonthlySummary({
      params: monthlySummaryParams,
      enabled: hasEnvironments
    });

  /**
   * Params for time range based queries
   */
  const timeRangeParams = useMemo(() => {
    const { startAt, endAt } = computeTimeRange(
      filters.timeRange,
      filters.customStartAt,
      filters.customEndAt
    );

    return {
      environmentIds: filters.environmentId
        ? [filters.environmentId]
        : filteredEnvironments.map(e => e.id),
      sourceIds: filters.sourceId ? [filters.sourceId] : ALL_SOURCE_IDS,
      apiIds: filters.apiId ? [filters.apiId] : ALL_API_IDS,
      startAt,
      endAt
    };
  }, [
    filters.environmentId,
    filters.sourceId,
    filters.apiId,
    filters.timeRange,
    filters.customEndAt,
    filters.customStartAt,
    filteredEnvironments
  ]);

  if (projectsLoading) {
    return <PageLayout.LoadingState />;
  }

  const userProjectIds = new Set(userProjects.map(p => p.id));
  const visibleProjects = (projectsData?.projects ?? []).filter(p =>
    userProjectIds.has(p.id)
  );

  return (
    <PageContent
      projects={visibleProjects}
      environments={filteredEnvironments}
      monthlySummary={monthlySummary}
      monthlySummaryLoading={monthlySummaryLoading}
      timeRangeParams={timeRangeParams}
      filters={filters}
      onFiltersChange={setFilters}
      queriesEnabled={hasEnvironments}
    />
  );
};

export default PageLoader;
