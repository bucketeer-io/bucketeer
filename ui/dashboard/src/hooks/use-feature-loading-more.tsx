import { useState, useMemo, useEffect, useCallback } from 'react';
import { useQueryFeatures } from '@queries/features';
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
  pageSize = 50,
  selectedFlagIds,
  currentFeatureId,
  filterSelected = false
}: UseFeatureFlagsLoaderParams) => {
  const [cursor, setCursor] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [flags, setFlags] = useState<Feature[]>([]);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  /**
   * Cache selected flags so labels persist across searches
   */
  const [selectedFlagsCache, setSelectedFlagsCache] = useState<
    Map<string, Feature>
  >(new Map());

  const { data, isLoading } = useQueryFeatures({
    params: {
      cursor: String(cursor),
      pageSize,
      searchKeyword: searchQuery,
      environmentId,
      archived: false
    }
  });

  /**
   * Debounced search
   */
  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        setCursor(0);
        setSearchQuery(value);
      }, 300),
    []
  );

  const onSearchChange = useCallback(
    (value: string) => {
      if (!value) {
        debouncedSearch.cancel();
        setCursor(0);
        setSearchQuery('');
      } else {
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  const loadMore = useCallback(() => {
    setIsLoadingMore(true);
    setCursor(prev => prev + pageSize);
  }, [pageSize]);

  /**
   * Accumulate / reset flags
   */
  useEffect(() => {
    if (!data?.features) return;

    setFlags(prev => {
      if (cursor === 0) {
        return data.features;
      }

      const existingIds = new Set(prev.map(f => f.id));
      const newFlags = data.features.filter(f => !existingIds.has(f.id));

      return newFlags.length ? [...prev, ...newFlags] : prev;
    });
  }, [data, cursor]);

  useEffect(() => {
    const total = Number(data?.totalCount ?? 0);
    // Stop loading more when cursor >= totalCount
    const nextCursor = cursor + pageSize;
    setHasMore(nextCursor < total);
  }, [data, cursor, pageSize]);
  /**
   * Cache flags for selected IDs
   */
  useEffect(() => {
    setIsLoadingMore(false);
    if (!flags.length) return;
    setSelectedFlagsCache(prev => {
      const hasNewFlags = flags.some(f => !prev.has(f.id));
      if (!hasNewFlags) return prev; // Skip if no changes

      const cache = new Map(prev);
      flags.forEach(flag => cache.set(flag.id, flag));
      return cache;
    });
  }, [flags]);

  /**
   * Merge cached selected flags with API flags
   */
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

  const filteredFlags = useMemo(() => {
    if (!filterSelected) return allAvailableFlags;

    const apiReturnedIds = new Set(data?.features.map(f => f.id));
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
    data,
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

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const isInitialLoading =
    isLoading && !searchQuery && cursor === 0 && flags.length === 0;

  return {
    flags,
    allAvailableFlags,
    isLoading,
    hasMore,
    isLoadingMore,
    isInitialLoading,
    data,
    remainingFlagOptions,
    filteredFlags,
    searchQuery,
    loadMore,
    onSearchChange
  };
};
