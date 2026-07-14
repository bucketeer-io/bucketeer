import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import PageContent from './page-content';

const PageLoader = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  return (
    <PageContent disabled={!editable} environmentId={currentEnvironment.id} />
  );
};

export default PageLoader;
