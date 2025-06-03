import { useCallback, useMemo } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { useLocation } from '@tanstack/react-router';
// import { useLocation, useNavigate } from 'react-router-dom';
import queryString, { ParsedQuery } from 'query-string';

export type SearchParams = ParsedQuery<string>;

export function useSearchParams() {
  const navigate = useNavigate();
  const location = useLocation();
  console.log(location)
  const searchOptions = useMemo<SearchParams>((): SearchParams => {
    return queryString.parse(location.search);
  }, [location.search]);

  const onChangSearchParams = useCallback(
    (options: Record<string, string | number | boolean | string[]>) => {
      navigate({
        to: `${location.pathname}?${decodeURIComponent(stringifyParams(options))}`,
        replace: true
      });
    },
    [navigate, location]
  );

  return { searchOptions, onChangSearchParams };
}

export const stringifyParams = queryString.stringify;
