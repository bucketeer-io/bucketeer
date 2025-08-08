import { useTranslation } from 'i18n';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type SegmentUploadingModalProps = {
  isOpen: boolean;
  onClose: () => void;
};

const SegmentUploadingModal = ({
  isOpen,
  onClose
}: SegmentUploadingModalProps) => {
  const { t } = useTranslation(['common', 'table']);

  return (
    <DialogModal
      className="w-[500px]"
      title={t('table:user-segment.segment-uploading-title')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="typo-para-medium text-gray-700 w-full">
          {t('table:user-segment.segment-uploading-desc')}
        </div>
      </div>

      <ButtonBar
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`close`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default SegmentUploadingModal;
