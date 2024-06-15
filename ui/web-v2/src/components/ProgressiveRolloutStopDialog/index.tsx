import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Modal } from '../Modal';

interface ProgressiveRolloutStopDialogProps {
  open: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

export const ProgressiveRolloutStopDialog: FC<ProgressiveRolloutStopDialogProps> =
  ({ open, onConfirm, onClose }) => {
    const { formatMessage: f } = useIntl();

    return (
      <Modal open={open} onClose={onClose}>
        <Dialog.Title
          as="h3"
          className="text-lg font-medium leading-6 text-gray-700"
        >
          {f(messages.autoOps.stopProgressiveRollout)}
        </Dialog.Title>
        <div className="mt-2">
          <p className="text-sm text-red-500">
            {f(messages.autoOps.stopProgressiveRolloutDescription)}
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
