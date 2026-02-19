import { useState } from 'react';
import { useQueryScheduledFlagChanges } from '@queries/scheduled-flag-changes';
import { useTranslation } from 'i18n';
import { ScheduledFlagChange, ScheduledFlagChangeStatuses } from '@types';
import SlideModal from 'components/modal/slide';
import ApplyNowDialog from '../scheduled-changes-banner/apply-now-dialog';
import CancelScheduleDialog from '../scheduled-changes-banner/cancel-schedule-dialog';
import EditScheduleDialog from '../scheduled-changes-banner/edit-schedule-dialog';
import ScheduleCard, { ScheduleCardAction } from './schedule-card';

const PANEL_STATUSES = [
  ScheduledFlagChangeStatuses.PENDING,
  ScheduledFlagChangeStatuses.CONFLICT
];

interface ScheduledChangesPanelProps {
  featureId: string;
  environmentId: string;
  isOpen: boolean;
  onClose: () => void;
}

const ScheduledChangesPanel = ({
  featureId,
  environmentId,
  isOpen,
  onClose
}: ScheduledChangesPanelProps) => {
  const { t } = useTranslation(['form']);
  const [activeDialog, setActiveDialog] = useState<{
    type: ScheduleCardAction;
    schedule: ScheduledFlagChange;
  } | null>(null);

  const { data: listData, isLoading } = useQueryScheduledFlagChanges({
    params: {
      environmentId,
      featureId,
      statuses: [...PANEL_STATUSES],
      orderBy: 'SCHEDULED_AT',
      orderDirection: 'ASC',
      pageSize: 50
    },
    enabled: isOpen && !!featureId && !!environmentId
  });

  const schedules = listData?.scheduledFlagChanges ?? [];

  const handleAction = (
    action: ScheduleCardAction,
    schedule: ScheduledFlagChange
  ) => {
    setActiveDialog({ type: action, schedule });
  };

  const handleCloseDialog = () => setActiveDialog(null);

  return (
    <>
      <SlideModal
        title={t('feature-flags.scheduled-changes-panel-title')}
        isOpen={isOpen}
        onClose={onClose}
      >
        <div className="flex flex-col gap-y-4 p-4">
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <p className="typo-para-small text-gray-400">
                {t('feature-flags.loading')}
              </p>
            </div>
          ) : schedules.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <p className="typo-para-small text-gray-400">
                {t('feature-flags.no-scheduled-changes')}
              </p>
            </div>
          ) : (
            schedules.map(schedule => (
              <ScheduleCard
                key={schedule.id}
                schedule={schedule}
                onAction={handleAction}
              />
            ))
          )}
        </div>
      </SlideModal>

      {activeDialog?.type === 'EDIT_SCHEDULE' && (
        <EditScheduleDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
      {activeDialog?.type === 'CANCEL_SCHEDULE' && (
        <CancelScheduleDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
      {activeDialog?.type === 'APPLY_NOW' && (
        <ApplyNowDialog
          schedule={activeDialog.schedule}
          isOpen
          onClose={handleCloseDialog}
        />
      )}
    </>
  );
};

export default ScheduledChangesPanel;
