import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from '@/constants/routing';
import { useCurrentEnvironment } from '@/modules/me';
import { Dialog, Transition } from '@headlessui/react';
import {
  InformationCircleIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/solid';
import { Fragment, FC } from 'react';
import { useIntl } from 'react-intl';
import { Link } from 'react-router-dom';

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
  const currentEnvironment = useCurrentEnvironment();

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
          <div className="rounded-md bg-yellow-50 p-4 mt-2">
            <div className="flex">
              <div className="flex-shrink-0">
                <ExclamationCircleIcon
                  className="h-5 w-5 text-yellow-400"
                  aria-hidden="true"
                />
              </div>
              <div className="ml-3 flex-1">
                <p className="text-sm text-yellow-700">
                  {f(messages.segment.confirm.cannotDelete, {
                    segmentName: <strong>{`${segment.name}`}</strong>,
                    length: segment.featuresList.length,
                  })}
                </p>
                <div className="mt-2 text-sm text-yellow-700">
                  <ul className="list-disc space-y-1 pl-5">
                    {segment.featuresList.map((feature) => (
                      <li key={feature.id}>
                        <Link
                          className="link text-left"
                          to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}`}
                        >
                          <p className="truncate w-60">{feature.name}</p>
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-red-500">
            {f(messages.segment.confirm.deleteDescription, {
              segmentName: <strong>{segment && segment.name}</strong>,
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
