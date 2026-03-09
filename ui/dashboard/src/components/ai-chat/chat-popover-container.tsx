import { useEffect, useMemo, useState } from 'react';
import { useQuerySuggestions } from '@queries/suggestions';
import { getCurrentEnvIdStorage } from 'storage/environment';
import ChatPopover from './chat-popover';
import { usePageContext } from './use-page-context';
import { useSSEChat } from './use-sse-chat';

interface ChatPopoverContainerProps {
  onClose: () => void;
}

const ChatPopoverContainer = ({ onClose }: ChatPopoverContainerProps) => {
  const pageContext = usePageContext();
  const environmentId = getCurrentEnvIdStorage() || '';

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

  const suggestionsParams = useMemo(
    () => ({ environmentId, pageContext: effectivePageContext }),
    [environmentId, effectivePageContext]
  );

  const { data } = useQuerySuggestions({
    params: suggestionsParams
  });

  const { messages, isStreaming, errorKey, sendMessage, clearMessages } =
    useSSEChat({ pageContext: effectivePageContext });

  return (
    <ChatPopover
      messages={messages}
      suggestions={data?.suggestions ?? []}
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
