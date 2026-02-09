import { useTranslation } from 'i18n';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  disabled?: boolean;
  isOpen: boolean;
  title: string;
  description: React.ReactElement | string;
  loading?: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

const ConfirmModal = ({
  disabled,
  isOpen,
  title,
  description,
  loading,
  onClose,
  onSubmit
}: ConfirmModalProps) => {
  const { t } = useTranslation(['common']);

  return (
    <DialogModal
      className="max-w-[500px]"
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
          <Button loading={loading} onClick={onSubmit} disabled={disabled}>
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
