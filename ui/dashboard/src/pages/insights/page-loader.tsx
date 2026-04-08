import { useCallback, useMemo } from 'react';
import { useQueryInsightsMonthlySummary } from '@queries/insights';
import { useQueryProjects } from '@queries/projects';
import {
  getCurrentEnvironment,
  getCurrentProject,
  getUniqueProjects,
  useAuth
} from 'auth';
import { ALL, ALL_API_IDS, ALL_SOURCE_IDS } from 'constants/insight';
import { InsightApiId, InsightSourceId } from '@types';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';
import { isAll, useInsightsFilters } from './use-insights-filters';
import { computeTimeRange, normalizeEnvId } from './utils';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const roles = consoleAccount!.environmentRoles;

  const userProjects = useMemo(() => getUniqueProjects(roles), [roles]);
  const userEnvironments = useMemo(
    () => roles.map(r => r.environment),
    [roles]
  );

  const currentEnv = getCurrentEnvironment(consoleAccount!);
  const currentProject = getCurrentProject(roles, currentEnv.id);
  const initialProjectId = currentProject?.id ?? userProjects[0]?.id ?? '';
  const initialEnvironmentId = normalizeEnvId(currentEnv.id);

  const { filters, setFilters } = useInsightsFilters(
    initialProjectId,
    initialEnvironmentId
  );

  const organizationId = roles[0]?.environment.organizationId ?? '';

  const { data: projectsData, isLoading: projectsLoading } = useQueryProjects({
    params: { cursor: '', organizationId }
  });

  const filteredEnvironments = useMemo(
    () =>
      !isAll(filters.projectId)
        ? userEnvironments.filter(env => env.projectId === filters.projectId)
        : userEnvironments,
    [userEnvironments, filters.projectId]
  );

  const handleProjectChange = useCallback(
    (projectId: string) => {
      const firstEnv = !isAll(projectId)
        ? userEnvironments.find(env => env.projectId === projectId)
        : userEnvironments[0];
      setFilters({
        ...filters,
        projectId: projectId || ALL,
        environmentId: normalizeEnvId(firstEnv?.id ?? '')
      });
    },
    [userEnvironments, filters, setFilters]
  );

  const hasEnvironments = filteredEnvironments.length > 0;

  const environmentIds = useMemo(
    () =>
      !isAll(filters.environmentId)
        ? [filters.environmentId]
        : filteredEnvironments.map(e => e.id),
    [filters.environmentId, filteredEnvironments]
  );

  const monthlySummaryParams = useMemo(
    () => ({
      environmentIds,
      sourceIds: !isAll(filters.sourceId)
        ? [filters.sourceId as InsightSourceId]
        : ALL_SOURCE_IDS
    }),
    [environmentIds, filters.sourceId]
  );

  const { data: monthlySummary, isLoading: monthlySummaryLoading } =
    useQueryInsightsMonthlySummary({
      params: monthlySummaryParams,
      enabled: hasEnvironments
    });

  const timeRangeParams = useMemo(() => {
    const { startAt, endAt } = computeTimeRange(
      filters.timeRange,
      filters.customStartAt,
      filters.customEndAt
    );
    return {
      environmentIds,
      sourceIds: !isAll(filters.sourceId)
        ? [filters.sourceId as InsightSourceId]
        : ALL_SOURCE_IDS,
      apiIds: !isAll(filters.apiId)
        ? [filters.apiId as InsightApiId]
        : ALL_API_IDS,
      startAt,
      endAt
    };
  }, [
    environmentIds,
    filters.sourceId,
    filters.apiId,
    filters.timeRange,
    filters.customStartAt,
    filters.customEndAt
  ]);

  const visibleProjects = useMemo(() => {
    const userProjectIds = new Set(userProjects.map(p => p.id));
    return (projectsData?.projects ?? []).filter(p => userProjectIds.has(p.id));
  }, [userProjects, projectsData]);

  if (projectsLoading) {
    return <PageLayout.LoadingState />;
  }

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
