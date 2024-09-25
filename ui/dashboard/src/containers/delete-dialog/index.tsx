import { PropsWithChildren } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconTrashSpecial } from '@icons';
import { Button } from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

type Props = PropsWithChildren & {
  title: string;
  isOpen: boolean;
  className?: string;
  onClose: () => void;
  onSubmit: () => void;
};

const DeleteMemberDialog = ({
  title,
  isOpen,
  children,
  className,
  onClose,
  onSubmit
}: Props) => {
  const { t } = useTranslation(['common']);
  return (
    <DialogModal title={title} isOpen={isOpen} onClose={onClose}>
      <div className="flex-center flex-col flex-1 px-5 py-8 gap-y-8">
        <div className="flex-center size-20">
          <Icon icon={IconTrashSpecial} size={'fit'} />
        </div>
        <div className={cn('flex-center px-[41.5px]', className)}>
          {children}
        </div>
      </div>
      <ButtonBar
        primaryButton={
          <Button variant="secondary" onClick={onClose}>
            {t('cancel')}
          </Button>
        }
        secondaryButton={
          <Button variant={'negative'} onClick={onSubmit}>
            {t('delete')}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default DeleteMemberDialog;
