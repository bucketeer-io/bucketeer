export type APIKeyRole =
  | 'UNKNOWN'
  | 'SDK_CLIENT'
  | 'SDK_SERVER'
  | 'PUBLIC_API_READ_ONLY'
  | 'PUBLIC_API_WRITE'
  | 'PUBLIC_API_ADMIN';

export interface APIKey {
  id: string;
  name: string;
  role: APIKeyRole;
  disabled: true;
  createdAt: string;
  updatedAt: string;
}

export interface APIKeyCollection {
  apiKeys: Array<APIKey>;
  cursor: string;
  totalCount: string;
}
