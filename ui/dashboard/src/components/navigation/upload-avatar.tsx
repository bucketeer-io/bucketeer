import { IconAddOutlined } from 'react-icons-material-design';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useTranslation } from 'i18n';
import { AvatarImage } from 'components/avatar';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type UploadAvatarProps = {
  isOpen: boolean;
  onClose: () => void;
  onUploadPhoto: () => void;
};

const UploadAvatarModal = ({
  isOpen,
  onClose,
  onUploadPhoto
}: UploadAvatarProps) => {
  const { t } = useTranslation(['common']);

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
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
          <AvatarImage image={primaryAvatar} size="xl" alt="user-avatar" />
        </div>

        <Button onClick={onUploadPhoto} variant="text" size="sm">
          <Icon icon={IconAddOutlined} size="sm" />
          {t(`upload-photo`)}
        </Button>
      </div>

      <ButtonBar
        secondaryButton={<Button>{t(`save`)}</Button>}
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
