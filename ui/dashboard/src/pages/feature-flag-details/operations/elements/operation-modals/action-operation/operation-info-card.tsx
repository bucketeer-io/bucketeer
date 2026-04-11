import { InfoIcon } from 'lucide-react';
import Icon from 'components/icon';

type OperationInfoCardProps = {
  title: string;
  description: string;
  className?: string;
};

const OperationInfoCard = ({
  title,
  description,
  className
}: OperationInfoCardProps) => {
  return (
    <div
      className={`w-full rounded-lg border-l-[8px] border-primary-500 px-4 py-3 shadow-card ${className ?? ''}`}
    >
      <div className="flex items-start gap-4 typo-para-medium">
        <Icon
          icon={InfoIcon}
          size="xxs"
          className="mt-[5px] text-primary-500"
        />
        <div>
          <p className="font-bold text-primary-500">{title}</p>
          <p className="typo-para-medium text-gray-500 w-full mt-2">
            {description}
          </p>
        </div>
      </div>
    </div>
  );
};

export default OperationInfoCard;
