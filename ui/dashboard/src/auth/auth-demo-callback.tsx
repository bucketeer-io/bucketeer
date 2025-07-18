import { FC, useEffect, memo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { exchangeDemoToken, ExchangeTokenPayload } from '@api/auth';
import { urls } from 'configs';
import { PAGE_PATH_DEMO_SITE } from 'constants/routing';
import { getCookieState } from 'cookie';
import { useSubmit, useToast } from 'hooks';
import queryString from 'query-string';
import { AppLoading } from 'app';
import { useAuth } from './auth-context';

export const AuthDemoCallbackPage: FC = memo(() => {
  const { errorNotify } = useToast();
  const { setIsGoogleAuthError, setIsInitialLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const query = location.search;

  const { onSubmit: onGoogleDemoHandler } = useSubmit(
    async (payload: ExchangeTokenPayload) => {
      try {
        const response = await exchangeDemoToken(payload);
        if (response.token) {
          navigate(`${PAGE_PATH_DEMO_SITE}?token=${response.token}`);
        }
      } catch (error) {
        setIsGoogleAuthError(true);
        setIsInitialLoading(false);
        errorNotify(error);
        navigate(PAGE_PATH_DEMO_SITE);
      }
    }
  );

  useEffect(() => {
    const { code, state } = queryString.parse(query);
    const cookieState = getCookieState();
    setIsInitialLoading(true);

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
