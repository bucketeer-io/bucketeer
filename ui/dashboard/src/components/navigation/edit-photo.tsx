import { useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { readImageFile } from 'utils/files';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import PhotoResize, { PhotoResizeHandle } from 'components/photo-resize';
import PhotoSelect, { PhotoSelectRef } from 'components/photo-select';

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10 Mb

const ASPECT_RATIO = {
  WIDTH: 3,
  HEIGHT: 1
};

interface FormValues {
  image: string | undefined;
}

type ProcessingStatus = 'select' | 'resize' | 'ready-to-upload';

const formSchema = yup.object().shape({
  image: yup.string().required()
});

interface Props {
  onUpload: (image: string) => Promise<any>;
  onCompleted: () => void;
}

export type EditPhottoProfileProps = {
  isOpen: boolean;
  onClose: () => void;
};

const EditPhottoProfileModal = ({
  isOpen,
  onClose
}: EditPhottoProfileProps) => {
  const { t } = useTranslation(['common']);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      image: ''
    }
  });

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
      <div className="p-5 flex flex-col items-center">
        <div className={progressingStatus === 'select' ? 'w-full' : 'hidden'}>
          <PhotoSelect
            ref={photoSelectRef}
            maxFileSize={MAX_FILE_SIZE}
            onChange={file => {
              readImageFile(file).then(result => {
                form.setValue('image', result);
              });
              setProcessingStatus('resize');
            }}
          />
        </div>

        {form.getValues('image') && progressingStatus === 'resize' && (
          <div className="flex w-full flex-col items-center">
            <div className="h-[400px] w-full">
              <PhotoResize
                ref={photoResizeRef}
                aspect={ASPECT_RATIO.WIDTH / ASPECT_RATIO.HEIGHT}
                value={form.getValues('image')}
                onChange={value => {
                  form.setValue('image', value);
                  setProcessingStatus('ready-to-upload');
                }}
              />
            </div>
            <div className="mt-16">
              <Button
                variant="secondary"
                type="button"
                disabled={form.formState.isSubmitting}
                onClick={() => {
                  setProcessingStatus('select');
                  photoSelectRef.current?.click();
                }}
              >
                {`Change Picture`}
              </Button>
            </div>
          </div>
        )}

        {form.getValues('image') && progressingStatus === 'ready-to-upload' && (
          <img
            className="mx-auto h-[400px] object-contain"
            alt=""
            src={form.getValues('image')}
          />
        )}
      </div>
    </DialogModal>
  );
};

export default EditPhottoProfileModal;
