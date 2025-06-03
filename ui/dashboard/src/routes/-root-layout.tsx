import React, { useCallback, useEffect } from 'react';
import {
  Outlet,
  useLocation,
  useNavigate,
  useParams
} from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from 'constants/routing';
import { pickBy } from 'lodash';
import {
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { getTokenStorage } from 'storage/token';
import { isNotEmpty } from 'utils/data-type';
import { stringifyParams, useSearchParams } from 'utils/search-params';
import { AppLoading } from 'app';

const RootLayout = () => {
  const authToken = getTokenStorage();
  const { consoleAccount } = useAuth();

  const { pathname } = useLocation();

  const navigate = useNavigate();
  const { envUrlCode, ...params } = useParams({
    from: ''
  });
  const { searchOptions } = useSearchParams();

  // const editable = hasEditable(account);
  const currentEnv = getCurrentEnvironment(consoleAccount!);

  const handleCheckEnvCodeOnInit = useCallback(() => {
    const isExistEnv = consoleAccount?.environmentRoles?.find(
      item => item.environment.urlCode === envUrlCode
    );
    console.log(consoleAccount?.environmentRoles);
    console.log({ isExistEnv });
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
    const path = params['*'] ? `/${params['*']}` : '';

    setCurrentEnvIdStorage(environment.id || environment.urlCode);
    return navigate({
      to: `${PAGE_PATH_ROOT}${environment.urlCode}${path}${queryParams}`
    });
  }, [envUrlCode, currentEnv, params, searchOptions, consoleAccount]);

  useEffect(() => {
    if (consoleAccount && envUrlCode) handleCheckEnvCodeOnInit();
  }, [consoleAccount, envUrlCode]);

  if (!consoleAccount && pathname !== '/v3/' && !!authToken?.accessToken)
    return <AppLoading />;

  return (
    <div>
      <Outlet />
    </div>
  );
};

export default RootLayout;
