import { DomainEventEntityType, DomainEventType } from '@types';

export const getActionText = (type: DomainEventType) => {
  switch (type.split('_').at(-1)) {
    case 'CREATED':
      return 'created';
    case 'DELETED':
      return 'deleted';
    case 'CLONED':
      return 'cloned';
    case 'ARCHIVED':
      return 'archived';
    case 'UNARCHIVED':
      return 'unarchived';
    case 'ENABLED':
      return 'enabled';
    case 'DISABLED':
      return 'disabled';
    case 'UPDATED':
    default:
      return 'updated';
  }
};

export const getEntityTypeText = (entityType: DomainEventEntityType) => {
  switch (entityType) {
    case 'APIKEY':
      return 'api key';
    case 'AUTOOPS_RULE':
      return 'operations rule';
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
