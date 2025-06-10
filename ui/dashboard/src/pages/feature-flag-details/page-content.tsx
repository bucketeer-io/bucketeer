// import { Navigate, Route, Routes, useParams } from 'react-router-dom';
import { useParams, Outlet, useLocation } from '@tanstack/react-router';
import { getCurrentEnvironment, useAuth } from 'auth';
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
import { Tabs, TabsList, TabsContent, TabsLink } from 'components/tabs-link';
import PageLayout from 'elements/page-layout';
import { TabItem } from './types';

const PageContent = () => {
  const { t } = useTranslation(['table', 'common']);
  const { featureId } = useParams({
    strict: false
  });
  const { pathname } = useLocation();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const url = `/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}`;

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
            <TabsLink
              key={index}
              to={`${url}${item.to}`}
              replace={true}
              className={pathname.includes(item.to) ? 'border-primary-500' : ''}
            >
              {item.title}
            </TabsLink>
          ))}
        </TabsList>

        <TabsContent className="pt-2">
          <Outlet />
        </TabsContent>
      </Tabs>
    </PageLayout.Content>
  );
};

export default PageContent;
