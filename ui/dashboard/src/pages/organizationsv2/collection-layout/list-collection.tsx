import { ListTableCollection } from './list-table-collection';
import type { CollectionProps } from './types';

export const ListCollection = ({ organizations }: CollectionProps) => {
  return <ListTableCollection organizations={organizations} />;
};
