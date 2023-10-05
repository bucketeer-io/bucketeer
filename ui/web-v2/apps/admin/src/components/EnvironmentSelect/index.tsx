import { Listbox, Transition } from '@headlessui/react';
import { CheckIcon, SelectorIcon } from '@heroicons/react/solid';
import { unwrapUndefinable } from 'option-t/lib/Undefinable/unwrap';
import React, {
  Fragment,
  FC,
  useState,
  memo,
  useCallback,
  useEffect,
} from 'react';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { useHistory } from 'react-router-dom';

import { PAGE_PATH_ROOT } from '../../constants/routing';
import { messages } from '../../lang/messages';
import {
  setCurrentEnvironment,
  useCurrentEnvironment,
  useEnvironments,
} from '../../modules/me';
import { EnvironmentV2 } from '../../proto/environment/environment_pb';
import { AppDispatch } from '../../store';
import { classNames } from '../../utils/css';

export const EnvironmentSelect: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const history = useHistory();
  const currenEnvironment = useCurrentEnvironment();
  const environments = useEnvironments();
  const initialEnvironment = unwrapUndefinable(
    environments.find((env) => env.id == currenEnvironment.id)
  );
  const [selected, setSelected] =
    useState<EnvironmentV2.AsObject>(initialEnvironment);

  const handleChange = useCallback(
    (value: string) => {
      const env = environments.find((env) => env.id == value);
      dispatch(setCurrentEnvironment(env.id));
      history.push(PAGE_PATH_ROOT);
    },
    [setCurrentEnvironment]
  );

  useEffect(() => {
    setSelected(currenEnvironment);
  }, [setSelected, currenEnvironment]);

  return (
    <div className="">
      <p className="mx-3">
        <span className="text-xs text-white">
          {f(messages.environment.select.label)}
        </span>
      </p>
      <Listbox value={selected.id} onChange={handleChange}>
        <div className="relative mt-1">
          <Listbox.Button
            className={classNames(
              'w-full py-2 pl-3 pr-10',
              'text-left bg-white',
              'rounded-md',
              'border',
              'cursor-default',
              'focus:outline-none',
              'focus-visible:ring-2',
              'focus-visible:ring-opacity-75',
              'focus-visible:ring-white focus-visible:ring-offset-orange-300',
              'focus-visible:ring-offset-2 focus-visible:border-indigo-500',
              'text-sm text-gray-700'
            )}
          >
            <span className="block truncate">
              {' '}
              {`(${selected.projectId}) ${selected.name}`}
            </span>
            <span
              className={classNames(
                'absolute inset-y-0 right-0',
                'flex items-center',
                'pr-2 pointer-events-none'
              )}
            >
              <SelectorIcon
                className="w-5 h-5 text-gray-400"
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
                'absolute min-w-full max-w-fit py-1 mt-1',
                'overflow-auto text-sm text-gray-700 bg-white rounded-md',
                'shadow-lg max-h-60 ring-1 ring-black z-20',
                'ring-opacity-5 focus:outline-none whitespace-nowrap'
              )}
            >
              {environments.map((env, envIdx) => (
                <Listbox.Option
                  key={envIdx}
                  value={env.id}
                  className={classNames(
                    'cursor-default select-none relative',
                    'py-2 pl-10 pr-4 hover:bg-gray-100'
                  )}
                >
                  {({ selected }) => (
                    <>
                      {selected ? (
                        <span
                          className={classNames(
                            'absolute inset-y-0 left-0',
                            'flex items-center pl-3'
                          )}
                        >
                          <CheckIcon className="w-5 h-5" aria-hidden="true" />
                        </span>
                      ) : null}
                      <span className="block text-sm">
                        {`(${env.projectId}) ${env.name}`}
                      </span>
                    </>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </Transition>
        </div>
      </Listbox>
    </div>
  );
});
