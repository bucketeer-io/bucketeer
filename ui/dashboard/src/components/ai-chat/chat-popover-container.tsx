import { useEffect, useMemo, useState } from 'react';
import { Suggestion } from '@types';
import { useSSEChat } from './use-sse-chat';
import { usePageContext } from './use-page-context';
import ChatPopover from './chat-popover';

interface ChatPopoverContainerProps {
  suggestions: Suggestion[];
  onClose: () => void;
}

const ChatPopoverContainer = ({
  suggestions,
  onClose
}: ChatPopoverContainerProps) => {
  const pageContext = usePageContext();

  const isFeaturePage =
    pageContext.pageType === 'feature_flags' ||
    pageContext.pageType === 'targeting';

  const [selectedFlagId, setSelectedFlagId] = useState<string | undefined>(
    pageContext.featureId
  );

  // Sync with URL-derived featureId when navigating between targeting pages
  useEffect(() => {
    setSelectedFlagId(pageContext.featureId);
  }, [pageContext.featureId]);

  const effectivePageContext = useMemo(
    () =>
      isFeaturePage
        ? { ...pageContext, featureId: selectedFlagId }
        : pageContext,
    [pageContext, selectedFlagId, isFeaturePage]
  );

  const { messages, isStreaming, errorKey, sendMessage, clearMessages } =
    useSSEChat({ pageContext: effectivePageContext });

  return (
    <ChatPopover
      messages={messages}
      suggestions={suggestions}
      isStreaming={isStreaming}
      errorKey={errorKey}
      selectedFlagId={selectedFlagId}
      showFlagSelector={isFeaturePage}
      onSelectFlag={setSelectedFlagId}
      onSend={sendMessage}
      onClear={clearMessages}
      onClose={onClose}
    />
  );
};

export default ChatPopoverContainer;
