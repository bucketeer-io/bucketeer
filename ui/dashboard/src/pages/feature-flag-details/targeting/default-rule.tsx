import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import Card from '../elements/card';
import ServeDropdown from './serve-dropdown';

const DefaultRule = () => {
  const { t } = useTranslation(['table']);
  return (
    <Card>
      <div>
        <div className="flex items-center gap-x-2">
          <p className="typo-para-medium leading-4 text-gray-700">
            {t('feature-flags.default-rule')}
            <span className="text-accent-red-400">*</span>
          </p>
          <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
        </div>
        <p className="typo-para-small text-gray-600 mt-4">
          {t('feature-flags.default-rule-desc')}
        </p>
      </div>
      <ServeDropdown
        label={t('feature-flags.variation')}
        isExpand
        serveValue={1}
        onChangeServe={() => {}}
      />
    </Card>
  );
};

export default DefaultRule;
