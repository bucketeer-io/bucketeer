import { useState, useMemo, useEffect, useCallback } from 'react';
import { useQueryFeatures } from '@queries/features';
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
  const [cursor, setCursor] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [flags, setFlags] = useState<Feature[]>([]);
  const [hasMore, setHasMore] = useState(true);
  const [isTyping, setIsTyping] = useState(false);

  const [selectedFlagsCache, setSelectedFlagsCache] = useState<
    Map<string, Feature>
  >(new Map());

  const { data, isLoading, isFetching } = useQueryFeatures({
    params: {
      cursor: String(cursor),
      pageSize,
      searchKeyword: searchQuery,
      environmentId,
      archived: false
    }
  });

  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        setFlags([]);
        setHasMore(true);
        setCursor(0);
        setSearchQuery(value);
      }, 300),
    []
  );

  const onSearchChange = useCallback(
    (value: string) => {
      if (!value) {
        debouncedSearch.cancel();
        setIsTyping(false);
        setCursor(0);
        setSearchQuery('');
      } else {
        setIsTyping(true);
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  const loadMore = useCallback(() => {
    setCursor(prev => prev + pageSize);
  }, [pageSize]);

  useEffect(() => {
    if (!data?.features) return;

    setIsTyping(false);
    setFlags(prev => {
      if (cursor === 0) {
        return data.features;
      }
      const existingIds = new Set(prev.map(f => f.id));
      const newFlags = data.features.filter(f => !existingIds.has(f.id));
      return newFlags.length ? [...prev, ...newFlags] : prev;
    });

    if (data.totalCount != null) {
      const total = Number(data.totalCount);
      setHasMore(cursor + pageSize < total);
    }
  }, [data, cursor, pageSize]);

  useEffect(() => {
    if (!flags.length) return;
    setSelectedFlagsCache(prev => {
      const hasNewFlags = flags.some(f => !prev.has(f.id));
      if (!hasNewFlags) return prev;
      const cache = new Map(prev);
      flags.forEach(flag => cache.set(flag.id, flag));
      return cache;
    });
  }, [flags]);

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

  const isInitialLoading =
    isLoading && !searchQuery && cursor === 0 && flags.length === 0;

  const isSearching = isTyping || (isFetching && cursor === 0 && !!searchQuery);

  return {
    flags,
    allAvailableFlags,
    isLoading,
    hasMore,
    isLoadingMore: isFetching && cursor > 0,
    isInitialLoading,
    isSearching,
    data,
    remainingFlagOptions,
    filteredFlags,
    searchQuery,
    loadMore,
    onSearchChange
  };
};
