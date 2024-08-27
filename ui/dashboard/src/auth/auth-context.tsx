import React, {
  createContext,
  useContext,
  useEffect,
  ReactNode,
  useState
} from 'react';
import { useNavigate } from 'react-router-dom';
import { accountOrganizationFetcher } from '@api/account';
import { accountMeFetcher } from '@api/account';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { Undefinable } from 'option-t/undefinable';
import {
  getCurrentEnvIdStorage,
  setCurrentEnvIdStorage
} from 'storage/environment';
import { clearOrgIdStorage, getOrgIdStorage } from 'storage/organization';
import {
  clearTokenStorage,
  getTokenStorage,
  setTokenStorage
} from 'storage/token';
import { AuthToken, ConsoleAccount, Organization } from '@types';

interface AuthContextType {
  syncSignIn: (authToken: AuthToken) => void;
  logout: () => void;
  isLogin: boolean;
  isInitialLoading: boolean;
  consoleAccount: Undefinable<ConsoleAccount>;
  myOrganizations: Array<Organization>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children
}) => {
  const navigate = useNavigate();
  const authToken: AuthToken | null = getTokenStorage();
  const organizationId = getOrgIdStorage();
  const environmentId = getCurrentEnvIdStorage();

  const [isInitialLoading, setIsInitialLoading] = useState(false);
  const [isLogin, setIsLogin] = useState<boolean>(false);
  const [consoleAccount, setConsoleAccount] =
    useState<Undefinable<ConsoleAccount>>();

  const [myOrganizations, setMyOrganizations] = useState<Array<Organization>>(
    []
  );

  const onMeFetcher = (id: string) => {
    accountMeFetcher({ organizationId: id })
      .then(response => {
        setConsoleAccount(response.account);
        setIsLogin(true);
        if (!environmentId) {
          setCurrentEnvIdStorage(
            response.account.environmentRoles[0].environment.id
          );
        }
      })
      .finally(() => setIsInitialLoading(false));
  };

  const onSyncAuthentication = () => {
    setIsInitialLoading(true);
    if (organizationId) {
      onMeFetcher(organizationId);
    } else {
      accountOrganizationFetcher().then(response => {
        const organizationsList = response.organizations || [];
        if (organizationsList.length === 1) {
          onMeFetcher(organizationsList[0].id);
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
        logout,
        isLogin,
        isInitialLoading,
        consoleAccount,
        myOrganizations
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
