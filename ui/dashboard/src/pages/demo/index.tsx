import { useState } from 'react';
import { authenticationUrl } from '@api/auth';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { urls } from 'configs';
import { useSubmit } from 'hooks';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconGoogle } from '@icons';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Icon from 'components/icon';
import FormLoading from 'elements/form-loading';
import DemoForm from './demo-form';

const AccessDemoPage = () => {
  const { t } = useTranslation(['auth', 'common', 'form', 'message']);

  const { data: demoSiteStatusData, isLoading } = useQueryDemoSiteStatus();
  const isDemoSiteEnabled = !demoSiteStatusData?.isDemoSiteEnabled;

  const [isAuthenticated] = useState(false);

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    const state = `${Date.now()}`;

    return authenticationUrl({
      state,
      redirectUrl: urls.AUTH_REDIRECT,
      type: 2 // Google auth type
    }).then(response => {
      if (response.url) {
        window.location.href = response.url;
      }
    });
  });

  return (
    <AuthWrapper>
      {isLoading ? (
        <FormLoading />
      ) : (
        <>
          <h1 className="text-gray-900 typo-head-bold-huge">
            {t(isAuthenticated ? 'privacy-notice' : 'demo')}
          </h1>
          <div
            className={cn('text-gray-600 typo-para-medium mt-6', {
              'text-accent-red-500': !isDemoSiteEnabled
            })}
            dangerouslySetInnerHTML={{
              __html: t(
                isAuthenticated
                  ? 'message:demo-privacy-description'
                  : isDemoSiteEnabled
                    ? 'message:demo-available'
                    : 'message:demo-not-available'
              )
            }}
          />
          {isDemoSiteEnabled && (
            <>
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
                <DemoForm isDemoSiteEnabled={isDemoSiteEnabled} />
              )}
            </>
          )}
        </>
      )}
    </AuthWrapper>
  );
};

export default AccessDemoPage;
