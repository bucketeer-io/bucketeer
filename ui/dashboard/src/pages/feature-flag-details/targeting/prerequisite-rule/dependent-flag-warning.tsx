import { Fragment } from 'react';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import InfoMessage from 'components/info-message';

export const DependentFlagWarning = ({
  dependentFlags
}: {
  dependentFlags: Feature[];
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  return (
    <InfoMessage
      typeOfIcon="warning"
      isToggleable={false}
      title={t('form:targeting.prerequisite-flags-used-by')}
      description={
        <Fragment>
          <span className="block">
            {t('form:targeting.prerequisite-flags-used-by-desc')}
          </span>
          <span className="block">
            {t('form:targeting.prerequisite-flags-used-by-desc-two')}
          </span>
        </Fragment>
      }
      linkElements={dependentFlags.map((item, index) => (
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
