import { useState } from 'react';
import { useTranslation } from 'i18n';
import { Tabs, TabsList, TabsTrigger, TabsContent } from 'components/tabs';
import PageLayout from 'elements/page-layout';
import OrganizationProjects from './projects';

const PageContent = () => {
  const { t } = useTranslation(['common']);
  const [selectedTab] = useState('PROJECTS');

  return (
    <PageLayout.Content className="pt-4">
      <Tabs defaultValue={selectedTab} className="flex-1 flex h-full flex-col">
        <TabsList>
          <TabsTrigger value="PROJECTS">{t(`projects`)}</TabsTrigger>
          <TabsTrigger value="USERS">{t(`users`)}</TabsTrigger>
          <TabsTrigger value="SETTINGS">{t(`settings`)}</TabsTrigger>
        </TabsList>
        <TabsContent value={selectedTab} className="pt-2">
          <OrganizationProjects />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
