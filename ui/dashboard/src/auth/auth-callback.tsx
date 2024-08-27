import { FC, useEffect, memo } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { exchangeToken, ExchangeTokenPayload } from '@api/auth';
import { urls } from 'configs';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { getCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import queryString from 'query-string';
import { setTokenStorage } from 'storage/token';
import Spinner from 'components/spinner';
import { useAuth } from './auth-context';

export const AuthCallbackPage: FC = memo(() => {
  const { syncSignIn } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const query = location.search;

  const { onSubmit: onGoogleLoginHandler } = useSubmit(
    (payload: ExchangeTokenPayload) =>
      exchangeToken(payload).then(response => {
        if (response.token) {
          setTokenStorage(response.token);
          syncSignIn(response.token);
          navigate(PAGE_PATH_ROOT);
        }
      })
  );

  useEffect(() => {
    const { code, state } = queryString.parse(query);
    const cookieState = getCookieState();
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

  return (
    <div className="mt-20 flex justify-center w-full">
      <Spinner />
    </div>
  );
});
