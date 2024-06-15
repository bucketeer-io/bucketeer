import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Modal } from '../Modal';

interface TriggerDeleteDialogProps {
  open: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

export const TriggerDeleteDialog: FC<TriggerDeleteDialogProps> = ({
  open,
  onConfirm,
  onClose,
}) => {
  const { formatMessage: f } = useIntl();

  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-700"
      >
        {f(messages.trigger.deleteTriggerDialogTitle)}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-red-500">
          {f(messages.trigger.deleteTriggerDialogMessage)}
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
            {f(messages.trigger.deleteTriggerDialogBtnLabel)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
