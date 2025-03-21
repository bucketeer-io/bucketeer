import { useTranslation } from 'i18n';
import { IconChecked, IconNotifications } from '@icons';
import Icon from 'components/icon';

interface NotificationItemProps {
  title: string;
  date: string;
  description: string;
  highlightText?: string;
}

const NotificationItem = ({
  title,
  date,
  description,
  highlightText
}: NotificationItemProps) => {
  const rowCls = 'flex items-center w-full justify-between gap-x-2';
  return (
    <div className="flex flex-col w-full gap-y-4 pb-4 border-b border-gray-200 last:border-transparent last:pb-0">
      <div className={rowCls}>
        <p className="typo-head-bold-medium leading-[18px]">{title}</p>
        <p className="typo-para-medium leading-4 text-gray-700">{date}</p>
      </div>
      <div className={rowCls}>
        <div className="flex items-center gap-x-1 typo-para-medium leading-4 text-gray-700">
          {description}
          {highlightText && (
            <p className="text-primary-500 underline">{highlightText}</p>
          )}
        </div>
        <Icon
          icon={IconChecked}
          size={'sm'}
          color="gray-500"
          className="flex-center"
        />
      </div>
    </div>
  );
};

const Notifications = () => {
  const { t } = useTranslation(['common', 'table']);
  const isEmpty = false;
  return (
    <div className="flex flex-col w-full max-h-[500px]">
      <div className="flex items-center w-full p-5 border-b border-gray-200">
        <h1 className="typo-head-bold-huge leading-6">{t('notifications')}</h1>
      </div>
      <div className="flex flex-col w-full py-8 px-5 gap-y-4 overflow-auto">
        {isEmpty ? (
          <div className="flex flex-col items-center gap-y-4 px-[59.5px]">
            <Icon icon={IconNotifications} size={'fit'} />
            <p className="typo-para-big text-gray-700 text-center">
              {t('table:empty.notifications')}
            </p>
          </div>
        ) : (
          <>
            <NotificationItem
              title="Experiment Updated"
              date="2d"
              description="Changes were made to the"
              highlightText="“Experiment Name 2”"
            />
            <NotificationItem
              title="Changes Saved"
              date="2d"
              description="Changes were saved successfully!"
            />
            <NotificationItem
              title="Experiment Updated"
              date="2d"
              description="Changes were made to the"
              highlightText="“Experiment Name 2”"
            />
            <NotificationItem
              title="Changes Saved"
              date="2d"
              description="Changes were saved successfully!"
            />
          </>
        )}
      </div>
    </div>
  );
};

export default Notifications;
