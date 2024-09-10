import * as Popover from '@radix-ui/react-popover';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useAuth } from 'auth';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import MenuItemComponent from './menu-item';

const UserMenu = () => {
  const { t } = useTranslation(['common']);
  const { logout, myOrganizations, consoleAccount } = useAuth();

  const menuItems = compact([
    {
      label: t(`navigation.user-profile`),
      icon: IconUser,
      onClick: () => {}
    },
    myOrganizations.length > 1 && {
      label: consoleAccount?.organization?.name || '',
      icon: IconBuilding,
      actIcon: IconChevronRight,
      onClick: () => {}
    },
    {
      label: t(`navigation.logout`),
      icon: IconLogout,
      onClick: logout
    }
  ]);

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
