import { cn } from 'utils/style';

interface Props {
  status: string;
  className?: string;
}

const Status = ({ status, className }: Props) => {
  return (
    <div
      className={cn(
        'flex-center w-fit px-2 py-1.5 typo-para-small leading-[14px] rounded-[3px] text-gray-600 bg-gray-100 capitalize',
        className
      )}
    >
      {status}
    </div>
  );
};

export default Status;
