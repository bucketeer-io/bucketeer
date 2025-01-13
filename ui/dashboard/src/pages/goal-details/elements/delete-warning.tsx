import { useTranslation } from 'i18n';
import { IconInfoFilled } from '@icons';
import Icon from 'components/icon';

const DeleteWarning = () => {
  const { t } = useTranslation(['form']);

  return (
    <div className="flex items-center w-full p-4 gap-x-2 rounded border-l-4 border-accent-blue-500 bg-accent-blue-50">
      <Icon icon={IconInfoFilled} size={'xxs'} />
      <p className="typo-para-small leading-[14px] text-accent-blue-500">
        {t('goal-details.delete-warning-desc')}
      </p>
    </div>
  );
};

export default DeleteWarning;
