import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { IconDelete } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type DeleteMemberProps = {
  onSubmit: () => void;
  isOpen: boolean;
  onClose: () => void;
  member: Account;
  loading: boolean;
};

const DeleteMemberModal = ({
  onSubmit,
  isOpen,
  onClose,
  member,
  loading
}: DeleteMemberProps) => {
  const { t } = useTranslation(['common']);

  return (
    <DialogModal
      className="max-w-[500px]"
      title={t(`delete-member`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
        <IconDelete />
        <div className="typo-para-medium sm:typo-para-big text-gray-700 text-center">
          <Trans
            i18nKey="table:members.delete-member-desc"
            values={{ email: member.email }}
            components={{ bold: <strong /> }}
          />
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            variant="negative"
            loading={loading}
            className="w-24"
            onClick={onSubmit}
          >
            {t(`delete`)}
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

export default DeleteMemberModal;
