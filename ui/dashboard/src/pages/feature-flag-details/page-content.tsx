import { Navigate, Route, Routes, useParams } from 'react-router-dom';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURE_VARIATION,
  PAGE_PATH_FEATURE_SETTING,
  PAGE_PATH_FEATURE_EVALUATION,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURE_TRIGGER,
  PAGE_PATH_FEATURE_HISTORY,
  PAGE_PATH_FEATURE_CODE_REFS
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import NotFoundPage from 'pages/not-found';
import { Tabs, TabsList, TabsContent, TabsLink } from 'components/tabs-link';
import PageLayout from 'elements/page-layout';
import CodeReferencesPage from './code-refs';
import EvaluationPage from './evaluation';
import HistoryPage from './history';
import Operations from './operations';
import SettingsPage from './settings';
import TargetingPage from './targeting';
import TriggerPage from './trigger';
import { TabItem } from './types';
import Variation from './variation';

const PageContent = ({
  feature,
  refetchFeature
}: {
  feature: Feature;
  refetchFeature: () => void;
}) => {
  const { t } = useTranslation(['table', 'common']);
  const { flagId } = useParams();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);

  const url = `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${flagId}`;

  const featureFlagTabs: Array<TabItem> = [
    {
      title: t(`feature-flags.targeting`),
      to: PAGE_PATH_FEATURE_TARGETING
    },
    {
      title: t(`feature-flags.evaluations`),
      to: PAGE_PATH_FEATURE_EVALUATION
    },
    {
      title: t(`feature-flags.variations`),
      to: PAGE_PATH_FEATURE_VARIATION
    },
    {
      title: t(`feature-flags.operations`),
      to: PAGE_PATH_FEATURE_AUTOOPS
    },
    {
      title: t(`feature-flags.triggers`),
      to: PAGE_PATH_FEATURE_TRIGGER
    },
    {
      title: t(`feature-flags.code-references`),
      to: PAGE_PATH_FEATURE_CODE_REFS
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
    <PageLayout.Content className="pt-4">
      <Tabs>
        <TabsList className="px-6 w-fit min-w-full">
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
              path={`${PAGE_PATH_FEATURE_HISTORY}/*`}
              element={<HistoryPage feature={feature} />}
            />
            <Route
              path={PAGE_PATH_FEATURE_SETTING}
              element={<SettingsPage feature={feature} editable={editable} />}
            />
            <Route
              path={PAGE_PATH_FEATURE_VARIATION}
              element={<Variation feature={feature} editable={editable} />}
            />
            <Route
              path={`${PAGE_PATH_FEATURE_EVALUATION}/*`}
              element={<EvaluationPage feature={feature} />}
            />
            <Route
              path={`${PAGE_PATH_FEATURE_CODE_REFS}/*`}
              element={<CodeReferencesPage feature={feature} />}
            />
            <Route
              path={`${PAGE_PATH_FEATURE_TRIGGER}/*`}
              element={<TriggerPage feature={feature} editable={editable} />}
            />
            <Route
              path={`${PAGE_PATH_FEATURE_AUTOOPS}/*`}
              element={
                <Operations
                  feature={feature}
                  refetchFeature={refetchFeature}
                  editable={editable}
                />
              }
            />
            <Route
              path={PAGE_PATH_FEATURE_TARGETING}
              element={<TargetingPage feature={feature} editable={editable} />}
            />

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
