import { Trans } from 'react-i18next';
import { useScreen } from 'hooks';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import { IconDelete } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type DeleteMemberProps = {
  onSubmit: () => void;
  isOpen: boolean;
  onClose: () => void;
  goal: Goal;
  loading: boolean;
  disabled?: boolean;
};

const DeleteGoalModal = ({
  onSubmit,
  isOpen,
  onClose,
  goal,
  loading,
  disabled
}: DeleteMemberProps) => {
  const { t } = useTranslation(['common']);
  const { fromMobileScreen } = useScreen();
  const sizeIcon = fromMobileScreen ? 'fit' : '3xl';
  return (
    <DialogModal
      className="max-w-[500px]"
      title={t(`delete-goal`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
        <Icon icon={IconDelete} size={sizeIcon} />
        <div className="typo-para-medium sm:typo-para-big text-gray-700 text-center">
          <Trans
            i18nKey="table:goals.delete-goal-desc"
            values={{ name: goal.name }}
            components={{ bold: <strong /> }}
          />
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            disabled={disabled}
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

export default DeleteGoalModal;
