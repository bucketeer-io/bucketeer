import { memo } from 'react';
import { I18nextProvider } from 'react-i18next';
import { QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider } from '@tanstack/react-router';
import { Navigate } from '@tanstack/react-router';
import { AuthProvider, useAuth, getCurrentEnvironment } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { i18n } from 'i18n';
import { queryClient, router } from 'router';
import { getTokenStorage } from 'storage/token';
import SignInPage from 'pages/signin';
import SelectOrganizationPage from 'pages/signin/organization';
import Spinner from 'components/spinner';

export const AppLoading = () => (
  <div className="flex items-center justify-center h-screen w-full">
    <Spinner size="md" />
  </div>
);

function InnerApp() {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

function App() {
  return (
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <InnerApp />
        </AuthProvider>
      </QueryClientProvider>
    </I18nextProvider>
  );
}

export const Root = memo(() => {
  const authToken = getTokenStorage();
  const { isInitialLoading, isLogin, consoleAccount, myOrganizations } =
    useAuth();

  if (isInitialLoading) {
    return <AppLoading />;
  }

  if (isLogin && consoleAccount) {
    const currentEnvironment = getCurrentEnvironment(consoleAccount!);
    return (
      <Navigate
        to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`}
        replace
      />
    );
  }

  if (!!authToken && myOrganizations.length > 1) {
    return <SelectOrganizationPage />;
  }

  return <SignInPage />;
});

export default App;
