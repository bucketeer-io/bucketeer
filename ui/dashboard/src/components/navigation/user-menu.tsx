import * as Popover from '@radix-ui/react-popover';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useAuth } from 'auth';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import LogoutConfirmModal from './logout-confirm';
import MenuItemComponent from './menu-item';
import UserProfileModal from './user-profile';

const UserMenu = () => {
  const { t } = useTranslation(['common']);
  const { logout, myOrganizations, consoleAccount } = useAuth();

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

  const [openProfileModal, onOpenProfileModal, onCloseProfileModal] =
    useToggleOpen(false);

  const avatar = consoleAccount?.avatarUrl
    ? consoleAccount.avatarUrl
    : primaryAvatar;

  const onHandleLogout = () => {
    logout();
    onCloseConfirmModal();
  };

  const menuItems = compact([
    {
      label: t(`navigation.user-profile`),
      icon: IconUser,
      onClick: onOpenProfileModal
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
      onClick: onOpenConfirmModal
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
        <AvatarImage image={avatar} size="sm" alt="user-avatar" />
      </Popover.Trigger>

      {openConfirmModal && (
        <LogoutConfirmModal
          isOpen={openConfirmModal}
          onClose={onCloseConfirmModal}
          onSubmit={onHandleLogout}
        />
      )}

      {openProfileModal && (
        <UserProfileModal
          isOpen={openProfileModal}
          onClose={onCloseProfileModal}
        />
      )}
    </Popover.Root>
  );
};

export default UserMenu;
