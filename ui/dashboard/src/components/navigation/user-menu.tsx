import { useCallback, useMemo, useRef, useState } from 'react';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { AccountAvatar, accountUpdater } from '@api/account/account-updater';
import * as Popover from '@radix-ui/react-popover';
import defaultAvatar from 'assets/avatars/default.svg';
import { getCurrentEnvironment, useAuth } from 'auth';
import { urls } from 'configs';
import { useToast, useToggleOpen } from 'hooks';
import { getLanguage, Language, setLanguage, useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { setConsoleVersion } from 'storage/console';
import { onChangeFontWithLocalized } from 'utils/function';
import {
  IconBucketWhite,
  IconBuilding,
  IconChevronRight,
  IconLogout,
  IconUser
} from '@icons';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import { AvatarImage } from 'components/avatar';
import { DropdownOption } from 'components/dropdown';
import EditPhotoProfileModal from './edit-photo';
import MenuItemComponent from './menu-item';
import UploadAvatarModal from './upload-avatar';
import UserProfileModal from './user-profile';

const UserMenu = ({ onOpenSwitchOrg }: { onOpenSwitchOrg: () => void }) => {
  const { t } = useTranslation(['common']);
  const { myOrganizations, consoleAccount, logout, onMeFetcher } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { errorNotify } = useToast();
  const popoverCloseRef = useRef<HTMLButtonElement>(null);

  const [selectedAvatar, setSelectedAvatar] = useState<AccountAvatar | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(false);
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

  const handleUpdateLanguage = useCallback(
    async (value: string) => {
      try {
        if (value === consoleAccount?.language)
          return popoverCloseRef?.current?.click();
        setIsLoading(true);
        const resp = await accountUpdater({
          organizationId: currentEnvironment.organizationId,
          email: consoleAccount!.email,
          language: value
        });
        if (resp) {
          const i18nLanguage = getLanguage();
          if (value !== i18nLanguage) setLanguage(value as Language);
          onChangeFontWithLocalized(value === Language.JAPANESE);
          popoverCloseRef?.current?.click();
          await onMeFetcher({
            organizationId: currentEnvironment.organizationId
          });
          setIsLoading(false);
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [consoleAccount, currentEnvironment]
  );

  const menuItems = compact([
    {
      label: t('old-console'),
      icon: IconBucketWhite,
      actIcon: IconLaunchOutlined,
      onClick: () => {
        setConsoleVersion({
          version: 'old',
          email: consoleAccount!.email
        });
        window.location.href = urls.OLD_CONSOLE_ENDPOINT;
      }
    },
    !isHiddenProfileMenu && {
      label: t(`navigation.user-profile`),
      icon: IconUser,
      onClick: onOpenProfileModal
    },
    !isHiddenProfileMenu && {
      label:
        languageList.find(item => item.value === consoleAccount?.language)
          ?.label || '',
      icon: languageList.find(item => item.value === consoleAccount?.language)
        ?.icon,
      actIcon: IconChevronRight,
      loading: isLoading,
      options: languageList as DropdownOption[],
      onSelectOption: handleUpdateLanguage
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
        <Popover.Close ref={popoverCloseRef} className="hidden" />
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
