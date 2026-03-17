import { memo } from 'react';
import { SortingState } from '@tanstack/react-table';
import { getCurrentEnvironment, useAuth } from 'auth';
import { sortingListFields } from 'constants/collection';
import { useScreen } from 'hooks';
import { UserSegment } from '@types';
import { isNotEmpty } from 'utils/data-type';
import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import TableListContent from 'elements/table-list-content';
import { CardCollection } from '../collection-layout/card-collection';
import { useColumns } from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';
import { UserSegmentsActionsType, UserSegmentsFilters } from '../types';
import { useFetchSegments } from './use-fetch-segment';

const CollectionLoader = memo(
  ({
    getUploadingStatus,
    onAdd,
    filters,
    setFilters,
    onActionHandler,
    onClearFilters
  }: {
    getUploadingStatus: (segment: UserSegment) => boolean | undefined;
    onAdd?: () => void;
    filters: UserSegmentsFilters;
    setFilters: (values: Partial<UserSegmentsFilters>) => void;
    organizationIds?: string[];
    onActionHandler: (
      value: UserSegment,
      type: UserSegmentsActionsType
    ) => void;
    onClearFilters: () => void;
  }) => {
    const columns = useColumns({ getUploadingStatus, onActionHandler });
    const { consoleAccount } = useAuth();
    const { fromMobileScreen } = useScreen();
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);

    const {
      data: collection,
      isLoading,
      refetch,
      isError
    } = useFetchSegments({
      ...filters,
      environmentId: currentEnvironment.id
    });

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

    const userSegments = collection?.segments || [];
    const totalCount = Number(collection?.totalCount) || 0;

    const emptyState = (
      <CollectionEmpty
        data={userSegments}
        isFilter={isNotEmpty(filters?.isInUseStatus)}
        searchQuery={filters.searchQuery as string}
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
            data={userSegments}
            columns={columns}
            onSortingChange={onSortingChangeHandler}
            emptyCollection={emptyState}
          />
        ) : (
          <CardCollection
            isLoading={isLoading}
            data={userSegments}
            getUploadingStatus={getUploadingStatus}
            onActions={onActionHandler}
            emptyCollection={emptyState}
          />
        )}
        {!isLoading && (
          <Pagination
            page={filters.page as number}
            totalCount={totalCount}
            onChange={page => setFilters({ page })}
          />
        )}
      </TableListContent>
    );
  }
);

export default CollectionLoader;
