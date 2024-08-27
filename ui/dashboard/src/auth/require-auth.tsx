import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from 'auth';
import Spinner from 'components/spinner';

const Loading = () => (
  <div className="mt-20 flex justify-center w-full">
    <Spinner />
  </div>
);

const RequireAuth = () => {
  const { isLogin, isInitialLoading } = useAuth();

  if (isInitialLoading) {
    return <Loading />;
  }

  if (!isLogin) {
    return null;
  }

  return isLogin ? <Outlet /> : <Navigate to="/" />;
};

export default RequireAuth;
