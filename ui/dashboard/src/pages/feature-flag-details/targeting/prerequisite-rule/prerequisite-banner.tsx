import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import InfoMessage from 'components/info-message';

const PrerequisiteBanner = ({
  hasPrerequisiteFlags
}: {
  hasPrerequisiteFlags: Feature[];
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <InfoMessage
      title={
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
      description={
        <Trans
          i18nKey={`form:targeting.prerequisite-flags-desc`}
          values={{
            text: t(
              hasPrerequisiteFlags.length > 1 ? 'table:flags' : 'common:flag'
            )?.toLowerCase()
          }}
        />
      }
      linkElements={hasPrerequisiteFlags.map((item, index) => (
        <li
          key={index}
          className="typo-para-small text-primary-500 underline w-fit max-w-full truncate"
        >
          <Link
            to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.id}${PAGE_PATH_FEATURE_TARGETING}`}
          >
            {item.name}
          </Link>
        </li>
      ))}
    />
  );
};

export default PrerequisiteBanner;
