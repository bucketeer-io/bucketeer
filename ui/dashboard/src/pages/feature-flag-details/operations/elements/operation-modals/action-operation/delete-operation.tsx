import { useTranslation } from 'i18n';
import { InfoIcon } from 'lucide-react';
import { Environment, Feature } from '@types';
import { OpsTypeMap } from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import OperationActiveModal from './operation-active';

export type DeleteOperationModalProps = {
  operationType: OpsTypeMap;
  editable: boolean;
  isRunning?: boolean;
  isOpen: boolean;
  environment: Environment;
  feature: Feature;
  loading?: boolean;
  refetchFeature: () => void;
  onClose: () => void;
  onSubmit: () => void;
};

export const DeleteOperationModal = ({
  editable,
  operationType,
  feature,
  environment,
  isOpen,
  loading,
  isRunning = false,
  refetchFeature,
  onClose,
  onSubmit
}: DeleteOperationModalProps) => {
  const { t } = useTranslation(['common', 'table', 'form']);

  const isRolloutType = operationType === OpsTypeMap.ROLLOUT;
  const isScheduleType = operationType === OpsTypeMap.SCHEDULE;

  const transKey = `table:popover.delete-${isRolloutType ? 'rollout' : isScheduleType ? 'operation' : 'kill-switch'}`;

  const infoTitleKey =
    operationType === OpsTypeMap.SCHEDULE
      ? 'form:operation.confirm-delete-schedule-title'
      : operationType === OpsTypeMap.EVENT_RATE
        ? 'form:operation.confirm-delete-event-rate-title'
        : 'form:operation.confirm-delete-rollout-title';

  const infoDescKey =
    operationType === OpsTypeMap.SCHEDULE
      ? 'form:operation.confirm-delete-schedule-desc'
      : operationType === OpsTypeMap.EVENT_RATE
        ? 'form:operation.confirm-delete-event-rate-desc'
        : 'form:operation.confirm-delete-rollout-desc';

  return (
    <DialogModal
      className="max-w-[600px]"
      title={t(transKey)}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isRunning ? (
        <OperationActiveModal
          isDeleting
          refetchFeature={refetchFeature}
          onClose={onClose}
          onActionOperation={onSubmit}
          editable={editable}
          feature={feature}
          environment={environment}
          loading={loading}
        />
      ) : (
        <>
          <div className="flex flex-col w-full items-start px-5 py-4">
            <div className="w-full rounded-lg border-l-[8px] border-primary-500 px-4 py-3 shadow-card">
              <div className="flex items-start gap-4 typo-para-medium">
                <Icon
                  icon={InfoIcon}
                  size="xxs"
                  className="mt-[5px] text-primary-500"
                />
                <div>
                  <p className="font-bold text-primary-500">
                    {t(infoTitleKey)}
                  </p>
                  <p className="typo-para-medium text-gray-500 w-full mt-2">
                    {t(infoDescKey)}
                  </p>
                </div>
              </div>
            </div>
          </div>
          <ButtonBar
            secondaryButton={
              <Button loading={loading} onClick={onSubmit} disabled={!editable}>
                {t(`submit`)}
              </Button>
            }
            primaryButton={
              <Button onClick={onClose} variant="secondary">
                {t(`cancel`)}
              </Button>
            }
          />
        </>
      )}
    </DialogModal>
  );
};
