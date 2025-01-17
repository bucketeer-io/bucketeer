import { ExperimentStatus } from '@types';
import { cn } from 'utils/style';

interface Props {
  text: string;
  status?: ExperimentStatus;
  isInUseStatus?: boolean;
  className?: string;
}

const Status = ({ text, status, isInUseStatus = false, className }: Props) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] rounded-[3px] text-gray-600 bg-gray-100 capitalize',
        {
          'bg-accent-orange-50 text-accent-orange-500': status === 'WAITING',
          'bg-accent-green-50 text-accent-green-500':
            status === 'RUNNING' || isInUseStatus,
          'bg-accent-red-50 text-accent-red-500': [
            'STOPPED',
            'FORCE_STOPPED'
          ].includes(status || '')
        },
        className
      )}
    >
      {text}
    </div>
  );
};

export default Status;
