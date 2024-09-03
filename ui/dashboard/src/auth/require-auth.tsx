import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from 'auth';
import { AppLoading } from 'app';

const RequireAuth = () => {
  const { isLogin, isInitialLoading } = useAuth();

  if (isInitialLoading) {
    return <AppLoading />;
  }

  if (!isLogin) {
    return null;
  }

  return isLogin ? <Outlet /> : <Navigate to="/" />;
};

export default RequireAuth;
