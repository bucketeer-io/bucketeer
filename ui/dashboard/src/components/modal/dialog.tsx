import { ReactNode, useCallback } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Dialog from '@radix-ui/react-dialog';
import { cn } from 'utils/style';
import { Button } from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';

export type ModalSize = 'sm' | 'md';

const widthBySize: Record<ModalSize, number> = { sm: 496, md: 780 };

const DialogModal = ({
  size = 'sm',
  title,
  isOpen,
  onClose,
  closeOnPressEscape = true,
  closeOnClickOutside = true,
  children
}: {
  size?: ModalSize;
  title: string;
  isOpen: boolean;
  onClose: () => void;
  closeOnPressEscape?: boolean;
  closeOnClickOutside?: boolean;
  children: ReactNode;
}) => {
  const onOpenChange = useCallback((v: boolean) => {
    if (v === false) onClose();
  }, []);

  return (
    <Dialog.Root open={isOpen} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay
          className={cn(
            'fixed inset-0 grid h-full w-full animate-fade',
            'place-items-center overflow-y-auto bg-overlay'
          )}
        >
          <Dialog.Content
            className="relative mx-4 my-8 animate-zoom rounded-lg bg-gray-50"
            style={{
              maxWidth: widthBySize[size]
            }}
            onEscapeKeyDown={
              closeOnPressEscape ? undefined : event => event.preventDefault()
            }
            onPointerDownOutside={
              closeOnClickOutside ? undefined : event => event.preventDefault()
            }
          >
            <div className="z-10 flex-initial shadow-header">
              <div
                className={cn('flex items-center justify-between px-4 py-3.5')}
              >
                <Dialog.Title className="typo-head-bold-huge">
                  {title}
                </Dialog.Title>
                <Dialog.Close asChild>
                  <Button size="icon-sm" variant="grey" onClick={onClose}>
                    <Icon icon={IconCloseRound} />
                  </Button>
                </Dialog.Close>
              </div>
            </div>
            <Divider />
            {children}
          </Dialog.Content>
        </Dialog.Overlay>
      </Dialog.Portal>
    </Dialog.Root>
  );
};

export default DialogModal;
