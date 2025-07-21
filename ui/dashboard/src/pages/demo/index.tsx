import { authenticationUrl } from '@api/auth';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { urls } from 'configs';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { useTranslation } from 'i18n';
import { getDemoTokenStorage } from 'storage/demo-token';
import { cn } from 'utils/style';
import { IconGoogle } from '@icons';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Icon from 'components/icon';
import FormLoading from 'elements/form-loading';
import DemoForm from './demo-form';

const AccessDemoPage = () => {
  const { t } = useTranslation(['common', 'auth', 'message']);

  const authToken = getDemoTokenStorage();
  const { data: demoSiteStatusData, isLoading } = useQueryDemoSiteStatus();

  const isDemoSiteEnabled = demoSiteStatusData?.isDemoSiteEnabled;
  const isAuthenticated = Boolean(authToken);

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    const state = `${Date.now()}`;
    setCookieState(state);

    return authenticationUrl({
      state,
      redirectUrl: urls.AUTH_DEMO_REDIRECT,
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
            {t(isAuthenticated ? 'auth:privacy-notice' : 'auth:demo')}
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
