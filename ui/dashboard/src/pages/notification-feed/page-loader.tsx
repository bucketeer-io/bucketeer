import { useAuth } from 'auth';
import PageContent from './page-content';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  // Creating, editing, publishing, and deleting notifications is system
  // admin only (see proto/notification/service.proto); environment editor
  // role does not grant this.
  const isSystemAdmin = !!consoleAccount?.isSystemAdmin;

  return (
    <PageContent disabled={!isSystemAdmin} isSystemAdmin={isSystemAdmin} />
  );
};

export default PageLoader;
