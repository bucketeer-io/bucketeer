import { useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router';
import logo from 'assets/logos/logo-white.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import * as ROUTING from 'constants/routing';
import { useScreen, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import flatMapDeep from 'lodash/flatMapDeep';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { cn } from 'utils/style';
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import SectionMenu from './menu-section';
import MyProjects from './my-projects';
import SwitchOrganization from './switch-organization';
import UserMenu from './user-menu';

const LogoIcon = () => (
  <span className="mt-4 block w-full bg-primary-600 p-[5px] rounded-md text-2xl font-bold text-white text-center">
    B
  </span>
);

const Navigation = ({
  onClickNavLink,
  forceExpanded = false
}: {
  onClickNavLink: () => void;
  forceExpanded?: boolean;
}) => {
  const { t } = useTranslation(['common']);
  const { pathname } = useLocation();
  const navigate = useNavigate();
  const { consoleAccount } = useAuth();

  const { fromTabletScreen, isMobile } = useScreen();
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
          href: `/${envUrlCode}${ROUTING.PAGE_PATH_FEATURES}`
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
  ).map(item => item.href);

  const [isOpenSetting, onOpenSetting, onCloseSetting] = useToggleOpen(
    settingPaths.includes(pathname)
  );

  const [isOpenSwitchOrg, onOpenSwitchOrg, onCloseSwitchOrg] =
    useToggleOpen(false);

  const [_isExpanded, setIsExpanded, setIsCloseExpand] = useToggleOpen(false);
  const isExpanded = forceExpanded || _isExpanded;

  useEffect(() => {
    if (forceExpanded) return;
    if (fromTabletScreen) {
      setIsExpanded();
    } else if (isMobile) {
      setIsCloseExpand();
    }
    // Run only on mount to set the initial collapsed/expanded state.
    // Subsequent resizes should not override a user's manual toggle.
  }, []);
  return (
    <div
      className={cn(
        'fixed h-screen bg-primary-500 z-50 transition-[width] duration-300 ease-in-out',
        isOpenSwitchOrg ? 'overflow-visible' : 'overflow-hidden',
        !fromTabletScreen && isExpanded
          ? 'w-[248px] px-3 sm:px-6 py-4'
          : isMobile
            ? 'w-[248px] px-3 sm:px-6 py-4'
            : 'w-[60px] md:w-[248px] py-4 px-2 md:px-6 md:py-8'
      )}
    >
      <div
        className={cn(
          'flex flex-col size-full relative pt-4 md:pt-0',
          isOpenSwitchOrg ? 'overflow-visible' : 'overflow-hidden'
        )}
      >
        {!fromTabletScreen && !forceExpanded && (
          <button
            onClick={() =>
              !_isExpanded ? setIsExpanded() : setIsCloseExpand()
            }
            className={cn(
              'hidden sm:flex items-center justify-center hover:bg-primary-300 hover:rounded-full w-6 h-6',
              isExpanded
                ? 'absolute right-0 top-4'
                : 'absolute left-1/2 -translate-x-1/2 top-0'
            )}
          >
            <Icon
              icon={isExpanded ? ChevronLeft : ChevronRight}
              size="sm"
              color="primary-50"
            />
          </button>
        )}
        <Link to={ROUTING.PAGE_PATH_ROOT} onClick={onCloseSetting}>
          {!fromTabletScreen ? (
            isExpanded ? (
              <img src={logo} alt="Bucketer" />
            ) : (
              <LogoIcon />
            )
          ) : (
            <img src={logo} alt="Bucketer" />
          )}
        </Link>

        <div className="flex flex-col flex-1 items-center pt-6">
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -right-[100%]',
              { 'right-0': isOpenSetting }
            )}
          >
            <div
              className={cn(
                'w-full flex',
                !fromTabletScreen
                  ? isExpanded
                    ? 'justify-start'
                    : 'justify-center'
                  : 'justify-start'
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
                <span
                  className={cn(
                    isExpanded ? 'inline-block' : 'hidden md:inline-block'
                  )}
                >
                  {t(`navigation.back-to-main`)}
                </span>
              </button>
            </div>
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {settingMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                isExpanded={isExpanded}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
                onClickNavLink={onClickNavLink}
              />
            ))}
          </div>
          <div
            className={cn(
              'w-full absolute ease-in-out transition-all duration-500 -left-[100%]',
              { 'left-0': !isOpenSetting }
            )}
          >
            <div
              className={cn(
                'px-3 opacity-80 uppercase typo-head-bold-tiny text-primary-50 mb-3',
                isExpanded ? 'block' : 'hidden md:block'
              )}
            >
              {t(`environment`)}
            </div>
            <MyProjects isExpanded={isExpanded} />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            {mainMenuSections.map((item, index) => (
              <SectionMenu
                key={index}
                isExpanded={isExpanded}
                className="first:mt-0 mt-4"
                title={item.title}
                items={item.menus}
                onClickNavLink={onClickNavLink}
              />
            ))}
          </div>
        </div>

        <Divider className="mb-3 bg-primary-50 opacity-10" />

        <div
          className={cn(
            'flex gap-5 items-center justify-between',
            isExpanded ? 'flex-row' : 'flex-col md:flex-row'
          )}
        >
          <UserMenu onOpenSwitchOrg={onOpenSwitchOrg} />
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
      {!isMobile && isOpenSwitchOrg ? (
        <SwitchOrganization
          isExpanded={isExpanded}
          isOpen={isOpenSwitchOrg}
          onCloseSwitchOrg={onCloseSwitchOrg}
          onCloseSetting={onCloseSetting}
        />
      ) : (
        <DialogModal
          className="w-[290px]"
          title=""
          isOpen={isOpenSwitchOrg}
          onClose={onCloseSwitchOrg}
          overlayCls="!z-[500]"
        >
          <SwitchOrganization
            isExpanded={isExpanded}
            isOpen={isOpenSwitchOrg}
            onCloseSwitchOrg={onCloseSwitchOrg}
            onCloseSetting={onCloseSetting}
          />
        </DialogModal>
      )}
    </div>
  );
};

export default Navigation;
