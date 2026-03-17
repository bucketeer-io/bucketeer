import { useState } from 'react';
import { useScreen } from 'hooks/use-screen';
import PaginationCell from './pagination-cell';
import PaginationGroup from './pagination-group';

export type PaginationActionsProps = {
  pageIndex: number;
  totalItems: number;
  itemsPerPage: number;
  onPageChange?: (page: number) => void;
};

const PaginationActions = ({
  totalItems,
  itemsPerPage,
  pageIndex,
  onPageChange
}: PaginationActionsProps) => {
  const { fromMobileScreen } = useScreen();
  const [currentPage, setCurrentPage] = useState(pageIndex);
  const totalPages = Math.ceil(totalItems / itemsPerPage);
  const maxVisibleButtons = fromMobileScreen ? 5 : 3;

  const cells = () => {
    let startPage, endPage;

    if (totalPages <= maxVisibleButtons) {
      startPage = 1;
      endPage = totalPages;
    } else {
      const maxPagesBeforeCurrentPage = Math.floor(maxVisibleButtons / 2);
      const maxPagesAfterCurrentPage = Math.ceil(maxVisibleButtons / 2) - 1;

      if (currentPage <= maxPagesBeforeCurrentPage) {
        startPage = 1;
        endPage = maxVisibleButtons;
      } else if (currentPage + maxPagesAfterCurrentPage >= totalPages) {
        startPage = totalPages - maxVisibleButtons + 1;
        endPage = totalPages;
      } else {
        startPage = currentPage - maxPagesBeforeCurrentPage;
        endPage = currentPage + maxPagesAfterCurrentPage;
      }
    }

    const pages = [];

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    if (startPage > 1) {
      pages.unshift('...');
      pages.unshift(1);
    }

    if (endPage < totalPages) {
      pages.push('...');
      pages.push(totalPages);
    }

    return pages;
  };

  const handleNext = () => handlePageChange(currentPage + 1);

  const handlePrevious = () => handlePageChange(currentPage - 1);

  const handlePageChange = (page?: number) => {
    if (page) {
      if (page < 1 || page > totalPages) return;
      setCurrentPage(page);
      if (onPageChange) onPageChange(page);
    }
  };

  const handleLast = () => handlePageChange(totalPages);

  const handleFirst = () => handlePageChange(1);

  const renderCell = cells();

  return (
    <div className="flex gap-4">
      <PaginationGroup>
        <PaginationCell
          variant="first"
          onClick={handleFirst}
          disabled={currentPage === 1}
        />
        <PaginationCell
          variant="previous"
          onClick={handlePrevious}
          disabled={currentPage === 1}
        />
      </PaginationGroup>
      <PaginationGroup>
        {renderCell.map((value, index) =>
          typeof value === 'string' ? (
            <span key={`${value}${index}`} className="text-gray-500">
              ...
            </span>
          ) : (
            <PaginationCell
              key={value}
              value={value}
              checked={currentPage === value}
              onClick={handlePageChange}
            />
          )
        )}
      </PaginationGroup>
      <PaginationGroup>
        <PaginationCell
          variant="next"
          onClick={handleNext}
          disabled={currentPage === totalPages}
        />
        <PaginationCell
          variant="last"
          onClick={handleLast}
          disabled={currentPage === totalPages}
        />
      </PaginationGroup>
    </div>
  );
};

export default PaginationActions;
