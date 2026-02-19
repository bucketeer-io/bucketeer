import { useState } from 'react';
import { useUpdateScheduledFlagChange } from '@queries/scheduled-flag-changes';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { ScheduledFlagChange } from '@types';
import { IconWatch } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { ReactDatePicker } from 'components/date-time-picker';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

interface EditScheduleDialogProps {
  schedule: ScheduledFlagChange;
  isOpen: boolean;
  onClose: () => void;
}

const EditScheduleDialog = ({
  schedule,
  isOpen,
  onClose
}: EditScheduleDialogProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { notify, errorNotify } = useToast();
  const updateMutation = useUpdateScheduledFlagChange();

  const [scheduleAt, setScheduleAt] = useState(schedule.scheduledAt);
  const scheduleDate = scheduleAt ? new Date(Number(scheduleAt) * 1000) : null;

  const handleSave = async () => {
    try {
      const resp = await updateMutation.mutateAsync({
        environmentId: schedule.environmentId,
        id: schedule.id,
        scheduledAt: scheduleAt
      });
      if (resp) {
        notify({
          message: t('form:feature-flags.schedule-updated')
        });
        onClose();
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:feature-flags.edit-schedule')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="flex gap-x-4">
          <div>
            <p className="typo-para-small text-gray-700 mb-1">
              {t('form:feature-flags.update-date')}
              <span className="text-accent-red-500 ml-0.5">*</span>
            </p>
            <ReactDatePicker
              dateFormat="yyyy/MM/dd"
              minDate={new Date()}
              selected={scheduleDate}
              showTimeSelect={false}
              className="w-[186px]"
              onChange={date => {
                if (date) {
                  if (scheduleDate) {
                    date.setHours(
                      scheduleDate.getHours(),
                      scheduleDate.getMinutes(),
                      0,
                      0
                    );
                  }
                  setScheduleAt(
                    String(Math.floor(date.getTime() / 1000))
                  );
                }
              }}
            />
          </div>
          <div>
            <p className="typo-para-small text-gray-700 mb-1">
              {t('form:feature-flags.update-time')}
              <span className="text-accent-red-500 ml-0.5">*</span>
            </p>
            <ReactDatePicker
              dateFormat="HH:mm"
              timeFormat="HH:mm"
              selected={scheduleDate}
              showTimeSelectOnly={true}
              className="w-[124px]"
              onChange={date => {
                if (date) {
                  setScheduleAt(
                    String(Math.floor(date.getTime() / 1000))
                  );
                }
              }}
              icon={
                <Icon icon={IconWatch} className="flex-center" />
              }
            />
          </div>
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            loading={updateMutation.isPending}
            onClick={handleSave}
          >
            {t('submit')}
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

export default EditScheduleDialog;
