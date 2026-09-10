import { useCallback, useMemo } from 'react';
import { userSegmentFetcher } from '@api/user-segment';
import { SEGMENT_DETAILS_QUERY_KEY } from '@queries/user-segment-details';
import { useInfiniteQueryUserSegments } from '@queries/user-segments';
import { useQueries } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import uniq from 'lodash/uniq';
import { UserSegment } from '@types';

type UseUserSegmentsLoaderParams = {
  environmentId: string;
  pageSize?: number;
  enabled?: boolean;
  selectedSegmentIds?: string[];
};

export const useUserSegmentsLoader = ({
  environmentId,
  pageSize = LIST_PAGE_SIZE,
  enabled = true,
  selectedSegmentIds = []
}: UseUserSegmentsLoaderParams) => {
  const uniqueSelectedSegmentIds = useMemo(
    () => uniq(selectedSegmentIds),
    [selectedSegmentIds]
  );

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } =
    useInfiniteQueryUserSegments({
      params: {
        pageSize,
        environmentId
      },
      enabled
    });

  const fetchedSegments = useMemo<UserSegment[]>(
    () => data?.pages.flatMap(page => page.segments) ?? [],
    [data]
  );

  const fetchedSegmentIds = useMemo(
    () => new Set(fetchedSegments.map(segment => segment.id)),
    [fetchedSegments]
  );

  // Selected ids not present on already-loaded pages are resolved with a
  // direct point lookup instead of paging through the whole collection, so a
  // stale/deleted id (or a large environment) can't force every page to load.
  const unresolvedSelectedIds = useMemo(
    () => uniqueSelectedSegmentIds.filter(id => !fetchedSegmentIds.has(id)),
    [uniqueSelectedSegmentIds, fetchedSegmentIds]
  );

  const preloadResults = useQueries({
    queries: unresolvedSelectedIds.map(id => ({
      queryKey: [SEGMENT_DETAILS_QUERY_KEY, { id, environmentId }],
      queryFn: () => userSegmentFetcher({ id, environmentId }),
      enabled
    }))
  });

  const preloadedSegments = useMemo(
    () =>
      preloadResults
        .map(result => result.data?.segment)
        .filter((segment): segment is UserSegment => !!segment),
    [preloadResults.map(r => r.dataUpdatedAt).join(',')]
  );

  const isPreloadingSelection = preloadResults.some(result => result.isLoading);

  const { userSegments, selectedSegments } = useMemo(() => {
    const map = new Map(fetchedSegments.map(segment => [segment.id, segment]));
    preloadedSegments.forEach(segment => map.set(segment.id, segment));
    const resolvedSelected = uniqueSelectedSegmentIds
      .map(id => map.get(id))
      .filter((segment): segment is UserSegment => !!segment);
    return {
      userSegments: Array.from(map.values()),
      selectedSegments: resolvedSelected
    };
  }, [fetchedSegments, preloadedSegments, uniqueSelectedSegmentIds]);

  const loadMore = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const isInitialLoading = isLoading && userSegments.length === 0;
  const isResolvingSelection = isPreloadingSelection;

  const totalCount = Number(data?.pages[0]?.totalCount ?? 0);
  const hasNoSegments = !isLoading && totalCount === 0;

  return {
    userSegments,
    selectedSegments,
    isLoading,
    hasMore: !!hasNextPage,
    isLoadingMore: isFetchingNextPage,
    isInitialLoading,
    isResolvingSelection,
    hasNoSegments,
    loadMore
  };
};
