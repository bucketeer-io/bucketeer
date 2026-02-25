import { Navigate, Route, Routes, useParams } from 'react-router-dom';
import {
  PAGE_PATH_MEMBERS,
  PAGE_PATH_NEW,
  PAGE_PATH_ORGANIZATIONS,
  PAGE_PATH_PROJECTS,
  PAGE_PATH_SETTINGS
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import { Tabs, TabsList, TabsContent, TabsLink } from 'components/tabs-link';
import PageLayout from 'elements/page-layout';
import OrganizationUsers from './members';
import OrganizationProjects from './projects';
import OrganizationSettings from './settings';
import { TabItem } from './types';

const PageContent = ({ organization }: { organization: Organization }) => {
  const { t } = useTranslation(['common']);
  const { organizationId } = useParams();
  const url = `${PAGE_PATH_ORGANIZATIONS}/${organizationId}`;

  const organizationTabs: Array<TabItem> = [
    {
      title: t(`projects`),
      to: PAGE_PATH_PROJECTS
    },
    {
      title: t(`members`),
      to: PAGE_PATH_MEMBERS
    },
    {
      title: t(`settings`),
      to: PAGE_PATH_SETTINGS
    }
  ];

  return (
    <PageLayout.Content className="pt-4">
      <Tabs>
        <TabsList className="px-3 sm:px-6">
          {organizationTabs.map((item, index) => (
            <TabsLink key={index} to={`${url}${item.to}`}>
              {item.title}
            </TabsLink>
          ))}
        </TabsList>

        <TabsContent className="pt-2 w-full">
          <Routes>
            <Route
              index
              element={<Navigate to={`${url}${PAGE_PATH_PROJECTS}`} replace />}
            />
            <Route
              path={PAGE_PATH_PROJECTS}
              element={<OrganizationProjects />}
            />
            <Route
              path={`${PAGE_PATH_PROJECTS}${PAGE_PATH_NEW}`}
              element={<OrganizationProjects />}
            />
            <Route
              path={`${PAGE_PATH_PROJECTS}/:projectId`}
              element={<OrganizationProjects />}
            />
            <Route path={PAGE_PATH_MEMBERS} element={<OrganizationUsers />} />
            <Route
              path={PAGE_PATH_SETTINGS}
              element={<OrganizationSettings organization={organization} />}
            />
          </Routes>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
