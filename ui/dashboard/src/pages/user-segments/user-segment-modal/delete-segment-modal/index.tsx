import { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { IconDelete } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import DialogModal from 'components/modal/dialog';
import SegmentWarning from '../edit-segment-modal/segment-warning';

export type DeleteUserSegmentProps = {
  onSubmit: () => void;
  isOpen: boolean;
  onClose: () => void;
  userSegment: UserSegment;
  loading: boolean;
};

const DeleteUserSegmentModal = ({
  onSubmit,
  isOpen,
  onClose,
  userSegment,
  loading
}: DeleteUserSegmentProps) => {
  const { t } = useTranslation(['common']);

  const isInUseSegment = useMemo(
    () => userSegment.isInUseStatus || userSegment.features.length > 0,
    [userSegment]
  );

  return (
    <DialogModal
      className="w-[500px]"
      title={t(`delete-user-segment`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
        {isInUseSegment ? (
          <SegmentWarning features={userSegment.features} className="mt-0" />
        ) : (
          <>
            <IconDelete />
            <div className="typo-para-big text-gray-700 text-center">
              <Trans
                i18nKey="table:user-segment.delete-user-segment-desc"
                values={{ name: userSegment.name }}
                components={{ bold: <strong /> }}
              />
            </div>
          </>
        )}
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            variant="negative"
            loading={loading}
            disabled={isInUseSegment || loading}
            className="w-24"
            onClick={onSubmit}
          >
            {t(`delete`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default DeleteUserSegmentModal;
