import { getCurrentEnvironment, useAuth } from 'auth';
import PageContent from './page-content';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  // Creating, editing, publishing, and deleting notifications is system
  // admin only (see proto/notification/service.proto); environment editor
  // role does not grant this.
  const isSystemAdmin = !!consoleAccount?.isSystemAdmin;

  return (
    <PageContent
      disabled={!isSystemAdmin}
      isSystemAdmin={isSystemAdmin}
      environmentId={currentEnvironment.id}
    />
  );
};

export default PageLoader;
