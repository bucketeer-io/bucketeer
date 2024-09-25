/* eslint-disable @typescript-eslint/no-explicit-any */
import { useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

export function useQuery() {
  const { search } = useLocation();
  return useMemo(() => new URLSearchParams(search), [search]);
}

export function useAddQuery() {
  const navigate = useNavigate();

  const addQuery = (query: URLSearchParams, params: Record<string, any>) => {
    const searchParams = new URLSearchParams(query);
    Object.keys(params).forEach(key => {
      searchParams.set(key, params[key]);
    });
    navigate(`?${searchParams.toString()}`, { replace: true });
  };

  return {
    addQuery
  };
}
