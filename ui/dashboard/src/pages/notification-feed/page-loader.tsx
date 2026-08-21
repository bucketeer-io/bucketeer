import { useAuth } from 'auth';
import PageContent from './page-content';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const isSystemAdmin = !!consoleAccount?.isSystemAdmin;

  return (
    <PageContent disabled={!isSystemAdmin} isSystemAdmin={isSystemAdmin} />
  );
};

export default PageLoader;
