import { useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router';
import logo from 'assets/logos/logo-white.svg';
import { useAuth, getCurrentEnvironment } from 'auth';
import * as ROUTING from 'constants/routing';
import { WALKTHROUGH_TARGETS } from 'constants/walkthrough';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import flatMapDeep from 'lodash/flatMapDeep';
import { setNavigationCollapsedStorage } from 'storage/navigation';
import { cn } from 'utils/style';
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import SectionMenu from './menu-section';
import MyProjects from './my-projects';
import NotificationBell from './notification-bell';
import SwitchOrganization from './switch-organization';
import UserMenu from './user-menu';
import logoIcon from '/img/bucketeer-logo-icon.png';

type NavigationProps = {
  onClickNavLink: () => void;
  isCollapsed: boolean;
  onToggleCollapsed: (value: boolean) => void;
};

const Navigation = ({
  onClickNavLink,
  isCollapsed,
  onToggleCollapsed
}: NavigationProps) => {
  const { t } = useTranslation(['common']);
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();

  const toggleCollapsed = () => {
    const next = !isCollapsed;
    setNavigationCollapsedStorage(next);
    onToggleCollapsed(next);
  };

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
          label: t(`navigation.notifications`),
          icon: IconSystem.IconNotifications,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_NOTIFICATION_FEED}`
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
          label: t(`navigation.feature-flags`),
          icon: IconSystem.IconSwitch,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_FEATURES}`,
          tourId: WALKTHROUGH_TARGETS.FEATURE_FLAGS_MENU
        },
        {
          label: t(`navigation.user-segment`),
          icon: IconSystem.IconUser,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_USER_SEGMENTS}`
        },
        {
          label: t(`navigation.insights`),
          icon: IconSystem.IconUsage,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_INSIGHTS}`
        },
        {
          label: t(`navigation.debugger`),
          icon: IconSystem.IconDebugger,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_DEBUGGER}`
        },
        {
          label: t(`navigation.audit-logs`),
          icon: IconSystem.IconLogs,
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_AUDIT_LOGS}`
        }
      ]
    },
    {
      title: t('experimentation'),
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
  )
    .map(item => item.href)
    .filter((href): href is string => !!href);

  const [isOpenSetting, onOpenSetting, onCloseSetting, setIsOpenSetting] =
    useToggleOpen(settingPaths.includes(pathname));

  useEffect(() => {
    setIsOpenSetting(settingPaths.includes(pathname));
  }, [pathname]);

  // Keep the sliding settings panel in sync with programmatic navigation
  // (e.g. the walkthrough moving between the main and settings areas).
  useEffect(() => {
    if (settingPaths.some(path => pathname.startsWith(path))) {
      onOpenSetting();
    } else {
      onCloseSetting();
    }
  }, [pathname]);

  const [isOpenSwitchOrg, onOpenSwitchOrg, onCloseSwitchOrg] =
    useToggleOpen(false);

  return (
    <div
      className={cn(
        'fixed h-screen bg-primary-500 z-50 py-8 transition-all duration-300 ease-in-out',
        isCollapsed ? 'w-[60px] px-2' : 'w-[248px] px-6'
      )}
    >
      <div className="flex flex-col size-full relative overflow-hidden">
        <div
          className={cn('group relative flex items-center', {
            'w-full justify-center': isCollapsed
          })}
        >
          <Link
            to={ROUTING.PAGE_PATH_ROOT}
            onClick={onCloseSetting}
            className={cn(
              'overflow-hidden',
              isCollapsed && 'flex-center w-full'
            )}
          >
            {isCollapsed ? (
              <img
                src={logoIcon}
                alt="Bucketeer"
                className="w-11 h-11 shrink-0"
              />
            ) : (
              <img src={logo} alt="Bucketer" />
            )}
          </Link>
          <button
            type="button"
            onClick={toggleCollapsed}
            className={cn(
              'flex-center size-6 rounded-md text-primary-50 shrink-0',
              'opacity-0 group-hover:opacity-100 focus-visible:opacity-100 transition-opacity',
              isCollapsed
                ? 'absolute inset-0 mx-auto w-11 h-11 bg-primary-500/80 hover:bg-primary-400 pointer-events-none group-hover:pointer-events-auto focus-visible:pointer-events-auto'
                : 'ml-auto hover:bg-primary-400'
            )}
            aria-label={t(
              isCollapsed ? `navigation.expand` : `navigation.collapse`
            )}
          >
            <Icon
              icon={IconSystem.IconChevronRight}
              color="primary-50"
              size="xs"
              className={cn({ 'rotate-180': !isCollapsed })}
            />
          </button>
        </div>

        <div className="flex flex-col flex-1 items-center pt-6">
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -right-[100%]',
              { 'right-0': isOpenSetting }
            )}
          >
            <Tooltip
              hidden={!isCollapsed}
              content={t(`navigation.back-to-main`)}
              side="right"
              trigger={
                <button
                  onClick={() => {
                    onCloseSetting();
                    navigate(`/${envUrlCode}${ROUTING.PAGE_PATH_FEATURES}`);
                  }}
                  aria-label={
                    isCollapsed ? t(`navigation.back-to-main`) : undefined
                  }
                  className={cn(
                    'flex items-center gap-x-2 text-primary-50 rounded-lg',
                    isCollapsed
                      ? 'justify-center w-full py-2 hover:bg-primary-400'
                      : 'px-3'
                  )}
                >
                  <Icon icon={IconSystem.IconBackspace} />
                  {!isCollapsed && <span>{t(`navigation.back-to-main`)}</span>}
                </button>
              }
            />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {settingMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
                isCollapsed={isCollapsed}
              />
            ))}
          </div>
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -left-[100%]',
              { 'left-0': !isOpenSetting }
            )}
          >
            {!isCollapsed && (
              <div className="px-3 opacity-80 uppercase typo-head-bold-tiny text-primary-50 mb-3">
                {t(`environment`)}
              </div>
            )}
            <MyProjects isCollapsed={isCollapsed} />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {mainMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
                onClickNavLink={onClickNavLink}
                isCollapsed={isCollapsed}
              />
            ))}
          </div>
        </div>

        <Divider className="mb-3 bg-primary-50 opacity-10" />

        <div
          className={cn('flex items-center justify-between', {
            'flex-col gap-y-3': isCollapsed
          })}
        >
          <UserMenu onOpenSwitchOrg={onOpenSwitchOrg} />
          <div
            className={cn('flex items-center justify-center gap-2', {
              'flex-col': isCollapsed
            })}
          >
            <NotificationBell envUrlCode={envUrlCode} />
            <Tooltip
              hidden={!isCollapsed}
              content={t(`settings`)}
              side="right"
              trigger={
                <button
                  type="button"
                  aria-label={isCollapsed ? t(`settings`) : undefined}
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
              }
            />
          </div>
        </div>
      </div>
      <SwitchOrganization
        isOpen={isOpenSwitchOrg}
        onCloseSwitchOrg={onCloseSwitchOrg}
        onCloseSetting={onCloseSetting}
        isCollapsed={isCollapsed}
      />
    </div>
  );
};

export default Navigation;
