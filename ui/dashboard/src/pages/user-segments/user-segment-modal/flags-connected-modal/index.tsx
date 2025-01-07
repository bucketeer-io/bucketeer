import { useTranslation } from 'i18n';
import { IconFlagConnected } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type DeleteMemberProps = {
  isOpen: boolean;
  onClose: () => void;
};

const FlagsConnectedModal = ({ isOpen, onClose }: DeleteMemberProps) => {
  const { t } = useTranslation(['common']);

  return (
    <DialogModal
      className="w-[496px]"
      title={t(`flags-connected`)}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="py-8 px-5 flex-center flex-col gap-6">
        <Icon icon={IconFlagConnected} size={'fit'} />
        <div className="typo-para-big text-gray-700 text-center">
          {t('flags-connected-desc')}
        </div>
      </div>

      <ButtonBar
        primaryButton={
          <Button onClick={onClose} variant="primary">
            {t(`close`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default FlagsConnectedModal;
