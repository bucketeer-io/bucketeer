import { memo } from 'react';
import { Outlet, useLocation } from '@tanstack/react-router';
import { useAuth } from 'auth';
import { getTokenStorage } from 'storage/token';
import { AppLoading } from 'app';

const RootLayout = memo(() => {
  const authToken = getTokenStorage();
  const { consoleAccount } = useAuth();

  const { pathname } = useLocation();

  if (!consoleAccount && pathname !== '/v3/' && !!authToken?.accessToken)
    return <AppLoading />;

  return (
    <div>
      <Outlet />
    </div>
  );
});

export default RootLayout;
