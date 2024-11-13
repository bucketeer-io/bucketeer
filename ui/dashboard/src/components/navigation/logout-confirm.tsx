import { useTranslation } from 'i18n';
import { IconLogoutConfirm } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type LogoutConfirmModalProps = {
  onSubmit: () => void;
  isOpen: boolean;
  onClose: () => void;
};

const LogoutConfirmModal = ({
  onSubmit,
  isOpen,
  onClose
}: LogoutConfirmModalProps) => {
  const { t } = useTranslation(['common', 'auth']);

  return (
    <DialogModal
      className="w-[500px]"
      title={t(`auth:logout-title`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
        <IconLogoutConfirm />
        <div className="typo-para-big text-gray-700 text-center">
          {t(`auth:logout-description`)}
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button variant="negative" className="w-24" onClick={onSubmit}>
            {t(`yes`)}
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

export default LogoutConfirmModal;
