import { useCallback, useMemo, useState } from 'react';
import { useQueryInsightsMonthlySummary } from '@queries/insights';
import { useQueryProjects } from '@queries/projects';
import { getUniqueProjects, useAuth } from 'auth';
import { ALL_API_IDS, ALL_SOURCE_IDS } from 'constants/insight';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';
import { InsightsFilters, computeTimeRange, normalizeEnvId } from './utils';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const roles = consoleAccount!.environmentRoles;

  const userProjects = useMemo(() => getUniqueProjects(roles), [roles]);
  const userEnvironments = useMemo(
    () => roles.map(r => r.environment),
    [roles]
  );

  const initialProjectId = userProjects[0]?.id ?? '';
  const initialEnvironmentId = normalizeEnvId(
    userEnvironments.find(e => e.projectId === initialProjectId)?.id ?? ''
  );

  const [filters, setFilters] = useState<InsightsFilters>({
    projectId: initialProjectId,
    environmentId: initialEnvironmentId,
    sourceId: '',
    apiId: '',
    timeRange: '24h'
  });

  const organizationId = roles[0]?.environment.organizationId ?? '';

  const { data: projectsData, isLoading: projectsLoading } = useQueryProjects({
    params: { cursor: '', organizationId }
  });

  const filteredEnvironments = useMemo(() => {
    return filters.projectId
      ? userEnvironments.filter(env => env.projectId === filters.projectId)
      : userEnvironments;
  }, [userEnvironments, filters.projectId]);

  const handleProjectChange = useCallback(
    (projectId: string) => {
      const firstEnv = projectId
        ? userEnvironments.find(env => env.projectId === projectId)
        : userEnvironments[0];
      setFilters(prev => ({
        ...prev,
        projectId,
        environmentId: normalizeEnvId(firstEnv?.id ?? '')
      }));
    },
    [userEnvironments]
  );

  const hasEnvironments = filteredEnvironments.length > 0;

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
      onProjectChange={handleProjectChange}
      queriesEnabled={hasEnvironments}
    />
  );
};

export default PageLoader;
