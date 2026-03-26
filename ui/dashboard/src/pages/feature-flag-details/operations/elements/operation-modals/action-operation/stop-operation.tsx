import { useTranslation } from 'i18n';
import { Environment, Feature } from '@types';
import { OpsTypeMap } from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import OperationActiveModal from './operation-active';
import OperationInfoCard from './operation-info-card';

export type StopOperationModalProps = {
  editable: boolean;
  operationType: OpsTypeMap;
  isRunning?: boolean;
  isOpen: boolean;
  environment: Environment;
  feature: Feature;
  loading?: boolean;
  refetchFeatures: () => void;
  onClose: () => void;
  onSubmit: () => void;
};

export const StopOperationModal = ({
  editable,
  operationType,
  feature,
  environment,
  isOpen,
  loading,
  isRunning = false,
  refetchFeatures,
  onClose,
  onSubmit
}: StopOperationModalProps) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const transKey =
    operationType === OpsTypeMap.SCHEDULE
      ? 'schedule'
      : operationType === OpsTypeMap.EVENT_RATE
        ? 'kill-switch'
        : 'rollout';

  const infoTitleKey =
    operationType === OpsTypeMap.SCHEDULE
      ? 'form:operation.confirm-stop-schedule-title'
      : operationType === OpsTypeMap.EVENT_RATE
        ? 'form:operation.confirm-stop-event-rate-title'
        : 'form:operation.confirm-stop-title';

  const infoDescKey =
    operationType === OpsTypeMap.SCHEDULE
      ? 'form:operation.confirm-stop-schedule-desc'
      : operationType === OpsTypeMap.EVENT_RATE
        ? 'form:operation.confirm-stop-event-rate-desc'
        : 'form:operation.confirm-stop-rollout-desc';

  return (
    <DialogModal
      className="max-w-[600px]"
      title={t(`table:stop-${transKey}`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isRunning ? (
        <OperationActiveModal
          refetchFeature={refetchFeatures}
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
            <OperationInfoCard
              title={t(infoTitleKey)}
              description={t(infoDescKey)}
            />
          </div>
          <ButtonBar
            secondaryButton={
              <Button loading={loading} onClick={onSubmit} disabled={!editable}>
                {t(`stop`)}
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
