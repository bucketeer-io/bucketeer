import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';
import { PAGE_PATH_EXPERIMENTS } from 'constants/routing';
import { Experiment } from '@types';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';

export type ExperimentDetailsTab = 'results' | 'settings';

const PageContent = ({ experiment }: { experiment: Experiment }) => {
  const { t } = useTranslation(['common']);
  const { tab: currentTab, envUrlCode, experimentId } = useParams();
  const navigate = useNavigate();

  return (
    <PageLayout.Content>
      <Tabs
        className="flex-1 flex h-full flex-col"
        value={currentTab}
        onValueChange={value =>
          navigate(
            `/${envUrlCode}${PAGE_PATH_EXPERIMENTS}/${experimentId}/${value}`
          )
        }
      >
        <TabsList>
          <TabsTrigger value="results">{t(`results`)}</TabsTrigger>
          <TabsTrigger value="settings">{t(`settings`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={currentTab as ExperimentDetailsTab}>
          <CollectionLoader
            currentTab={currentTab as ExperimentDetailsTab}
            experiment={experiment}
          />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
