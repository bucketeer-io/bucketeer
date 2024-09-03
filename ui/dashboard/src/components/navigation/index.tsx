import { Link } from 'react-router-dom';
import logo from 'assets/logos/logo-white.svg';
import { useAuth, useCurrentEnvironment } from 'auth';
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
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import SectionMenu from './menu-section';
import ProjectList from './project-list';
import UserMenu from './user-menu';

const Navigation = ({ onClickNavLink }: { onClickNavLink: () => void }) => {
  const [isShowSetting, onShowSetting, onCloseSetting] = useToggleOpen(false);
  const { consoleAccount } = useAuth();
  const currentEnvironment = useCurrentEnvironment(consoleAccount!);
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
              <span>{`Back to Main`}</span>
            </button>
            <Divider className="my-5 bg-primary-50 opacity-10" />
            <SectionMenu
              title="General"
              items={[
                {
                  label: 'Projects',
                  icon: IconSystem.IconFolder,
                  href: '/projects'
                },
                {
                  label: 'Organizations',
                  icon: IconSystem.IconBuilding,
                  href: '/organizations'
                },
                {
                  label: 'Members',
                  icon: IconSystem.IconMember,
                  href: '/members'
                },
                {
                  label: 'Usage',
                  icon: IconSystem.IconUsage,
                  href: '/usage'
                }
              ]}
            />
            <SectionMenu
              title="INTEGRATIONS"
              className="mt-5"
              items={[
                {
                  label: 'Integrations',
                  icon: IconSystem.IconIntegration,
                  href: '/integrations'
                },
                {
                  label: 'API Keys',
                  icon: IconSystem.IconKey,
                  href: '/api-keys'
                }
              ]}
            />
          </div>
        ) : (
          <div className="flex flex-col flex-1 pt-6">
            <div className="px-3 uppercase typo-head-bold-tiny text-primary-50 mb-3">
              {`Environment`}
            </div>
            <ProjectList />
            <Divider className="my-5 bg-primary-50 opacity-10" />
            <SectionMenu
              title={`Management`}
              items={[
                {
                  label: 'Audit Logs',
                  icon: IconSystem.IconLogs,
                  href: `/${environmentUrlCode}${PAGE_PATH_AUDIT_LOGS}`
                },
                {
                  label: 'Feature Flags',
                  icon: IconSystem.IconSwitch,
                  href: `/${environmentUrlCode}${PAGE_PATH_FEATURES}`
                },
                {
                  label: 'User Segment',
                  icon: IconSystem.IconUser,
                  href: `/${environmentUrlCode}${PAGE_PATH_USER_SEGMENTS}`
                },
                {
                  label: 'Debugger',
                  icon: IconSystem.IconDebugger,
                  href: `/${environmentUrlCode}${PAGE_PATH_DEBUGGER}`
                }
              ]}
              onClickNavLink={onClickNavLink}
            />
            <SectionMenu
              title={`Analysis`}
              className="mt-4"
              items={[
                {
                  label: 'Goals',
                  icon: IconSystem.IconNote,
                  href: `/${environmentUrlCode}${PAGE_PATH_GOALS}`
                },
                {
                  label: 'Experiments',
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
