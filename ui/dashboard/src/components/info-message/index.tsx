import { ReactNode, useEffect, useState } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconChevronDown, IconInfoFilled, IconToastWarning } from '@icons';
import Icon from 'components/icon';

const ICON_MAP = {
  info: <Icon icon={IconInfoFilled} size="xxs" color="accent-blue-500" />,
  warning: <Icon icon={IconToastWarning} size="xxs" color="accent-yellow-500" />
};

const TEXT_COLOR_MAP = {
  info: 'text-accent-blue-500',
  warning: 'text-accent-yellow-500'
};

const CONTAINER_COLOR_MAP = {
  info: 'border-accent-blue-500 bg-accent-blue-50',
  warning: 'border-accent-yellow-500 bg-accent-yellow-50'
};

type InfoMessageProps = {
  title: ReactNode;
  description?: ReactNode;
  linkElements?: ReactNode;
  className?: string;
  typeOfIcon?: 'info' | 'warning';
  isToggleable?: boolean;
};

const InfoMessage = ({
  title,
  description,
  linkElements,
  className,
  typeOfIcon = 'info',
  isToggleable = true
}: InfoMessageProps) => {
  const { t } = useTranslation(['common']);

  const [isExpanded, setIsExpanded] = useState(!isToggleable);

  useEffect(() => {
    if (!isToggleable) {
      setIsExpanded(true);
    }
  }, [isToggleable]);

  const hasLinks = !!linkElements;
  const shouldShowDescription = !!description && (!hasLinks || isExpanded);

  return (
    <div
      className={cn(
        'flex flex-col w-full gap-y-3 rounded border-l-4 p-4',
        CONTAINER_COLOR_MAP[typeOfIcon],
        className
      )}
    >
      <div className="flex items-center w-full gap-x-2">
        {ICON_MAP[typeOfIcon]}
        <p
          className={cn(
            'typo-para-small leading-[14px]',
            TEXT_COLOR_MAP[typeOfIcon]
          )}
        >
          {title}
        </p>
      </div>

      {hasLinks && isToggleable && (
        <button
          type="button"
          onClick={() => setIsExpanded(prev => !prev)}
          className="flex w-fit max-w-full items-center gap-x-2 pl-6"
        >
          <p className="typo-para-small text-gray-700">
            {t(isExpanded ? 'close' : 'see-more')}
          </p>
          <Icon
            icon={IconChevronDown}
            size="xxs"
            color="gray-600"
            className={cn(
              'transition-transform duration-200',
              isExpanded && 'rotate-180'
            )}
          />
        </button>
      )}

      {shouldShowDescription && (
        <p
          className={cn(
            'typo-para-small pl-10',
            typeOfIcon === 'warning'
              ? 'text-accent-yellow-500'
              : 'text-gray-700',
            {
              'pl-[30px]': !isToggleable
            }
          )}
        >
          {description}
        </p>
      )}

      {isExpanded && hasLinks && (
        <ul
          className={cn(
            'flex flex-col gap-y-2 list-decimal list-inside pl-10 w-full max-w-full',
            !isToggleable && 'pl-[30px]'
          )}
        >
          {linkElements}
        </ul>
      )}
    </div>
  );
};

export default InfoMessage;
