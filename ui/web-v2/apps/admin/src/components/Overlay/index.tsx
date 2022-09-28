import { Dialog, Transition } from '@headlessui/react';
import React, { Fragment, FC, memo } from 'react';

import { classNames } from '../../utils/css';

export interface OverlayProps {
  open: boolean;
  onClose: () => void;
}

export const Overlay: FC<OverlayProps> = memo(({ onClose, open, children }) => {
  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog
        as="div"
        static
        className="z-30 fixed inset-0 overflow-hidden"
        open={open}
        onClose={onClose}
      >
        <div className="absolute inset-0 overflow-hidden">
          <Transition.Child
            as={Fragment}
            enter="ease-in-out duration-500"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="ease-in-out duration-500"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <Dialog.Overlay
              className={classNames(
                'absolute inset-0',
                'bg-gray-500 bg-opacity-75 transition-opacity'
              )}
            />
          </Transition.Child>
          <div className="fixed inset-y-0 max-w-full right-0 flex">
            <Transition.Child
              as={Fragment}
              enter="transform transition ease-in-out duration-300"
              enterFrom="translate-x-full"
              enterTo="translate-x-0"
              leave="transform transition ease-in-out duration-300"
              leaveFrom="translate-x-0"
              leaveTo="translate-x-full"
            >
              <div className="max-w-3xl">
                <div className="h-full bg-white shadow-xl overflow-y-scroll">
                  {children}
                </div>
              </div>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  );
});
