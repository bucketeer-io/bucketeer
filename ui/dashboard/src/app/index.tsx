import { memo, useCallback, useEffect, useState } from 'react';
import { I18nextProvider } from 'react-i18next';
import {
  BrowserRouter,
  Route,
  Routes,
  useLocation,
  useNavigate,
  useParams
} from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import logo from 'assets/logos/logo-white.svg';
import {
  AuthCallbackPage,
  AuthDemoCallbackPage,
  AuthProvider,
  getCurrentEnvironment,
  hasEditable,
  useAuth
} from 'auth';
import { ENVIRONMENT_WITH_EMPTY_ID } from 'constants/app';
import {
  PAGE_PATH_APIKEYS,
  PAGE_PATH_AUDIT_LOGS,
  PAGE_PATH_AUTH_CALLBACK,
  PAGE_PATH_AUTH_DEMO_CALLBACK,
  PAGE_PATH_AUTH_SIGNIN,
  PAGE_PATH_DEBUGGER,
  PAGE_PATH_DEMO_SITE,
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
import { useScreen } from 'hooks';
import { ConfirmProvider } from 'hooks/use-unsaved-leave-page';
import { i18n } from 'i18n';
import pickBy from 'lodash/pickBy';
import {
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { getIsLoginFirstTimeStorage } from 'storage/login';
import {
  getCurrentProjectEnvironmentStorage,
  setCurrentProjectEnvironmentStorage
} from 'storage/project-environment';
import { getTokenStorage } from 'storage/token';
import { v4 as uuid } from 'uuid';
import { ConsoleAccount, EnvironmentRole } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId } from 'utils/function';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import { IconMenu } from '@icons';
import AccessDeniedPage from 'pages/access-denied';
import APIKeysPage from 'pages/api-keys';
import AuditLogsPage from 'pages/audit-logs';
import DebuggerPage from 'pages/debugger';
import AccessDemoPage from 'pages/demo';
import CreateDemoPage from 'pages/demo/demo-create';
import NotFoundPage from 'pages/not-found';
import NotificationsPage from 'pages/notifications';
import PushesPage from 'pages/pushes';
import SettingsPage from 'pages/settings';
import SignInPage from 'pages/signin';
import SignInEmailPage from 'pages/signin/email';
import UserInformation from 'pages/signin/information';
import SelectOrganizationPage from 'pages/signin/organization';
import UserSegmentsPage from 'pages/user-segments';
import Button from 'components/button';
import Drawer from 'components/drawer';
import Icon from 'components/icon';
import Navigation from 'components/navigation';
import Spinner from 'components/spinner';
import {
  ExperimentsRoot,
  FeatureFlagsRoot,
  GoalsRoot,
  MemberRoot,
  OrganizationsRoot,
  ProjectsRoot
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
        <ConfirmProvider>
          <BrowserRouter>
            <AuthProvider>
              <Routes>
                <Route
                  path={PAGE_PATH_AUTH_CALLBACK}
                  element={<AuthCallbackPage />}
                />
                <Route
                  path={PAGE_PATH_AUTH_DEMO_CALLBACK}
                  element={<AuthDemoCallbackPage />}
                />
                <Route
                  path={PAGE_PATH_AUTH_SIGNIN}
                  element={<SignInEmailPage />}
                />
                <Route
                  path={PAGE_PATH_DEMO_SITE}
                  element={<AccessDemoPage />}
                />
                <Route
                  path={`${PAGE_PATH_DEMO_SITE}/new`}
                  element={<CreateDemoPage />}
                />
                <Route path={`${PAGE_PATH_ROOT}*`} element={<Root />} />
              </Routes>
            </AuthProvider>
          </BrowserRouter>
        </ConfirmProvider>
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
  const [showMenu, setShowMenu] = useState<boolean>(false);
  const { fromMobileScreen } = useScreen();
  const { isInitialLoading, isLogin, consoleAccount, myOrganizations } =
    useAuth();

  const handleChangePageKey = useCallback(() => {
    setPageKey(uuid());
  }, [setPageKey]);

  if (isInitialLoading) {
    return <AppLoading />;
  }

  if (isLogin && consoleAccount) {
    const isLoginFirstTime = getIsLoginFirstTimeStorage();
    if (isLoginFirstTime) {
      return <UserInformation />;
    }
    return (
      <div className="flex flex-col sm:flex-row w-full h-full">
        <div className="flex sticky z-20 top-0 left-0 gap-3 items-center justify-between sm:hidden w-full h-[50px] px-4 bg-primary-400">
          <img src={logo} alt="Bucketer" />
          <Button
            className="bg-transparent hover:bg-transparent p-0"
            onClick={() => setShowMenu(true)}
          >
            <Icon icon={IconMenu} />
          </Button>
        </div>
        {fromMobileScreen ? (
          <Navigation onClickNavLink={handleChangePageKey} />
        ) : (
          <Drawer
            side="left"
            open={showMenu}
            onClose={() => setShowMenu(false)}
          >
            <Navigation onClickNavLink={handleChangePageKey} />
          </Drawer>
        )}

        <div className="w-full ml-0 sm:ml-[60px] md:ml-[248px] shadow-lg overflow-y-auto">
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
    const { pathname } = useLocation();
    const { envUrlCode, ...params } = useParams();
    const { searchOptions } = useSearchParams();

    const editable = hasEditable(account);
    const currentEnv = getCurrentEnvironment(account);

    const handleCheckEnvCodeOnInit = useCallback(() => {
      const envIdStorage = getCurrentEnvIdStorage();
      const projectEnvironment = getCurrentProjectEnvironmentStorage();
      let isExistEnv: EnvironmentRole | undefined = undefined;
      if (projectEnvironment) {
        const projectEnvironmentId = checkEnvironmentEmptyId(
          projectEnvironment.environmentId
        );
        isExistEnv = account.environmentRoles.find(item => {
          const { environment, project } = item || {};
          return (
            environment.id === projectEnvironmentId &&
            project.id === projectEnvironment.projectId
          );
        });
        if (envUrlCode && isExistEnv?.environment.urlCode !== envUrlCode)
          isExistEnv = undefined;
      }
      if (!isExistEnv) {
        isExistEnv = account.environmentRoles.find(
          item => item.environment.urlCode === envUrlCode
        );
      }

      if (!envUrlCode || !isExistEnv) {
        setCurrentEnvIdStorage(currentEnv.id || ENVIRONMENT_WITH_EMPTY_ID);
        setCurrentProjectEnvironmentStorage({
          environmentId: currentEnv.id || ENVIRONMENT_WITH_EMPTY_ID,
          projectId: currentEnv.projectId
        });
        return navigate(`/${currentEnv.urlCode}${PAGE_PATH_FEATURES}`, {
          replace: true
        });
      }

      const { environment } = isExistEnv || {};
      if (environment.id === envIdStorage && environment.urlCode === envUrlCode)
        return;

      const stringifyQueryParams = stringifyParams(
        pickBy(searchOptions, v => isNotEmpty(v as string))
      );
      const queryParams = isNotEmpty(stringifyQueryParams)
        ? `?${stringifyQueryParams}`
        : '';
      const path = params['*'] ? `/${params['*']}` : '';

      setCurrentEnvIdStorage(environment.id || ENVIRONMENT_WITH_EMPTY_ID);
      setCurrentProjectEnvironmentStorage({
        environmentId: environment.id,
        projectId: environment.projectId
      });
      return navigate(`/${environment.urlCode}${path}${queryParams}`, {
        replace: true
      });
    }, [envUrlCode, currentEnv, params, searchOptions, account]);

    useEffect(() => {
      handleCheckEnvCodeOnInit();
    }, [account, envUrlCode]);

    if (pathname === '/') return <AppLoading />;

    return (
      <Routes>
        {!editable && (
          <Route
            path={`/:any${PAGE_PATH_NEW}`}
            element={<AccessDeniedPage />}
          />
        )}
        <Route
          path={`${PAGE_PATH_FEATURES}/*`}
          element={<FeatureFlagsRoot />}
        />
        <Route path={`${PAGE_PATH_SETTINGS}`} element={<SettingsPage />} />
        <Route path={`${PAGE_PATH_PROJECTS}/*`} element={<ProjectsRoot />} />
        <Route path={`${PAGE_PATH_APIKEYS}/*`} element={<APIKeysPage />} />
        <Route path={`${PAGE_PATH_MEMBERS}/*`} element={<MemberRoot />} />
        <Route
          path={`${PAGE_PATH_NOTIFICATIONS}/*`}
          element={<NotificationsPage />}
        />
        <Route path={`${PAGE_PATH_PUSHES}/*`} element={<PushesPage />} />
        <Route path={`${PAGE_PATH_GOALS}/*`} element={<GoalsRoot />} />
        <Route
          path={`${PAGE_PATH_USER_SEGMENTS}/*`}
          element={<UserSegmentsPage />}
        />
        <Route
          path={`${PAGE_PATH_EXPERIMENTS}/*`}
          element={<ExperimentsRoot />}
        />
        <Route path={`${PAGE_PATH_AUDIT_LOGS}/*`} element={<AuditLogsPage />} />
        <Route path={`${PAGE_PATH_DEBUGGER}/*`} element={<DebuggerPage />} />

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    );
  }
);

export default App;
