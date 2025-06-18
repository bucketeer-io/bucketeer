import { i18n } from 'i18n';
import { APIKeyOption } from './types';

const translate = i18n.t;

export const apiKeyOptions: APIKeyOption[] = [
  {
    id: 'client-sdk',
    label: translate('form:api-key.client-sdk'),
    description: translate('form:api-key.client-sdk-desc'),
    value: 'SDK_CLIENT'
  },
  {
    id: 'server-sdk',
    label: translate('form:api-key.server-sdk'),
    description: translate('form:api-key.server-sdk-desc'),
    value: 'SDK_SERVER'
  },
  {
    id: 'public-api-read-only',
    label: translate('form:api-key.public-api-read-only'),
    description: translate('form:api-key.public-api-read-only-desc'),
    value: 'PUBLIC_API_READ_ONLY'
  },
  {
    id: 'public-api-write',
    label: translate('form:api-key.public-api-write'),
    description: translate('form:api-key.public-api-write-desc'),
    value: 'PUBLIC_API_WRITE'
  },
  {
    id: 'public-api-admin',
    label: translate('form:api-key.public-api-admin'),
    description: translate('form:api-key.public-api-admin-desc'),
    value: 'PUBLIC_API_ADMIN'
  }
];
