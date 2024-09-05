import { Link } from 'react-router-dom';
import logo from 'assets/logos/logo-white.svg';
import { useAuth, getCurrentEnvironment } from 'auth';
import {
  PAGE_PATH_AUDIT_LOGS,
  PAGE_PATH_FEATURES,
  PAGE_PATH_ROOT,
  PAGE_PATH_USER_SEGMENTS,
  PAGE_PATH_DEBUGGER,
  PAGE_PATH_GOALS,
  PAGE_PATH_EXPERIMENTS
} from 'constants/routing';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import SectionMenu from './menu-section';
import MyProjects from './my-projects';
import UserMenu from './user-menu';

const Navigation = ({ onClickNavLink }: { onClickNavLink: () => void }) => {
  const { t } = useTranslation(['common']);
  const [isShowSetting, onShowSetting, onCloseSetting] = useToggleOpen(false);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const environmentUrlCode = currentEnvironment.urlCode;

  return (
    <div className="fixed h-screen w-[248px] bg-primary-500 py-8 px-6">
      <div className="flex flex-col h-full">
        <Link to={PAGE_PATH_ROOT}>
          <img src={logo} alt="Bucketer" />
        </Link>

        {isShowSetting ? (
          <div className="flex flex-col flex-1 pt-6">
            <button
              onClick={onCloseSetting}
              className="flex items-center gap-x-2 text-primary-50"
            >
              <Icon icon={IconSystem.IconBackspace} />
              <span>{t(`navigation.back-to-main`)}</span>
            </button>
            <Divider className="my-5 bg-primary-50 opacity-10" />
            <SectionMenu
              title={t('general')}
              items={[
                {
                  label: t(`organizations`),
                  icon: IconSystem.IconBuilding,
                  href: '/organizations'
                },
                {
                  label: t(`settings`),
                  icon: IconSystem.IconSetting,
                  href: '/settings'
                },
                {
                  label: t(`projects`),
                  icon: IconSystem.IconFolder,
                  href: '/projects'
                },

                {
                  label: t(`usage`),
                  icon: IconSystem.IconUsage,
                  href: '/usage'
                }
              ]}
            />
            <SectionMenu
              title={t(`access`)}
              className="mt-5"
              items={[
                {
                  label: t(`members`),
                  icon: IconSystem.IconMember,
                  href: '/members'
                },
                {
                  label: t(`API-keys`),
                  icon: IconSystem.IconKey,
                  href: '/api-keys'
                }
              ]}
            />
            <SectionMenu
              title={t(`integrations`)}
              className="mt-5"
              items={[
                {
                  label: `Slack`,
                  icon: IconSystem.IconSlack,
                  href: '/integrations/slack'
                },
                {
                  label: `FCM`,
                  icon: IconSystem.IconFCM,
                  href: '/integrations/fcm'
                }
              ]}
            />
          </div>
        ) : (
          <div className="flex flex-col flex-1 pt-6">
            <div className="px-3 uppercase typo-head-bold-tiny text-primary-50 mb-3">
              {t(`environment`)}
            </div>
            <MyProjects />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            <SectionMenu
              title={t(`management`)}
              items={[
                {
                  label: t(`navigation.audit-logs`),
                  icon: IconSystem.IconLogs,
                  href: `/${environmentUrlCode}${PAGE_PATH_AUDIT_LOGS}`
                },
                {
                  label: t(`navigation.feature-flags`),
                  icon: IconSystem.IconSwitch,
                  href: `/${environmentUrlCode}${PAGE_PATH_FEATURES}`
                },
                {
                  label: t(`navigation.user-segment`),
                  icon: IconSystem.IconUser,
                  href: `/${environmentUrlCode}${PAGE_PATH_USER_SEGMENTS}`
                },
                {
                  label: t(`navigation.debugger`),
                  icon: IconSystem.IconDebugger,
                  href: `/${environmentUrlCode}${PAGE_PATH_DEBUGGER}`
                }
              ]}
              onClickNavLink={onClickNavLink}
            />
            <SectionMenu
              title={t(`analysis`)}
              className="mt-4"
              items={[
                {
                  label: t(`navigation.goals`),
                  icon: IconSystem.IconNote,
                  href: `/${environmentUrlCode}${PAGE_PATH_GOALS}`
                },
                {
                  label: t(`navigation.experiments`),
                  icon: IconSystem.IconProton,
                  href: `/${environmentUrlCode}${PAGE_PATH_EXPERIMENTS}`
                }
              ]}
            />
          </div>
        )}
        <Divider className="mb-3 bg-primary-50 opacity-10" />

        <div className="flex items-center justify-between">
          <UserMenu />
          <button onClick={onShowSetting}>
            <Icon icon={IconSystem.IconSetting} color="primary-50" />
          </button>
        </div>
      </div>
    </div>
  );
};

export default Navigation;
