import { ReactNode } from 'react';
import { formatLongDateTime } from 'utils/date-time';
import { Tooltip } from 'components/tooltip';

type Props = {
  date: Date | string;
  trigger: ReactNode;
};

const DateTooltip = ({ date, trigger }: Props) => {
  const _date = date instanceof Date ? date : new Date(Number(date) * 1000);
  const dateFormatted = formatLongDateTime(_date);

  return (
    <Tooltip
      trigger={trigger}
      content={dateFormatted}
      className="bg-gray-800"
    />
  );
};

export default DateTooltip;
