import React, {
  createContext,
  useContext,
  useEffect,
  ReactNode,
  useState
} from 'react';
import { useNavigate } from 'react-router-dom';
import { accountOrganizationFetcher, MeFetcherParams } from '@api/account';
import { accountMeFetcher } from '@api/account';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { Undefinable } from 'option-t/undefinable';
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

interface AuthContextType {
  logout: () => void;
  isLogin: boolean;

  consoleAccount: Undefinable<ConsoleAccount>;
  myOrganizations: Array<Organization>;

  syncSignIn: (authToken: AuthToken) => void;
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
  const navigate = useNavigate();
  const authToken: AuthToken | null = getTokenStorage();
  const organizationId = getOrgIdStorage();
  const environmentId = getCurrentEnvIdStorage();

  const [isInitialLoading, setIsInitialLoading] = useState(
    !!authToken?.accessToken
  );
  const [isLogin, setIsLogin] = useState<boolean>(false);
  const [consoleAccount, setConsoleAccount] =
    useState<Undefinable<ConsoleAccount>>();

  const [myOrganizations, setMyOrganizations] = useState<Organization[]>([]);
  const [isGoogleAuthError, setIsGoogleAuthError] = useState(false);

  const onMeFetcher = (params: MeFetcherParams) => {
    return accountMeFetcher(params)
      .then(response => {
        const environmentRoles = response.account.environmentRoles;

        if (environmentRoles.length > 0) {
          setConsoleAccount(response.account);
          setIsLogin(true);
          if (!environmentId) {
            setCurrentEnvIdStorage(environmentRoles[0].environment.id);
          }
        } else logout();
      })
      .finally(() => setIsInitialLoading(false));
  };

  const onSyncAuthentication = () => {
    accountOrganizationFetcher().then(response => {
      const organizationsList = response.organizations || [];
      if (organizationId) {
        onMeFetcher({ organizationId });
      } else if (organizationsList.length === 1) {
        setOrgIdStorage(organizationsList[0].id);
        onMeFetcher({ organizationId: organizationsList[0].id });
      } else {
        setIsInitialLoading(false);
      }
      setMyOrganizations(organizationsList);
    });
  };

  const syncSignIn = (authToken: AuthToken) => {
    setTokenStorage(authToken);
    onSyncAuthentication();
  };

  const logout = () => {
    setConsoleAccount(undefined);
    setMyOrganizations([]);
    setIsLogin(false);
    clearOrgIdStorage();
    clearTokenStorage();
    clearCurrentEnvIdStorage();
    navigate(PAGE_PATH_ROOT);
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
