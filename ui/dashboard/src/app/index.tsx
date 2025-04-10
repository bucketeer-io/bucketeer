import { memo, useCallback, useEffect, useState } from 'react';
import { I18nextProvider } from 'react-i18next';
import {
  BrowserRouter,
  Route,
  Routes,
  useParams,
  useNavigate
} from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import {
  AuthCallbackPage,
  AuthProvider,
  useAuth,
  getCurrentEnvironment,
  hasEditable
} from 'auth';
import {
  PAGE_PATH_APIKEYS,
  PAGE_PATH_AUDIT_LOGS,
  PAGE_PATH_AUTH_CALLBACK,
  PAGE_PATH_AUTH_SIGNIN,
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURES,
  PAGE_PATH_GOALS,
  PAGE_PATH_MEMBERS,
  PAGE_PATH_NEW,
  PAGE_PATH_NOTIFICATIONS,
  PAGE_PATH_ORGANIZATIONS,
  PAGE_PATH_PROJECTS,
  PAGE_PATH_PUSHES,
  PAGE_PATH_ROOT,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_USER_SEGMENTS
} from 'constants/routing';
import { i18n } from 'i18n';
import { getTokenStorage } from 'storage/token';
import { v4 as uuid } from 'uuid';
import { ConsoleAccount } from '@types';
import APIKeysPage from 'pages/api-keys';
import AuditLogsPage from 'pages/audit-logs';
import MembersPage from 'pages/members';
import NotFoundPage from 'pages/not-found';
import NotificationsPage from 'pages/notifications';
import PushesPage from 'pages/pushes';
import SettingsPage from 'pages/settings';
import SignInPage from 'pages/signin';
import SignInEmailPage from 'pages/signin/email';
import SelectOrganizationPage from 'pages/signin/organization';
import UserSegmentsPage from 'pages/user-segments';
import Navigation from 'components/navigation';
import Spinner from 'components/spinner';
import {
  ExperimentsRoot,
  OrganizationsRoot,
  ProjectsRoot,
  GoalsRoot,
  FeatureFlagsRoot
} from './routers';

export const AppLoading = () => (
  <div className="flex items-center justify-center h-screen w-full">
    <Spinner size="md" />
  </div>
);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 60 * 1000 // Set the global stale time to 30 minutes
    }
  }
});

function App() {
  return (
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter basename="/v3">
          <AuthProvider>
            <Routes>
              <Route
                path={PAGE_PATH_AUTH_CALLBACK}
                element={<AuthCallbackPage />}
              />
              <Route
                path={PAGE_PATH_AUTH_SIGNIN}
                element={<SignInEmailPage />}
              />
              <Route path={`${PAGE_PATH_ROOT}*`} element={<Root />} />
            </Routes>
          </AuthProvider>
        </BrowserRouter>
        {/* {process.env.NODE_ENV === 'development' && (
          <ReactQueryDevtools initialIsOpen={false} />
        )} */}
      </QueryClientProvider>
    </I18nextProvider>
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
        <div className="w-full ml-[248px] shadow-lg overflow-y-auto">
          <Routes>
            {consoleAccount.isSystemAdmin && (
              <Route
                path={`${PAGE_PATH_ORGANIZATIONS}/*`}
                element={<OrganizationsRoot />}
              />
            )}
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
          <Route
            path={`/:any${PAGE_PATH_NEW}`}
            element={<h3>{`403 Access denied`}</h3>}
          />
        )}
        <Route
          path={`${PAGE_PATH_FEATURES}/*`}
          element={<FeatureFlagsRoot />}
        />
        <Route path={`${PAGE_PATH_SETTINGS}`} element={<SettingsPage />} />
        <Route path={`${PAGE_PATH_PROJECTS}/*`} element={<ProjectsRoot />} />
        <Route path={`${PAGE_PATH_APIKEYS}`} element={<APIKeysPage />} />
        <Route path={`${PAGE_PATH_MEMBERS}`} element={<MembersPage />} />
        <Route
          path={`${PAGE_PATH_NOTIFICATIONS}`}
          element={<NotificationsPage />}
        />
        <Route path={`${PAGE_PATH_PUSHES}`} element={<PushesPage />} />
        <Route path={`${PAGE_PATH_GOALS}/*`} element={<GoalsRoot />} />
        <Route
          path={`${PAGE_PATH_USER_SEGMENTS}`}
          element={<UserSegmentsPage />}
        />
        <Route
          path={`${PAGE_PATH_EXPERIMENTS}/*`}
          element={<ExperimentsRoot />}
        />
        <Route path={`${PAGE_PATH_AUDIT_LOGS}`} element={<AuditLogsPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    );
  }
);

export default App;
