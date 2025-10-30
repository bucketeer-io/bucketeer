import { ReactNode, useCallback } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Dialog from '@radix-ui/react-dialog';
import { cn } from 'utils/style';
import Button from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';

export type SliderType = 'create' | 'edit';

export type SliderProps = {
  direction?: 'slide-up' | 'slide-left';
  title: string;
  isOpen: boolean;
  children?: ReactNode;
  shouldCloseOnOverlayClick?: boolean;
  onClose: () => void;
};

const SlideModal = ({
  direction = 'slide-left',
  title,
  isOpen,
  onClose,
  children,
  shouldCloseOnOverlayClick = true
}: SliderProps) => {
  const onOpenChange = useCallback(
    (v: boolean) => {
      if (v === false && shouldCloseOnOverlayClick) onClose();
    },
    [shouldCloseOnOverlayClick]
  );

  return (
    <Dialog.Root open={isOpen} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 h-full w-full animate-fade bg-overlay" />
        <Dialog.Content>
          <div
            className={cn(
              'fixed z-50 flex h-full w-full flex-col rounded-l-lg bg-gray-50 max-w-[542px]',
              direction === 'slide-left' && 'right-0 top-0 animate-slide-left',
              direction === 'slide-up' && 'bottom-0 left-0 animate-slide-up'
            )}
          >
            <div className="z-10 flex-initial shadow-header">
              <div
                className={cn('flex items-center justify-between px-4 py-3.5')}
              >
                <Dialog.Title className="typo-head-bold-huge">
                  {title}
                </Dialog.Title>
                <Dialog.Description className="hidden" />
                <Dialog.Close asChild>
                  <Button
                    type="button"
                    size="icon-sm"
                    variant="grey"
                    onClick={onClose}
                  >
                    <Icon icon={IconCloseRound} />
                  </Button>
                </Dialog.Close>
              </div>
            </div>
            <Divider />
            <div className="flex-1 overflow-auto">{children}</div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
};

export default SlideModal;
