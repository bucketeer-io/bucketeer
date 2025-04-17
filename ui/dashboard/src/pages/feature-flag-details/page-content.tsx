import { Navigate, Route, Routes, useParams } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURE_VARIATION,
  PAGE_PATH_FEATURE_SETTING,
  PAGE_PATH_FEATURE_EXPERIMENTS,
  PAGE_PATH_FEATURE_EVALUATION,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURE_TRIGGER,
  PAGE_PATH_FEATURE_HISTORY
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import NotFoundPage from 'pages/not-found';
import { Tabs, TabsList, TabsContent, TabsLink } from 'components/tabs-link';
import PageLayout from 'elements/page-layout';
import { TabItem } from './types';
import Variation from './variation';

const PageContent = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['table', 'common']);
  const { flagId } = useParams();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const url = `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${flagId}`;

  const featureFlagTabs: Array<TabItem> = [
    {
      title: t(`feature-flags.targeting`),
      to: PAGE_PATH_FEATURE_TARGETING
    },
    {
      title: t(`feature-flags.variation`),
      to: PAGE_PATH_FEATURE_VARIATION
    },
    {
      title: t(`feature-flags.evaluation`),
      to: PAGE_PATH_FEATURE_EVALUATION
    },
    {
      title: t(`feature-flags.operations`),
      to: PAGE_PATH_FEATURE_AUTOOPS
    },
    {
      title: t(`feature-flags.trigger`),
      to: PAGE_PATH_FEATURE_TRIGGER
    },
    {
      title: t(`feature-flags.experiments`),
      to: PAGE_PATH_FEATURE_EXPERIMENTS
    },
    {
      title: t(`feature-flags.history`),
      to: PAGE_PATH_FEATURE_HISTORY
    },
    {
      title: t(`common:settings`),
      to: PAGE_PATH_FEATURE_SETTING
    }
  ];

  return (
    <PageLayout.Content className="p-6 pt-4">
      <Tabs>
        <TabsList>
          {featureFlagTabs.map((item, index) => (
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
                <Navigate to={`${url}${PAGE_PATH_FEATURE_TARGETING}`} replace />
              }
            />
            <Route
              path={PAGE_PATH_FEATURE_VARIATION}
              element={<Variation feature={feature} />}
            />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
