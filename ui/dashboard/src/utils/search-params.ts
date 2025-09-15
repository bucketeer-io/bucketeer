import { useCallback, useMemo } from 'react';
import { NavigateOptions, useLocation, useNavigate } from 'react-router-dom';
import queryString, { ParsedQuery } from 'query-string';

export type SearchParams = ParsedQuery<string>;

export function useSearchParams() {
  const navigate = useNavigate();
  const location = useLocation();

  const searchOptions = useMemo<SearchParams>((): SearchParams => {
    return queryString.parse(location.search);
  }, [navigate, location]);

  const onChangSearchParams = useCallback(
    (
      options: Record<string, string | number | boolean | string[]>,
      state?: NavigateOptions['state']
    ) => {
      navigate(
        `${location.pathname}?${decodeURIComponent(stringifyParams(options))}`,
        {
          replace: true,
          state
        }
      );
    },
    [navigate, location]
  );

  return { searchOptions, onChangSearchParams };
}

export const stringifyParams = queryString.stringify;
