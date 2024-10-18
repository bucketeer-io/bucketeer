import { useState } from 'react';
import { useTranslation } from 'i18n';
import { Organization } from '@types';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import PageLayout from 'elements/page-layout';
import OrganizationProjects from './projects';
import OrganizationSettings from './settings';
import OrganizationUsers from './users';

const PageContent = ({ organization }: { organization: Organization }) => {
  const { t } = useTranslation(['common']);
  const [selectedTab, setSelectedTab] = useState('PROJECTS');

  const getTabContent = () => {
    switch (selectedTab) {
      case 'PROJECTS':
        return <OrganizationProjects />;

      case 'USERS':
        return <OrganizationUsers />;

      case 'SETTINGS':
        return <OrganizationSettings organization={organization} />;

      default:
        return null;
    }
  };

  return (
    <PageLayout.Content className="pt-4">
      <Tabs
        defaultValue={selectedTab}
        onValueChange={setSelectedTab}
        className="flex-1 flex h-full flex-col"
      >
        <TabsList>
          <TabsTrigger value="PROJECTS">{t(`projects`)}</TabsTrigger>
          <TabsTrigger value="USERS">{t(`users`)}</TabsTrigger>
          <TabsTrigger value="SETTINGS">{t(`settings`)}</TabsTrigger>
        </TabsList>
        <TabsContent value={selectedTab} className="pt-2">
          {getTabContent()}
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
