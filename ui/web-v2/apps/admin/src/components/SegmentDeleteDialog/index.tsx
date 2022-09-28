import { Dialog, Transition } from '@headlessui/react';
import { Fragment, FC } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Segment } from '../../proto/feature/segment_pb';
import { Modal } from '../Modal';

interface SegmentDeleteDialogProps {
  open: boolean;
  segment: Segment.AsObject;
  onConfirm: () => void;
  onClose: () => void;
}

export const SegmentDeleteDialog: FC<SegmentDeleteDialogProps> = ({
  open,
  segment,
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
        {f(messages.segment.confirm.deleteTitle)}
      </Dialog.Title>
      <div className="mt-2">
        {segment && segment.isInUseStatus ? (
          <p className="text-sm text-gray-700">
            {f(messages.segment.confirm.cannotDelete, {
              segmentName: `${segment.name}`,
            })}
          </p>
        ) : (
          <p className="text-sm text-red-500">
            {f(messages.segment.confirm.deleteDescription, {
              segmentName: `${segment && segment.name}`,
            })}
          </p>
        )}
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
          <button
            type="button"
            className="btn-submit"
            disabled={segment && segment.isInUseStatus}
            onClick={onConfirm}
          >
            {f(messages.button.submit)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
