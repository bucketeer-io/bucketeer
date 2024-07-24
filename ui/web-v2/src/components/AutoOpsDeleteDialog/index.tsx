import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { OperationType, SelectedOperation } from '../FeatureAutoOpsRulesForm';
import { Modal } from '../Modal';

interface AutoOpsDeleteDialogProps {
  open: boolean;
  onConfirm: () => void;
  onClose: () => void;
  selectedOperation: SelectedOperation;
}

export const AutoOpsDeleteDialog: FC<AutoOpsDeleteDialogProps> = ({
  open,
  onConfirm,
  onClose,
  selectedOperation
}) => {
  const { formatMessage: f } = useIntl();

  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-700"
      >
        {selectedOperation.type === OperationType.SCHEDULE &&
          f(messages.autoOps.dialog.deleteScheduleTitle)}
        {selectedOperation.type === OperationType.EVENT_RATE &&
          f(messages.autoOps.dialog.deleteKillSwitchTitle)}
        {selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT &&
          f(messages.autoOps.dialog.deleteProgressiveRolloutTitle)}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-red-500">
          {selectedOperation.type === OperationType.SCHEDULE &&
            f(messages.autoOps.dialog.deleteScheduleDescription)}
          {selectedOperation.type === OperationType.EVENT_RATE &&
            f(messages.autoOps.dialog.deleteKillSwitchDescription)}
          {selectedOperation.type === OperationType.PROGRESSIVE_ROLLOUT &&
            f(messages.autoOps.dialog.deleteProgressiveRolloutDescription)}
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
            {f(messages.button.delete)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
