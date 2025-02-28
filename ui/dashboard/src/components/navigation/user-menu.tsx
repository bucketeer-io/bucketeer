import { useMemo, useState } from 'react';
import { AvatarCommand } from '@api/account/account-updater';
import * as Popover from '@radix-ui/react-popover';
import defaultAvatar from 'assets/avatars/default.svg';
import { useAuth } from 'auth';
import { useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { IconBuilding, IconChevronRight, IconLogout, IconUser } from '@icons';
import { AvatarImage } from 'components/avatar';
import EditPhotoProfileModal from './edit-photo';
import LogoutConfirmModal from './logout-confirm';
import MenuItemComponent from './menu-item';
import UploadAvatarModal from './upload-avatar';
import UserProfileModal from './user-profile';

const UserMenu = () => {
  const { t } = useTranslation(['common']);
  const { logout, myOrganizations, consoleAccount } = useAuth();

  const [selectedAvatar, setSelectedAvatar] = useState<AvatarCommand | null>(
    null
  );

  const [openConfirmModal, onOpenConfirmModal, onCloseConfirmModal] =
    useToggleOpen(false);

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

  const onSelectAvatar = (avatar: AvatarCommand | null, cb?: () => void) => {
    setSelectedAvatar(avatar);
    onOpenProfileModal();
    if (cb) cb();
  };

  const onHandleLogout = () => {
    logout();
    onCloseConfirmModal();
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
        <AvatarImage image={avatarSrc} size="sm" alt="user-avatar" />
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
          onClose={onCloseUploadPhotoModal}
        />
      )}
    </Popover.Root>
  );
};

export default UserMenu;
