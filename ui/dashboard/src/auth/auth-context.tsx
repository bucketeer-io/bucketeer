import React, {
  createContext,
  useContext,
  useEffect,
  ReactNode,
  useState,
  useCallback
} from 'react';
import { useNavigate } from 'react-router-dom';
import { accountOrganizationFetcher, MeFetcherParams } from '@api/account';
import { accountMeFetcher } from '@api/account';
import { urls } from 'configs';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToast } from 'hooks';
import { getLanguage, Language, setLanguage, useTranslation } from 'i18n';
import { isNil } from 'lodash';
import { Undefinable } from 'option-t/undefinable';
import { clearConsoleVersion, getConsoleVersion } from 'storage/console';
import {
  clearCurrentEnvIdStorage,
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { setIsLoginFirstTimeStorage } from 'storage/login';
import {
  clearOrgIdStorage,
  getOrgIdStorage,
  setOrgIdStorage
} from 'storage/organization';
import {
  clearCurrentProjectEnvironmentStorage,
  setCurrentProjectEnvironmentStorage
} from 'storage/project-environment';
import {
  clearTokenStorage,
  getTokenStorage,
  setTokenStorage
} from 'storage/token';
import { AuthToken, ConsoleAccount, Organization } from '@types';
import { onChangeFontWithLocalized } from 'utils/function';
import { useSearchParams } from 'utils/search-params';
import { getAccountAccess } from './utils';

interface AuthContextType {
  logout: () => void;
  isLogin: boolean;

  consoleAccount: Undefinable<ConsoleAccount>;
  myOrganizations: Array<Organization>;

  syncSignIn: (authToken: AuthToken) => Promise<void>;
  onMeFetcher: (params: MeFetcherParams) => Promise<void>;

  isInitialLoading: boolean;
  setIsInitialLoading: (v: boolean) => void;

  isGoogleAuthError: boolean;
  setIsGoogleAuthError: (v: boolean) => void;
  isAccessDemoPage: boolean;
  setIsAccessDemoPage: (v: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children
}) => {
  const { t } = useTranslation(['message']);
  const navigate = useNavigate();
  const authToken: AuthToken | null = getTokenStorage();
  const organizationId = getOrgIdStorage();
  const environmentId = getCurrentEnvIdStorage();
  const { errorNotify } = useToast();
  const { searchOptions } = useSearchParams();

  const [isInitialLoading, setIsInitialLoading] = useState(
    !!authToken?.accessToken
  );
  const [isLogin, setIsLogin] = useState<boolean>(false);

  const [consoleAccount, setConsoleAccount] =
    useState<Undefinable<ConsoleAccount>>();

  const [myOrganizations, setMyOrganizations] = useState<Organization[]>([]);
  const [isGoogleAuthError, setIsGoogleAuthError] = useState(false);
  const [isAccessDemoPage, setIsAccessDemoPage] = useState(true);

  const clearOrgAndEnvStorage = () => {
    clearOrgIdStorage();
    clearCurrentEnvIdStorage();
    clearCurrentProjectEnvironmentStorage();
  };

  const handleCheckConsoleVersion = useCallback(
    (account: ConsoleAccount) => {
      const backFromOldConsole = !!searchOptions?.fromOldConsole;
      const consoleVersion = getConsoleVersion();
      const isRedirectToOldConsole =
        consoleVersion?.version === 'old' &&
        consoleVersion?.email === account.email;
      if (!isRedirectToOldConsole || backFromOldConsole)
        return clearConsoleVersion();

      if (isRedirectToOldConsole)
        return (window.location.href = urls.OLD_CONSOLE_ENDPOINT as string);
    },
    [searchOptions]
  );

  const onMeFetcher = async (params: MeFetcherParams) => {
    try {
      const response = await accountMeFetcher(params);
      const environmentRoles = response.account.environmentRoles;
      if (!environmentRoles?.length) {
        clearOrgAndEnvStorage();
        errorNotify(null, t('message:env-are-empty'));
        return logout();
      }
      handleCheckConsoleVersion(response.account);
      setConsoleAccount(response.account);
      setIsLogin(true);
      if (response.account.lastSeen === '0' || !response.account.lastSeen)
        return setIsLoginFirstTimeStorage(true);
      const isJapanese = response.account.language === Language.JAPANESE;
      onChangeFontWithLocalized(isJapanese);

      if (response.account.language !== getLanguage()) {
        await setLanguage(response.account.language as Language);
      }
      if (isNil(environmentId)) {
        const environment = environmentRoles[0].environment;
        setCurrentEnvIdStorage(environment.id);
        setCurrentProjectEnvironmentStorage({
          environmentId: environment.id,
          projectId: environment.projectId
        });
      }
    } catch (error) {
      clearOrgAndEnvStorage();
      errorNotify(error, t('message:org-not-found'));
    } finally {
      setIsInitialLoading(false);
    }
  };

  const onSyncAuthentication = async () => {
    try {
      const response = await accountOrganizationFetcher();
      const organizationsList = response.organizations || [];
      const isExistOrg = organizationsList.find(
        item => item.id === organizationId
      );
      if (organizationId && isExistOrg) {
        await onMeFetcher({ organizationId });
      } else if (organizationsList.length === 1) {
        setOrgIdStorage(organizationsList[0].id);
        await onMeFetcher({ organizationId: organizationsList[0].id });
      } else {
        setIsInitialLoading(false);
      }
      setMyOrganizations(organizationsList);
    } catch (error) {
      errorNotify(error);
    }
  };

  const syncSignIn = async (authToken: AuthToken) => {
    setTokenStorage(authToken);
    onSyncAuthentication();
  };

  const logout = () => {
    setConsoleAccount(undefined);
    setMyOrganizations([]);
    setIsLogin(false);
    clearTokenStorage();
    navigate(PAGE_PATH_ROOT);
  };

  useEffect(() => {
    if (authToken) {
      onSyncAuthentication();
    }
  }, []);

  useEffect(() => {
    const handleTokenRefreshed = () => {
      setIsInitialLoading(false);
    };
    window.addEventListener('tokenRefreshed', handleTokenRefreshed);
    window.addEventListener('unauthenticated', () => {
      logout();
      clearOrgIdStorage();
    });
    return () => {
      window.removeEventListener('tokenRefreshed', handleTokenRefreshed);
      window.removeEventListener('unauthenticated', () => {
        logout();
        clearOrgIdStorage();
      });
    };
  }, []);

  return (
    <AuthContext.Provider
      value={{
        isLogin,
        logout,

        consoleAccount,
        myOrganizations,

        syncSignIn,
        onMeFetcher,

        isInitialLoading,
        setIsInitialLoading,

        isGoogleAuthError,
        setIsGoogleAuthError,

        isAccessDemoPage,
        setIsAccessDemoPage
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  const { t } = useTranslation(['message']);

  if (!context) {
    throw new Error(t('auth-context-error'));
  }
  return context;
};

export const useAuthAccess = () => {
  const { consoleAccount } = useAuth();
  return getAccountAccess(consoleAccount!);
};
