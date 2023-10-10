import { Header } from '@/components/Header';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import {
  NavLink,
  Route,
  Switch,
  Redirect,
  useRouteMatch,
} from 'react-router-dom';

import {
  PAGE_PATH_SETTINGS,
  PAGE_PATH_PUSHES,
  PAGE_PATH_NOTIFICATIONS,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT,
  PAGE_PATH_WEBHOOKS,
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useCurrentEnvironment } from '../../modules/me';
import { NotificationIndexPage } from '../notification';
import { PushIndexPage } from '../push';
import { WebhookIndexPage } from '../webhook';

export const SettingsIndexPage: FC = memo(() => {
  const { url } = useRouteMatch();
  const { formatMessage: f } = useIntl();
  const currentEnvironment = useCurrentEnvironment();

  return (
    <div className="">
      <div className="bg-white border-b border-gray-300">
        <div className="">
          <Header
            title={f(messages.settings.list.header.title)}
            description={f(messages.settings.list.header.description)}
          />
        </div>
        <div className="px-10 hidden sm:block">
          <nav className="-mb-px flex" aria-label="Tabs">
            {createTabs().map((tab, idx) => (
              <NavLink
                key={idx}
                className="
                    tab-item
                    border-transparent
                    text-gray-500
                    hover:text-gray-700
                    whitespace-nowrap py-4 px-5 border-b-2
                    font-medium text-sm"
                to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_SETTINGS}${tab.to}`}
              >
                {tab.message}
              </NavLink>
            ))}
          </nav>
        </div>
      </div>
      <div className="my-10">
        <Switch>
          <Route
            exact
            path={`${url}`}
            component={() => <Redirect to={`${url}${PAGE_PATH_PUSHES}`} />}
          />
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_PUSHES}`,
              `${url}${PAGE_PATH_PUSHES}/:pushId`,
            ]}
          >
            <PushIndexPage />
          </Route>
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_NOTIFICATIONS}`,
              `${url}${PAGE_PATH_NOTIFICATIONS}/:notificationId`,
            ]}
          >
            <NotificationIndexPage />
          </Route>
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_WEBHOOKS}`,
              `${url}${PAGE_PATH_WEBHOOKS}/:webhookId`,
            ]}
          >
            <WebhookIndexPage />
          </Route>
        </Switch>
      </div>
    </div>
  );
});

export interface TabItem {
  readonly message: string;
  readonly to: string;
}

const createTabs = (): Array<TabItem> => {
  return [
    {
      message: intl.formatMessage(messages.settings.tab.pushes),
      to: PAGE_PATH_PUSHES,
    },
    {
      message: intl.formatMessage(messages.settings.tab.notifications),
      to: PAGE_PATH_NOTIFICATIONS,
    },
    {
      message: intl.formatMessage(messages.settings.tab.webhooks),
      to: PAGE_PATH_WEBHOOKS,
    },
  ];
};
