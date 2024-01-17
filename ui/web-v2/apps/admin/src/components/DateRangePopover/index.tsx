import 'react-date-range/dist/styles.css'; // main css file
import 'react-date-range/dist/theme/default.css'; // theme css file
import './dateRangeStyles.css';
import { isLanguageJapanese } from '@/lang/getSelectedLanguage';
import { AuditLogSearchOptions } from '@/types/auditLog';
import { Popover, Transition } from '@headlessui/react';
import { XIcon, ChevronDownIcon } from '@heroicons/react/solid';
import en from 'date-fns/locale/en-US';
import ja from 'date-fns/locale/ja';
import dayjs from 'dayjs';
import React, { FC, Fragment, memo, useEffect, useRef, useState } from 'react';
import {
  DateRangePicker,
  defaultStaticRanges,
  createStaticRanges,
} from 'react-date-range';
import { FormattedDate, useIntl } from 'react-intl';
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

      const isSameDay = from.isSame(to, 'day');

      if (isSameDay) {
        return (
          <FormattedDate
            value={from.toDate()}
            year="numeric"
            month="short"
            day="numeric"
          />
        );
      }

      return (
        <>
          <FormattedDate
            value={from.toDate()}
            year="numeric"
            month="short"
            day="numeric"
          />
          <span className="mx-1">-</span>
          <FormattedDate
            value={to.toDate()}
            year="numeric"
            month="short"
            day="numeric"
          />
        </>
      );
    };

    const handleClear = (e) => {
      console.log(e);
      e.stopPropagation();
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
                  'group',
                  'rounded-md inline-flex items-center',
                  'h-10',
                  'text-sm',
                  'border border-gray-300'
                )}
              >
                {isDateSelected ? (
                  <div className="flex items-center">
                    <div className="pl-3 flex">
                      <span>{f(messages.show)}:&nbsp;</span>
                      {getSelectedDate()}
                      <button
                        onClick={handleClear}
                        className="px-3 text-gray-500 hover:text-gray-600 mt-[1px]"
                      >
                        <XIcon className="h-4 w-4" aria-hidden="true" />
                      </button>
                    </div>
                    <div className="h-6 bg-gray-300 w-[1px]" />
                    <div className="px-2">
                      <ChevronDownIcon
                        className="w-5 h-5 text-gray-400"
                        aria-hidden="true"
                      />
                    </div>
                  </div>
                ) : (
                  <div className="flex items-center">
                    <div className="px-3">
                      <span>{f(messages.show)}:</span>
                      <span className="text-[#808080] ml-2">
                        {f(messages.mostRecent)}
                      </span>
                    </div>
                    <div className="h-6 bg-gray-300 w-[1px]" />
                    <div className="px-2">
                      <ChevronDownIcon
                        className="w-5 h-5 text-gray-400"
                        aria-hidden="true"
                      />
                    </div>
                  </div>
                )}
              </div>
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
