import { useCallback, useEffect, useMemo } from 'react';
import { ALL } from 'constants/insight';
import { InsightApiId, InsightSourceId } from '@types';
import { useSearchParams } from 'utils/search-params';
import { InsightsFilters } from './utils';

const str = (v: unknown): string | undefined =>
  typeof v === 'string' ? v : undefined;

export const isAll = (v: string) => v === ALL;

export const useInsightsFilters = (
  initialProjectId: string,
  initialEnvironmentId: string
) => {
  const { searchOptions, onChangSearchParams } = useSearchParams();

  useEffect(() => {
    if (!window.location.search.includes('projectId=')) {
      onChangSearchParams({
        projectId: initialProjectId || ALL,
        environmentId: initialEnvironmentId || ALL,
        sourceId: ALL,
        apiId: ALL,
        timeRange: '24h'
      });
    }
  }, []);

  const filters: InsightsFilters = useMemo(
    () => ({
      projectId: str(searchOptions.projectId) ?? initialProjectId,
      environmentId: str(searchOptions.environmentId) ?? initialEnvironmentId,
      sourceId: (str(searchOptions.sourceId) ?? ALL) as
        | InsightSourceId
        | typeof ALL,
      apiId: (str(searchOptions.apiId) ?? ALL) as InsightApiId | typeof ALL,
      timeRange:
        (str(searchOptions.timeRange) as InsightsFilters['timeRange']) ?? '24h',
      customStartAt: str(searchOptions.customStartAt),
      customEndAt: str(searchOptions.customEndAt)
    }),
    [searchOptions, initialProjectId, initialEnvironmentId]
  );

  const setFilters = useCallback(
    (next: InsightsFilters) => {
      const params: Record<string, string> = {
        projectId: next.projectId || ALL,
        environmentId: next.environmentId || ALL,
        sourceId: next.sourceId,
        apiId: next.apiId,
        timeRange: next.timeRange
      };
      if (next.customStartAt) params.customStartAt = next.customStartAt;
      if (next.customEndAt) params.customEndAt = next.customEndAt;
      onChangSearchParams(params);
    },
    [onChangSearchParams]
  );

  return { filters, setFilters };
};
