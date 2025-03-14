// import { LIST_PAGE_SIZE } from 'constants/app';
// import Pagination from 'components/pagination';
import CollectionEmpty from 'elements/collection/collection-empty';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from '../collection-layout/empty-collection';
import GridViewCollection from '../collection-layout/grid-view-collection';
import { FlagActionType, FlagsTemp } from '../types';

export const mockFlags: FlagsTemp[] = [
  {
    id: 'flag-1',
    name: 'Flag using boolean',
    type: 'boolean',
    status: 'active',
    tags: ['Android'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-2',
    name: 'Flag using string',
    type: 'string',
    status: 'no_activity',
    tags: ['Web'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-3',
    name: 'Flag using number',
    type: 'number',
    status: 'new',
    tags: ['Android'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-4',
    name: 'Flag using json',
    type: 'json',
    status: 'no_activity',
    tags: ['IOS'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  }
];

const CollectionLoader = ({
  onHandleActions
}: {
  onHandleActions: (item: FlagsTemp, type: FlagActionType) => void;
}) => {
  const isError = false;

  const emptyState = (
    <CollectionEmpty
      data={mockFlags}
      searchQuery={''}
      onClear={() => {}}
      empty={<EmptyCollection onAdd={() => {}} />}
    />
  );

  return isError ? (
    <PageLayout.ErrorState onRetry={() => {}} />
  ) : (
    <>
      <GridViewCollection data={mockFlags} onActions={onHandleActions} />

      {/* {totalCount > LIST_PAGE_SIZE && !isLoading && (
        <Pagination
          page={filters.page}
          totalCount={totalCount}
          onChange={page => setFilters({ page })}
        />
      )} */}
    </>
  );
};

export default CollectionLoader;
