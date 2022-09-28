import { Listbox, Transition } from '@headlessui/react';
import { CheckIcon, SelectorIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';

export interface SortItem {
  readonly key: string;
  readonly message?: string;
}

export interface SortSelectProps {
  sortKey: string;
  sortItems: SortItem[];
  onChange: (sort: string) => void;
}

export const SortSelect: FC<SortSelectProps> = memo(
  ({ sortKey, sortItems, onChange }) => {
    const { formatMessage: f } = useIntl();
    const selected = sortItems.find((item) => item.key == sortKey);

    return (
      <Listbox
        value={selected.key}
        onChange={(value) => {
          onChange(value);
        }}
      >
        <div className="relative">
          <Listbox.Button
            className={classNames(
              'w-full h-10 py-2 pl-3 pr-10',
              'text-left bg-white',
              'cursor-default',
              'rounded-md',
              'hover:bg-gray-100',
              'text-gray-700',
              'focus:outline-none focus-visible:ring-2 focus-visible:ring-opacity-75 focus-visible:ring-white focus-visible:ring-offset-orange-300 focus-visible:ring-offset-2 focus-visible:border-indigo-500',
              'text-sm'
            )}
          >
            <span className="block truncate">
              {`${f(messages.sort)}: `}
              {selected.message}
            </span>
            <span
              className={classNames(
                'absolute',
                'inset-y-0 right-0',
                'flex items-center',
                'pr-2 pointer-events-none'
              )}
            >
              <SelectorIcon
                className="w-5 h-5 text-gray-700"
                aria-hidden="true"
              />
            </span>
          </Listbox.Button>
          <Transition
            as={Fragment}
            leave="transition ease-in duration-100"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Listbox.Options
              className={classNames(
                'absolute',
                'w-full',
                'py-1 mt-1 overflow-auto',
                'bg-white text-sm z-50',
                'rounded-md',
                'shadow-lg max-h-60 ring-1 ring-black',
                'ring-opacity-5 focus:outline-none',
                'text-gray-700'
              )}
            >
              {sortItems.map((item, idx) => (
                <Listbox.Option
                  key={idx}
                  value={item.key}
                  className={classNames(
                    'cursor-default select-none relative py-2 pl-10 pr-4',
                    'text-sm hover:bg-gray-100',
                    'text-gray-700'
                  )}
                >
                  {({ selected }) => (
                    <>
                      {selected ? (
                        <span className="absolute inset-y-0 left-0 flex items-center pl-3">
                          <CheckIcon className="w-5 h-5" aria-hidden="true" />
                        </span>
                      ) : null}
                      <span className="block truncate">{item.message}</span>
                    </>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </Transition>
        </div>
      </Listbox>
    );
  }
);
