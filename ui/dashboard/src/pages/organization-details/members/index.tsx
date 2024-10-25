import { useParams } from 'react-router-dom';
import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickBy';
import { OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import Filter from 'elements/filter';
import { OrganizationMembersFilters } from '../types';
import CollectionLoader from './collection-loader';

const OrganizationMembers = () => {
  const { organizationId } = useParams();
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [filters, setFilters] = usePartialState<OrganizationMembersFilters>({
    page: Number(searchOptions.page) || 1,
    orderBy: (searchOptions.orderBy as OrderBy) || 'CREATED_AT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'DESC',
    searchQuery: (searchOptions.searchQuery as string) || ''
  });

  const onChangeFilters = (values: Partial<OrganizationMembersFilters>) => {
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

export default OrganizationMembers;
