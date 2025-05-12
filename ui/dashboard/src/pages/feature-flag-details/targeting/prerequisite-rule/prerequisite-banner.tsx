import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import ToastMessage from 'components/toast';

const PrerequisiteBanner = ({
  hasPrerequisiteFlags
}: {
  hasPrerequisiteFlags: Feature[];
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <ToastMessage
      message={
        <Trans
          i18nKey={'form:targeting.prerequisite-flags'}
          values={{
            quantity: hasPrerequisiteFlags.length,
            flag: t(
              hasPrerequisiteFlags.length > 1 ? 'table:flags' : 'common:flag'
            )?.toLowerCase()
          }}
        />
      }
      messageType="info"
      toastType="prerequisite-message"
      className="!w-full !max-w-full"
      toastChildren={
        <div className="flex flex-col w-full gap-y-3">
          <p className="typo-para-medium text-gray-600">
            <Trans
              i18nKey={`form:targeting.prerequisite-flags-desc`}
              values={{
                text: t(
                  hasPrerequisiteFlags.length > 1
                    ? 'table:flags'
                    : 'common:flag'
                )?.toLowerCase()
              }}
            />
          </p>
          <ul className="flex flex-col w-full gap-y-2">
            {hasPrerequisiteFlags.map((item, index) => (
              <li
                className="flex items-center gap-x-2 typo-para-medium"
                key={index}
              >
                <p className="text-gray-700">{index + 1}.</p>
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.id}${PAGE_PATH_FEATURE_TARGETING}`}
                  className="underline text-primary-500"
                >
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      }
    />
  );
};

export default PrerequisiteBanner;
