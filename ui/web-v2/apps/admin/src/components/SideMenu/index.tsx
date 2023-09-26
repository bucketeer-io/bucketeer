import MUAccountCircleIcon from '@material-ui/icons/AccountCircle';
import MUBarChart from '@material-ui/icons/BarChart';
import MUFlagIcon from '@material-ui/icons/Flag';
import MUNotificationsIcon from '@material-ui/icons/Notifications';
import MUOpenInNew from '@material-ui/icons/OpenInNew';
import MUPeopleIcon from '@material-ui/icons/People';
import MUPermIdentityIcon from '@material-ui/icons/PermIdentity';
import MURemoveCircleIcon from '@material-ui/icons/RemoveCircle';
import MUSettingsIcon from '@material-ui/icons/Settings';
import MUSettingsApplications from '@material-ui/icons/SettingsApplications';
import MUShowChartIcon from '@material-ui/icons/ShowChart';
import MUSubjectIcon from '@material-ui/icons/Subject';
import MUToggleOnIcon from '@material-ui/icons/ToggleOn';
import MUVpnKeyIcon from '@material-ui/icons/VpnKey';
import { FC, ReactNode, memo, useCallback } from 'react';
import { useIntl } from 'react-intl';
import { useDispatch } from 'react-redux';
import { useHistory, Link, NavLink } from 'react-router-dom';

import {
  PAGE_PATH_ADMIN,
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_FEATURES,
  PAGE_PATH_GOALS,
  PAGE_PATH_APIKEYS,
  PAGE_PATH_USER_SEGMENTS,
  PAGE_PATH_USERS,
  PAGE_PATH_ACCOUNTS,
  PAGE_PATH_ROOT,
  PAGE_PATH_AUDIT_LOGS,
  PAGE_PATH_DOCUMENTATION,
  PAGE_PATH_SETTINGS,
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { clearToken } from '../../modules/auth';
import { clearMe, useCurrentEnvironment, useMe } from '../../modules/me';
import { AppDispatch } from '../../store';
import { EnvironmentSelect } from '../EnvironmentSelect';

export interface MenuItem {
  readonly messageComponent?: ReactNode;
  readonly path: string;
  readonly external: boolean;
  readonly target: string;
  readonly iconElement: JSX.Element;
}

interface Divider {
  messageComponent: null;
}

const createMenuItems = (
  isAdmin: boolean,
  environmentUrlCode: string
): Array<MenuItem | Divider> => {
  const items: Array<MenuItem | Divider> = [];
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.featureFlags)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_FEATURES}`,
    external: null,
    target: null,
    iconElement: <MUToggleOnIcon />,
  });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.goals)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_GOALS}`,
    external: null,
    target: null,
    iconElement: <MUFlagIcon />,
  });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.experiments)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_EXPERIMENTS}`,
    external: null,
    target: null,
    iconElement: <MUBarChart />,
  });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.userSegments)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_USER_SEGMENTS}`,
    external: null,
    target: null,
    iconElement: <MUPeopleIcon />,
  });
  // items.push({ TODO: User implementation
  //   messageComponent: <span>{intl.formatMessage(messages.sideMenu.user)}</span>,
  //   path: `/${environmentUrlCode}${PAGE_PATH_USERS}`,
  //   external: null,
  //   target: null,
  //   iconElement: <MUPermIdentityIcon />,
  // });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.auditLog)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_AUDIT_LOGS}`,
    external: null,
    target: null,
    iconElement: <MUNotificationsIcon />,
  });
  items.push({ messageComponent: null });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.accounts)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_ACCOUNTS}`,
    external: null,
    target: null,
    iconElement: <MUAccountCircleIcon />,
  });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.apiKeys)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_APIKEYS}`,
    external: null,
    target: null,
    iconElement: <MUVpnKeyIcon />,
  });
  items.push({
    messageComponent: (
      <span>
        {intl.formatMessage(messages.sideMenu.documentation)}
        <MUOpenInNew fontSize="small" />
      </span>
    ),
    path: `${PAGE_PATH_DOCUMENTATION}`,
    external: true,
    target: '_blank',
    iconElement: <MUSubjectIcon />,
  });
  items.push({
    messageComponent: (
      <span>{intl.formatMessage(messages.sideMenu.settings)}</span>
    ),
    path: `/${environmentUrlCode}${PAGE_PATH_SETTINGS}`,
    external: null,
    target: null,
    iconElement: <MUSettingsIcon />,
  });
  if (isAdmin) {
    items.push({
      messageComponent: (
        <span>{intl.formatMessage(messages.sideMenu.adminSettings)}</span>
      ),
      path: PAGE_PATH_ADMIN,
      external: null,
      target: null,
      iconElement: <MUSettingsApplications />,
    });
  }
  return items;
};

function isMenuItem(item: MenuItem | Divider): item is MenuItem {
  return item.messageComponent !== null;
}

interface Props {
  onClickNavLink: () => void;
}

export const SideMenu: FC<Props> = memo(({ onClickNavLink }) => {
  const dispatch = useDispatch<AppDispatch>();
  const history = useHistory();
  const { formatMessage: f } = useIntl();
  const me = useMe();
  const currentEnvironment = useCurrentEnvironment();

  const handleLogout = useCallback(async () => {
    dispatch(clearMe());
    dispatch(clearToken());
    history.push(PAGE_PATH_ROOT);
  }, [dispatch]);

  if (!me.isLogin) {
    return null;
  }
  return (
    <div className="flex flex-col w-full h-full bg-primary shadow-lg">
      <div className="p-4">
        <Link to={PAGE_PATH_ROOT}>
          <img src="/assets/logo.png" alt="Bucketer" />
        </Link>
      </div>
      <div className="w-full px-3 pb-1">
        <EnvironmentSelect />
      </div>
      <div className="flex-grow">
        {createMenuItems(me.isAdmin, currentEnvironment.urlCode).map((item, i) =>
          isMenuItem(item) ? (
            <div key={i} className="py-1">
              <SideMenuItem item={item} onClick={onClickNavLink} />
            </div>
          ) : (
            <div
              key={i}
              className="py-1 mb-2 shadow-md border-b border-purple-600"
            />
          )
        )}
      </div>
      <div className="bg-purple-600 h-12 items-center">
        <Link
          to={PAGE_PATH_ROOT}
          onClick={handleLogout}
          className="flex px-5 py-2.5"
        >
          <div className="flex justify-content items-center text-white ml-3">
            <div className="w-5 h-6 mr-2">{<MURemoveCircleIcon />}</div>
            <div>{f(messages.sideMenu.logout)}</div>
          </div>
        </Link>
      </div>
    </div>
  );
});

interface SideMenuItemProps {
  item: MenuItem;
  onClick: () => void;
}

const SideMenuItem: FC<SideMenuItemProps> = ({ item, onClick }) => {
  const history = useHistory();
  return item.external ? (
    <a href={item.path} target={item.target}>
      <div className="px-3">
        <div className="sidemenu-item flex px-5 py-2.5 rounded-md">
          <div className="flex justify-content items-center">
            <div className="w-5 h-6 mr-2">{item.iconElement}</div>
            <div>{item.messageComponent}</div>
          </div>
        </div>
      </div>
    </a>
  ) : (
    <div className="px-3">
      <NavLink
        to={item.path}
        target={item.target}
        onClick={onClick}
        className="sidemenu-item flex px-5 py-2.5 rounded-md"
      >
        <div className="flex justify-content items-center">
          <div className="w-5 h-6 mr-2">{item.iconElement}</div>
          <div>{item.messageComponent}</div>
        </div>
      </NavLink>
    </div>
  );
};
