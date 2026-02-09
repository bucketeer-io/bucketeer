import { memo, useMemo } from 'react';
import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { useScreen } from 'hooks';
import { Project } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useFetchMembers } from 'pages/members/collection-loader/use-fetch-members';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { CardCollection } from '../collection-layout/card-collection';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { ProjectFilters } from '../types';
import { useFetchProjects } from './use-fetch-projects';

const CollectionLoader = memo(
  ({
    onAdd,
    filters,
    setFilters,
    organizationId,
    onActionHandler,
    onClearFilters
  }: {
    onAdd?: () => void;
    filters: ProjectFilters;
    setFilters: (values: Partial<ProjectFilters>) => void;
    organizationId?: string;
    onActionHandler: (value: Project) => void;
    onClearFilters: () => void;
  }) => {
    const { consoleAccount } = useAuth();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    const { data: accountCollection } = useFetchMembers({
      organizationId: currentEnvironment.organizationId,
      pageSize: 0
    });

    const accounts = useMemo(
      () => accountCollection?.accounts || [],
      [accountCollection]
    );
    const columns = useColumns({
      organizationId,
      accounts,
      currentEnvironment,
      onActionHandler
    });
    const { fromMobileScreen } = useScreen();
    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchProjects({ ...filters, organizationId });

    const onSortingChangeHandler = (sorting: SortingState) => {
      const updateOrderBy =
        sorting.length > 0
          ? sortingListFields[sorting[0].id]
          : sortingListFields.default;

      setFilters({
        orderBy: updateOrderBy,
        orderDirection: sorting[0]?.desc ? 'DESC' : 'ASC'
      });
    };

    const projects = collection?.projects || [];
    const totalCount = Number(collection?.totalCount) || 0;

    const emptyState = (
      <CollectionEmpty
        data={projects}
        searchQuery={filters.searchQuery}
        isFilter={isNotEmpty(filters?.disabled)}
        onClear={onClearFilters}
        empty={<EmptyCollection onAdd={onAdd} />}
      />
    );

    return isError ? (
      <PageLayout.ErrorState onRetry={refetch} />
    ) : (
      <TableListContent>
        {fromMobileScreen ? (
          <DataTable
            isLoading={isLoading}
            data={projects}
            columns={columns}
            onSortingChange={onSortingChangeHandler}
            emptyCollection={emptyState}
          />
        ) : (
          <CardCollection
            isLoading={isLoading}
            accounts={accounts}
            currentEnvironment={currentEnvironment}
            emptyCollection={emptyState}
            organizationId={organizationId}
            data={projects}
            onActions={onActionHandler}
          />
        )}
        {!isLoading && (
          <Pagination
            page={filters.page}
            totalCount={totalCount}
            onChange={page => setFilters({ page })}
          />
        )}
      </TableListContent>
    );
  }
);

export default CollectionLoader;
