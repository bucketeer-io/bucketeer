import { ExperimentStatus, OperationStatus } from '@types';
import { cn } from 'utils/style';

interface Props {
  text: string;
  status?: ExperimentStatus | OperationStatus;
  isInUseStatus?: boolean;
  className?: string;
}

const Status = ({ text, status, isInUseStatus = false, className }: Props) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] rounded-[3px] text-gray-600 bg-gray-100 dark:text-dark-gray-200 dark:bg-dark-black-700 capitalize',
        {
          'bg-accent-orange-50 dark:bg-accent-orange-900/30 text-accent-orange-500 dark:text-accent-orange-400':
            status === 'WAITING',
          'bg-accent-green-50 dark:bg-accent-green-900/30 text-accent-green-500 dark:text-accent-green-400':
            status === 'RUNNING' || isInUseStatus,
          'bg-gray-50 dark:bg-dark-black-700 text-gray-500 dark:text-dark-gray-200':
            status === 'FINISHED',
          'bg-accent-red-50 dark:bg-accent-red-900/30 text-accent-red-500 dark:text-accent-red-400':
            ['STOPPED', 'FORCE_STOPPED'].includes(status || '')
        },
        className
      )}
    >
      {text}
    </div>
  );
};

export default Status;
