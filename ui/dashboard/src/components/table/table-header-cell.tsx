import { useMemo } from 'react';
import { OrderDirection, TableHeaderCellProps } from '@types';
import { cn } from 'utils/style';
import { IconAngleDown, IconAngleUp } from '@icons';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';

type SortedType = OrderDirection & 'DEFAULT';

const TableHeaderCell = <T,>({
  isSelectAllRows,
  header,
  sortingState,
  spreadColumn,
  handleToggleSelectAllRows,
  onSortingTable
}: TableHeaderCellProps<T>) => {
  const { columnDef } = useMemo(() => spreadColumn(header), [header]);

  const renderType = () => {
    switch (columnDef?.headerCellType) {
      case 'checkbox':
        return (
          <Checkbox
            checked={isSelectAllRows}
            onCheckedChange={handleToggleSelectAllRows}
          />
        );

      case 'empty':
        return <></>;
      case 'title':
      default:
        return (
          <div
            className="text-gray-500 typo-para-small flex items-center gap-3 whitespace-nowrap uppercase select-none"
            onClick={() => {
              if (!columnDef.sorting) return;
              header.column.toggleSorting();
            }}
          >
            {columnDef?.header}
            {columnDef.sorting && (
              <SortIcon
                sortedType={
                  (sortingState?.id === columnDef.accessorKey
                    ? sortingState?.orderDirection
                    : 'DEFAULT') as SortedType
                }
              />
            )}
          </div>
        );
    }
  };

  return (
    <th
      style={{
        width: columnDef?.size,
        minWidth: columnDef?.minSize
      }}
      className={cn('h-[60px] px-4 py-5 cursor-pointer')}
      onClick={() =>
        columnDef?.sorting &&
        onSortingTable &&
        onSortingTable(
          columnDef?.accessorKey || columnDef?.id || '',
          columnDef?.sortingKey
        )
      }
    >
      {renderType()}
    </th>
  );
};

const sortIconCls = ({
  sortedType,
  selfType
}: {
  sortedType?: SortedType;
  selfType: SortedType;
}) => {
  return cn('flex-center size-fit text-gray-300', {
    'text-gray-500': sortedType === selfType
  });
};

const SortIcon = ({ sortedType }: { sortedType?: SortedType }) => {
  return (
    <div className="flex flex-col gap-y-0.5">
      <div
        className={sortIconCls({ sortedType, selfType: 'ASC' as SortedType })}
      >
        <Icon icon={IconAngleUp} size="fit" />
      </div>
      <div
        className={sortIconCls({ sortedType, selfType: 'DESC' as SortedType })}
      >
        <Icon icon={IconAngleDown} size="fit" />
      </div>
    </div>
  );
};

export default TableHeaderCell;
