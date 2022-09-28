import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import {
  NavLink,
  Redirect,
  Route,
  Switch,
  useRouteMatch,
} from 'react-router-dom';

import {
  PAGE_PATH_ACCOUNTS,
  PAGE_PATH_ADMIN,
  PAGE_PATH_AUDIT_LOGS,
  PAGE_PATH_ENVIRONMENTS,
  PAGE_PATH_NOTIFICATIONS,
  PAGE_PATH_PROJECTS,
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AdminAccountIndexPage } from '../admin/account';
import { AdminAuditLogIndexPage } from '../admin/auditLog';
import { AdminNotificationIndexPage } from '../admin/notification';
import { AdminProjectIndexPage } from '../admin/projects';

import { AdminEnvironmentIndexPage } from './environment';

export const AdminIndexPage: FC = memo(() => {
  const { url } = useRouteMatch();
  const { formatMessage: f } = useIntl();
  return (
    <div className="">
      <div className="bg-white border-b border-gray-300">
        <div className="">
          <div className="py-5 px-10 text-gray-700">
            <p className="text-xl">
              {f(messages.adminSettings.list.header.title)}
            </p>
            <p className="text-sm">
              {f(messages.adminSettings.list.header.description)}
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
                to={`${PAGE_PATH_ADMIN}${tab.to}`}
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
            component={() => <Redirect to={`${url}${PAGE_PATH_PROJECTS}`} />}
          />
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_PROJECTS}`,
              `${url}${PAGE_PATH_PROJECTS}/:projectId`,
            ]}
          >
            <AdminProjectIndexPage />
          </Route>
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_ENVIRONMENTS}`,
              `${url}${PAGE_PATH_ENVIRONMENTS}/:environmentId`,
            ]}
          >
            <AdminEnvironmentIndexPage />
          </Route>
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_ACCOUNTS}`,
              `${url}${PAGE_PATH_ACCOUNTS}/:accountId`,
            ]}
          >
            <AdminAccountIndexPage />
          </Route>
          <Route
            exact
            path={[
              `${url}${PAGE_PATH_NOTIFICATIONS}`,
              `${url}${PAGE_PATH_NOTIFICATIONS}/:notificationId`,
            ]}
          >
            <AdminNotificationIndexPage />
          </Route>
          <Route exact path={`${url}${PAGE_PATH_AUDIT_LOGS}`}>
            <AdminAuditLogIndexPage />
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
      message: intl.formatMessage(messages.adminSettings.tab.projects),
      to: PAGE_PATH_PROJECTS,
    },
    {
      message: intl.formatMessage(messages.adminSettings.tab.environments),
      to: PAGE_PATH_ENVIRONMENTS,
    },
    {
      message: intl.formatMessage(messages.adminSettings.tab.account),
      to: PAGE_PATH_ACCOUNTS,
    },
    {
      message: intl.formatMessage(messages.adminSettings.tab.notifications),
      to: PAGE_PATH_NOTIFICATIONS,
    },
    {
      message: intl.formatMessage(messages.adminSettings.tab.auditLogs),
      to: PAGE_PATH_AUDIT_LOGS,
    },
  ];
};
