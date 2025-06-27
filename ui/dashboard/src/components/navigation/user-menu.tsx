import { useMemo, useState } from 'react';
import { AccountAvatar } from '@api/account/account-updater';
import * as Popover from '@radix-ui/react-popover';
import defaultAvatar from 'assets/avatars/default.svg';
import { useAuth } from 'auth';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import EditPhotoProfileModal from './edit-photo';
import MenuItemComponent from './menu-item';
import UploadAvatarModal from './upload-avatar';
import UserProfileModal from './user-profile';

const UserMenu = ({ onOpenSwitchOrg }: { onOpenSwitchOrg: () => void }) => {
  const { t } = useTranslation(['common']);
  const { logout, myOrganizations, consoleAccount } = useAuth();

  const [selectedAvatar, setSelectedAvatar] = useState<AccountAvatar | null>(
    null
  );

  const [openProfileModal, onOpenProfileModal, onCloseProfileModal] =
    useToggleOpen(false);

  const [
    openUploadAvatarModal,
    onOpenUploadAvatarModal,
    onCloseUploadAvatarModal
  ] = useToggleOpen(false);

  const [
    openUploadPhotoModal,
    onOpenUploadPhotoModal,
    onCloseUploadPhotoModal
  ] = useToggleOpen(false);

  const avatar = consoleAccount?.avatarImage;
  const isHiddenProfileMenu =
    consoleAccount?.isSystemAdmin && !consoleAccount?.organization.systemAdmin;

  const avatarSrc = useMemo(
    () =>
      avatar
        ? `data:${consoleAccount?.avatarFileType};base64,${avatar}`
        : defaultAvatar,
    [avatar, defaultAvatar]
  );

  const onSelectAvatar = (avatar: AccountAvatar | null, cb?: () => void) => {
    setSelectedAvatar(avatar);
    onOpenProfileModal();
    if (cb) cb();
  };

  const menuItems = compact([
    !isHiddenProfileMenu && {
      label: t(`navigation.user-profile`),
      icon: IconUser,
      onClick: onOpenProfileModal
    },
    myOrganizations.length > 1 && {
      label: consoleAccount?.organization?.name || '',
      icon: IconBuilding,
      actIcon: IconChevronRight,
      onClick: onOpenSwitchOrg
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
        <div className="bg-primary-600 rounded-lg min-w-[200px] max-w-[220px] mb-2">
          {menuItems.map((item, index) => (
            <MenuItemComponent {...item} key={index} />
          ))}
        </div>
      </Popover.Content>
      <Popover.Trigger>
        <AvatarImage image={avatarSrc} size="sm" alt="user-avatar" />
      </Popover.Trigger>

      {openProfileModal && (
        <UserProfileModal
          selectedAvatar={selectedAvatar}
          isOpen={openProfileModal}
          onClose={() => {
            onCloseProfileModal();
            setSelectedAvatar(null);
          }}
          onEditAvatar={() => {
            onCloseProfileModal();
            onOpenUploadAvatarModal();
          }}
        />
      )}
      {openUploadAvatarModal && (
        <UploadAvatarModal
          isOpen={openUploadAvatarModal}
          onClose={() => {
            onCloseUploadAvatarModal();
            onOpenProfileModal();
            setSelectedAvatar(null);
          }}
          onUploadPhoto={() => {
            onCloseUploadAvatarModal();
            onOpenUploadPhotoModal();
          }}
          onSelectAvatar={avatar =>
            onSelectAvatar(avatar, onCloseUploadAvatarModal)
          }
        />
      )}
      {openUploadPhotoModal && (
        <EditPhotoProfileModal
          onUpload={avatar => onSelectAvatar(avatar, onCloseUploadPhotoModal)}
          isOpen={openUploadPhotoModal}
          onClose={() => {
            onOpenProfileModal();
            onCloseUploadPhotoModal();
          }}
        />
      )}
    </Popover.Root>
  );
};

export default UserMenu;
