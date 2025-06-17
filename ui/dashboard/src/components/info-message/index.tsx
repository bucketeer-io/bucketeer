import { ReactNode, useState } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconChevronDown, IconInfoFilled } from '@icons';
import Icon from 'components/icon';

const InfoMessage = ({
  title,
  description,
  linkElements,
  className
}: {
  title: ReactNode;
  description?: ReactNode;
  linkElements?: ReactNode;
  className?: string;
}) => {
  const { t } = useTranslation(['common']);
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div
      className={cn(
        'flex flex-col w-full p-4 gap-y-3 rounded border-l-4 border-accent-blue-500 bg-accent-blue-50',
        className
      )}
    >
      <div className="flex items-center w-full gap-x-2">
        <Icon icon={IconInfoFilled} size={'xxs'} color="accent-blue-500" />
        <p className="typo-para-small leading-[14px] text-accent-blue-500">
          {title}
        </p>
      </div>
      {linkElements && (
        <div
          className="flex w-fit max-w-full items-center gap-x-2 pl-6 cursor-pointer"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          <p className="typo-para-small text-gray-700">
            {t(isExpanded ? 'close' : 'see-more')}
          </p>
          <Icon
            icon={IconChevronDown}
            size="xxs"
            color="gray-600"
            className={cn('flex-center rotate-0 transition-all duration-200', {
              'rotate-180': isExpanded
            })}
          />
        </div>
      )}
      {(linkElements ? isExpanded && description : description) && (
        <p className="typo-para-small text-gray-700 pl-10">{description}</p>
      )}
      {isExpanded && (
        <ul className="flex flex-col gap-y-2 list-decimal list-inside pl-10 w-full max-w-full">
          {linkElements}
        </ul>
      )}
    </div>
  );
};

export default InfoMessage;
