import { ReactNode, useCallback } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Dialog from '@radix-ui/react-dialog';
import { cn } from 'utils/style';
import Button from 'components/button';
import Divider from 'components/divider';
import Icon from 'components/icon';

export type ModalSize = 'sm' | 'md';
export type ModalProps = {
  size?: ModalSize;
  title: string | ReactNode;
  closeContent?: ReactNode;
  isOpen: boolean;
  onClose: () => void;
  closeOnPressEscape?: boolean;
  closeOnClickOutside?: boolean;
  isShowHeader?: boolean;
  children?: ReactNode;
  className?: string;
  overlayCls?: string;
};

const DialogModal = ({
  closeContent,
  title,
  isOpen,
  onClose,
  closeOnPressEscape = true,
  closeOnClickOutside = true,
  isShowHeader = true,
  children,
  className,
  overlayCls
}: ModalProps) => {
  const onOpenChange = useCallback((v: boolean) => {
    if (v === false) onClose();
  }, []);

  return (
    <div className="w-full p-6">
      <Dialog.Root open={isOpen} onOpenChange={onOpenChange}>
        <Dialog.Portal>
          <Dialog.Overlay
            className={cn(
              'fixed inset-0 grid h-full w-full animate-fade z-50',
              'place-items-center overflow-y-auto bg-overlay',
              'p-6',
              overlayCls
            )}
          >
            <Dialog.Content
              className={cn(
                'relative w-full mx-4 my-8 animate-zoom rounded-lg bg-gray-50',
                className
              )}
              onEscapeKeyDown={
                closeOnPressEscape ? undefined : event => event.preventDefault()
              }
              onPointerDownOutside={
                closeOnClickOutside
                  ? undefined
                  : event => event.preventDefault()
              }
            >
              <div
                className={cn('z-10 flex-initial shadow-header', {
                  hidden: !isShowHeader
                })}
              >
                <div
                  className={cn(
                    'flex items-center justify-between px-4 py-3.5'
                  )}
                >
                  {title && (
                    <Dialog.Title className="typo-head-bold-small sm:typo-head-bold-huge">
                      {title}
                    </Dialog.Title>
                  )}
                  <Dialog.Description className="hidden" />
                  <Dialog.Close asChild>
                    {closeContent ?? (
                      <Button size="icon-sm" variant="grey" onClick={onClose}>
                        <Icon icon={IconCloseRound} />
                      </Button>
                    )}
                  </Dialog.Close>
                </div>
              </div>
              <Divider />
              {children}
            </Dialog.Content>
          </Dialog.Overlay>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
};

export default DialogModal;
