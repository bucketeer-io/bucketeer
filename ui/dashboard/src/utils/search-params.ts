import { useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import queryString, { ParsedQuery } from 'query-string';

export function useSearchParams(): SearchParams {
  const location = useLocation();

  return useMemo<SearchParams>((): SearchParams => {
    return queryString.parse(location.search);
  }, [location.search]);
}

export const stringifySearchParams = queryString.stringify;

export type SearchParams = ParsedQuery<string>;
