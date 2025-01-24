import { useQueryTags } from '@queries/tags';
import { EntityType } from '@types';

export const useFetchTags = ({
  organizationId,
  entityType
}: {
  organizationId: string;
  entityType?: EntityType;
}) => {
  return useQueryTags({
    params: {
      cursor: String(0),
      pageSize: 9999,
      organizationId,
      entityType: entityType || 'ACCOUNT'
    }
  });
};
