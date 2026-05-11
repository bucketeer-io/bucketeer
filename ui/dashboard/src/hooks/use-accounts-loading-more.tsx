import { useCallback, useEffect, useMemo, useState } from 'react';
import { accountsFetcher } from '@api/account';
import {
  ACCOUNTS_QUERY_KEY,
  useInfiniteQueryAccounts
} from '@queries/accounts';
import { useQueries } from '@tanstack/react-query';
import { LIST_PAGE_SIZE } from 'constants/app';
import { debounce } from 'lodash';
import { Account } from '@types';

type UseAccountsLoaderParams = {
  organizationId: string;
  environmentId?: string;
  environmentRole?: number;
  pageSize?: number;
  enabled?: boolean;
  preloadEmails?: string[];
};

export const useAccountsLoader = ({
  organizationId,
  environmentId,
  environmentRole,
  pageSize = LIST_PAGE_SIZE,
  enabled = true,
  preloadEmails = []
}: UseAccountsLoaderParams) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [isTyping, setIsTyping] = useState(false);

  const { data, isLoading, isFetchingNextPage, hasNextPage, fetchNextPage } =
    useInfiniteQueryAccounts({
      params: {
        pageSize,
        searchKeyword: searchQuery,
        organizationId,
        environmentId,
        environmentRole
      },
      enabled
    });

  const accumulatedAccounts = useMemo<Account[]>(
    () => data?.pages.flatMap(page => page.accounts) ?? [],
    [data]
  );

  const uniquePreloadEmails = useMemo(
    () => Array.from(new Set(preloadEmails.filter(Boolean))),
    [preloadEmails.join(',')]
  );

  const preloadResults = useQueries({
    queries: uniquePreloadEmails.map(email => ({
      queryKey: [
        ACCOUNTS_QUERY_KEY,
        {
          searchKeyword: email,
          organizationId,
          environmentId,
          environmentRole,
          cursor: '0',
          pageSize: 1
        }
      ],
      queryFn: () =>
        accountsFetcher({
          cursor: '0',
          pageSize: 1,
          searchKeyword: email,
          organizationId,
          environmentId,
          environmentRole
        })
    }))
  });

  const preloadedAccounts = useMemo(
    () => preloadResults.flatMap(r => r.data?.accounts ?? []),
    [preloadResults.map(r => r.dataUpdatedAt).join(',')]
  );

  const allAccounts = useMemo(() => {
    if (!preloadedAccounts.length) return accumulatedAccounts;
    const emailsInList = new Set(accumulatedAccounts.map(a => a.email));
    const missing = preloadedAccounts.filter(a => !emailsInList.has(a.email));
    return missing.length
      ? [...accumulatedAccounts, ...missing]
      : accumulatedAccounts;
  }, [accumulatedAccounts, preloadedAccounts]);

  const debouncedSearch = useMemo(
    () =>
      debounce((value: string) => {
        setSearchQuery(value);
        setIsTyping(false);
      }, 300),
    []
  );

  useEffect(() => () => debouncedSearch.cancel(), [debouncedSearch]);

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

  const getAccountLabel = useCallback(
    (email: string) => {
      const account = allAccounts.find(a => a.email === email);
      if (account?.firstName || account?.lastName)
        return `${account.firstName} ${account.lastName}`.trim();
      return account?.email ?? email;
    },
    [allAccounts]
  );

  const emailOptions = useMemo(
    () =>
      accumulatedAccounts.map(account => ({
        label: account.email,
        value: account.email
      })),
    [accumulatedAccounts]
  );

  const isInitialLoading =
    enabled && isLoading && accumulatedAccounts.length === 0;

  const isSearching =
    isTyping || (enabled && isLoading && accumulatedAccounts.length === 0);

  return {
    accounts: allAccounts,
    emailOptions,
    isLoading,
    hasMore: !!hasNextPage,
    isLoadingMore: isFetchingNextPage,
    isInitialLoading,
    isSearching,
    loadMore,
    onSearchChange,
    getAccountLabel
  };
};
