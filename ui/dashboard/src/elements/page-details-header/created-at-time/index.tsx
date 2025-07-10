import { IconAccessTimeOutlined } from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Icon from 'components/icon';
import DateTooltip from 'elements/date-tooltip';

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
    <DateTooltip
      align="end"
      alignOffset={-40}
      trigger={
        <div
          className={cn(
            'flex items-center w-fit h-6 gap-1.5 text-gray-500 whitespace-nowrap -mb-1',
            className
          )}
        >
          <Icon icon={IconAccessTimeOutlined} size={'xxs'} />
          {Number(createdAt) === 0 ? (
            t('never')
          ) : (
            <p className="typo-para-small">
              {t('created-at-time', {
                time: formatDateTime(createdAt)
              })}
            </p>
          )}
        </div>
      }
      date={Number(createdAt) === 0 ? null : createdAt}
    />
  );
};

export default CreatedAtTime;
