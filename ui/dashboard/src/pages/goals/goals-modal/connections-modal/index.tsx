import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Goal } from '@types';
import { IconExperimentsConnected, IconOperationsConnected } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type ConfirmModalProps = {
  goal: Goal;
  isOpen: boolean;
  onClose: () => void;
};

const ConnectionsModal = ({ goal, isOpen, onClose }: ConfirmModalProps) => {
  const { t } = useTranslation(['common']);
  const connectionType = useMemo(() => goal?.connections?.type, [goal]);

  return (
    <DialogModal
      className="w-[500px]"
      title={
        connectionType === 'experiments'
          ? t(`experiments-connected`)
          : t(`operations-connected`)
      }
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex-center flex-col w-full items-start px-5 py-8 gap-y-8">
        <div className="flex-center w-full">
          <Icon
            icon={
              goal?.connections?.type === 'experiments'
                ? IconExperimentsConnected
                : IconOperationsConnected
            }
            size={'fit'}
          />
        </div>

        <div className="flex-center flex-col w-full gap-y-5">
          <div className="flex-center w-full text-center px-[67px] text-gray-700">
            <Trans
              i18nKey="goal-connected-desc"
              values={{ type: goal?.connections?.type }}
            />
          </div>
          <div className="flex flex-col w-full p-4 gap-y-5 rounded bg-gray-100">
            {goal?.connections?.data.map((item, index) => (
              <div
                key={index}
                className="flex items-center gap-x-2 typo-para-medium leading-4 text-primary-500"
              >
                <p>{index + 1}.</p>
                <p className="underline">{item.name}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      <ButtonBar
        primaryButton={<Button onClick={onClose}>{t(`close`)}</Button>}
      />
    </DialogModal>
  );
};

export default ConnectionsModal;
