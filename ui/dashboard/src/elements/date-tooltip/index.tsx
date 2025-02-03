import { ReactNode } from 'react';
import { formatLongDateTime } from 'utils/date-time';
import { Tooltip } from 'components/tooltip';

type Props = {
  date: string;
  trigger: ReactNode;
};

const DateTooltip = ({ date, trigger }: Props) => {
  const dateFormatted = formatLongDateTime({ value: date });

  return (
    <Tooltip
      trigger={trigger}
      content={dateFormatted}
      className="bg-gray-800"
    />
  );
};

export default DateTooltip;
