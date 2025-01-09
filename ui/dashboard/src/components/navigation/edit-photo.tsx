import { useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { AvatarCommand } from '@api/account/account-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { readImageFile } from 'utils/files';
import { cn } from 'utils/style';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import PhotoResize, { PhotoResizeHandle } from 'components/photo-resize';
import PhotoSelect, { PhotoSelectRef } from 'components/photo-select';

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10 Mb

type ProcessingStatus = 'select' | 'resize';

const formSchema = yup.object().shape({
  avatarImage: yup.string().required(),
  avatarFileType: yup.string().required()
});

export type EditPhotoProfileProps = {
  isOpen: boolean;
  onClose: () => void;
  onUpload: (avatar: AvatarCommand) => void;
};

const EditPhotoProfileModal = ({
  isOpen,
  onClose,
  onUpload
}: EditPhotoProfileProps) => {
  const { t } = useTranslation(['common']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      avatarImage: '',
      avatarFileType: ''
    }
  });
  const avatarImage = form.watch('avatarImage');
  const avatarFileType = form.watch('avatarFileType');

  const [progressingStatus, setProcessingStatus] =
    useState<ProcessingStatus>('select');

  const photoSelectRef = useRef<PhotoSelectRef>(null);
  const photoResizeRef = useRef<PhotoResizeHandle>(null);

  return (
    <DialogModal
      className="w-[466px]"
      title={t('edit-photo-profile')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div
        className={cn(
          'p-5',
          progressingStatus === 'select' ? 'w-full' : 'hidden'
        )}
      >
        <PhotoSelect
          ref={photoSelectRef}
          maxFileSize={MAX_FILE_SIZE}
          onChange={file => {
            readImageFile(file).then(result => {
              form.setValue('avatarImage', result);
              form.setValue('avatarFileType', file.type);
            });
            setProcessingStatus('resize');
          }}
        />
      </div>

      {avatarImage && progressingStatus === 'resize' && (
        <>
          <div className="flex w-full flex-col items-center p-5">
            <div className="h-[220px] w-full">
              <PhotoResize
                ref={photoResizeRef}
                aspect={2 / 2}
                value={avatarImage}
                onChange={value => form.setValue('avatarImage', value)}
              />
            </div>
          </div>
          <ButtonBar
            secondaryButton={
              <Button
                onClick={() =>
                  onUpload({
                    avatarImage: avatarImage?.split(',')[1] || '',
                    avatarFileType
                  })
                }
                loading={form.formState.isSubmitting}
              >
                {t(`submit`)}
              </Button>
            }
            primaryButton={
              <Button onClick={onClose} variant="secondary">
                {t(`cancel`)}
              </Button>
            }
          />
        </>
      )}
    </DialogModal>
  );
};

export default EditPhotoProfileModal;
