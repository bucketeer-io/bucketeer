import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Feature, FeaturePrerequisite } from '@types';
import ToastMessage from 'components/toast';

const PrerequisiteBanner = ({
  prerequisite,
  features
}: {
  prerequisite: FeaturePrerequisite[];
  features: Feature[];
}) => {
  const { t } = useTranslation(['form']);
  return (
    <ToastMessage
      message={
        <Trans
          i18nKey={'form:.feature-flags.prerequisite-flags'}
          values={{
            quantity: prerequisite.length,
            flag: prerequisite.length > 1 ? 'flags' : 'flag'
          }}
        />
      }
      messageType="info"
      toastType="prerequisite-message"
      className="!w-full !max-w-full"
      toastChildren={
        <div className="flex flex-col w-full gap-y-3">
          <p className="typo-para-medium text-gray-600">
            {t('feature-flags.prerequisite-flags-desc')}
          </p>
          <ul className="flex flex-col w-full gap-y-5">
            {prerequisite.map((item, index) => (
              <li
                className="flex items-center gap-x-2 typo-para-medium"
                key={index}
              >
                <p className="text-gray-700">{index + 1}.</p>
                <p className="text-primary-500">
                  {
                    features.find(feature => feature.id === item.featureId)
                      ?.name
                  }
                </p>
              </li>
            ))}
          </ul>
        </div>
      }
    />
  );
};

export default PrerequisiteBanner;
