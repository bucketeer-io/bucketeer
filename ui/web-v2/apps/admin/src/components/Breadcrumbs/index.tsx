import { HomeIcon } from '@heroicons/react/solid';
import React, { FC, memo } from 'react';
import { Link } from 'react-router-dom';

import { PAGE_PATH_ROOT } from '../../constants/routing';
import { classNames } from '../../utils/css';

export interface Pages {
  name: string;
  path: string;
  current: boolean;
}

export interface BreadcrumbsProps {
  pages: Pages[];
}

export const Breadcrumbs: FC<BreadcrumbsProps> = memo(({ pages }) => {
  return (
    <nav className="flex text-xs" aria-label="Breadcrumb">
      <ol className="flex items-center">
        <li>
          <div>
            <Link to={PAGE_PATH_ROOT} className="">
              <HomeIcon
                className="flex-shrink-0 h-4 w-4 text-gray-500"
                aria-hidden="true"
              />
              <span className="sr-only">Home</span>
            </Link>
          </div>
        </li>
        {pages.map((page) => (
          <li key={page.name}>
            <div className="flex items-center space-x-3">
              <div
                className={classNames('ml-4 text-gray-500 hover:text-gray-700')}
              >
                {'/'}
              </div>
              <Link
                to={page.path}
                className={`${
                  page.current
                    ? 'text-gray-500 hover:text-gray-700'
                    : 'text-blue-500 underline'
                } ml-4`}
              >
                {page.name}
              </Link>
            </div>
          </li>
        ))}
      </ol>
    </nav>
  );
});
