import { useState, useMemo, useEffect, useCallback } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { LIST_PAGE_SIZE } from 'constants/app';
import { debounce } from 'lodash';

type UseAccountsLoaderParams = {
  organizationId: string;
  environmentId?: string;
  environmentRole?: number;
  pageSize?: number;
};

export const useAccountsLoader = ({
  organizationId,
  environmentId,
  environmentRole,
  pageSize = LIST_PAGE_SIZE
}: UseAccountsLoaderParams) => {
  const [cursor, setCursor] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [emails, setEmails] = useState<string[]>([]);
  const [hasMore, setHasMore] = useState(true);

  const { data, isLoading, isFetching } = useQueryAccounts({
    params: {
      cursor: String(cursor),
      pageSize,
      searchKeyword: searchQuery,
      organizationId,
      environmentId,
      environmentRole
    }
  });

  const resetSearchState = useCallback(() => {
    setEmails([]);
    setHasMore(true);
  }, []);

  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        resetSearchState();
        setCursor(0);
        setSearchQuery(value);
      }, 300),
    [resetSearchState]
  );

  const onSearchChange = useCallback(
    (value: string) => {
      if (!value) {
        debouncedSearch.cancel();
        resetSearchState();
        setCursor(0);
        setSearchQuery('');
      } else {
        debouncedSearch(value);
      }
    },
    [debouncedSearch, resetSearchState]
  );

  const loadMore = useCallback(() => {
    setCursor(prev => prev + pageSize);
  }, [pageSize]);

  useEffect(() => {
    if (!data?.accounts) return;

    setEmails(prev => {
      if (cursor === 0) {
        return data.accounts.map(a => a.email);
      }
      const existingSet = new Set(prev);
      const newEmails = data.accounts
        .map(a => a.email)
        .filter(e => !existingSet.has(e));
      return newEmails.length ? [...prev, ...newEmails] : prev;
    });
  }, [data]);

  useEffect(() => {
    const total = Number(data?.totalCount ?? 0);
    const nextCursor = cursor + pageSize;
    setHasMore(nextCursor < total);
  }, [data, cursor, pageSize]);

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const emailOptions = useMemo(
    () => emails.map(email => ({ label: email, value: email })),
    [emails]
  );

  const isInitialLoading =
    isLoading && !searchQuery && cursor === 0 && emails.length === 0;

  return {
    emails,
    emailOptions,
    isLoading,
    hasMore,
    isLoadingMore: isFetching && cursor > 0,
    isInitialLoading,
    loadMore,
    onSearchChange
  };
};
