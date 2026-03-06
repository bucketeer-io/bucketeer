import { useCallback, useEffect, useRef, useState } from 'react';
import { chatStreamer, CHAT_ERROR, type ChatErrorCode } from '@api/ai-chat';
import { getCurrentEnvIdStorage } from 'storage/environment';
import { AIChatMessage, PageContext } from '@types';

export { type ChatErrorCode } from '@api/ai-chat';

export const MAX_INPUT_LENGTH = 2000;
const MAX_HISTORY_MESSAGES = 50;

interface UseSSEChatOptions {
  pageContext: PageContext;
}

interface UseSSEChatReturn {
  messages: AIChatMessage[];
  isStreaming: boolean;
  errorKey: ChatErrorCode | null;
  sendMessage: (content: string) => void;
  clearMessages: () => void;
}

export const useSSEChat = ({
  pageContext
}: UseSSEChatOptions): UseSSEChatReturn => {
  const [messages, setMessages] = useState<AIChatMessage[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [errorKey, setErrorKey] = useState<ChatErrorCode | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const isStreamingRef = useRef(false);
  const pageContextRef = useRef(pageContext);
  const messagesRef = useRef(messages);

  // Sync refs on every render
  pageContextRef.current = pageContext;
  isStreamingRef.current = isStreaming;
  messagesRef.current = messages;

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  const sendMessage = useCallback(async (content: string) => {
    if (isStreamingRef.current) return;

    const trimmed = content.slice(0, MAX_INPUT_LENGTH);
    const userMessage: AIChatMessage = {
      id: crypto.randomUUID(),
      role: 'user',
      content: trimmed
    };
    const assistantMessage: AIChatMessage = {
      id: crypto.randomUUID(),
      role: 'assistant',
      content: ''
    };

    const currentMessages = [
      ...messagesRef.current.filter(m => m.content),
      userMessage
    ];

    setMessages(prev => {
      const updated = [...prev, userMessage, assistantMessage];
      return updated.slice(-MAX_HISTORY_MESSAGES);
    });
    isStreamingRef.current = true;
    setIsStreaming(true);
    setErrorKey(null);

    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    try {
      await chatStreamer(
        {
          messages: currentMessages.map(m => ({
            role: m.role,
            content: m.content
          })),
          pageContext: pageContextRef.current,
          environmentId: getCurrentEnvIdStorage() || ''
        },
        chunk => {
          if (abortController.signal.aborted) return;
          if (chunk.error) {
            setErrorKey(chunk.error as ChatErrorCode);
            return;
          }
          if (chunk.content) {
            setMessages(prev => {
              const updated = [...prev];
              const last = updated[updated.length - 1];
              if (last?.role === 'assistant') {
                updated[updated.length - 1] = {
                  ...last,
                  content: last.content + chunk.content
                };
              }
              return updated;
            });
          }
        },
        abortController.signal
      );
    } catch (err) {
      if ((err as Error).name !== 'AbortError') {
        const msg = (err as Error).message;
        const knownErrors = Object.values(CHAT_ERROR) as string[];
        setErrorKey(
          knownErrors.includes(msg)
            ? (msg as ChatErrorCode)
            : CHAT_ERROR.UNKNOWN
        );
      }
    } finally {
      isStreamingRef.current = false;
      setIsStreaming(false);
      abortControllerRef.current = null;
    }
  }, []);

  const clearMessages = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    setMessages([]);
    setErrorKey(null);
    setIsStreaming(false);
  }, []);

  return { messages, isStreaming, errorKey, sendMessage, clearMessages };
};
