import { useInfiniteQueryAccounts } from '~/@queries/accounts';
import type { AppFiltersType } from '~/services/app-filters/types';
import { useAppFilters } from '~/services/app-filters/use-app-filters';

export const useFetch = ({
  filtersTypes,
  searchQuery
}: {
  filtersTypes?: AppFiltersType[];
  searchQuery?: string;
} = {}) => {
  const { filtersParams } = useAppFilters(filtersTypes);

  return useInfiniteQueryAccounts({
    params: {
      term: searchQuery,
      page_size: 15,
      acc_type: 5,
      ...(filtersTypes
        ? { service_type: filtersParams.actor_service_id }
        : undefined)
    }
  });
};
