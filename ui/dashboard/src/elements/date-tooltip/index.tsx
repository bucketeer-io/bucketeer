import { ReactNode } from 'react';
import { formatLongDateTime } from 'utils/date-time';
import { Tooltip, TooltipProps } from 'components/tooltip';

interface Props extends TooltipProps {
  date: string | null;
  trigger: ReactNode;
}

const DateTooltip = ({ date, trigger, ...props }: Props) => {
  const dateFormatted = date ? formatLongDateTime({ value: date }) : null;

  return (
    <Tooltip
      trigger={trigger}
      content={dateFormatted}
      className="bg-gray-800"
      {...props}
    />
  );
};

export default DateTooltip;
