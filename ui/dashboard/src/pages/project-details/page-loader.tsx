import { useNavigate, useParams } from 'react-router-dom';
import { useQueryProjectDetails } from '@queries/project-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import { useFormatDateTime } from 'utils/date-time';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();
  const { projectId } = useParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data, isLoading, refetch, isError } = useQueryProjectDetails({
    params: { id: projectId! }
  });

  const project = data?.project;
  const isErrorState = isError || !project;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader
            title={project.name}
            description={`Created ${formatDateTime(project.createdAt)}`}
            onBack={() =>
              navigate(`/${currentEnvironment.urlCode}/${PAGE_PATH_PROJECTS}`)
            }
          />
          <PageContent project={project} />
        </>
      )}
    </>
  );
};

export default PageLoader;
