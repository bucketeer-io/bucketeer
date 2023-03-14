import 'react-date-range/dist/styles.css'; // main css file
import 'react-date-range/dist/theme/default.css'; // theme css file
import { isLanguageJapanese } from '@/lang/getSelectedLanguage';
import { AuditLogSearchOptions } from '@/types/auditLog';
import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon } from '@heroicons/react/solid';
import en from 'date-fns/locale/en-US';
import ja from 'date-fns/locale/ja';
import React, { FC, Fragment, memo, useEffect, useRef, useState } from 'react';
import { DateRangePicker, defaultStaticRanges } from 'react-date-range';
import { useIntl } from 'react-intl';
import { usePopper } from 'react-popper';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';

export interface DateRangePopoverProps {
  options: AuditLogSearchOptions;
  onChange: (from: number, to: number) => void;
}

export const DateRangePopover: FC<DateRangePopoverProps> = memo(
  ({ options, onChange }) => {
    const { formatMessage: f } = useIntl();
    const referenceElement = useRef<HTMLButtonElement | null>(null);
    const popperElement = useRef<HTMLDivElement | null>(null);
    const popper = usePopper(referenceElement.current, popperElement.current, {
      placement: 'bottom-start',
    });
    const [ranges, setRanges] = useState([
      {
        startDate: options.from ? new Date(options.from * 1000) : new Date(),
        endDate: options.to ? new Date(options.to * 1000) : new Date(),
        key: 'selection',
      },
    ]);

    const handleApply = () => {
      const { startDate, endDate } = ranges[0];

      onChange(
        Math.round(startDate.getTime() / 1000),
        Math.round(endDate.getTime() / 1000)
      );
    };

    let staticRanges = defaultStaticRanges;
    if (isLanguageJapanese) {
      staticRanges = defaultStaticRanges.map((range) => ({
        ...range,
        label: 'Japanese text',
      }));
    }

    return (
      <Popover>
        {({ open }) => (
          <>
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
                    'max-w-lg px-4 mt-3',
                    'transform sm:px-0 lg:max-w-5xl shadow-lg'
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
                          <DateRangePicker
                            onChange={(item: any) =>
                              setRanges([item.selection])
                            }
                            showSelectionPreview={true}
                            moveRangeOnFirstSelection={false}
                            months={2}
                            ranges={ranges}
                            direction="horizontal"
                            rangeColors={['#5d3597']}
                            inputRanges={[]} // hide input ranges
                            locale={isLanguageJapanese ? ja : en}
                            staticRanges={staticRanges}
                          />
                        </div>
                        <div className="flex justify-end mt-4 space-x-2">
                          <button
                            type="button"
                            className="btn-cancel"
                            onClick={() => close()}
                          >
                            Cancel
                          </button>
                          <button
                            type="button"
                            className="btn-submit"
                            onClick={() => {
                              close();
                              handleApply();
                            }}
                          >
                            Apply
                          </button>
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
