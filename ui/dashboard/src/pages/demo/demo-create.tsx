import { useEffect } from 'react';
import { Trans } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { PAGE_PATH_DEMO_SITE } from 'constants/routing';
import { useTranslation } from 'i18n';
import { getDemoTokenStorage } from 'storage/demo-token';
import { IconBackspace } from '@icons';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Icon from 'components/icon';
import DemoForm from './demo-form';

const CreateDemoPage = () => {
  const navigate = useNavigate();
  const { t } = useTranslation(['common', 'auth', 'message']);

  const demoAuthToken = getDemoTokenStorage();
  const { data: demoSiteStatusData, isSuccess } = useQueryDemoSiteStatus();

  const isDemoSiteEnabled = demoSiteStatusData?.isDemoSiteEnabled;

  useEffect(() => {
    if (!demoAuthToken || (isSuccess && !isDemoSiteEnabled)) {
      navigate(PAGE_PATH_DEMO_SITE);
    }
  }, [isDemoSiteEnabled]);

  return (
    <AuthWrapper>
      <div className="-mt-2">
        <Button
          variant="secondary-2"
          onClick={() => navigate(PAGE_PATH_DEMO_SITE)}
          className="p-2 h-auto"
        >
          <Icon icon={IconBackspace} size="sm" />
        </Button>
        <h1 className="text-gray-900 typo-head-bold-huge mt-8">
          {t('auth:demo-organization')}
        </h1>

        <h3 className="text-gray-900 typo-head-light-medium mt-6">
          {t('auth:privacy-notice')}
        </h3>
        <div className="text-gray-600 typo-para-medium mt-2">
          <Trans
            i18nKey="message:demo-privacy-description"
            components={{
              comp: (
                <Link
                  target="_blank"
                  to={`https://app.slack.com/client/T08PSQ7BQ/C043026BME1`}
                  className="text-primary-500 underline"
                />
              )
            }}
          />
        </div>
        <DemoForm isDemoSiteEnabled={isDemoSiteEnabled} />
      </div>
    </AuthWrapper>
  );
};

export default CreateDemoPage;
