import * as Popover from '@radix-ui/react-popover';
import {
  IconBuildingOutlined,
  IconChevronRightOutlined,
  IconLogoutOutlined,
  IconUserOutlined
} from '@icons';
import { AvatarImage } from 'components/avatar';
import MenuItemComponent from './menu-item';

const UserMenu = () => {
  const menuItems = [
    {
      label: 'User Profile',
      icon: IconUserOutlined,
      onClick: () => {}
    },
    {
      label: 'Polaris Edge',
      icon: IconBuildingOutlined,
      actIcon: IconChevronRightOutlined,
      onClick: () => {}
    },
    {
      label: 'Logout',
      icon: IconLogoutOutlined,
      onClick: () => {}
    }
  ];

  return (
    <Popover.Root>
      <Popover.Content align="start" className="border-none p-0">
        <div className="bg-primary-600 rounded-lg w-[200px] mb-2">
          {menuItems.map((item, index) => (
            <MenuItemComponent {...item} key={index} />
          ))}
        </div>
      </Popover.Content>
      <Popover.Trigger>
        <AvatarImage size="sm" image="./assets/avatars/primary.svg" />
      </Popover.Trigger>
    </Popover.Root>
  );
};

export default UserMenu;
