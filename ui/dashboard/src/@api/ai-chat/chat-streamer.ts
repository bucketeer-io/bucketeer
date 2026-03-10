import { urls } from 'configs';
import { getTokenStorage } from 'storage/token';
import { AIChatMessage, PageContext } from '@types';

const MAX_BUFFER_SIZE = 1024 * 1024; // 1MB

export interface ChatStreamerParams {
  messages: Pick<AIChatMessage, 'role' | 'content'>[];
  pageContext: PageContext;
  environmentId: string;
}

export interface ChatStreamChunk {
  content?: string;
  error?: string;
  done: boolean;
}

/** Error codes used as i18n keys in ai-chat:errors.{code} */
export const CHAT_ERROR = {
  NOT_AUTHENTICATED: 'not-authenticated',
  RATE_LIMIT: 'rate-limit',
  NETWORK: 'network',
  UNKNOWN: 'unknown'
} as const;

export type ChatErrorCode = (typeof CHAT_ERROR)[keyof typeof CHAT_ERROR];

interface SSEErrorPayload {
  error: string;
}

interface SSEContentPayload {
  content: string;
}

function isSSEError(v: unknown): v is SSEErrorPayload {
  return (
    typeof v === 'object' &&
    v !== null &&
    'error' in v &&
    typeof (v as Record<string, unknown>).error === 'string'
  );
}

function isSSEContent(v: unknown): v is SSEContentPayload {
  return (
    typeof v === 'object' &&
    v !== null &&
    'content' in v &&
    typeof (v as Record<string, unknown>).content === 'string'
  );
}

/**
 * Process a single SSE data payload and invoke the appropriate callback.
 * Returns 'done' if the stream should be terminated, 'continue' otherwise.
 */
function processSSEData(
  data: string,
  onChunk: (chunk: ChatStreamChunk) => void
): 'done' | 'continue' {
  if (data === '[DONE]') {
    onChunk({ done: true });
    return 'done';
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(data);
  } catch {
    return 'continue';
  }

  if (isSSEError(parsed)) {
    onChunk({ error: CHAT_ERROR.UNKNOWN, done: true });
    return 'done';
  }

  if (isSSEContent(parsed) && parsed.content) {
    onChunk({ content: parsed.content, done: false });
  }

  return 'continue';
}

/**
 * Streams chat responses via SSE using native fetch.
 * Native fetch is required instead of axiosClient because axios does not
 * support ReadableStream for server-sent events.
 */
export const chatStreamer = async (
  params: ChatStreamerParams,
  onChunk: (chunk: ChatStreamChunk) => void,
  signal?: AbortSignal
): Promise<void> => {
  const authToken = getTokenStorage();
  if (!authToken) {
    throw new Error(CHAT_ERROR.NOT_AUTHENTICATED);
  }

  const baseUrl = urls.WEB_API_ENDPOINT || '';
  let response: Response;
  try {
    response = await fetch(`${baseUrl}/v1/aichat/chat`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${authToken.accessToken}`
      },
      body: JSON.stringify({
        messages: params.messages,
        pageContext: params.pageContext,
        environmentId: params.environmentId
      }),
      signal
    });
  } catch {
    throw new Error(CHAT_ERROR.NETWORK);
  }

  if (!response.ok) {
    if (response.status === 401) {
      document.dispatchEvent(
        new CustomEvent('unauthenticated', { bubbles: true })
      );
      throw new Error(CHAT_ERROR.NOT_AUTHENTICATED);
    }
    if (response.status === 429) {
      throw new Error(CHAT_ERROR.RATE_LIMIT);
    }
    throw new Error(CHAT_ERROR.NETWORK);
  }

  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error(CHAT_ERROR.NETWORK);
  }

  const decoder = new TextDecoder();
  let buffer = '';

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      if (buffer.length > MAX_BUFFER_SIZE) {
        throw new Error(CHAT_ERROR.NETWORK);
      }

      const lines = buffer.split('\n');
      buffer = lines.pop() || '';

      for (const line of lines) {
        const trimmedLine = line.replace(/\r$/, '');
        if (!trimmedLine.startsWith('data:')) continue;
        const data = trimmedLine.startsWith('data: ')
          ? trimmedLine.slice(6)
          : trimmedLine.slice(5);

        if (processSSEData(data, onChunk) === 'done') return;
      }
    }

    // Flush remaining bytes from TextDecoder
    buffer += decoder.decode();

    // Process remaining buffer
    const remaining = buffer.trim();
    if (remaining.startsWith('data:')) {
      const data = remaining.startsWith('data: ')
        ? remaining.slice(6).trim()
        : remaining.slice(5).trim();
      if (data) {
        processSSEData(data, onChunk);
      }
    }

    // Always send done if [DONE] was not received
    onChunk({ done: true });
  } finally {
    reader.releaseLock();
  }
};
