import * as Popover from '@radix-ui/react-popover';
import primaryAvatar from 'assets/avatars/primary.svg';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import MenuItemComponent from './menu-item';

const UserMenu = () => {
  const menuItems = [
    {
      label: 'User Profile',
      icon: IconUser,
      onClick: () => {}
    },
    {
      label: 'Polaris Edge',
      icon: IconBuilding,
      actIcon: IconChevronRight,
      onClick: () => {}
    },
    {
      label: 'Logout',
      icon: IconLogout,
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
        <AvatarImage size="sm" image={primaryAvatar} />
      </Popover.Trigger>
    </Popover.Root>
  );
};

export default UserMenu;
