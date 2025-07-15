import { useState } from 'react';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { useSubmit } from 'hooks';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconGoogle } from '@icons';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Icon from 'components/icon';
import FormLoading from 'elements/form-loading';
import CreateDemoOrganizationForm from './create-demo-org-form';

const AccessDemoPage = () => {
  const { t } = useTranslation(['auth', 'common', 'form', 'message']);

  const { data: demoSiteStatusData, isLoading } = useQueryDemoSiteStatus();
  const isDemoSiteEnabled = demoSiteStatusData?.isDemoSiteEnabled;

  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    // call the API to authenticate with Google

    return new Promise(rs => {
      setTimeout(() => {
        rs(setIsAuthenticated(true));
      }, 1000);
    });
  });

  return (
    <AuthWrapper>
      {isLoading ? (
        <FormLoading />
      ) : (
        <>
          <h1 className="text-gray-900 typo-head-bold-huge">{t('demo')}</h1>
          <p
            className={cn('text-gray-600 typo-para-medium mt-6', {
              'text-accent-red-500': !isDemoSiteEnabled
            })}
          >
            {t(
              `message:${isDemoSiteEnabled ? 'demo-available' : 'demo-not-available'}`
            )}
          </p>
          {!isDemoSiteEnabled && (
            <div>
              {!isAuthenticated ? (
                <Button
                  loading={submitting}
                  onClick={onGoogleLoginHandler}
                  variant={'secondary-2'}
                  className="w-full mt-8"
                >
                  <Icon icon={IconGoogle} />
                  {`Sign in With Google`}
                </Button>
              ) : (
                <CreateDemoOrganizationForm
                  isDemoSiteEnabled={isDemoSiteEnabled}
                />
              )}
            </div>
          )}
        </>
      )}
    </AuthWrapper>
  );
};

export default AccessDemoPage;
