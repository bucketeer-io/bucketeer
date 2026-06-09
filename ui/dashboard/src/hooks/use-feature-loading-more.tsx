import { useState, useMemo, useEffect, useCallback } from 'react';
import { useInfiniteQueryFeatures } from '@queries/features';
import { LIST_PAGE_SIZE } from 'constants/app';
import { debounce } from 'lodash';
import { Feature } from '@types';

type UseFeatureFlagsLoaderParams = {
  environmentId: string;
  pageSize?: number;
  selectedFlagIds: readonly string[];
  currentFeatureId?: string;
  filterSelected?: boolean;
};

export const useFeatureFlagsLoader = ({
  environmentId,
  pageSize = LIST_PAGE_SIZE,
  selectedFlagIds,
  currentFeatureId,
  filterSelected = false
}: UseFeatureFlagsLoaderParams) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [isTyping, setIsTyping] = useState(false);

  const [selectedFlagsCache, setSelectedFlagsCache] = useState<
    Map<string, Feature>
  >(new Map());

  const {
    data,
    isLoading,
    isFetching,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage
  } = useInfiniteQueryFeatures({
    environmentId,
    searchKeyword: searchQuery,
    pageSize,
    archived: false
  });

  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        setIsTyping(false);
        setSearchQuery(value);
      }, 300),
    []
  );

  const onSearchChange = useCallback(
    (value: string) => {
      if (!value) {
        debouncedSearch.cancel();
        setIsTyping(false);
        setSearchQuery('');
      } else {
        setIsTyping(true);
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  const loadMore = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const flags = useMemo(
    () => data?.pages.flatMap(p => p.features) ?? [],
    [data]
  );

  useEffect(() => {
    if (!data) return;
    if (!flags.length) return;
    setSelectedFlagsCache(prev => {
      const hasNewFlags = flags.some(f => !prev.has(f.id));
      if (!hasNewFlags) return prev;
      const cache = new Map(prev);
      flags.forEach(flag => cache.set(flag.id, flag));
      return cache;
    });
  }, [data]);

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const allAvailableFlags = useMemo(() => {
    const map = new Map(
      flags.filter(f => f.id !== currentFeatureId).map(flag => [flag.id, flag])
    );

    selectedFlagIds.forEach(id => {
      if (!map.has(id)) {
        const cached = selectedFlagsCache.get(id);
        if (cached) {
          map.set(id, cached);
        }
      }
    });

    return Array.from(map.values());
  }, [flags, selectedFlagIds, selectedFlagsCache, currentFeatureId]);

  const lastPageFeatures = data?.pages.at(-1)?.features;

  const filteredFlags = useMemo(() => {
    if (!filterSelected) return allAvailableFlags;

    const apiReturnedIds = new Set(lastPageFeatures?.map(f => f.id));
    return allAvailableFlags.filter(flag => {
      if (flag.id === currentFeatureId) return false;
      if (
        searchQuery &&
        selectedFlagIds.includes(flag.id) &&
        !apiReturnedIds.has(flag.id)
      ) {
        return false;
      }
      return true;
    });
  }, [
    allAvailableFlags,
    lastPageFeatures,
    searchQuery,
    selectedFlagIds,
    currentFeatureId,
    filterSelected
  ]);

  const remainingFlagOptions = useMemo(() => {
    const selectedSet = new Set(selectedFlagIds);
    return filteredFlags.map(item => ({
      label: item.name,
      value: item.id,
      enabled: item.enabled,
      disabled: selectedSet.has(item.id)
    }));
  }, [filteredFlags, selectedFlagIds]);

  const isInitialLoading = isLoading && !searchQuery && flags.length === 0;
  const isSearching =
    !isTyping && isFetching && !isFetchingNextPage && !!searchQuery;

  return {
    flags,
    allAvailableFlags,
    isLoading,
    hasMore: !!hasNextPage,
    isLoadingMore: isFetchingNextPage,
    isInitialLoading,
    isSearching,
    data: data?.pages.at(-1),
    remainingFlagOptions,
    filteredFlags,
    searchQuery,
    loadMore,
    onSearchChange
  };
};
