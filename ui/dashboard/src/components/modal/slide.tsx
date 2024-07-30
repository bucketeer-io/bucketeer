import { ReactNode, useCallback } from 'react';
import { IconCloseRound } from 'react-icons-material-design';
import * as Dialog from '@radix-ui/react-dialog';
import { cn } from 'utils/style';
import { Button } from 'components/button';
import Icon from 'components/icon';

export type SlideModalHeaderSize = 'xs' | 'sm' | 'md';

const headerSizes = {
  xs: 'md:py-2',
  sm: 'md:py-4',
  md: 'md:py-4 lg:py-6'
};

const SlideModal = ({
  direction = 'slide-left',
  title,
  isOpen,
  onClose,
  children,
  header = 'md',
  shouldCloseOnOverlayClick = true
}: {
  direction?: 'slide-up' | 'slide-left';
  title: string;
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
  header?: SlideModalHeaderSize;
  shouldCloseOnOverlayClick?: boolean;
}) => {
  const onOpenChange = useCallback((v: boolean) => {
    if (v === false && shouldCloseOnOverlayClick) onClose();
  }, []);

  return (
    <Dialog.Root open={isOpen} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 h-full w-full animate-fade bg-overlay" />
        <Dialog.Content>
          <div
            className={cn(
              'fixed flex h-full w-full flex-col overflow-y-hidden bg-light-50 lg:w-2/3 lg:max-w-[808px]',
              direction === 'slide-left' && 'right-0 top-0 animate-slide-left',
              direction === 'slide-up' && 'bottom-0 left-0 animate-slide-up'
            )}
          >
            <div className="z-10 flex-initial px-4 shadow-header">
              <div
                className={cn(
                  'flex items-center justify-between px-2 py-3',
                  headerSizes[header]
                )}
              >
                <h2 className="typo-label-big">{title}</h2>
                <Button size="icon-sm" variant="grey" onClick={onClose}>
                  <Icon icon={IconCloseRound} />
                </Button>
              </div>
            </div>

            <div className="flex-1 overflow-hidden">{children}</div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
};

export default SlideModal;
