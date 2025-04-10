import { ReactDatePickerCustomHeaderProps } from 'react-datepicker';
import { getYear } from 'date-fns';
import { range } from 'lodash';
import { cn } from 'utils/style';
import { IconChevronRight } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

const years = range(2020, getYear(new Date()) + 10, 1);
const months = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December'
];

const CustomizeHeader = ({
  date,
  changeYear,
  changeMonth,
  decreaseMonth,
  increaseMonth,
  prevMonthButtonDisabled,
  nextMonthButtonDisabled,
  customHeaderCount
}: ReactDatePickerCustomHeaderProps) => {
  console.log({ customHeaderCount });
  return (
    <div
      className={cn('flex flex-col gap-y-5 w-full divide-y divide-gray-200')}
    >
      <div
        className={cn(
          'flex items-center justify-between w-full relative z-[1000] pb-5',
          
        )}
      >
        <Button
          type="button"
          size={'icon-sm'}
          variant={'secondary-2'}
          onClick={decreaseMonth}
          disabled={prevMonthButtonDisabled}
          className="size-8"
        >
          <Icon
            icon={IconChevronRight}
            className="rotate-180"
            size={'xs'}
            color="gray-500"
          />
        </Button>
        <div className="flex items-center gap-x-2">
          <select
            value={months[date.getMonth()]}
            onChange={({ target: { value } }) =>
              changeMonth(months.indexOf(value))
            }
            className="border-none font-sofia-pro typo-para-medium text-gray-700"
          >
            {months.map(option => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
          <select
            value={getYear(date)}
            onChange={({ target: { value } }) => changeYear(+value)}
            className="border-none font-sofia-pro typo-para-medium text-gray-700"
          >
            {years.map(option => (
              <option key={option} value={option}>
                {option}
              </option>
            ))}
          </select>
        </div>
        <Button
          type="button"
          size={'icon-sm'}
          variant={'secondary-2'}
          onClick={increaseMonth}
          disabled={nextMonthButtonDisabled}
          className="size-8"
        >
          <Icon icon={IconChevronRight} size={'xs'} color="gray-500" />
        </Button>
      </div>
    </div>
  );
};

export default CustomizeHeader;
