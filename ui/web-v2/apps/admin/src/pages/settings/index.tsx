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
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { useCurrentEnvironment } from '../../modules/me';
import { NotificationIndexPage } from '../notification';
import { PushIndexPage } from '../push';

export const SettingsIndexPage: FC = memo(() => {
  const { url } = useRouteMatch();
  const { formatMessage: f } = useIntl();
  const currentEnvironment = useCurrentEnvironment();

  return (
    <div className="">
      <div className="bg-white border-b border-gray-300">
        <div className="">
          <div className="py-5 px-10 text-gray-700">
            <p className="text-xl">{f(messages.settings.list.header.title)}</p>
            <p className="text-sm">
              {f(messages.settings.list.header.description)}
            </p>
          </div>
        </div>
        <div className="-mt-4 -ml-5 px-10 hidden sm:block">
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
                to={`${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_SETTINGS}${tab.to}`}
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
  ];
};
