import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { OpsTypeMap } from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';

export type StopOperationModalProps = {
  operationType: OpsTypeMap;
  isOpen: boolean;
  loading?: boolean;
  onClose: () => void;
  onSubmit: () => void;
};

const StopOperationModal = ({
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
          ? 'event-rate'
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
          <Trans
            i18nKey={'table:stop-operation-type-desc'}
            values={{
              type: t(`form:feature-flags.${transKey}`)
            }}
          />
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button loading={loading} onClick={onSubmit}>
            {t(`submit`)}
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
