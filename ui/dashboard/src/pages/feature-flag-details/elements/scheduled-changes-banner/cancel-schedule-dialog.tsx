import { Trans } from 'react-i18next';
import { useDeleteScheduledFlagChange } from '@queries/scheduled-flag-changes';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { ScheduledFlagChange } from '@types';
import { IconCalendarCancel } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import { formatScheduledDate } from './utils';

interface CancelScheduleDialogProps {
  schedule: ScheduledFlagChange;
  isOpen: boolean;
  onClose: () => void;
}

const CancelScheduleDialog = ({
  schedule,
  isOpen,
  onClose
}: CancelScheduleDialogProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();
  const deleteMutation = useDeleteScheduledFlagChange();

  const handleConfirm = async () => {
    try {
      await deleteMutation.mutateAsync({
        environmentId: schedule.environmentId,
        id: schedule.id
      });
      notify({
        message: t('form:feature-flags.schedule-cancelled')
      });
      onClose();
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:feature-flags.cancel-schedule')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-center px-5 py-8 gap-y-6">
        <IconCalendarCancel />
        <p className="typo-para-big text-gray-700 text-center">
          <Trans
            i18nKey="form:feature-flags.cancel-schedule-confirm"
            values={{ datetime: formatScheduledDate(schedule.scheduledAt) }}
            components={{ bold: <strong /> }}
          />
        </p>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            variant="negative"
            loading={deleteMutation.isPending}
            onClick={handleConfirm}
          >
            {t('confirm')}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t('cancel')}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default CancelScheduleDialog;
