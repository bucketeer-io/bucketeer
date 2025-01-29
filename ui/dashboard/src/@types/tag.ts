export type EntityType = 'UNKNOWN' | 'FEATURE_FLAG';

export interface Tag {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  entityType: EntityType;
  environmentId: string;
}

export interface TagCollection {
  tags: Tag[];
  cursor: string;
  totalCount: string;
}
