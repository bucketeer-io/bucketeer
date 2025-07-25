import { FC, useEffect, memo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { exchangeDemoToken, ExchangeTokenPayload } from '@api/auth';
import { urls } from 'configs';
import { PAGE_PATH_DEMO_SITE } from 'constants/routing';
import { getCookieState } from 'cookie';
import { ServerErrorType, useSubmit } from 'hooks';
import { useTranslation } from 'i18n';
import queryString from 'query-string';
import { setDemoTokenStorage } from 'storage/demo-token';
import { AppLoading } from 'app';
import { useAuth } from './auth-context';

export const AuthDemoCallbackPage: FC = memo(() => {
  const { t } = useTranslation(['common', 'message']);
  const { setDemoGoogleAuthError } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const query = location.search;

  const { onSubmit: onGoogleDemoHandler } = useSubmit(
    async (payload: ExchangeTokenPayload) => {
      try {
        const response = await exchangeDemoToken(payload);
        if (response.demoCreationToken) {
          setDemoTokenStorage(response.demoCreationToken);
          navigate(`${PAGE_PATH_DEMO_SITE}/new`);
        }
      } catch (error) {
        setDemoGoogleAuthError(
          (error as ServerErrorType)?.response?.data?.message ||
            t('something-went-wrong')
        );
        navigate(PAGE_PATH_DEMO_SITE);
      }
    }
  );

  useEffect(() => {
    const { code, state } = queryString.parse(query);
    const cookieState = getCookieState();

    if (!!code && typeof code === 'string' && state === cookieState) {
      onGoogleDemoHandler({
        code,
        redirectUrl: urls.AUTH_DEMO_REDIRECT,
        type: 2 // Google auth type
      });
    } else {
      throw new Error('Invalid authentication.');
    }
  }, [query]);

  return <AppLoading />;
});
