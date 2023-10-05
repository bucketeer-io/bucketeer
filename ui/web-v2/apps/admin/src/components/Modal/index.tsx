import { Dialog, Transition } from '@headlessui/react';
import { Fragment, FC, useCallback } from 'react';

import { classNames } from '../../utils/css';

interface ModalProps {
  open: boolean;
  onClose: () => void;
  overflowVisible?: boolean;
}

export const Modal: FC<ModalProps> = ({
  open,
  onClose,
  children,
  overflowVisible,
}) => {
  const handleClose = useCallback((): void => {
    onClose();
  }, [onClose]);

  return (
    <Transition appear show={open} as={Fragment}>
      <Dialog
        as="div"
        className="fixed inset-0 z-10 overflow-y-auto"
        onClose={handleClose}
      >
        <div className="min-h-screen px-4 text-center">
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Dialog.Overlay
              className={classNames(
                'absolute inset-0 bg-gray-500 bg-opacity-75 transition-opacity'
              )}
            />
          </Transition.Child>

          {/* This element is to trick the browser into centering the modal contents. */}
          <span
            className="inline-block h-screen align-middle"
            aria-hidden="true"
          >
            &#8203;
          </span>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300"
            enterFrom="opacity-0 scale-95"
            enterTo="opacity-100 scale-100"
            leave="ease-in duration-200"
            leaveFrom="opacity-100 scale-100"
            leaveTo="opacity-0 scale-95"
          >
            <div
              className={classNames(
                'inline-block w-full max-w-md p-6 my-8',
                'text-left align-middle',
                'transition-all transform',
                'bg-white shadow-xl rounded-2xl',
                overflowVisible ? 'overflow-visible' : 'overflow-hidden'
              )}
            >
              {children}
            </div>
          </Transition.Child>
        </div>
      </Dialog>
    </Transition>
  );
};
