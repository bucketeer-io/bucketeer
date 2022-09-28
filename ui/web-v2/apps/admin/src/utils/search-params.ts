import { parse, stringify, ParsedQuery } from 'query-string';
import { useMemo } from 'react';
import { useLocation } from 'react-router-dom';

export function useSearchParams(): SearchParams {
  const location = useLocation();
  return useMemo<SearchParams>((): SearchParams => {
    return parse(location.search);
  }, [location.search]);
}

export const stringifySearchParams = stringify;

export type SearchParams = ParsedQuery<string>;
