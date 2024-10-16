import { useParams } from 'react-router-dom';
import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickBy';
import { OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Filter from 'elements/filter';
import { OrganizationUsersFilters } from '../types';
import CollectionLoader from './collection-loader';

const OrganizationUsers = () => {
  const { organizationId } = useParams();
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [filters, setFilters] = usePartialState<OrganizationUsersFilters>({
    page: Number(searchOptions.page) || 1,
    orderBy: (searchOptions.orderBy as OrderBy) || 'DEFAULT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'ASC',
    searchQuery: (searchOptions.searchQuery as string) || ''
  });

  const onChangeFilters = (values: Partial<OrganizationUsersFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const filterParams = { ...filters, organizationId };

  return (
    <>
      <Filter
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <CollectionLoader filters={filterParams} setFilters={onChangeFilters} />
    </>
  );
};

export default OrganizationUsers;
