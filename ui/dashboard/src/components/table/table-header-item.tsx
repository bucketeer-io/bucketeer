import { useMemo } from 'react';
import { SortedObjType } from '@types';
import { SortedType } from 'utils/sort';
import { cn } from 'utils/style';
import { IconAngleDown, IconAngleUp } from '@icons';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';

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

const TableHeaderItem = ({
  text,
  sort,
  type = 'title',
  width,
  isSelectAllRows,
  colIndex,
  sortedObj,
  fieldName,
  handleToggleSelectAllRows,
  handleSortedData
}: TableHeaderItemProps) => {
  const renderType = useMemo(() => {
    switch (type) {
      case 'title':
        return (
          <div
            className="text-gray-500 typo-para-small flex items-center gap-3 whitespace-nowrap select-none"
            onClick={() => {
              if (sort && handleSortedData) {
                handleSortedData(colIndex, fieldName);
              }
            }}
          >
            {text}
            {sort && (
              <SortIcon
                sortedType={
                  (sortedObj?.colIndex === colIndex && sortedObj?.sortedType) ||
                  ''
                }
              />
            )}
          </div>
        );
      case 'checkbox':
        return (
          <Checkbox
            checked={isSelectAllRows}
            onCheckedChange={handleToggleSelectAllRows}
          />
        );

      case 'empty':
        return <></>;
    }
  }, [type, isSelectAllRows, sortedObj]);

  return (
    <th
      style={{
        width
      }}
      className="h-[60px] px-4 py-5 cursor-pointer"
    >
      {renderType}
    </th>
  );
};

const sortIconCls = ({
  sortedType,
  selfType
}: {
  sortedType: SortedType;
  selfType: SortedType;
}) => {
  return cn('flex-center size-fit text-gray-300', {
    'text-gray-500': sortedType === selfType
  });
};

const SortIcon = ({ sortedType }: { sortedType: SortedType }) => {
  return (
    <div className="flex flex-col gap-y-0.5">
      <div className={sortIconCls({ sortedType, selfType: 'asc' })}>
        <Icon icon={IconAngleUp} className="size-fit" />
      </div>
      <div className={sortIconCls({ sortedType, selfType: 'desc' })}>
        <Icon icon={IconAngleDown} className="size-fit" />
      </div>
    </div>
  );
};

export default TableHeaderItem;
