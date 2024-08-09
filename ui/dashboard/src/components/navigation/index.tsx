import { Link } from 'react-router-dom';
import logo from 'assets/logo.svg';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToggleOpen } from 'hooks';
import * as IconSystem from '@icons';
import Divider from 'components/divider';
import Icon from 'components/icon';
import SectionMenu from './menu-section';
import ProjectList from './project-list';
import UserMenu from './user-menu';

const Navigation = () => {
  const [isShowSetting, onShowSetting, onCloseSetting] = useToggleOpen(false);

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
                { label: 'Projects', icon: IconSystem.IconFolder, href: '/' },
                {
                  label: 'Organizations',
                  icon: IconSystem.IconBuilding,
                  href: '/'
                },
                {
                  label: 'Members',
                  icon: IconSystem.IconMember,
                  href: '/'
                },
                {
                  label: 'Usage',
                  icon: IconSystem.IconUsage,
                  href: '/'
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
                  href: '/'
                },
                {
                  label: 'API Keys',
                  icon: IconSystem.IconKey,
                  href: '/'
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
                  href: '/'
                },
                {
                  label: 'Feature Flags',
                  icon: IconSystem.IconSwitch,
                  href: '/'
                },
                {
                  label: 'User Segment',
                  icon: IconSystem.IconUser,
                  href: '/'
                },
                {
                  label: 'Debugger',
                  icon: IconSystem.IconDebugger,
                  href: '/'
                }
              ]}
            />
            <SectionMenu
              title={`Analysis`}
              className="mt-4"
              items={[
                {
                  label: 'Goals',
                  icon: IconSystem.IconNote,
                  href: '/'
                },
                {
                  label: 'Experiments',
                  icon: IconSystem.IconProton,
                  href: '/'
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
