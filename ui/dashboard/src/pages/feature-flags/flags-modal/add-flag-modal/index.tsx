import { useTranslation } from 'i18n';
import SlideModal from 'components/modal/slide';
import CreateFlagForm from './create-flag-form';

interface AddFlagModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddFlagModal = ({ isOpen, onClose }: AddFlagModalProps) => {
  const { t } = useTranslation(['common', 'form']);

  return (
    <SlideModal title={t('new-flag')} isOpen={isOpen} onClose={onClose}>
      <CreateFlagForm onClose={onClose} />
    </SlideModal>
  );
};

export default AddFlagModal;
