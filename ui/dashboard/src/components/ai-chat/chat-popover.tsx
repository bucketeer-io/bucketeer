import { FormEvent, memo, useCallback, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'i18n';
import { IconClose } from '@icons';
import Icon from 'components/icon';
import { ScrollArea } from 'components/scroll-area';
import Button from 'components/button';
import { AIChatMessage, Suggestion } from '@types';
import type { ChatErrorCode } from '@api/ai-chat';
import { cn } from 'utils/style';
import { MAX_INPUT_LENGTH } from './use-sse-chat';
import SuggestionCard from './suggestion-card';
import PresetQuestions from './preset-questions';
import FlagSelector from './flag-selector';

interface ChatPopoverProps {
  messages: AIChatMessage[];
  suggestions: Suggestion[];
  isStreaming: boolean;
  errorKey: ChatErrorCode | null;
  selectedFlagId?: string;
  showFlagSelector?: boolean;
  onSelectFlag?: (flagId: string | undefined) => void;
  onSend: (content: string) => void;
  onClear: () => void;
  onClose: () => void;
}

const ChatPopover = memo(
  ({
    messages,
    suggestions,
    isStreaming,
    errorKey,
    selectedFlagId,
    showFlagSelector,
    onSelectFlag,
    onSend,
    onClear,
    onClose
  }: ChatPopoverProps) => {
    const { t } = useTranslation(['ai-chat']);
    const [input, setInput] = useState('');
    const bottomRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
      bottomRef.current?.scrollIntoView({
        behavior: isStreaming ? 'instant' : 'smooth'
      });
    }, [messages, isStreaming]);

    const handleSubmit = useCallback(
      (e: FormEvent) => {
        e.preventDefault();
        const trimmed = input.trim();
        if (!trimmed || isStreaming) return;
        onSend(trimmed);
        setInput('');
      },
      [input, isStreaming, onSend]
    );

    const handlePresetSelect = useCallback(
      (question: string) => {
        if (isStreaming) return;
        onSend(question);
      },
      [isStreaming, onSend]
    );

    const isEmpty = messages.length === 0;

    return (
      <div className="flex h-[500px] w-[380px] flex-col rounded-xl border border-gray-200 bg-white shadow-lg">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 px-4 py-3">
          <h3 className="typo-para-medium font-semibold text-gray-800">
            {t('ai-chat:title')}
          </h3>
          <div className="flex items-center gap-2">
            {messages.length > 0 && (
              <button
                type="button"
                className="rounded typo-para-tiny text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-300"
                onClick={onClear}
              >
                {t('ai-chat:clear')}
              </button>
            )}
            <button
              type="button"
              className="rounded text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-300"
              onClick={onClose}
              aria-label={t('ai-chat:close')}
            >
              <Icon icon={IconClose} size="sm" />
            </button>
          </div>
        </div>

        {/* Flag Selector */}
        {showFlagSelector && onSelectFlag && (
          <div className="border-b border-gray-200 px-4 py-2">
            <FlagSelector
              selectedFlagId={selectedFlagId}
              onSelectFlag={onSelectFlag}
            />
          </div>
        )}

        {/* Messages */}
        <ScrollArea className="flex-1 px-4 py-3">
          {isEmpty ? (
            <div className="flex flex-col gap-4">
              <p className="typo-para-small text-gray-500">
                {t('ai-chat:welcome')}
              </p>
              {suggestions.length > 0 && (
                <div className="flex flex-col gap-2">
                  <p className="typo-para-tiny font-medium text-gray-600">
                    {t('ai-chat:suggestions-title')}
                  </p>
                  {suggestions.map(s => (
                    <SuggestionCard
                      key={s.id}
                      suggestion={s}
                      onClick={handlePresetSelect}
                    />
                  ))}
                </div>
              )}
              <div className="flex flex-col gap-2">
                <p className="typo-para-tiny font-medium text-gray-600">
                  {t('ai-chat:quick-questions')}
                </p>
                <PresetQuestions onSelect={handlePresetSelect} />
              </div>
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              {messages.map(msg => (
                <div
                  key={msg.id}
                  className={cn(
                    'max-w-[90%] rounded-lg px-3 py-2',
                    msg.role === 'user'
                      ? 'ml-auto bg-primary-500 text-white'
                      : 'mr-auto bg-gray-100 text-gray-700'
                  )}
                >
                  <p className="typo-para-small whitespace-pre-wrap">
                    {msg.content ||
                      (isStreaming &&
                      msg.role === 'assistant' &&
                      msg === messages[messages.length - 1] ? (
                        <span className="animate-pulse motion-reduce:animate-none">
                          {t('ai-chat:thinking')}
                        </span>
                      ) : (
                        ''
                      ))}
                  </p>
                </div>
              ))}
              <div ref={bottomRef} />
            </div>
          )}
          {errorKey && (
            <div
              className="mt-2 rounded-lg bg-accent-red-50 px-3 py-2"
              role="alert"
            >
              <p className="typo-para-tiny text-accent-red-500">
                {t(`ai-chat:errors.${errorKey}`)}
              </p>
            </div>
          )}
        </ScrollArea>

        {/* Streaming status for screen readers */}
        <div aria-live="polite" className="sr-only">
          {isStreaming ? t('ai-chat:thinking') : ''}
        </div>

        {/* Input */}
        <form
          onSubmit={handleSubmit}
          className="flex items-center gap-2 border-t border-gray-200 px-4 py-3"
        >
          <input
            type="text"
            data-chat-input
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder={t('ai-chat:input-placeholder')}
            aria-label={t('ai-chat:input-placeholder')}
            disabled={isStreaming}
            maxLength={MAX_INPUT_LENGTH}
            className={cn(
              'flex-1 rounded-lg border border-gray-200 px-3 py-2',
              'typo-para-small text-gray-700 placeholder:text-gray-400',
              'focus:border-primary-300 focus:outline-none focus:ring-1 focus:ring-primary-300',
              'disabled:cursor-not-allowed disabled:bg-gray-50'
            )}
          />
          <Button
            type="submit"
            variant="primary"
            size="sm"
            disabled={!input.trim() || isStreaming}
            loading={isStreaming}
          >
            {t('ai-chat:send')}
          </Button>
        </form>
      </div>
    );
  }
);

ChatPopover.displayName = 'ChatPopover';

export default ChatPopover;
