import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon, XIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, useEffect, useState } from 'react';

import { classNames } from '../../utils/css';

export interface FilterChipProps {
  label: string;
  onRemove: () => void;
}

export const FilterChip: FC<FilterChipProps> = memo(({ label, onRemove }) => {
  return (
    <div className="flex items-center bg-gray-200 text-gray-700 text-sm px-1">
      {label}
      <button
        type="button"
        className={classNames(
          'border rounded-sm',
          'px-1 text-sm flex items-center',
          'text-gray-700'
        )}
        disabled={false}
        onClick={onRemove}
      >
        <XIcon className="ml-1 h-4 w-4" aria-hidden="true" />
      </button>
    </div>
  );
});
