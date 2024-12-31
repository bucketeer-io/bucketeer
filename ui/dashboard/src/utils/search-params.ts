import { useCallback, useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import queryString, { ParsedQuery } from 'query-string';

export type SearchParams = ParsedQuery<string>;

export function useSearchParams() {
  const navigate = useNavigate();
  const location = useLocation();

  const searchOptions = useMemo<SearchParams>((): SearchParams => {
    return queryString.parse(location.search);
  }, [location.search]);

  const onChangSearchParams = useCallback(
    (options: Record<string, string | number | boolean>) => {
      navigate(`${location.pathname}?${stringifyParams(options)}`, {
        replace: true
      });
    },
    [navigate]
  );

  return { searchOptions, onChangSearchParams };
}

export const stringifyParams = queryString.stringify;
