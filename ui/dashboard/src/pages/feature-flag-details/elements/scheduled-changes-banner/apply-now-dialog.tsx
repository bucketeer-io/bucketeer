import { useExecuteScheduledFlagChange } from '@queries/scheduled-flag-changes';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { ScheduledFlagChange } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

const formatScheduledDate = (timestamp: string | number): string => {
  const date = new Date(Number(timestamp) * 1000);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}/${month}/${day} ${hours}:${minutes}`;
};

interface ApplyNowDialogProps {
  schedule: ScheduledFlagChange;
  isOpen: boolean;
  onClose: () => void;
}

const ApplyNowDialog = ({
  schedule,
  isOpen,
  onClose
}: ApplyNowDialogProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const executeMutation = useExecuteScheduledFlagChange();

  const handleConfirm = async () => {
    try {
      await executeMutation.mutateAsync({
        environmentId: schedule.environmentId,
        id: schedule.id
      });
      notify({
        message: t('form:feature-flags.schedule-applied-now')
      });
      invalidateFeature(queryClient);
      invalidateFeatures(queryClient);
      onClose();
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:feature-flags.apply-now')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-center px-5 py-8 gap-y-4">
        <p className="typo-para-medium text-gray-600 text-center">
          {t('form:feature-flags.apply-now-confirm', {
            datetime: formatScheduledDate(schedule.scheduledAt)
          })}
        </p>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            loading={executeMutation.isPending}
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

export default ApplyNowDialog;
