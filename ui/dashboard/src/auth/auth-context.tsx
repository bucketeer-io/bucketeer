import React, {
  createContext,
  useContext,
  useEffect,
  ReactNode,
  useState
} from 'react';
import { useNavigate } from 'react-router-dom';
import { accountOrganizationFetcher, MeFetcherPayload } from '@api/account';
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
  syncSignIn: (authToken: AuthToken) => void;
  onMeFetcher: (payload: MeFetcherPayload) => Promise<void>;
  logout: () => void;
  isLogin: boolean;
  isInitialLoading: boolean;
  consoleAccount: Undefinable<ConsoleAccount>;
  myOrganizations: Array<Organization>;

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

  const onMeFetcher = (payload: MeFetcherPayload) => {
    return accountMeFetcher(payload).then(response => {
      setConsoleAccount(response.account);
      setIsLogin(true);
      if (!environmentId) {
        setCurrentEnvIdStorage(
          response.account.environmentRoles[0].environment.id
        );
      }
      setIsInitialLoading(false);
    });
  };

  const onSyncAuthentication = () => {
    if (organizationId) {
      onMeFetcher({ organizationId });
    } else {
      accountOrganizationFetcher().then(response => {
        const organizationsList = response.organizations || [];
        if (organizationsList.length === 1) {
          setOrgIdStorage(organizationsList[0].id);
          onMeFetcher({ organizationId: organizationsList[0].id });
        } else {
          setMyOrganizations(organizationsList);
          setIsInitialLoading(false);
        }
      });
    }
  };

  const syncSignIn = (authToken: AuthToken) => {
    setTokenStorage(authToken);
    onSyncAuthentication();
  };

  const logout = () => {
    setConsoleAccount(undefined);
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
        syncSignIn,
        onMeFetcher,
        logout,
        isLogin,
        isInitialLoading,
        consoleAccount,
        myOrganizations,

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
