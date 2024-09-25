import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQueryProjectDetails } from '@queries/project-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_PROJECTS } from 'constants/routing';
import PageDetailHeader from 'containers/page-details-header';
import { ProjectDetailsContent } from 'containers/pages/project-details/';
import { environmentTab, settingTab } from 'helpers/tab';
import { useAddQuery, useQuery } from 'hooks';
import { useFormatDateTime } from 'utils/date-time';
import Spinner from 'components/spinner';

const tabs = [environmentTab, settingTab];

const ProjectDetailsPage = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const query = useQuery();
  const { addQuery } = useAddQuery();
  const tab = query.get('tab');

  const params = useParams();
  const { projectId } = params;

  const formatDateTime = useFormatDateTime();

  const { data, isLoading } = useQueryProjectDetails({
    params: {
      id: projectId as string
    },
    enabled: !!projectId
  });

  const [targetTab, setTargetTab] = useState(tab || tabs[0].value);

  const handleChangeTab = (value: string) => {
    setTargetTab(value);
    addQuery(query, { tab: value });
  };

  return (
    <div className="flex flex-col size-full overflow-auto">
      {isLoading ? (
        <div className="pt-20 flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <>
          <PageDetailHeader
            title={data?.project?.name || ''}
            description={
              data?.project?.createdAt
                ? formatDateTime(data?.project?.createdAt)
                : ''
            }
            navigateRoute={`/${currentEnvironment.id}${PAGE_PATH_PROJECTS}`}
            tabs={tabs}
            targetTab={targetTab}
            status="new"
            onSelectTab={handleChangeTab}
          />
          <ProjectDetailsContent
            targetTab={targetTab}
            projectId={projectId}
            projectData={data?.project}
          />
        </>
      )}
    </div>
  );
};

export default ProjectDetailsPage;
