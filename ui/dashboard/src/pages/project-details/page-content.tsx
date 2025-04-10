import { Navigate, Route, Routes, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_PROJECTS,
  PAGE_PATH_ENVIRONMENTS,
  PAGE_PATH_SETTINGS
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Project } from '@types';
import { Tabs, TabsList, TabsContent, TabsLink } from 'components/tabs-link';
import PageLayout from 'elements/page-layout';
import ProjectEnvironments from './environments';
import ProjectSettings from './settings';
import { TabItem } from './types';

const PageContent = ({ project }: { project: Project }) => {
  const { t } = useTranslation(['common']);
  const { projectId } = useParams();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const url = `/${currentEnvironment.urlCode}${PAGE_PATH_PROJECTS}/${projectId}`;

  const projectTabs: Array<TabItem> = [
    {
      title: t(`environments`),
      to: PAGE_PATH_ENVIRONMENTS
    },
    {
      title: t(`settings`),
      to: PAGE_PATH_SETTINGS
    }
  ];

  return (
    <PageLayout.Content className="pt-4">
      <Tabs>
        <TabsList className="px-6">
          {projectTabs.map((item, index) => (
            <TabsLink key={index} to={`${url}${item.to}`}>
              {item.title}
            </TabsLink>
          ))}
        </TabsList>

        <TabsContent className="pt-2">
          <Routes>
            <Route
              index
              element={
                <Navigate to={`${url}${PAGE_PATH_ENVIRONMENTS}`} replace />
              }
            />
            <Route
              path={PAGE_PATH_ENVIRONMENTS}
              element={<ProjectEnvironments />}
            />
            <Route
              path={PAGE_PATH_SETTINGS}
              element={<ProjectSettings project={project} />}
            />
          </Routes>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
