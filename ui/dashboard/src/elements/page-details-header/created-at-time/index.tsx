import { IconAccessTimeOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Icon from 'components/icon';

const CreatedAtTime = ({
  createdAt,
  className
}: {
  createdAt: string;
  className?: string;
}) => {
  const { t } = useTranslation(['table']);
  const formatDateTime = useFormatDateTime();

  return (
    <div
      className={cn('flex items-center h-6 text-gray-500 gap-1.5', className)}
    >
      <Icon icon={IconAccessTimeOutlined} size="xxs" />
      <p className="typo-para-small">
        {t('created-at-time', {
          time: formatDateTime(createdAt)
        })}
      </p>
    </div>
  );
};

export default CreatedAtTime;
