import { Dialog, Transition } from '@headlessui/react';
import { Fragment, FC } from 'react';
import { useIntl } from 'react-intl';

import { Modal } from '../Modal';

interface ConfirmDialogProps {
  open: boolean;
  onConfirm: () => void;
  onClose: () => void;
  title: string;
  description: string;
  onCloseButton: string;
  onConfirmButton: string;
}

export const ConfirmDialog: FC<ConfirmDialogProps> = ({
  open,
  onConfirm,
  onClose,
  title,
  description,
  onCloseButton,
  onConfirmButton,
}) => {
  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-700"
      >
        {title}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-gray-700">{description}</p>
      </div>
      <div className="pt-5">
        <div className="flex justify-end">
          <button
            type="button"
            className="btn-cancel mr-3"
            disabled={false}
            onClick={onClose}
          >
            {onCloseButton}
          </button>
          <button
            type="button"
            className="btn-submit"
            disabled={false}
            onClick={onConfirm}
          >
            {onConfirmButton}
          </button>
        </div>
      </div>
    </Modal>
  );
};
