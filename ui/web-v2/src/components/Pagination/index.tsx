import {
  ArrowNarrowLeftIcon,
  ArrowNarrowRightIcon,
} from '@heroicons/react/solid';
import { FC, memo } from 'react';

import { classNames } from '../../utils/css';

interface PagingProps {
  maxPage: number;
  currentPage: number;
  onChange: (page: number) => void;
}

type Page = number | '...';

export const Pagination: FC<PagingProps> = memo(
  ({ maxPage, currentPage, onChange }) => {
    const pages: Page[] = [];
    pages.push(1);
    if (currentPage == 4) pages.push(2);
    if (currentPage >= 5) pages.push('...');
    if (currentPage >= 3) pages.push(currentPage - 1);
    if (currentPage >= 2 && currentPage <= maxPage - 1) pages.push(currentPage);
    if (currentPage <= maxPage - 2) pages.push(currentPage + 1);
    if (maxPage >= 7 && currentPage <= 4 && currentPage >= maxPage - 4)
      pages.push('...');
    if (currentPage <= maxPage - 4) pages.push('...');
    if (currentPage == maxPage - 3) pages.push(maxPage - 1);
    if (maxPage >= 2) pages.push(maxPage);

    const handleOnChange = (page) => {
      onChange(page);
      document.getElementById('header') &&
        document.getElementById('header').scrollIntoView();
    };
    return (
      <nav className="py-4 flex items-center justify-between">
        <div className="flex-1 flex justify-end">
          <button
            type="button"
            className={classNames(
              'disabled:opacity-50 disabled:cursor-auto disabled:pointer-events-none',
              'pb-2 pr-1 inline-flex items-center text-sm border-b-2 border-transparent',
              'hover:text-gray-700 hover:text-gray-700 hover:border-gray-300 hover:border-b-2',
              'font-medium text-gray-500'
            )}
            disabled={currentPage == 1}
            onClick={() => handleOnChange(currentPage - 1)}
          >
            <ArrowNarrowLeftIcon
              className="h-5 w-10 text-gray-400"
              aria-hidden="true"
            />
          </button>
        </div>
        <div className="flex">
          {pages.map((page: Page, idx: number) => {
            return (
              <button
                key={idx}
                type="button"
                className={`${
                  page == '...'
                    ? 'disabled:cursor-auto disabled:pointer-events-none'
                    : page == currentPage
                    ? 'text-primary border-b-2 border-primary cursor-auto pointer-events-none'
                    : 'text-gray-500 border-b-2 border-transparent hover:text-gray-700 hover:border-gray-300 hover:border-b-2'
                } pb-2 px-4 inline-flex items-center text-sm font-medium`}
                disabled={page == '...'}
                onClick={() => handleOnChange(page)}
              >
                {page}
              </button>
            );
          })}
        </div>
        <div className="flex-1 flex">
          <button
            type="button"
            className={classNames(
              'disabled:opacity-50 disabled:cursor-auto disabled:pointer-events-none',
              'pb-2 pl-1 inline-flex items-center text-sm border-b-2 border-transparent',
              'hover:text-gray-700 hover:text-gray-700 hover:border-gray-300 hover:border-b-2',
              'font-medium text-gray-500'
            )}
            disabled={currentPage == maxPage}
            onClick={() => handleOnChange(currentPage + 1)}
          >
            <ArrowNarrowRightIcon
              className="h-5 w-10 text-gray-400"
              aria-hidden="true"
            />
          </button>
        </div>
      </nav>
    );
  }
);
