import { ReactNode } from 'react';
import { formatLongDateTime } from 'utils/date-time';
import { Tooltip } from 'components/tooltip';

type Props = {
  date: string | null;
  trigger: ReactNode;
};

const DateTooltip = ({ date, trigger }: Props) => {
  const dateFormatted = date ? formatLongDateTime({ value: date }) : null;

  return (
    <Tooltip
      trigger={trigger}
      content={dateFormatted}
      className="bg-gray-800"
    />
  );
};

export default DateTooltip;
