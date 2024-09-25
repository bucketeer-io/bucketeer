import { useMemo } from 'react';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { ColumnType, SpreadColumn } from 'hooks/use-table';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import Checkbox from 'components/checkbox';
import { Popover } from 'components/popover';
import { Switch } from 'components/switch';
import { TagGroup } from 'components/tag';
import OperationTag from 'components/tag/operation-tag';
import StatusTag, { StatusTagType } from 'components/tag/status-tag';
import Tag from 'components/tag/tag';
import { Tooltip } from 'components/tooltip';
import VariationGroup from 'components/variation/variation-group';
import Flag from './table-row-items/flag';
import Member from './table-row-items/member';
import Text from './table-row-items/text';
import Title from './table-row-items/title';

type StatusTagItem = {
  label: string;
  value: StatusTagType;
};

const statusTags: StatusTagItem[] = [
  { label: 'New', value: 'new' },
  { label: 'Activity', value: 'activity' },
  { label: 'No Activity', value: 'noActivity' },
  { label: 'Waiting', value: 'waiting' },
  { label: 'In Use', value: 'inUse' }
];

const TableRowItem = <T,>({
  cell,
  cellType = 'text',
  text,
  description,
  status,
  operators = [],
  flagIconType,
  variations = [],
  tagType,
  tagVariant,
  expandable,
  width,
  options = [],
  tooltip,
  disabled,
  addonSlot,
  rowId,
  rowsSelected = [],
  descriptionKey,
  statusKey,
  onClickCell,
  onClickPopover,
  handleSelectRow,
  spreadColumn,
  renderFunc
}: ColumnType<T>) => {
  const formatDateTime = useFormatDateTime();
  const { columnDef } = useMemo(() => {
    if (spreadColumn && cell) return spreadColumn(cell);
    return {} as SpreadColumn<T>;
  }, [spreadColumn, cell]);

  const accessorKey = useMemo(() => columnDef?.accessorKey, [columnDef]);
  const originalRow = useMemo(() => cell?.row?.original, [cell]);

  const mainDescription = useMemo(() => {
    if (originalRow) {
      if (description) return description;
      if (descriptionKey)
        return originalRow[descriptionKey as keyof T] as string;
      return '';
    }
    return '';
  }, [originalRow, description, descriptionKey]);

  const mainStatus: StatusTagType | undefined = useMemo(() => {
    if (originalRow) {
      if (statusKey) return originalRow[statusKey as keyof T] as StatusTagType;
      if (status) return status;
      return undefined;
    }
    return undefined;
  }, [originalRow, statusKey, status]);

  const formattedText = (cellValue: unknown) => {
    if (cellValue instanceof Date || accessorKey === 'createdAt')
      return formatDateTime(String(cellValue));
    return cellValue as string;
  };

  const renderType = () => {
    switch (cellType) {
      case 'title':
        return (
          <Title
            text={String(cell?.getValue()) || text}
            description={mainDescription}
          />
        );
      case 'flag':
        return (
          <Flag
            text={String(cell?.getValue()) || text}
            description={mainDescription}
            status={mainStatus}
            flagIconType={flagIconType}
          />
        );
      case 'member':
        return (
          <Member
            text={(originalRow && String(cell?.getValue())) || text}
            description={mainDescription}
            status={
              (originalRow &&
                (originalRow['status' as keyof T] as StatusTagType)) ||
              undefined
            }
            statusLabel={
              (originalRow &&
                (originalRow['statusLabel' as keyof T] as string)) ||
              ''
            }
          />
        );
      case 'operation':
        return (
          <TagGroup>
            {operators.map((i, index) => (
              <OperationTag key={index} type={i} />
            ))}
          </TagGroup>
        );
      case 'toggle':
        return <Switch />;

      case 'icon':
        return (
          <Popover
            options={options}
            icon={IconMoreHorizOutlined}
            addonSlot={addonSlot}
            disabled={disabled}
            onClick={value =>
              onClickPopover && onClickPopover(value, originalRow)
            }
            align="end"
          />
        );
      case 'tag':
        return (
          <Tag
            text={String(cell?.getValue()) || text}
            type={tagType}
            variant={tagVariant}
          />
        );
      case 'variation':
        return <VariationGroup variations={variations} />;
      case 'checkbox':
        return (
          <Checkbox
            title={text}
            description={mainDescription}
            checked={(rowId && rowsSelected.includes(rowId)) || false}
            onCheckedChange={() => handleSelectRow && handleSelectRow(rowId)}
          />
        );

      case 'empty':
        return <></>;
      case 'text':
        return (
          <Text
            text={formattedText(cell?.getValue() ?? text)}
            description={mainDescription}
          />
        );

      case 'status':
        return (
          <Popover
            trigger={
              <TagGroup expandable={expandable}>
                <StatusTag
                  variant={(cell?.getValue() || '') as StatusTagType}
                />
              </TagGroup>
            }
            options={statusTags}
          />
        );
    }
  };

  return (
    <td
      className={cn(
        'h-[60px] px-4 py-1.5 first:rounded-l-md last:rounded-r-md bg-white',
        {
          'cursor-pointer': !!onClickCell
        }
      )}
      onClick={() => onClickCell && onClickCell(originalRow)}
      style={{
        width
      }}
    >
      <Tooltip
        content={tooltip}
        trigger={renderFunc ? renderFunc(originalRow) : renderType()}
      />
    </td>
  );
};

export default TableRowItem;
