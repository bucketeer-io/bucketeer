import { DomainEventEntityType, DomainEventType } from '@types';

export const getActionText = (
  type: DomainEventType,
  isLanguageJapanese: boolean
) => {
  switch (type.split('_').at(-1)) {
    case 'CREATED':
      return isLanguageJapanese ? '作成しました' : 'created';
    case 'DELETED':
      return isLanguageJapanese ? '削除しました' : 'deleted';
    case 'CLONED':
      return isLanguageJapanese ? '複製しました' : 'cloned';
    case 'ARCHIVED':
      return isLanguageJapanese ? 'アーカイブしました' : 'archived';
    case 'UNARCHIVED':
      return isLanguageJapanese ? 'アンアーカイブドしました' : 'unarchived';
    case 'ENABLED':
      return isLanguageJapanese ? '有効にしました' : 'enabled';
    case 'DISABLED':
      return isLanguageJapanese ? '無効にしました' : 'disabled';
    case 'UPDATED':
    default:
      return isLanguageJapanese ? '更新しました' : 'updated';
  }
};

export const getEntityTypeText = (entityType: DomainEventEntityType) => {
  switch (entityType) {
    case 'APIKEY':
      return 'api key';
    case 'AUTOOPS_RULE':
      return 'operation rule';
    case 'CODEREF':
      return 'code ref';
    case 'FEATURE':
      return 'feature flag';
    case 'ADMIN_ACCOUNT':
    case 'ADMIN_SUBSCRIPTION':
    case 'FLAG_TRIGGER':
    case 'PROGRESSIVE_ROLLOUT':
      return entityType.replace('_', ' ').toLowerCase();
    default:
      return entityType.toLowerCase();
  }
};

export const convertJSONToRender = (json: string) => {
  if (typeof json != 'string') {
    json = JSON.stringify(json, null, 4);
  }
  json = json
    .replace(
      /("(\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*"(?=\s*:))/g,
      '<span style="color: #e439ac">$1</span>'
    )
    .replace(
      /(:\s*)("(\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*")/g,
      '$1<span style="color: #40BF42">$2</span>'
    )
    .replace(/(:\s*)(\d+)/g, '$1<span style="color: #64748B">$2</span>'); // numeric values
  return json;
};

export const formatJSONWithIndent = (json: string) => {
  const parsedJSON = JSON.parse(json) || {};
  return JSON.stringify(parsedJSON, null, 4);
};

export const truncNumber = (num: number) => Math.trunc(num);

export const getPathName = (id: string, entityType: DomainEventEntityType) => {
  switch (entityType) {
    case 'APIKEY':
      return '/api-keys';
    case 'PROGRESSIVE_ROLLOUT':
    case 'AUTOOPS_RULE':
      return `/features/${id}/operations`;
    case 'FEATURE':
      return `/features/${id}/targeting`;
    case 'ADMIN_ACCOUNT':
    case 'ACCOUNT':
      return '/members';
    case 'FLAG_TRIGGER':
      return `/features/${id}/trigger`;
    case 'GOAL':
      return `/goals/${id}`;
    case 'PROJECT':
      return `/projects/${id}`;
    case 'ORGANIZATION':
      return `/organizations/${id}`;
    case 'PUSH':
      return `/pushes`;
    case 'EXPERIMENT':
      return `/experiments/${id}/results`;
    case 'SEGMENT':
      return '/segments';
    default:
      return null;
  }
};
