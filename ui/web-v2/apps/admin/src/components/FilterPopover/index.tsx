import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, useEffect, useRef, useState } from 'react';
import { useIntl } from 'react-intl';
import { usePopper } from 'react-popper';
import ReactSelect from 'react-select';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { FilterTypes } from '../FeatureList';

export type FilterType = 'maintainer' | 'hasExperiment' | 'enabled';

export interface Option {
  value: string;
  label: string;
}

export interface FilterPopoverProps {
  keys: Option[];
  values: Option[];
  onChangeKey: (key: string) => void;
  onAdd: (key: string, value?: string) => void;
  onAddMulti?: (key: string, value?: string[]) => void;
}

export const FilterPopover: FC<FilterPopoverProps> = memo(
  ({ keys, values, onChangeKey, onAdd, onAddMulti }) => {
    const { formatMessage: f } = useIntl();
    const [selectedFilterType, setSelectedFilterType] = useState<Option>(null);
    const referenceElement = useRef<HTMLButtonElement | null>(null);
    const popperElement = useRef<HTMLDivElement | null>(null);
    const popper = usePopper(referenceElement.current, popperElement.current, {
      placement: 'bottom-start',
    });
    const [valueOption, setValue] = useState<Option>(values[0]);
    const [multiValueOption, setMultiValue] = useState<Option[]>([]);

    const isMultiFilter = selectedFilterType?.value === FilterTypes.TAGS;

    const handleKeyChange = (o: Option) => {
      setSelectedFilterType(o);
      onChangeKey(o.value);
    };

    const handleOnClickAdd = () => {
      if (isMultiFilter) {
        onAddMulti(
          selectedFilterType.value,
          multiValueOption.map((o) => o.value)
        );
      } else {
        onAdd(selectedFilterType.value, valueOption.value);
      }
      setValue(values[0]);
    };

    const onPopoverClose = (isPopoverClose: boolean) => {
      if (isPopoverClose) {
        setSelectedFilterType(null);
      }
    };

    useEffect(() => {
      setValue(values[0]);
    }, [values, setValue]);

    return (
      <Popover>
        {({ open }) => (
          <>
            {onPopoverClose(open === false)}
            <Popover.Button
              ref={referenceElement}
              className={classNames(
                'group px-3 py-2',
                'rounded-md inline-flex items-center',
                'hover:bg-gray-100',
                'h-10',
                'text-gray-700',
                'focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75'
              )}
            >
              <span className="text-sm">{f(messages.filter.filter)}</span>
              <SelectorIcon
                className="w-5 h-5 text-gray-400"
                aria-hidden="true"
              />
            </Popover.Button>
            <div
              ref={popperElement}
              style={popper.styles.popper}
              {...popper.attributes.popper}
            >
              <Transition
                as={Fragment}
                enter="transition ease-out duration-200"
                enterFrom="opacity-0 translate-y-1"
                enterTo="opacity-100 translate-y-0"
                leave="transition ease-in duration-150"
                leaveFrom="opacity-100 translate-y-0"
                leaveTo="opacity-0 translate-y-1"
              >
                <Popover.Panel
                  className={classNames(
                    'absolute z-10',
                    'max-w-sm px-4 mt-3',
                    'transform sm:px-0 lg:max-w-3xl shadow-lg'
                  )}
                >
                  {({ close }) => (
                    <div
                      className={classNames(
                        'overflow-hidden rounded-lg',
                        'ring-1 ring-black ring-opacity-5'
                      )}
                    >
                      <div className="p-4 bg-gray-100">
                        <div className="flex">
                          <ReactSelect
                            value={selectedFilterType}
                            className={classNames(
                              'w-60 z-10 text-sm text-gray-700'
                            )}
                            classNamePrefix="react-select"
                            options={keys}
                            menuPortalTarget={document.body}
                            placeholder={f(messages.filter.add)}
                            onChange={handleKeyChange}
                            isSearchable={false}
                            styles={{
                              menuPortal: (base) => ({
                                ...base,
                                zIndex: 9999,
                              }),
                            }}
                          />
                          {selectedFilterType && (
                            <div className="flex">
                              <div className="mx-3 pt-[6px]">
                                {f(messages.feature.clause.operator.equal)}
                              </div>
                              <ReactSelect
                                placeholder={f(
                                  messages.feature.filter.tagsPlaceholder
                                )}
                                closeMenuOnSelect={isMultiFilter ? false : true}
                                className={classNames(
                                  `${
                                    isMultiFilter
                                      ? 'min-w-[270px]'
                                      : 'min-w-max'
                                  } text-sm text-gray-700 focus-visible:ring-white`
                                )}
                                classNamePrefix="react-select"
                                options={values}
                                menuPortalTarget={document.body}
                                styles={{
                                  menuPortal: (base) => ({
                                    ...base,
                                    zIndex: 9999,
                                  }),
                                  menu: ({ width, ...css }) => ({
                                    width: 'max-content',
                                    minWidth: '100%',
                                    ...css,
                                  }),
                                  singleValue: ({
                                    maxWidth,
                                    position,
                                    top,
                                    transform,
                                    ...otherStyles
                                  }) => ({ ...otherStyles }),
                                }}
                                value={
                                  isMultiFilter ? multiValueOption : valueOption
                                }
                                onChange={(o) => {
                                  if (isMultiFilter) {
                                    setMultiValue(o);
                                  } else {
                                    setValue(o);
                                  }
                                }}
                                isSearchable={false}
                                isMulti={isMultiFilter}
                              />
                              <div className={classNames('flex-none ml-4')}>
                                <button
                                  type="button"
                                  className="btn-submit"
                                  disabled={false}
                                  onClick={() => {
                                    handleOnClickAdd();
                                    close();
                                  }}
                                >
                                  {f(messages.button.add)}
                                </button>
                              </div>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  )}
                </Popover.Panel>
              </Transition>
            </div>
          </>
        )}
      </Popover>
    );
  }
);
