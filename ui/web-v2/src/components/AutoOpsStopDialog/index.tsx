import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Modal } from '../Modal';
import { OperationType, SelectedOperation } from '../FeatureAutoOpsRulesForm';

interface AutoOpsStopDialogProps {
  selectedOperation: SelectedOperation;
  open: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

export const AutoOpsStopDialog: FC<AutoOpsStopDialogProps> = ({
  selectedOperation,
  open,
  onConfirm,
  onClose
}) => {
  const { formatMessage: f } = useIntl();

  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-700"
      >
        {selectedOperation.type === OperationType.SCHEDULE &&
          f(messages.autoOps.stopSchedule)}
        {selectedOperation.type === OperationType.EVENT_RATE &&
          f(messages.autoOps.killSwitch)}
        {selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT &&
          f(messages.autoOps.stopProgressiveRollout)}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-red-500">
          {selectedOperation.type === OperationType.SCHEDULE &&
            f(messages.autoOps.stopScheduleDescription)}
          {selectedOperation.type === OperationType.EVENT_RATE &&
            f(messages.autoOps.stopKillSwitchDescription)}
          {selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT &&
            f(messages.autoOps.stopProgressiveRolloutDescription)}
        </p>
      </div>
      <div className="pt-5">
        <div className="flex justify-end">
          <button
            type="button"
            className="btn-cancel mr-3"
            disabled={false}
            onClick={onClose}
          >
            {f(messages.button.cancel)}
          </button>
          <button type="button" className="btn-submit" onClick={onConfirm}>
            {f(messages.button.stop)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
