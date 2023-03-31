import 'react-date-range/dist/styles.css'; // main css file
import 'react-date-range/dist/theme/default.css'; // theme css file
import './dateRangeStyles.css';
import { isLanguageJapanese } from '@/lang/getSelectedLanguage';
import { AuditLogSearchOptions } from '@/types/auditLog';
import { Popover, Transition } from '@headlessui/react';
import { SelectorIcon, XIcon } from '@heroicons/react/solid';
import en from 'date-fns/locale/en-US';
import ja from 'date-fns/locale/ja';
import dayjs from 'dayjs';
import React, { FC, Fragment, memo, useEffect, useRef, useState } from 'react';
import {
  DateRangePicker,
  defaultStaticRanges,
  createStaticRanges,
} from 'react-date-range';
import { useIntl } from 'react-intl';
import { usePopper } from 'react-popper';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';

const defaultEnRangesLabels = [
  'Today',
  'Yesterday',
  'This week',
  'Last week',
  'This month',
  'Last month',
];

const defaultJaRangesLabels = ['今日', '昨日', '今週', '先週', '今月', '先月'];

const extraEnRangesLabels = [
  { label: 'Last 3 months', month: 3 },
  { label: 'Last 6 months', month: 6 },
  { label: 'Last 12 months', month: 12 },
];

const extraJaRangesLabels = [
  { label: '過去3ヶ月', month: 3 },
  { label: '過去6ヶ月', month: 6 },
  { label: '過去12ヶ月', month: 12 },
];

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
        startDate: new Date(),
        endDate: new Date(),
        key: 'selection',
      },
    ]);
    const [isDateSelected, setIsDateSelected] = useState(false);

    useEffect(() => {
      if (options.from && options.to) {
        setIsDateSelected(true);
        setRanges([
          {
            ...ranges[0],
            startDate: new Date(options.from * 1000),
            endDate: new Date(options.to * 1000),
          },
        ]);
      } else if (isDateSelected === true) {
        setIsDateSelected(false);
      }
    }, [options]);

    const getSelectedDate = () => {
      const from = dayjs(new Date(options.from * 1000));
      const to = dayjs(new Date(options.to * 1000));

      const isSameYear = from.isSame(to, 'year');
      const isSameMonth = from.isSame(to, 'month');
      const isSameDay = from.isSame(to, 'day');

      if (!isSameYear) {
        return `${from.format('MMM D, YYYY')} - ${to.format('MMM D, YYYY')}`;
      }

      if (!isSameMonth) {
        return `${from.format('MMM D')} - ${to.format('MMM D, YYYY')}`;
      }

      if (!isSameDay) {
        return `${from.format('MMM D')} - ${to.format('D, YYYY')}`;
      }
      return from.format('MMM D, YYYY');
    };

    const handleClear = () => {
      setRanges([
        {
          startDate: new Date(),
          endDate: new Date(),
          key: 'selection',
        },
      ]);
      onChange(null, null);
    };

    const handleApply = () => {
      const { startDate, endDate } = ranges[0];

      onChange(
        Math.trunc(startDate.getTime() / 1000),
        Math.trunc(endDate.getTime() / 1000)
      );
    };

    let staticRanges = [];

    const labels = isLanguageJapanese
      ? defaultJaRangesLabels
      : defaultEnRangesLabels;

    const extraLabels = isLanguageJapanese
      ? extraJaRangesLabels
      : extraEnRangesLabels;

    staticRanges = defaultStaticRanges.map((range, rangeIdx) => ({
      ...range,
      label: labels[rangeIdx],
    }));

    staticRanges = [
      ...staticRanges,
      ...createStaticRanges(
        extraLabels.map(({ label, month }) => ({
          label,
          range: () => ({
            startDate: dayjs().subtract(month, 'month').toDate(),
            endDate: new Date(),
          }),
        }))
      ),
    ];

    return (
      <Popover>
        {({ open }) => (
          <>
            <Popover.Button ref={referenceElement}>
              <div
                className={classNames(
                  'group pl-3 pr-2 py-2',
                  'rounded-md inline-flex items-center',
                  'hover:bg-gray-100',
                  'h-10',
                  'text-sm text-gray-700',
                  'focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75',
                  `${isDateSelected && 'border'}`
                )}
              >
                {isDateSelected ? (
                  <div className="flex">
                    <span>
                      {f(messages.auditLog.filter.dates)}: {getSelectedDate()}
                    </span>
                    <SelectorIcon
                      className="w-5 h-5 text-gray-400 ml-2"
                      aria-hidden="true"
                    />
                  </div>
                ) : (
                  <>
                    <span className="text-sm">{f(messages.filter.filter)}</span>
                    <SelectorIcon
                      className="w-5 h-5 text-gray-400"
                      aria-hidden="true"
                    />
                  </>
                )}
              </div>
            </Popover.Button>
            {isDateSelected && (
              <button
                type="button"
                className="inline-flex items-center ml-2 rounded-md bg-white py-2.5 px-3.5 text-sm text-gray-900 shadow-sm ring-1 ring-inset ring-gray-200 hover:bg-gray-50"
                onClick={handleClear}
              >
                Clear
                <XIcon className="ml-1 h-4 w-4" aria-hidden="true" />
              </button>
            )}
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
                      <div className="bg-gray-100">
                        <div className="flex">
                          <DateRangePicker
                            onChange={(item: any) => {
                              setRanges([item.selection]);
                            }}
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
                        <div className="flex justify-end py-2 pr-2 space-x-2 border-t border-gray-300">
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
