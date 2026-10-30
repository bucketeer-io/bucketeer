import { APIKey } from '@types';
import { truncateTextCenter } from 'utils/converts';

export const API_KEY_LABEL_ELLIPSIS = '…';

// Label shown for a key in listings. Keys carry a short stored prefix, which is
// shown followed by an ellipsis. Keys created before the prefix was stored have
// none yet and fall back to the masked key returned by the API.
export const formatAPIKeyLabel = (
  apiKey: Pick<APIKey, 'apiKey' | 'apiKeyPrefix'>
): string => {
  const prefix = apiKey.apiKeyPrefix?.trim();
  if (prefix) {
    return `${prefix}${API_KEY_LABEL_ELLIPSIS}`;
  }
  return truncateTextCenter(apiKey.apiKey ?? '');
};
