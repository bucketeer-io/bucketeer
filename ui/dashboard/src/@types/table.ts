import { ReactNode } from 'react';
import { PopoverOption, PopoverValue } from 'components/popover';
import { FlagProps } from 'components/table/table-row-items/flag';
import { TextProps } from 'components/table/table-row-items/text';
import { TitleProps } from 'components/table/table-row-items/title';
import { OperationType } from 'components/tag/operation-tag';
import { StatusTagType } from 'components/tag/status-tag';
import { TagType, TagVariant } from 'components/tag/tag';
import { VariationGroupProps } from 'components/variation/variation-group';
import { AddonSlot } from './app';

export type SortedType = 'asc' | 'desc' | '';
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

export type SortedObjType = {
  colIndex: number;
  sortedType: SortedType;
};

export type CellValueType = string | number | Date | boolean;
export type TableSignature = {
  [key: string]: CellValueType;
};
export type TableHeaderItemProps = {
  text?: string;
  sort?: boolean;
  type?: 'title' | 'checkbox' | 'empty';
  defaultSortedType?: SortedType;
  sortedType?: SortedType;
  width?: string;
  isSelectAllRows?: boolean;
  colIndex?: number;
  sortedObj?: SortedObjType;
  fieldName?: string;
  handleToggleSelectAllRows?: () => void;
  handleSortedData?: (colIndex?: number, fieldName?: string) => void;
};

export type TableRowItemAdditionalProps = {
  type?: TableRowItemType;
  tagType?: TagType;
  tagVariant?: TagVariant;
  operators?: OperationType[];
  statusTags?: StatusTagType[];
  expandable?: boolean;
  width?: string;
  options?: PopoverOption<PopoverValue>[];
  tooltip?: string;
  addonSlot?: AddonSlot;
  disabled?: boolean;
  rowIndex?: number;
  rowsSelected?: number[];
  tableRows?: TableRows;
  handleSelectRow?: (rowIndex?: number) => void;
  onClick?: () => void;
  onClickPopover?: (value: PopoverValue) => void;
};

export type TableRowItemProps = TitleProps &
  TextProps &
  FlagProps &
  VariationGroupProps &
  TableRowItemAdditionalProps;

export type TableProps<T> = {
  headers: TableHeaders;
  rows: TableRows;
  elementEmpty?: ReactNode;
  originalData: T[];
  rowsData: T[];
  setRowsData: (data: T[]) => void;
};

export type TableHeaderProps = {
  data: TableHeaderItemProps[];
  isSelectAllRows?: boolean;
  sortedObj?: SortedObjType;
  handleToggleSelectAllRows?: () => void;
  handleSortedData?: (colIndex?: number, fieldName?: string) => void;
};

export type TableRowProps = {
  tableRows: TableRows;
  data: TableRowItemProps[];
  rowIndex: number;
  rowsSelected: number[];
  handleSelectRow: (rowIndex?: number) => void;
};

export type TableHeaders = TableHeaderItemProps[];
export type TableRows = TableRowItemProps[][];
