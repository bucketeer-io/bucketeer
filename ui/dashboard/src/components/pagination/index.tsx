import PaginationActions, {
  PaginationActionsProps
} from './pagination-actions';
import PaginationCell from './pagination-cell';
import PaginationCount from './pagination-count';
import PaginationGroup from './pagination-group';

export type PaginationProps = PaginationActionsProps;

const Pagination = ({ totalItems, itemsPerPage }: PaginationProps) => {
  return (
    <div className="flex items-center justify-between">
      <PaginationCount totalItems={totalItems} />
      <PaginationActions totalItems={totalItems} itemsPerPage={itemsPerPage} />
    </div>
  );
};

Pagination.Cell = PaginationCell;
Pagination.Group = PaginationGroup;
Pagination.Actions = PaginationActions;
Pagination.Count = PaginationCount;

export default Pagination;
