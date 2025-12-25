import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { OpsTypeMap } from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type StopOperationModalProps = {
  editable: boolean;
  operationType: OpsTypeMap;
  isOpen: boolean;
  loading?: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

const StopOperationModal = ({
  editable,
  operationType,
  isOpen,
  loading,
  onClose,
  onSubmit
}: StopOperationModalProps) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const transKey = useMemo(
    () =>
      operationType === OpsTypeMap.SCHEDULE
        ? 'schedule'
        : operationType === OpsTypeMap.EVENT_RATE
          ? 'kill-switch'
          : 'rollout',
    [operationType]
  );

  return (
    <DialogModal
      className="w-[500px]"
      title={t(`table:stop-${transKey}`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="typo-para-medium text-accent-red-500 w-full">
          <p>{t('table:stop-operation-rollout.title')}</p>
          <p>{t('table:stop-operation-rollout.desc')}</p>
        </div>
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
    </DialogModal>
  );
};

export default StopOperationModal;
