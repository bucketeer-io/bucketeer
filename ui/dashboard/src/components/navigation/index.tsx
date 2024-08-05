import { Link } from 'react-router-dom';
import { PAGE_PATH_ROOT } from 'constants/routing';
import { useToggleOpen } from 'hooks';
import {
  IconChevronRightOutlined,
  IconDebuggerOutlined,
  IconUsageOutlined,
  IconFolderOutlined,
  IconIntegrationOutlined,
  IconKeyOutlined,
  IconLogsOutlined,
  IconMemberOutlined,
  IconNoteOutlined,
  IconProtonOutlined,
  IconSwitchOutlined,
  IconUserOutlined,
  IconSettingOutlined,
  IconBackspaceOutlined,
  IconBuildingOutlined
} from '@icons';
import { AvatarImage } from 'components/avatar';
import Divider from 'components/divider';
import Icon from 'components/icon';
import SectionMenu from './menu-section';
import NavigationBottomAction from './navigation-bottom-action';
import NavigationUserMenu from './navigation-user-menu';

// import ProjectList from './project-list';

const Navigation = () => {
  const [isShowSetting, onShowSetting, onCloseSetting] = useToggleOpen(false);

  return (
    <div className="fixed h-screen w-[248px] bg-primary-500 py-8 px-6">
      <div className="flex flex-col h-full">
        <Link to={PAGE_PATH_ROOT}>
          <img src="./assets/logo.svg" alt="Bucketer" />
        </Link>

        {isShowSetting ? (
          <div className="flex flex-col flex-1 pt-6">
            <button
              onClick={onCloseSetting}
              className="flex items-center gap-x-2 text-primary-50"
            >
              <Icon icon={IconBackspaceOutlined} />
              <span>{`Back to Main`}</span>
            </button>

            <Divider className="my-5 bg-primary-50 opacity-10" />
            <SectionMenu
              title="General"
              items={[
                { label: 'Projects', icon: IconFolderOutlined, href: '/' },
                {
                  label: 'Organizations',
                  icon: IconBuildingOutlined,
                  href: '/'
                },
                {
                  label: 'Members',
                  icon: IconMemberOutlined,
                  href: '/'
                },
                {
                  label: 'Usage',
                  icon: IconUsageOutlined,
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
                  icon: IconIntegrationOutlined,
                  href: '/'
                },
                {
                  label: 'API Keys',
                  icon: IconKeyOutlined,
                  href: '/'
                }
              ]}
            />
          </div>
        ) : (
          <div className="flex flex-col flex-1 pt-6">
            <SectionMenu
              title={`Environment`}
              items={[
                {
                  icon: IconFolderOutlined,
                  actIcon: IconChevronRightOutlined,
                  label: `Abematv`,
                  onClick: () => {}
                }
              ]}
            />
            <Divider className="my-5 bg-primary-50 opacity-10" />

            <SectionMenu
              title={`Management`}
              items={[
                {
                  label: 'Audit Logs',
                  icon: IconLogsOutlined,
                  href: '/'
                },
                {
                  label: 'Feature Flags',
                  icon: IconSwitchOutlined,
                  href: '/'
                },
                {
                  label: 'User Segment',
                  icon: IconUserOutlined,
                  href: '/'
                },
                {
                  label: 'Debugger',
                  icon: IconDebuggerOutlined,
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
                  icon: IconNoteOutlined,
                  href: '/'
                },
                {
                  label: 'Experiments',
                  icon: IconProtonOutlined,
                  href: '/'
                }
              ]}
            />
          </div>
        )}
        <Divider className="mb-3 bg-primary-50 opacity-10" />

        <div className="flex items-center justify-between">
          <AvatarImage size="sm" image="./assets/avatars/primary.svg" />
          <button onClick={onShowSetting}>
            <Icon icon={IconSettingOutlined} />
          </button>
        </div>
      </div>
    </div>
  );
};

Navigation.BottomAction = NavigationBottomAction;
Navigation.UserMenu = NavigationUserMenu;

export default Navigation;
