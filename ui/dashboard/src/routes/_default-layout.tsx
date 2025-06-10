import { memo, useCallback, useEffect } from 'react';
import { createRoute, Outlet, useLocation } from '@tanstack/react-router';
import { useNavigate } from '@tanstack/react-router';
import { useParams } from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import DefaultEnvLayout from 'elements/default-env-layout';
import { Route as RootRoute } from './__root';

export const PathlessLayout = memo(() => {
  const { consoleAccount } = useAuth();
  const navigate = useNavigate();
  const { pathname, searchStr } = useLocation();
  const { env: envUrlCode, ...params } = useParams({
    strict: false
  });

  // const editable = hasEditable(account);
  const currentEnv = getCurrentEnvironment(consoleAccount!);

  const handleCheckEnvCodeOnInit = useCallback(() => {
    const isExistEnv = consoleAccount?.environmentRoles?.find(
      item => item.environment.urlCode === envUrlCode
    );
    if (!envUrlCode || !isExistEnv) return;

    const envIdStorage = getCurrentEnvIdStorage();
    if (envIdStorage === envUrlCode) return;
    const { environment } = isExistEnv;

    setCurrentEnvIdStorage(environment.id || environment.urlCode);
    return navigate({
      to: `${pathname}${searchStr}`
    });
  }, [envUrlCode, currentEnv, params, consoleAccount]);

  useEffect(() => {
    if (consoleAccount && envUrlCode) handleCheckEnvCodeOnInit();
  }, [consoleAccount, envUrlCode]);

  return (
    <DefaultEnvLayout>
      <Outlet />
    </DefaultEnvLayout>
  );
});

export const Route = createRoute({
  id: 'pathlessLayout',
  component: PathlessLayout,
  getParentRoute: () => RootRoute
});
