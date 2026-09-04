import { useCallback, useEffect, useMemo, useRef } from 'react';
import { useInfiniteQueryUserSegments } from '@queries/user-segments';
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
  const selectedSegmentsCacheRef = useRef<Map<string, UserSegment>>(new Map());
  const cachedEnvironmentIdRef = useRef(environmentId);

  const uniqueSelectedSegmentIds = useMemo(
    () => uniq(selectedSegmentIds),
    [selectedSegmentIds]
  );

  if (cachedEnvironmentIdRef.current !== environmentId) {
    cachedEnvironmentIdRef.current = environmentId;
    selectedSegmentsCacheRef.current = new Map();
  }

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

  const { userSegments, selectedSegments } = useMemo(() => {
    const map = new Map(fetchedSegments.map(segment => [segment.id, segment]));
    const cache = selectedSegmentsCacheRef.current;
    const resolvedSelected: UserSegment[] = [];
    uniqueSelectedSegmentIds.forEach(id => {
      const segment = map.get(id) ?? cache.get(id);
      if (segment) {
        map.set(id, segment);
        resolvedSelected.push(segment);
      }
    });
    return {
      userSegments: Array.from(map.values()),
      selectedSegments: resolvedSelected
    };
  }, [fetchedSegments, uniqueSelectedSegmentIds]);

  useEffect(() => {
    const cache = selectedSegmentsCacheRef.current;
    selectedSegments.forEach(segment => cache.set(segment.id, segment));
  }, [selectedSegments]);

  const loadMore = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const hasUnresolvedSelection =
    selectedSegments.length < uniqueSelectedSegmentIds.length;

  // Keep fetching in the background until every selected id is resolved,
  // so a selection beyond page 1 doesn't get silently dropped from the
  // label/count before the user ever opens the dropdown.
  useEffect(() => {
    if (hasUnresolvedSelection) loadMore();
  }, [hasUnresolvedSelection, loadMore]);

  const isInitialLoading = isLoading && userSegments.length === 0;
  const isResolvingSelection = isFetchingNextPage && hasUnresolvedSelection;

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
