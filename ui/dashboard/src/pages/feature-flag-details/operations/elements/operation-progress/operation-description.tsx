import { ReactElement } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';

export const OperationDescription = ({
  titleKey,
  value,
  isLastItem,
  className
}: {
  titleKey: string;
  value: string | number | ReactElement;
  isLastItem?: boolean;
  className?: string;
}) => {
  useTranslation(['form']);
  return (
    <div className={cn('flex items-center gap-x-2', className)}>
      <div className="flex items-center gap-1 typo-para-medium text-gray-600 dark:text-dark-gray-200">
        <Trans
          i18nKey={titleKey}
          values={{
            value
          }}
          components={{
            comp: <div className="inline-flex">{value}</div>,
            b: <span className="text-gray-700 dark:text-dark-gray-400" />
          }}
        />
      </div>
      {!isLastItem && (
        <p className="typo-para-medium text-gray-300 dark:text-dark-black-700 slash">
          |
        </p>
      )}
    </div>
  );
};
