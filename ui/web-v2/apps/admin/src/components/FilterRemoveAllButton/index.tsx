import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon, XIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, useEffect, useState } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';

export interface FilterRemoveAllButtonProps {
  onClick: () => void;
}

export const FilterRemoveAllButtonProps: FC<FilterRemoveAllButtonProps> = memo(
  ({ onClick }) => {
    const { formatMessage: f } = useIntl();

    return (
      <button
        type="button"
        className={classNames(
          'border rounded-sm',
          'px-1 text-sm flex items-center',
          'text-gray-600'
        )}
        disabled={false}
        onClick={onClick}
      >
        {f(messages.button.clearAll)}
        <XIcon className="ml-1 h-4 w-4" aria-hidden="true" />
      </button>
    );
  }
);
