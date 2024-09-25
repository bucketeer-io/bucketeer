import { useMemo } from 'react';
import PaginationActions from './pagination-actions';
import PaginationCell from './pagination-cell';
import PaginationCount from './pagination-count';
import PaginationGroup from './pagination-group';

export type PaginationProps = {
  cursor: number;
  pageSize: number;
  totalCount: number;
  setCursor: (cursor: number) => void;
  cb?: () => void;
};

export type Props = {
  paginationProps?: PaginationProps;
};

const Pagination = ({ paginationProps }: Props) => {
  const { cursor, pageSize, totalCount, setCursor, cb } = useMemo(() => {
    if (paginationProps) return paginationProps;
    return {} as PaginationProps;
  }, [paginationProps]);

  return (
    <div className="flex items-center justify-between">
      <PaginationCount
        totalItems={totalCount}
        value={totalCount < pageSize ? totalCount : pageSize}
      />
      <PaginationActions
        pageIndex={cursor === 0 ? 1 : cursor / pageSize + 1}
        totalItems={totalCount}
        itemsPerPage={pageSize}
        onPageChange={page => {
          setCursor(pageSize * (page - 1));
          if (cb) cb();
        }}
      />
    </div>
  );
};

Pagination.Cell = PaginationCell;
Pagination.Group = PaginationGroup;
Pagination.Actions = PaginationActions;
Pagination.Count = PaginationCount;

export default Pagination;
