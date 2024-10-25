import { memo, useCallback, useEffect, useState } from 'react';
import {
  BrowserRouter,
  Route,
  Routes,
  useParams,
  useNavigate
} from 'react-router-dom';
import {
  AuthCallbackPage,
  AuthProvider,
  useAuth,
  getCurrentEnvironment,
  hasEditable
} from 'auth';
import {
  PAGE_PATH_AUTH_CALLBACK,
  PAGE_PATH_AUTH_SIGNIN,
  PAGE_PATH_FEATURES,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
  PAGE_PATH_ROOT_ALL
} from 'constants/routing';
import { getTokenStorage } from 'storage/token';
import { v4 as uuid } from 'uuid';
import { ConsoleAccount } from '@types';
import DashboardPage from 'pages/dashboard';
import NotFoundPage from 'pages/not-found';
import SignInPage from 'pages/signin';
import SignInEmailPage from 'pages/signin/email';
import SelectOrganizationPage from 'pages/signin/organization';
import Navigation from 'components/navigation';
import Spinner from 'components/spinner';

export const AppLoading = () => (
  <div className="flex items-center justify-center h-screen w-full">
    <Spinner size="md" />
  </div>
);

function App() {
  return (
    <BrowserRouter basename="/v3">
      <AuthProvider>
        <Routes>
          <Route
            path={PAGE_PATH_AUTH_CALLBACK}
            element={<AuthCallbackPage />}
          />
          <Route path={PAGE_PATH_AUTH_SIGNIN} element={<SignInEmailPage />} />
          <Route path={PAGE_PATH_ROOT_ALL} element={<Root />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export const Root = memo(() => {
  const authToken = getTokenStorage();
  const [pageKey, setPageKey] = useState<string>(uuid());
  const { isInitialLoading, isLogin, consoleAccount, myOrganizations } =
    useAuth();

  const handleChangePageKey = useCallback(() => {
    setPageKey(uuid());
  }, [setPageKey]);

  if (isInitialLoading) {
    return <AppLoading />;
  }

  if (isLogin && consoleAccount) {
    return (
      <div className="flex flex-row w-full h-full">
        <Navigation onClickNavLink={handleChangePageKey} />
        <div className="flex-grow ml-[248px] shadow-lg overflow-y-auto">
          <Routes>
            <Route
              key={pageKey}
              path={'/:envUrlCode?/*'}
              element={<EnvironmentRoot account={consoleAccount} />}
            />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </div>
      </div>
    );
  }

  if (!!authToken && myOrganizations.length > 1) {
    return <SelectOrganizationPage />;
  }

  return <SignInPage />;
});

export const EnvironmentRoot = memo(
  ({ account }: { account: ConsoleAccount }) => {
    const navigate = useNavigate();
    const { envUrlCode } = useParams();

    const editable = hasEditable(account);
    const currentEnv = getCurrentEnvironment(account);

    useEffect(() => {
      if (!envUrlCode) {
        navigate(`${PAGE_PATH_ROOT}${currentEnv.urlCode}${PAGE_PATH_FEATURES}`);
      }
    }, [account, envUrlCode]);

    return (
      <Routes>
        {!editable && (
          <Route path={`/:any${PAGE_PATH_NEW}`}>
            <h3>{`403 Access denied`}</h3>
          </Route>
        )}
        <Route path={`${PAGE_PATH_FEATURES}`} element={<DashboardPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    );
  }
);

export default App;
