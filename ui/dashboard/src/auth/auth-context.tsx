import React, {
  createContext,
  useContext,
  useEffect,
  ReactNode,
  useState
} from 'react';
import { accountOrganizationFetcher, MeFetcherParams } from '@api/account';
import { accountMeFetcher } from '@api/account';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToast } from 'hooks';
import { Undefinable } from 'option-t/undefinable';
import { queryClient } from 'router';
import {
  clearCurrentEnvIdStorage,
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import {
  clearOrgIdStorage,
  getOrgIdStorage,
  setOrgIdStorage
} from 'storage/organization';
import {
  clearTokenStorage,
  getTokenStorage,
  setTokenStorage
} from 'storage/token';
import { AuthToken, ConsoleAccount, Organization } from '@types';

export interface AuthContextType {
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
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children
}) => {
  const authToken: AuthToken | null = getTokenStorage();
  const organizationId = getOrgIdStorage();
  const environmentId = getCurrentEnvIdStorage();
  const { errorNotify } = useToast();

  const [isInitialLoading, setIsInitialLoading] = useState(
    !!authToken?.accessToken
  );
  const [isLogin, setIsLogin] = useState<boolean>(false);
  const [consoleAccount, setConsoleAccount] =
    useState<Undefinable<ConsoleAccount>>();

  const [myOrganizations, setMyOrganizations] = useState<Organization[]>([]);
  const [isGoogleAuthError, setIsGoogleAuthError] = useState(false);

  const clearOrgAndEnvStorage = () => {
    clearOrgIdStorage();
    clearCurrentEnvIdStorage();
  };

  const onMeFetcher = async (params: MeFetcherParams) => {
    try {
      // const response = await accountMeFetcher(params);

      const response = await queryClient.ensureQueryData({
        queryKey: ['account', params.organizationId],
        queryFn: () =>
          accountMeFetcher({
            organizationId: params.organizationId
          })
      });

      const environmentRoles = response.account.environmentRoles;
      if (!environmentRoles.length) {
        clearOrgAndEnvStorage();
        errorNotify(null, 'The environments are empty.');
        return logout();
      }

      setConsoleAccount(response.account);
      setIsLogin(true);
      if (!environmentId) {
        setCurrentEnvIdStorage(environmentRoles[0].environment.id);
      }
    } catch (error) {
      clearOrgAndEnvStorage();
      errorNotify(error, 'The organization is not found.');
    } finally {
      setIsInitialLoading(false);
    }
  };

  const onSyncAuthentication = async () => {
    try {
      const response = await accountOrganizationFetcher();
      const organizationsList = response.organizations || [];
      if (organizationId) {
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
    window.location.replace(PAGE_PATH_ROOT);
  };

  useEffect(() => {
    if (authToken) {
      onSyncAuthentication();
    }
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
        setIsGoogleAuthError
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
