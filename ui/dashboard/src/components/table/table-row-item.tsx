import { useMemo } from 'react';
import { IconMoreHorizOutlined } from 'react-icons-material-design';
import { TableRowItemProps } from '@types';
import Checkbox from 'components/checkbox';
import Popover from 'components/popover';
import { Switch } from 'components/switch';
import { TagGroup } from 'components/tag';
import OperationTag from 'components/tag/operation-tag';
import StatusTag from 'components/tag/status-tag';
import Tag from 'components/tag/tag';
import Tooltip from 'components/tooltip';
import VariationGroup from 'components/variation/variation-group';
import Flag from './table-row-items/flag';
import Member from './table-row-items/member';
import Text from './table-row-items/text';
import Title from './table-row-items/title';

const TableRowItem = ({
  type = 'text',
  text,
  description,
  status,
  operators = [],
  flagIconType,
  variations = [],
  tagType,
  tagVariant,
  statusTags = [],
  expandable,
  width,
  options = [],
  tooltip,
  disabled,
  addonSlot,
  rowIndex,
  rowsSelected = [],
  tableRows,
  onClick,
  onClickPopover,
  handleSelectRow
}: TableRowItemProps) => {
  const renderType = useMemo(() => {
    switch (type) {
      case 'title':
        return <Title text={text} description={description} />;
      case 'flag':
        return (
          <Flag
            text={text}
            description={description}
            status={status}
            flagIconType={flagIconType}
          />
        );
      case 'member':
        return <Member text={text} description={description} />;
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
            onClick={onClickPopover}
          />
        );
      case 'tag':
        return <Tag text={text} type={tagType} variant={tagVariant} />;
      case 'variation':
        return <VariationGroup variations={variations} />;
      case 'checkbox':
        return (
          <Checkbox
            title={text}
            description={description}
            checked={
              (typeof rowIndex === 'number' &&
                rowsSelected.includes(rowIndex)) ||
              false
            }
            onCheckedChange={() => handleSelectRow && handleSelectRow(rowIndex)}
          />
        );

      case 'empty':
        return <></>;
      case 'text':
        return <Text text={text} description={description} />;
      case 'status':
        return (
          <TagGroup expandable={expandable}>
            {statusTags.map(i => (
              <StatusTag key={i} variant={i} />
            ))}
          </TagGroup>
        );
    }
  }, [type, rowsSelected, tableRows]);

  return (
    <td
      className="h-[60px] px-4 py-1.5 first:rounded-l-md last:rounded-r-md bg-white"
      onClick={onClick}
      style={{
        width
      }}
    >
      {tooltip ? (
        <Tooltip content={tooltip} trigger={renderType} />
      ) : (
        renderType
      )}
    </td>
  );
};

export default TableRowItem;
