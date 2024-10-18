import { IconErrorOutlineRound } from 'react-icons-material-design';
import { cn } from 'utils/style';
import Icon from 'components/icon';

export const InvalidMessage = ({
  className,
  children
}: {
  className?: string;
  children: string;
}) => {
  return (
    <div className={cn('h-full flex-grow flex-center', className)}>
      <div className="flex flex-col justify-center items-center">
        <Icon size="xl" color="accent-red-500" icon={IconErrorOutlineRound} />
        <div className="typo-head-semi-medium mt-2 text-center">{children}</div>
      </div>
    </div>
  );
};

export default InvalidMessage;
