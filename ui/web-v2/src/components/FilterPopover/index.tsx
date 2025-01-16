import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, useEffect, useRef, useState } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { FilterTypes } from '../FeatureList';
import { Select } from '../Select';

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
    const [valueOption, setValueOption] = useState<Option>(null);
    const [multiValueOption, setMultiValueOption] = useState<Option[]>([]);

    const isMultiFilter = selectedFilterType?.value === FilterTypes.TAGS;
    const isFilterTypeMaintainer =
      selectedFilterType?.value === FilterTypes.MAINTAINER;

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
        onAdd(selectedFilterType.value, valueOption?.value);
      }
    };

    useEffect(() => {
      if (isFilterTypeMaintainer) {
        setValueOption(null);
      } else {
        setValueOption(values[0]);
      }
    }, [values, setValueOption, isFilterTypeMaintainer]);

    return (
      <Popover>
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
          <SelectorIcon className="w-5 h-5 text-gray-400" aria-hidden="true" />
        </Popover.Button>
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
                  'rounded-lg',
                  'ring-1 ring-black ring-opacity-5'
                )}
              >
                <div className="p-4 bg-gray-100">
                  <div className="flex">
                    <Select
                      value={selectedFilterType}
                      className={classNames('w-60')}
                      options={keys}
                      placeholder={f(messages.filter.add)}
                      onChange={handleKeyChange}
                      isSearchable={false}
                    />
                    {selectedFilterType && (
                      <div className="flex">
                        <div className="mx-3 pt-[6px]">
                          {f(messages.feature.clause.operator.equal)}
                        </div>
                        <Select
                          placeholder={f(
                            selectedFilterType?.value === FilterTypes.TAGS
                              ? messages.tags.tagsPlaceholder
                              : messages.select
                          )}
                          closeMenuOnSelect={isMultiFilter ? false : true}
                          className={classNames(
                            isMultiFilter
                              ? 'min-w-[270px]'
                              : isFilterTypeMaintainer
                                ? 'min-w-[220px]'
                                : 'min-w-[106px]'
                          )}
                          options={values}
                          styles={{
                            menu: ({ ...css }) => ({
                              width: 'max-content',
                              minWidth: '100%',
                              ...css
                            }),
                            singleValue: ({ ...otherStyles }) => ({
                              ...otherStyles
                            })
                          }}
                          value={isMultiFilter ? multiValueOption : valueOption}
                          onChange={(o) => {
                            if (isMultiFilter) {
                              setMultiValueOption(o);
                            } else {
                              setValueOption(o);
                            }
                          }}
                          isSearchable={
                            selectedFilterType?.value === FilterTypes.TAGS ||
                            isFilterTypeMaintainer
                          }
                          isMulti={isMultiFilter}
                          clearable={isFilterTypeMaintainer}
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
      </Popover>
    );
  }
);
