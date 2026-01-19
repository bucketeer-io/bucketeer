import { ReactElement } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { TFunction } from 'i18next';
import { Environment, Feature } from '@types';
import { OpsTypeMap } from 'pages/feature-flag-details/operations/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import OperationActiveModal from './operation-active';

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

type StopOperationDefaultProps = {
  editable: boolean;
  loading?: boolean;
  description: ReactElement;
  trans: TFunction;
  onClose: () => void;
  onSubmit: () => void;
};

const StopOperationDefault = ({
  description,
  loading,
  editable,
  trans,
  onSubmit,
  onClose
}: StopOperationDefaultProps) => {
  return (
    <>
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="typo-para-medium text-accent-red-500 w-full">
          {description}
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button loading={loading} onClick={onSubmit} disabled={!editable}>
            {trans(`stop`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {trans(`cancel`)}
          </Button>
        }
      />
    </>
  );
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

  return (
    <DialogModal
      className={isRunning ? 'max-w-[600px]' : 'max-w-[500px]'}
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
        <StopOperationDefault
          description={
            <Trans
              i18nKey={'table:stop-operation-type-desc'}
              values={{
                type: t(`form:feature-flags.${transKey}`)
              }}
            />
          }
          trans={t}
          loading={loading}
          editable={editable}
          onSubmit={onSubmit}
          onClose={onClose}
        />
      )}
    </DialogModal>
  );
};
