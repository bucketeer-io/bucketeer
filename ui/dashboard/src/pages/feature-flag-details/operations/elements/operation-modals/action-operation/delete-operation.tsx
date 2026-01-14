import { ReactElement } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { TFunction } from 'i18next';
import { Environment, Feature } from '@types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import OperationActiveModal from './operation-active';

export type DeleteOperationModalProps = {
  isRolloutType: boolean;
  isScheduleType: boolean;
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

type DeleteOperationModalDefaultProps = {
  editable: boolean;
  loading?: boolean;
  description: ReactElement;
  trans: TFunction;
  onClose: () => void;
  onSubmit: () => void;
};

const DeleteOperationDefaultModal = ({
  loading,
  editable,
  description,
  trans,
  onSubmit,
  onClose
}: DeleteOperationModalDefaultProps) => {
  return (
    <>
      <div className="flex flex-col w-full items-start px-5 py-8">
        <div className="typo-para-medium text-gray-700 w-full">
          {description}
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button loading={loading} onClick={onSubmit} disabled={!editable}>
            {trans(`submit`)}
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

export const DeleteOperationModal = ({
  editable,
  isRolloutType,
  isScheduleType,
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
  const transKey = `table:popover.delete-${isRolloutType ? 'rollout' : isScheduleType ? 'operation' : 'kill-switch'}`;
  return (
    <DialogModal
      className={isRunning ? 'max-w-[600px]' : 'max-w-[500px]'}
      title={t(transKey)}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isRunning ? (
        <OperationActiveModal
          refetchFeature={refetchFeature}
          onClose={onClose}
          onActionOperation={onSubmit}
          editable={editable}
          feature={feature}
          environment={environment}
          loading={loading}
        />
      ) : (
        <DeleteOperationDefaultModal
          description={
            <Trans
              i18nKey={'table:operations.confirm-delete-operation'}
              values={{
                type: t(
                  `form:feature-flags.${isRolloutType ? 'rollout' : isScheduleType ? 'schedule' : 'kill-switch'}`
                )
              }}
              components={{
                bold: <strong />
              }}
            />
          }
          trans={t}
          onClose={onClose}
          onSubmit={onSubmit}
          editable={editable}
          loading={loading}
        />
      )}
    </DialogModal>
  );
};
