import { FC, useEffect, memo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { exchangeToken, ExchangeTokenPayload } from '@api/auth';
import { urls } from 'configs';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { getCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import queryString from 'query-string';
import { setTokenStorage } from 'storage/token';
import { AppLoading } from 'app';
import { useAuth } from './auth-context';

export const AuthCallbackPage: FC = memo(() => {
  const { syncSignIn, setIsGoogleAuthError, setIsInitialLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const query = location.search;

  const { onSubmit: onGoogleLoginHandler } = useSubmit(
    async (payload: ExchangeTokenPayload) => {
      try {
        const response = await exchangeToken(payload);
        if (response.token) {
          setTokenStorage(response.token);
          await syncSignIn(response.token);
          navigate(PAGE_PATH_ROOT);
        }
      } catch {
        setIsGoogleAuthError(true);
        setIsInitialLoading(false);
        navigate(PAGE_PATH_ROOT);
      }
    }
  );

  useEffect(() => {
    const { code, state } = queryString.parse(query);
    const cookieState = getCookieState();
    setIsInitialLoading(true);
    if (!!code && state === cookieState) {
      if (typeof code === 'string') {
        onGoogleLoginHandler({
          code,
          redirectUrl: urls.AUTH_REDIRECT,
          type: 2 // Google auth type
        });
      }
    } else {
      throw new Error('exchange token failed.');
    }
  }, [query]);

  return <AppLoading />;
});
