import { Dialog } from '@headlessui/react';
import { Fragment, FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Modal } from '../Modal';

interface SegmentUploadingDialogProps {
  open: boolean;
  onClose: () => void;
}

export const SegmentUploadingDialog: FC<SegmentUploadingDialogProps> = ({
  open,
  onClose,
}) => {
  const { formatMessage: f } = useIntl();

  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-700"
      >
        {f(messages.segment.uploading.title)}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-gray-600">
          {f(messages.segment.uploading.message)}
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
            {f(messages.close)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
