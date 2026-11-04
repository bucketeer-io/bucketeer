export type APIKeyRole =
  | 'UNKNOWN'
  | 'SDK_CLIENT'
  | 'SDK_SERVER'
  | 'PUBLIC_API_READ_ONLY'
  | 'PUBLIC_API_WRITE'
  | 'PUBLIC_API_ADMIN';

export interface APIKey {
  id: string;
  apiKey: string;
  // First characters of the key, kept for display. Empty for keys created
  // before the prefix was stored.
  apiKeyPrefix: string;
  name: string;
  role: APIKeyRole;
  disabled: true;
  description: string;
  environmentName: string;
  maintainer: string;
  createdAt: string;
  updatedAt: string;
  lastUsedAt: string;
}

export interface APIKeyCollection {
  apiKeys: Array<APIKey>;
  cursor: string;
  totalCount: string;
}
