import { Link, useLocation, useNavigate } from 'react-router-dom';
import logo from 'assets/logos/logo-white.svg';
import { useAuth, getCurrentEnvironment } from 'auth';
import * as ROUTING from 'constants/routing';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import flatMapDeep from 'lodash/flatMapDeep';
import { cn } from 'utils/style';
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import SectionMenu from './menu-section';
import MyProjects from './my-projects';
import SwitchOrganization from './switch-organization';
import UserMenu from './user-menu';

const Navigation = ({ onClickNavLink }: { onClickNavLink: () => void }) => {
  const { t } = useTranslation(['common']);
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();

  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const envUrlCode = currentEnvironment.urlCode;

  const settingMenuSections = [
    {
      title: t('general'),
      menus: compact([
        consoleAccount?.isSystemAdmin && {
          label: t(`organizations`),
          icon: IconSystem.IconBuilding,
          href: ROUTING.PAGE_PATH_ORGANIZATIONS
        },
        {
          label: t(`settings`),
          icon: IconSystem.IconSetting,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_SETTINGS}`
        },
        {
          label: t(`projects`),
          icon: IconSystem.IconFolder,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_PROJECTS}`
        },
        {
          label: t(`usage`),
          icon: IconSystem.IconUsage,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_USAGE}`
        }
      ])
    },
    {
      title: t('access'),
      menus: [
        {
          label: t(`members`),
          icon: IconSystem.IconMember,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_MEMBERS}`
        },
        {
          label: t(`api-keys`),
          icon: IconSystem.IconKey,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_APIKEYS}`
        }
      ]
    },
    {
      title: t('integrations'),
      menus: [
        {
          label: `Slack`,
          icon: IconSystem.IconSlack,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_NOTIFICATIONS}`
        },
        {
          label: `FCM`,
          icon: IconSystem.IconFCM,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_PUSHES}`
        }
      ]
    }
  ];

  const mainMenuSections = [
    {
      title: t('management'),
      menus: [
        {
          label: t(`navigation.audit-logs`),
          icon: IconSystem.IconLogs,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_AUDIT_LOGS}`
        },
        {
          label: t(`navigation.feature-flags`),
          icon: IconSystem.IconSwitch,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_FEATURES}`
        },
        {
          label: t(`navigation.user-segment`),
          icon: IconSystem.IconUser,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_USER_SEGMENTS}`
        },
        {
          label: t(`navigation.debugger`),
          icon: IconSystem.IconDebugger,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_DEBUGGER}`
        }
      ]
    },
    {
      title: t('analysis'),
      menus: [
        {
          label: t(`navigation.goals`),
          icon: IconSystem.IconNote,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_GOALS}`
        },
        {
          label: t(`navigation.experiments`),
          icon: IconSystem.IconProton,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_EXPERIMENTS}`
        }
      ]
    }
  ];

  const settingPaths = flatMapDeep(
    settingMenuSections.map(section => section.menus)
  ).map(item => item.href);

  const [isOpenSetting, onOpenSetting, onCloseSetting] = useToggleOpen(
    settingPaths.includes(pathname)
  );

  const [isOpenSwitchOrg, onOpenSwitchOrg, onCloseSwitchOrg] =
    useToggleOpen(false);

  return (
    <div className="fixed h-screen w-[248px] bg-primary-500 z-50 py-8 px-6">
      <div className="flex flex-col size-full relative overflow-hidden">
        <Link to={ROUTING.PAGE_PATH_ROOT}>
          <img src={logo} alt="Bucketer" />
        </Link>

        <div className="flex flex-col flex-1 items-center pt-6">
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -right-[100%]',
              { 'right-0': isOpenSetting }
            )}
          >
            <button
              onClick={() => {
                onCloseSetting();
                navigate(`/${envUrlCode}${ROUTING.PAGE_PATH_FEATURES}`);
              }}
              className="flex items-center gap-x-2 text-primary-50"
            >
              <Icon icon={IconSystem.IconBackspace} />
              <span>{t(`navigation.back-to-main`)}</span>
            </button>
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {settingMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
              />
            ))}
          </div>
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -left-[100%]',
              { 'left-0': !isOpenSetting }
            )}
          >
            <div className="px-3 opacity-80 uppercase typo-head-bold-tiny text-primary-50 mb-3">
              {t(`environment`)}
            </div>
            <MyProjects />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {mainMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
                onClickNavLink={onClickNavLink}
              />
            ))}
          </div>
        </div>

        <Divider className="mb-3 bg-primary-50 opacity-10" />

        <div className="flex items-center justify-between">
          <UserMenu
            onOpenSwitchOrg={
              isOpenSwitchOrg ? onCloseSwitchOrg : onOpenSwitchOrg
            }
          />
          <button
            type="button"
            onClick={() => {
              onOpenSetting();
              if (consoleAccount?.isSystemAdmin) {
                navigate(ROUTING.PAGE_PATH_ORGANIZATIONS);
              } else {
                navigate(`/${envUrlCode}${ROUTING.PAGE_PATH_SETTINGS}`);
              }
            }}
          >
            <Icon icon={IconSystem.IconSetting} color="primary-50" />
          </button>
        </div>
      </div>
      <SwitchOrganization
        isOpen={isOpenSwitchOrg}
        onCloseSwitchOrg={onCloseSwitchOrg}
      />
    </div>
  );
};

export default Navigation;
