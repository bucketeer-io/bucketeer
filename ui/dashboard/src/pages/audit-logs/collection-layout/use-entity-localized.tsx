import { useCallback } from 'react';
import { getLanguage, useTranslation } from 'i18n';
import { DomainEventEntityType } from '@types';

const useEntityLocalized = () => {
  const { t } = useTranslation(['common', 'table']);
  const language = getLanguage();
  const getEntityLocalized = useCallback(
    (entityType: DomainEventEntityType) => {
      let result = '';
      switch (entityType) {
        case 'FEATURE':
          return (result = t('source-type.feature-flag'));
        case 'APIKEY':
          return (result = t('source-type.api-key'));
        case 'GOAL':
        case 'EXPERIMENT':
        case 'PROJECT':
        case 'PUSH':
        case 'ACCOUNT':
        case 'SUBSCRIPTION':
          return (result = t(`source-type.${entityType.toLowerCase()}`));
        case 'ENVIRONMENT':
        case 'ORGANIZATION':
        case 'SEGMENT':
        case 'TAG':
          return (result = t(`${entityType.toLowerCase()}`));
        case 'ADMIN_ACCOUNT':
        case 'AUTOOPS_RULE':
        case 'ADMIN_SUBSCRIPTION':
        case 'PROGRESSIVE_ROLLOUT':
        case 'FLAG_TRIGGER': {
          const translationKey = entityType.replace('_', '-').toLowerCase();
          return (result = t(`${translationKey}`));
        }
        case 'CODEREF':
          return (result = t('table:feature-flags.code-references'));
        default:
          break;
      }
      return result.toLowerCase();
    },
    [language]
  );

  return { getEntityLocalized };
};

export default useEntityLocalized;
