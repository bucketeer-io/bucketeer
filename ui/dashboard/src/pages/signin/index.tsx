// import { useNavigate } from 'react-router-dom';
import { authenticationUrl } from '@api/auth';
import { useNavigate } from '@tanstack/react-router';
import { useAuth } from 'auth';
import { urls } from 'configs';
import { DEMO_SIGN_IN_ENABLED } from 'configs';
import { PAGE_PATH_AUTH_SIGNIN } from 'constants/routing';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { useTranslation } from 'i18n';
import { IconEmail, IconGoogle } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import AuthWrapper from './elements/auth-wrapper';

const SignIn = () => {
  const { t } = useTranslation(['auth']);
  const navigate = useNavigate();
  const { isGoogleAuthError, setIsGoogleAuthError } = useAuth();

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    const state = `${Date.now()}`;
    setCookieState(state);
    setIsGoogleAuthError(false);

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
      <div className="grid gap-6">
        <h1 className="text-gray-900 typo-head-bold-huge">
          {`Sign in to Bucketeer`}
        </h1>
        {isGoogleAuthError && (
          <p className="text-accent-red-500 typo-para-medium">
            {t(`error-message.invalid-google-auth`)}
          </p>
        )}
      </div>
      <div className="flex flex-col gap-4 mt-10">
        {DEMO_SIGN_IN_ENABLED && (
          <Button
            variant="secondary-2"
            onClick={() => {
              navigate({
                to: PAGE_PATH_AUTH_SIGNIN
              });
              setIsGoogleAuthError(false);
            }}
          >
            <Icon icon={IconEmail} />
            {`Sign in With Email`}
          </Button>
        )}
        <Button
          loading={submitting}
          onClick={onGoogleLoginHandler}
          variant={'secondary-2'}
        >
          <Icon icon={IconGoogle} />
          {`Sign in With Google`}
        </Button>
        {/* <Button variant={'secondary-2'}>
          <Icon icon={IconGithub} />
          {`Sign in With Github`}
        </Button>
        <Button variant={'secondary-2'}>
          <Icon icon={IconKey} />
          {`Sign in With SSO`}
        </Button> */}
      </div>
    </AuthWrapper>
  );
};

export default SignIn;
