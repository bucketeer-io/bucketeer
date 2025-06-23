import { DomainEventEntityType, DomainEventType } from '@types';

export const getActionText = (
  type: DomainEventType,
  isLanguageJapanese: boolean
) => {
  switch (type.split('_').at(-1)) {
    case 'CREATED':
      return isLanguageJapanese ? '作成' : 'created';
    case 'DELETED':
      return isLanguageJapanese ? '削除' : 'deleted';
    case 'CLONED':
      return isLanguageJapanese ? '複製' : 'cloned';
    case 'ARCHIVED':
      return isLanguageJapanese ? 'アーカイブ' : 'archived';
    case 'UNARCHIVED':
      return isLanguageJapanese ? 'アンアーカイブ' : 'unarchived';
    case 'ENABLED':
      return isLanguageJapanese ? '有効' : 'enabled';
    case 'DISABLED':
      return isLanguageJapanese ? '無効' : 'disabled';
    case 'UPDATED':
    default:
      return isLanguageJapanese ? '更新' : 'updated';
  }
};

export const convertJSONToRender = (json: string | null) => {
  if (!json) return null;
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
  try {
    const parsedJSON = JSON.parse(json) || {};
    return JSON.stringify(parsedJSON, null, 4);
  } catch {
    return null;
  }
};

export const truncNumber = (num: number) => Math.trunc(num);

export const getPathName = (
  id: string | undefined,
  entityType: DomainEventEntityType
) => {
  if (!id) return null;
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
