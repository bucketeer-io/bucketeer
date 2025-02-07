import { AnyObject } from 'yup';
import { CollectionStatusType } from '@types';

export type FlagsViewType = 'LIST_VIEW' | 'GRID_VIEW';
export type FlagTabType = CollectionStatusType & 'FAVORITES';
export type FlagStatusType = 'new' | 'no_activity' | 'active';
export type FlagDataType = 'number' | 'string' | 'json' | 'boolean';
export type FlagActionType = 'ARCHIVE' | 'UNARCHIVE' | 'CLONE' | 'ACTIVE';

export interface FlagsTemp {
  id: string;
  name: string;
  type: FlagDataType;
  status: FlagStatusType;
  tags: string[];
  variations: boolean | string | Array<AnyObject>;
  operations: Array<AnyObject>;
  disabled: boolean;
  updatedAt: string;
  createdAt: string;
}
