import { TagQueryOptions, useQueryTags } from '@queries/tags';
import { EntityType } from '@types';

export const useFetchTags = ({
  organizationId,
  environmentId,
  entityType,
  options
}: {
  organizationId?: string;
  environmentId?: string;
  entityType?: EntityType;
  options?: TagQueryOptions;
}) => {
  return useQueryTags({
    params: {
      cursor: String(0),
      pageSize: 9999,
      organizationId,
      environmentId,
      entityType: entityType || 'ACCOUNT'
    },
    ...options
  });
};
