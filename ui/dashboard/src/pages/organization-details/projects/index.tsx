import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickby';
import { OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Filter from 'elements/filter';
import { OrganizationProjectFilters } from '../types';

const OrganizationProjects = () => {
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [filters, setFilters] = usePartialState<OrganizationProjectFilters>({
    page: Number(searchOptions.page) || 1,
    orderBy: (searchOptions.orderBy as OrderBy) || 'DEFAULT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'ASC',
    searchQuery: (searchOptions.searchQuery as string) || ''
  });

  const onChangeFilters = (values: Partial<OrganizationProjectFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  return (
    <>
      <Filter
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
    </>
  );
};

export default OrganizationProjects;
