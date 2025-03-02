import { AppDispatch } from '../../store';
import { Dialog } from '@headlessui/react';
import { FC } from 'react';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { Modal } from '../Modal';
import RevertSvg from '../../assets/svg/revert.svg';
import { updateFeature } from '../../modules/features';

interface CancelScheduleDialogProps {
  open: boolean;
  onClose: () => void;
}

export const CancelScheduleDialog: FC<CancelScheduleDialogProps> = ({
  open,
  onClose
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();

  const handleCancel = () => {
    // dispatch(updateFeature({ id: 1, status: 'scheduled' }));
    onClose();
  };
  return (
    <Modal open={open} onClose={onClose}>
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-900"
      >
        Revert Updates
      </Dialog.Title>
      <div className="py-6 flex justify-center">
        <RevertSvg />
      </div>
      <div className="mt-2 text-center px-5">
        <p className="text-sm text-gray-500">
          The scheduled changes for <strong>01/31/2025 at 09:30</strong> are
          going to be reverted permanently. Are you sure you want to proceed?
        </p>
      </div>
      <div className="pt-10">
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
            className="btn bg-[#EB1726]"
            onClick={handleCancel}
          >
            Revert
          </button>
        </div>
      </div>
    </Modal>
  );
};
