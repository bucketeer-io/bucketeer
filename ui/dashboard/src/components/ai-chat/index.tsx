import { memo, useCallback, useState } from 'react';
import * as PopoverPrimitive from '@radix-ui/react-popover';
import { useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Suggestion } from '@types';
import { cn } from 'utils/style';
import ChatPopoverContainer from './chat-popover-container';

const EMPTY_SUGGESTIONS: Suggestion[] = [];

const ChatWidget = memo(() => {
  const { t } = useTranslation(['ai-chat']);
  const { isLogin } = useAuth();
  const [open, setOpen] = useState(false);

  const handleClose = useCallback(() => {
    setOpen(false);
  }, []);

  if (!isLogin) return null;

  return (
    <PopoverPrimitive.Root open={open} onOpenChange={setOpen}>
      <PopoverPrimitive.Trigger asChild>
        <button
          type="button"
          className={cn(
            'fixed bottom-6 right-6 z-50',
            'flex h-12 w-12 items-center justify-center',
            'rounded-full bg-primary-500 text-white shadow-lg',
            'transition-all duration-200 hover:bg-primary-700 hover:shadow-xl',
            'motion-reduce:transition-none',
            'focus:outline-none focus:ring-2 focus:ring-primary-300 focus:ring-offset-2'
          )}
          aria-label={t('ai-chat:open-chat')}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
        </button>
      </PopoverPrimitive.Trigger>

      {open && (
        <PopoverPrimitive.Portal>
          <PopoverPrimitive.Content
            side="top"
            align="end"
            sideOffset={12}
            className="z-50 animate-fade motion-reduce:animate-none"
            aria-label={t('ai-chat:title')}
            onOpenAutoFocus={e => {
              e.preventDefault();
              setTimeout(() => {
                document
                  .querySelector<HTMLInputElement>('[data-chat-input]')
                  ?.focus();
              }, 0);
            }}
          >
            <ChatPopoverContainer
              suggestions={EMPTY_SUGGESTIONS}
              onClose={handleClose}
            />
          </PopoverPrimitive.Content>
        </PopoverPrimitive.Portal>
      )}
    </PopoverPrimitive.Root>
  );
});

ChatWidget.displayName = 'ChatWidget';

export default ChatWidget;
