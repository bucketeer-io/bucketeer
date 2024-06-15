import { FC } from 'react';

import { classNames } from '../../../utils/css';
import { HelpTextTooltip } from '../../HelpTextTooltip';

export interface HeaderCell {
  id: string;
  label: string;
  helpText: string;
}

export interface TableHeaderProps {
  cells: HeaderCell[];
}

export const TableHeader: FC<TableHeaderProps> = ({ cells }) => {
  return (
    <thead className="bg-gray-50 rounded-md">
      <tr>
        {cells.map((headCell) => (
          <td
            key={headCell.id}
            className={classNames(
              'text-left text-xs rounded-md',
              'font-medium text-gray-500 uppercase tracking-wider border-b',
              'px-5'
            )}
          >
            <div className={classNames('flex flex-row items-center')}>
              {headCell.label}
              {headCell.helpText && (
                <HelpTextTooltip helpText={headCell.helpText} />
              )}
            </div>
          </td>
        ))}
      </tr>
    </thead>
  );
};

export interface TableProps {}

export const Table: FC<TableProps> = ({ children }) => {
  return (
    <div className={classNames('w-full rounded-md border')}>
      <table
        className={classNames(
          'min-w-full table-auto leading-normal rounded-md'
        )}
      >
        {children}
      </table>
    </div>
  );
};

export interface TableBodyProps {}

export const TableBody: FC<TableBodyProps> = ({ children }) => {
  return <tbody className="text-sm text-gray-600 rounded-md">{children}</tbody>;
};

export interface TableRowProps {}

export const TableRow: FC<TableRowProps> = ({ children }) => {
  return <tr className="rounded-md border-t">{children}</tr>;
};

export interface TableCellProps {
  textLeft?: boolean;
}

export const TableCell: FC<TableCellProps> = ({ children, textLeft }) => {
  return (
    <td
      className={classNames(
        'rounded-md w-[1%] px-5',
        textLeft ? 'text-left' : 'text-right'
      )}
    >
      {children}
    </td>
  );
};
