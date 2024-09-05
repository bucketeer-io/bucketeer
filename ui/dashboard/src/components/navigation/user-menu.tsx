import * as Popover from '@radix-ui/react-popover';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import MenuItemComponent from './menu-item';

const UserMenu = () => {
  const { t } = useTranslation(['common']);
  const { logout } = useAuth();

  const menuItems = [
    {
      label: t(`navigation.user-profile`),
      icon: IconUser,
      onClick: () => {}
    },
    {
      label: t(`navigation.polaris-edge`),
      icon: IconBuilding,
      actIcon: IconChevronRight,
      onClick: () => {}
    },
    {
      label: t(`navigation.logout`),
      icon: IconLogout,
      onClick: logout
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
