import { SearchIcon } from '@heroicons/react/solid';
import React, { FC, memo } from 'react';

import { classNames } from '../../utils/css';

export interface SearchInputProps {
  placeholder: string;
  onChange: (query: string) => void;
}

export const SearchInput: FC<SearchInputProps> = memo(
  ({ placeholder: placeHolder, onChange }) => {
    return (
      <div className="relative h-10">
        <div
          className={classNames(
            'absolute inset-y-0 left-0 pl-3 flex',
            'items-center pointer-events-none'
          )}
        >
          <SearchIcon className="h-5 w-5 text-gray-400" aria-hidden="true" />
        </div>
        <input
          type="text"
          name="search"
          id="search"
          className={classNames(
            'focus:ring-purple-600 focus:border-indigo-500',
            'block w-full h-full pl-10 sm:text-sm',
            'border border-gray-300 rounded-md',
            'text-gray-700'
          )}
          placeholder={placeHolder}
          onChange={(e) => onChange(e.target.value)}
        />
      </div>
    );
  }
);
