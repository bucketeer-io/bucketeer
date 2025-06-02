import { useTranslation } from 'react-i18next';
import { IconToastWarning } from '@icons';
import Icon from 'components/icon';

const ArchiveWarning = () => {
  const { t } = useTranslation(['table']);
  return (
    <div className="flex items-center w-full p-4 gap-x-2 rounded-xl bg-accent-yellow-50 text-accent-yellow-500 typo-para-medium">
      <Icon
        icon={IconToastWarning}
        size={'xs'}
        color="accent-yellow-500"
        className="flex-center"
      />
      <div>{t('feature-flags.archive-warning')}</div>
    </div>
  );
};

export default ArchiveWarning;
