import { memo, useCallback, useEffect } from 'react';
import { createRoute, Outlet } from '@tanstack/react-router';
import { useNavigate } from '@tanstack/react-router';
import { useParams } from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from 'constants/routing';
import { pickBy } from 'lodash';
import {
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import Navigation from 'components/navigation';
import { Route as RootRoute } from './__root';

export const PathlessLayout = memo(() => {
  const { consoleAccount } = useAuth();
  const navigate = useNavigate();
  const { envUrlCode, ...params } = useParams({
    strict: false
  });
  const { searchOptions } = useSearchParams();

  // const editable = hasEditable(account);
  const currentEnv = getCurrentEnvironment(consoleAccount!);

  const handleCheckEnvCodeOnInit = useCallback(() => {
    const isExistEnv = consoleAccount?.environmentRoles?.find(
      item => item.environment.urlCode === envUrlCode
    );
    if (!envUrlCode || !isExistEnv) {
      return navigate({
        to: `${PAGE_PATH_ROOT}${currentEnv?.urlCode}${PAGE_PATH_FEATURES}`
      });
    }

    const envIdStorage = getCurrentEnvIdStorage();
    if (envIdStorage === envUrlCode) return;
    const { environment } = isExistEnv;

    const stringifyQueryParams = stringifyParams(
      pickBy(searchOptions, v => isNotEmpty(v as string))
    );
    const queryParams = isNotEmpty(stringifyQueryParams)
      ? `?${stringifyQueryParams}`
      : '';

    console.log({ stringifyQueryParams });

    const path = params['*'] ? `/${params['*']}` : '';

    setCurrentEnvIdStorage(environment.id || environment.urlCode);
    return navigate({
      to: `${PAGE_PATH_ROOT}${environment.urlCode}${path}${queryParams}`
    });
  }, [envUrlCode, currentEnv, params, searchOptions, consoleAccount]);

  useEffect(() => {
    if (consoleAccount && envUrlCode) handleCheckEnvCodeOnInit();
  }, [consoleAccount, envUrlCode]);

  return (
    <div className="flex size-full">
      <Navigation onClickNavLink={() => {}} />
      <div className="w-full ml-[248px] shadow-lg overflow-y-auto">
        <Outlet />
      </div>
    </div>
  );
});

export const Route = createRoute({
  id: 'pathlessLayout',
  component: PathlessLayout,
  getParentRoute: () => RootRoute
});
