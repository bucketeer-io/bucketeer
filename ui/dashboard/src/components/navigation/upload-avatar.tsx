import { useCallback, useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { AvatarCommand } from '@api/account/account-updater';
import blueAvatar from 'assets/avatars/blue.svg';
import greenAvatar from 'assets/avatars/green.svg';
import orangeAvatar from 'assets/avatars/orange.svg';
import pinkAvatar from 'assets/avatars/pink.svg';
import primaryAvatar from 'assets/avatars/primary.svg';
import redAvatar from 'assets/avatars/red.svg';
import { useTranslation } from 'i18n';
import { readImageFile } from 'utils/files';
import { cn } from 'utils/style';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type UploadAvatarProps = {
  isOpen: boolean;
  onClose: () => void;
  onSelectAvatar: (avt: AvatarCommand | null) => void;
  onUploadPhoto: () => void;
};

type AvatarOption = {
  id: string;
  image: string;
};

enum AvatarColor {
  PRIMARY = 'primary',
  PINK = 'pink',
  GREEN = 'green',
  BLUE = 'blue',
  ORANGE = 'orange',
  RED = 'red'
}

const avatarOptions: AvatarOption[] = [
  {
    id: AvatarColor.PRIMARY,
    image: primaryAvatar
  },
  {
    id: AvatarColor.PINK,
    image: pinkAvatar
  },
  {
    id: AvatarColor.GREEN,
    image: greenAvatar
  },
  {
    id: AvatarColor.BLUE,
    image: blueAvatar
  },
  {
    id: AvatarColor.ORANGE,
    image: orangeAvatar
  },
  {
    id: AvatarColor.RED,
    image: redAvatar
  }
];

const UploadAvatarModal = ({
  isOpen,
  onClose,
  onSelectAvatar,
  onUploadPhoto
}: UploadAvatarProps) => {
  const { t } = useTranslation(['common']);
  const [currentAvatar, setCurrentAvatar] = useState<AvatarOption | null>(null);

  const getBorderColor = useCallback(
    (avatar: AvatarOption) => {
      if (!currentAvatar || currentAvatar.id !== avatar.id)
        return 'border-transparent';
      switch (currentAvatar.id) {
        case AvatarColor.PINK:
          return 'border-accent-pink-500';
        case AvatarColor.GREEN:
          return 'border-accent-green-500';
        case AvatarColor.BLUE:
          return 'border-accent-blue-500';
        case AvatarColor.ORANGE:
          return 'border-accent-orange-500';
        case AvatarColor.RED:
          return 'border-accent-red-500';
        case AvatarColor.PRIMARY:
        default:
          return 'border-primary-500';
      }
    },
    [currentAvatar]
  );

  const handleConvertImageToBase64 =
    useCallback(async (): Promise<AvatarCommand | null> => {
      if (!currentAvatar?.image) return null;
      const response = await fetch(currentAvatar.image);
      const blob = await response.blob();
      const base64String: string = await readImageFile(blob as File);
      const avatarImage = base64String?.split(',')[1] || '';
      return {
        avatarImage,
        avatarFileType: blob.type
      };
    }, [currentAvatar, readImageFile]);

  const handleSelectAvatar = useCallback(async () => {
    const avatar = await handleConvertImageToBase64();
    onSelectAvatar(avatar);
  }, [currentAvatar, onSelectAvatar]);

  return (
    <DialogModal
      className="w-[466px]"
      title={t('avatars')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-5 flex flex-col items-center">
        <div className="typo-para-big text-gray-700">{t(`upload-title`)}</div>
        <div className="flex items-center justify-center py-5 flex-wrap gap-6">
          {avatarOptions.map(avt => (
            <AvatarImage
              key={avt.id}
              image={avt.image}
              size="xl"
              alt="user-avatar"
              className={cn(
                'border-[3px] rounded-full cursor-pointer',
                getBorderColor(avt)
              )}
              onClick={() => setCurrentAvatar(avt)}
            />
          ))}
        </div>

        <Button onClick={onUploadPhoto} variant="text" size="sm">
          <Icon icon={IconAddOutlined} size="sm" />
          {t(`upload-photo`)}
        </Button>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={!currentAvatar} onClick={handleSelectAvatar}>
            {t(`save`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default UploadAvatarModal;
