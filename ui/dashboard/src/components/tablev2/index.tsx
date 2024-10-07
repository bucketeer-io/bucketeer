import * as React from 'react';
import { IconAngleDown, IconAngleUp } from '@icons';
import Icon from 'components/icon';

const TableRoot = React.forwardRef<
  HTMLTableElement,
  React.HTMLAttributes<HTMLTableElement>
>(({ ...props }, ref) => (
  <table
    className="border-separate border-spacing-y-3 w-full mb-6"
    ref={ref}
    {...props}
  />
));
TableRoot.displayName = 'Table';

const TableHeader = React.forwardRef<
  HTMLTableSectionElement,
  React.HTMLAttributes<HTMLTableSectionElement>
>(({ ...props }, ref) => <thead ref={ref} {...props} />);
TableHeader.displayName = 'TableHeader';

const TableBody = React.forwardRef<
  HTMLTableSectionElement,
  React.HTMLAttributes<HTMLTableSectionElement>
>(({ ...props }, ref) => <tbody ref={ref} {...props} />);
TableBody.displayName = 'TableBody';

const TableFooter = React.forwardRef<
  HTMLTableSectionElement,
  React.HTMLAttributes<HTMLTableSectionElement>
>(({ ...props }, ref) => <tfoot ref={ref} {...props} />);
TableFooter.displayName = 'TableFooter';

const TableRow = React.forwardRef<
  HTMLTableRowElement,
  React.HTMLAttributes<HTMLTableRowElement>
>(({ ...props }, ref) => <tr ref={ref} {...props} />);
TableRow.displayName = 'TableRow';

const TableHead = React.forwardRef<
  HTMLTableCellElement,
  React.ThHTMLAttributes<HTMLTableCellElement>
>(({ align, ...props }, ref) => (
  <th
    className="text-gray-500 p-4 typo-para-small uppercase"
    ref={ref}
    {...props}
    align={align ? align : 'left'}
  />
));
TableHead.displayName = 'TableHead';

const TableCell = React.forwardRef<
  HTMLTableCellElement,
  React.TdHTMLAttributes<HTMLTableCellElement>
>(({ ...props }, ref) => (
  <td
    className="h-[60px] px-4 py-1.5 first:rounded-l-md last:rounded-r-md bg-white"
    ref={ref}
    {...props}
  />
));
TableCell.displayName = 'TableCell';

const TableCaption = React.forwardRef<
  HTMLTableCaptionElement,
  React.HTMLAttributes<HTMLTableCaptionElement>
>(({ ...props }, ref) => <caption ref={ref} {...props} />);
TableCaption.displayName = 'TableCaption';

const TableHeadSort = () => {
  return (
    <div className="flex flex-col gap-y-0.5">
      <Icon icon={IconAngleUp} size="fit" color="gray-300" />
      <Icon icon={IconAngleDown} size="fit" color="gray-300" />
    </div>
  );
};
TableHeadSort.displayName = 'TableHeadSort';

const Table = {
  Root: TableRoot,
  Header: TableHeader,
  Head: TableHead,
  Body: TableBody,
  Row: TableRow,
  Footer: TableFooter,
  Cell: TableCell,
  Caption: TableCaption,
  HeadSort: TableHeadSort
};

export default Table;
