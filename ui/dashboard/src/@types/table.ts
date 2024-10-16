import { ReactNode } from 'react';
import type {
  Table as TanstackTableType,
  Cell,
  Header,
  Row
} from '@tanstack/react-table';
import { SortingType } from 'containers/pages';
import { ColumnData, SpreadColumn } from 'hooks/use-table';
import { PaginationProps } from 'components/pagination';
import { PopoverOption, PopoverValue } from 'components/popover';
import { FlagProps } from 'components/table/table-row-items/flag';
import { TextProps } from 'components/table/table-row-items/text';
import { TitleProps } from 'components/table/table-row-items/title';
import { OperationType } from 'components/tag/operation-tag';
import { TagType, TagVariant } from 'components/tag/tag';
import { VariationGroupProps } from 'components/variation/variation-group';
import { AddonSlot } from './app';
import { OrderBy } from './collection';

export type TableRowItemType =
  | 'title'
  | 'text'
  | 'tag'
  | 'variation'
  | 'operation'
  | 'toggle'
  | 'icon'
  | 'checkbox'
  | 'empty'
  | 'member'
  | 'flag'
  | 'status';

export type CellValueType = string | number | Date | boolean | string[];
export type TableSignature = {
  [key: string]: CellValueType;
};

export type TableHeaderCellProps<T> = {
  text?: string;
  sort?: boolean;
  width?: string;
  isSelectAllRows?: boolean;
  colIndex?: number;
  header: Header<T, unknown>;
  sortingState?: SortingType;
  spreadColumn: <T>(data: ColumnData<T>) => SpreadColumn<T>;
  handleToggleSelectAllRows?: () => void;
  onSortingTable?: (accessorKey: string, sortingKey?: OrderBy) => void;
};

export type TableRowItemAdditionalProps<T> = {
  cell?: Cell<T, unknown>;
  cellType?: TableRowItemType;
  tagType?: TagType;
  tagVariant?: TagVariant;
  operators?: OperationType[];
  expandable?: boolean;
  width?: string;
  options?: PopoverOption<PopoverValue>[];
  tooltip?: string;
  addonSlot?: AddonSlot;
  disabled?: boolean;
  rowId?: string;
  rowsSelected?: string[];
  descriptionKey?: string;
  statusKey?: string;
  handleSelectRow?: (rowId?: string) => void;
  onClickCell?: (row?: T) => void;
  onClickPopover?: (value: PopoverValue, row?: T) => void;
  spreadColumn?: (data: ColumnData<T>) => SpreadColumn<T>;
};

export type TableRowItemProps<T> = TitleProps &
  TextProps &
  FlagProps &
  VariationGroupProps &
  TableRowItemAdditionalProps<T>;

export type TableProps<T> = {
  table: TanstackTableType<T>;
  elementEmpty: ReactNode;
  paginationProps?: PaginationProps;
  rowsSelected?: string[];
  sortingState?: SortingType;
  onSortingTable?: (accessorKey: string, sortingKey?: OrderBy) => void;
  setRowsSelected?: (rows: string[]) => void;
  spreadColumn: <T>(data: ColumnData<T>) => SpreadColumn<T>;
};

export type TableHeaderProps<T> = {
  data: Header<T, unknown>[];
  isSelectAllRows?: boolean;
  sortingState?: SortingType;
  spreadColumn: <T>(data: ColumnData<T>) => SpreadColumn<T>;
  handleToggleSelectAllRows?: () => void;
  onSortingTable?: (accessorKey: string, sortingKey?: OrderBy) => void;
};

export type TableRowProps<T> = {
  row: Row<T>;
  rowsSelected: string[];
  handleSelectRow: (rowId?: string) => void;
  spreadColumn: <T>(data: ColumnData<T>) => SpreadColumn<T>;
};

export type TableHeaders<T> = TableHeaderCellProps<T>[];
