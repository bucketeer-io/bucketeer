import { useTranslation } from 'i18n';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  onSubmit: () => void;
  isOpen: boolean;
  title: string;
  description: React.ReactElement | string;
  loading?: boolean;
  onClose: () => void;
};

const ConfirmModal = ({
  onSubmit,
  isOpen,
  title,
  description,
  loading,
  onClose
}: ConfirmModalProps) => {
  const { t } = useTranslation(['common']);

  return (
    <DialogModal
      className="w-[500px]"
      title={title}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="typo-para-medium text-gray-700 w-full">
          {description}
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button loading={loading} onClick={onSubmit}>
            {t(`submit`)}
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

export default ConfirmModal;
