import { Transition } from '@headlessui/react';
import {
  CheckCircleIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/outline';
import { XIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, SyntheticEvent } from 'react';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { IToast, removeToast, selectAll } from '../../modules/toasts';
import { classNames } from '../../utils/css';

const AUTO_HIDE_DURATION = 5000;

export const Toasts: FC = memo(function Toasts() {
  const dispatch = useDispatch();
  const { formatMessage: f } = useIntl();
  const toasts = useSelector<AppState, IToast[]>((state) =>
    selectAll(state.toasts)
  );

  return (
    <>
      {toasts.map((item, idx) => {
        const handleClose = (_: SyntheticEvent<unknown, Event>): void => {
          dispatch(removeToast({ id: item.id }));
        };
        if (item.severity !== 'error') {
          setTimeout(
            () => dispatch(removeToast({ id: item.id })),
            AUTO_HIDE_DURATION
          );
        }
        return (
          <div
            key={idx}
            aria-live="assertive"
            className={classNames(
              'fixed inset-0 flex items-end px-4 py-6',
              'pointer-events-none sm:p-6 sm:items-start z-50'
            )}
          >
            <div
              className={classNames(
                'w-full flex flex-col items-center space-y-4 sm:items-end'
              )}
            >
              <Transition
                show={true}
                as={Fragment}
                enter="transform ease-out duration-300 transition"
                enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
                enterTo="translate-y-0 opacity-100 sm:translate-x-0"
                leave="transition ease-in duration-100"
                leaveFrom="opacity-100"
                leaveTo="opacity-0"
              >
                <div
                  className={classNames(
                    'border-l-4',
                    item.severity === 'error' && 'border-red-400',
                    item.severity === 'warning' && 'border-yellow-400',
                    item.severity === 'success' && 'border-green-400',
                    'max-w-sm w-full bg-white shadow-lg rounded-lg',
                    'pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden'
                  )}
                >
                  <div className="p-4">
                    <div className="flex items-center">
                      <div className="flex-shrink-0">
                        {item.severity === 'error' && (
                          <ExclamationCircleIcon
                            className="h-6 w-6 text-red-400"
                            aria-hidden="true"
                          />
                        )}
                        {item.severity === 'warning' && (
                          <ExclamationCircleIcon
                            className="h-6 w-6 text-yellow-400"
                            aria-hidden="true"
                          />
                        )}
                        {item.severity === 'success' && (
                          <CheckCircleIcon
                            className="h-6 w-6 text-green-400"
                            aria-hidden="true"
                          />
                        )}
                      </div>
                      <div className="ml-3 w-0 flex-1 pt-0.5">
                        <p className="text-sm font-medium">
                          <span className="text-red-400">
                            {item.severity === 'error' && f(messages.error)}
                          </span>
                          <span className="text-500-400">
                            {item.severity === 'success' && f(messages.success)}
                          </span>
                          <span className="text-green-400">
                            {item.severity === 'warning' && f(messages.warning)}
                          </span>
                        </p>
                        <p className="mt-1 text-sm text-gray-500">
                          {item.message}
                        </p>
                      </div>
                      <div className="ml-4 flex-shrink-0 flex">
                        <button
                          className={classNames(
                            'bg-white rounded-md inline-flex text-gray-400 ',
                            'hover:text-gray-500',
                            'focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                          )}
                          onClick={handleClose}
                        >
                          <span className="sr-only">Close</span>
                          <XIcon className="h-5 w-5" aria-hidden="true" />
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </Transition>
            </div>
          </div>
        );
      })}
    </>
  );
});
