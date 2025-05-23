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
  value: string | number;
  isLastItem?: boolean;
  className?: string;
}) => {
  useTranslation(['form']);
  return (
    <div className={cn('flex items-center gap-x-2', className)}>
      <p className="typo-para-medium text-gray-600">
        <Trans
          i18nKey={titleKey}
          values={{
            value
          }}
          components={{
            b: <span className="text-gray-700" />
          }}
        />
      </p>
      {!isLastItem && <p className="typo-para-medium text-gray-300 slash">|</p>}
    </div>
  );
};
