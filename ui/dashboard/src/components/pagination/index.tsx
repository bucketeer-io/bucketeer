import { useMemo } from 'react';
import { LIST_PAGE_SIZE } from 'constants/app';
import { cn } from 'utils/style';
import PaginationActions from './pagination-actions';
import PaginationCell from './pagination-cell';
import PaginationCount from './pagination-count';
import PaginationGroup from './pagination-group';

export type PaginationProps = {
  page: number;
  pageSize?: number;
  totalCount: number;
  onChange: (page: number) => void;
};

const Pagination = ({
  pageSize = LIST_PAGE_SIZE,
  totalCount,
  page,
  onChange
}: PaginationProps) => {
  const cursor = pageSize * (page - 1);
  const isShowPaginationAction = useMemo(
    () => totalCount > pageSize,
    [totalCount, pageSize]
  );

  return (
    <div
      className={cn('flex items-center justify-between', {
        'justify-end': !isShowPaginationAction
      })}
    >
      {totalCount > 0 && (
        <PaginationCount
          totalItems={totalCount}
          value={totalCount < pageSize ? totalCount : pageSize}
        />
      )}
      {isShowPaginationAction && (
        <PaginationActions
          pageIndex={cursor === 0 ? 1 : cursor / pageSize + 1}
          totalItems={totalCount}
          itemsPerPage={pageSize}
          onPageChange={onChange}
        />
      )}
    </div>
  );
};

Pagination.Cell = PaginationCell;
Pagination.Group = PaginationGroup;
Pagination.Actions = PaginationActions;
Pagination.Count = PaginationCount;

export default Pagination;
