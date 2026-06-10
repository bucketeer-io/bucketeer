import { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { IconCancelFilled } from 'react-icons-material-design';
import { cn } from 'utils/style';
import Button from 'components/button';
import Icon from 'components/icon';

interface DrawerHeaderProps {
  title: string;
  onClose?: () => void;
}
const DrawerHeader = ({ title, onClose }: DrawerHeaderProps) => {
  return (
    <div className="flex items-center justify-between p-4 border-b border-gray-200">
      <h2 className="text-lg font-medium">{title}</h2>
      {onClose && (
        <Button
          onClick={onClose}
          className="text-gray-500 hover:text-gray-700 focus:outline-none"
        >
          <Icon icon={IconCancelFilled} />
        </Button>
      )}
    </div>
  );
};

const DrawerBody = ({ children }: { children: React.ReactNode }) => {
  return <div className="p-4 flex-1 overflow-y-auto">{children}</div>;
};

const DrawerFooter = ({ children }: { children: React.ReactNode }) => {
  return <div className="p-4 border-t border-gray-200">{children}</div>;
};

type DrawerSide = 'left' | 'right' | 'top' | 'bottom';

interface DrawerProps {
  open: boolean;
  onClose: () => void;
  side?: DrawerSide;
  className?: string;
  children: React.ReactNode;
}

export const Drawer = ({
  open,
  onClose,
  side = 'left',
  className,
  children
}: DrawerProps) => {
  const drawerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
    };
  }, [open, onClose]);

  useEffect(() => {
    if (!open) return;
    const previous = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = previous;
    };
  }, [open]);

  if (!open) return null;

  return createPortal(
    <div
      role="dialog"
      aria-modal="true"
      className={cn('fixed inset-0 z-[200]', className)}
      style={{ pointerEvents: 'auto' }}
    >
      <div className="absolute inset-0 bg-overlay" onClick={onClose} />
      <div
        ref={drawerRef}
        className={cn(
          'absolute bg-white shadow-card',
          side === 'right' && 'right-0 top-0 h-full w-[248px]',
          side === 'left' && 'left-0 top-0 h-full w-[248px]',
          side === 'top' && 'top-0 left-0 w-full',
          side === 'bottom' &&
            'bottom-0 left-0 w-full rounded-t-2xl max-h-[90vh] overflow-y-auto'
        )}
      >
        {children}
      </div>
    </div>,
    document.body
  );
};

Drawer.Header = DrawerHeader;
Drawer.Body = DrawerBody;
Drawer.Footer = DrawerFooter;
export default Drawer;
