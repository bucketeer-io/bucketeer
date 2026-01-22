import { useCallback, useEffect, useMemo, useState } from 'react';
import { PAGE_SIZE } from 'constants/experiment';
import { debounce } from 'lodash';

export default function useLoadMore(
  initSearch = '',
  pageSize = PAGE_SIZE,
  debounceMs = 300
) {
  const [cursor, setCursor] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');

  const debouncedSearch = useMemo(
    () => debounce((v: string) => setSearchQuery(v), debounceMs),
    [debounceMs]
  );
  useEffect(() => {
    return () => debouncedSearch.cancel();
  }, [debouncedSearch]);

  const onSearchChange = useCallback(
    (value: string) => {
      setCursor(0);
      if (value === '') {
        debouncedSearch.cancel();
        setSearchQuery('');
      } else {
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  const loadMore = useCallback(() => {
    setCursor(prev => prev + pageSize);
  }, [pageSize]);

  const reset = useCallback(() => {
    debouncedSearch.cancel();
    setCursor(0);
    setSearchQuery(initSearch);
  }, [debouncedSearch]);

  return {
    cursor,
    pageSize,
    searchQuery,
    loadMore,
    onSearchChange,
    reset,
    setCursor,
    setSearchQuery
  } as const;
}
