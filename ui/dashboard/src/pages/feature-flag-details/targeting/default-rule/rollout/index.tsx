import { useTranslation } from 'i18n';

const DefaultRuleRollout = () => {
  const { t } = useTranslation(['form', 'table', 'common']);
  return (
    <div className="flex flex-col w-full gap-y-6">
      <p className="typo-para-medium text-gray-700">
        {t('common:source-type.progressive-rollout')}
      </p>
      <div className="flex items-center w-full gap-x-2">
        <p className="typo-para-medium text-gray-600">
          {t('table:results.variation')}
        </p>
      </div>
    </div>
  );
};

export default DefaultRuleRollout;
