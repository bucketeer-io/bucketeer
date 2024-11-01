import { useParams } from 'react-router-dom';
import { usePartialState } from 'hooks';
import pickBy from 'lodash/pickBy';
import { OrderBy, OrderDirection } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import CollectionLoader from 'pages/projects/collection-loader';
import { ProjectsFilters } from 'pages/projects/types';
import Filter from 'elements/filter';

const OrganizationProjects = () => {
  const { organizationId } = useParams();
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [filters, setFilters] = usePartialState<ProjectsFilters>({
    page: Number(searchOptions.page) || 1,
    orderBy: (searchOptions.orderBy as OrderBy) || 'CREATED_AT',
    orderDirection: (searchOptions.orderDirection as OrderDirection) || 'DESC',
    searchQuery: (searchOptions.searchQuery as string) || ''
  });

  const onChangeFilters = (values: Partial<ProjectsFilters>) => {
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
      <CollectionLoader
        filters={filters}
        organizationIds={[organizationId!]}
        setFilters={onChangeFilters}
        onActionHandler={() => {}}
      />
    </>
  );
};

export default OrganizationProjects;
