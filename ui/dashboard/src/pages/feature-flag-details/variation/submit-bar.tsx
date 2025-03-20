import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';

const SubmitBar = () => {
  const { t } = useTranslation(['common', 'table']);
  return (
    <div className="flex items-center justify-between w-full gap-x-6">
      <div className="flex items-center gap-x-2">
        <h3 className="typo-head-bold-small text-gray-800">
          {t('table:feature-flags.variation')}
        </h3>
        <Icon
          icon={IconInfo}
          size={'xxs'}
          color="gray-500"
          className="flex-center"
        />
      </div>
      <Button>{t('save-with-comment')}</Button>
    </div>
  );
};

export default SubmitBar;
