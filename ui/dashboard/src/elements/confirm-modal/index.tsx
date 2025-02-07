import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  onSubmit: () => void;
  isOpen: boolean;
  title: string;
  description: React.ReactElement | string;
  loading?: boolean;
  className?: string;
  onClose: () => void;
};

const ConfirmModal = ({
  onSubmit,
  isOpen,
  title,
  description,
  loading,
  className,
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
      <div
        className={cn('flex flex-col w-full items-start px-5 py-8', className)}
      >
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
