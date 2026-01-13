import * as React from 'react';
import { cn } from 'utils/style';

const TableRoot = React.forwardRef<
  HTMLTableElement,
  React.HTMLAttributes<HTMLTableElement>
>(({ ...props }, ref) => (
  <table
    className="table-fixed border-separate border-spacing-y-2 w-full mb-6"
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
>(({ align, className, ...props }, ref) => (
  <th
    className={cn(
      'text-gray-500 p-4 whitespace-nowrap typo-para-small uppercase',
      className
    )}
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
    className="px-4 py-2 h-[60px] min-h-[60px] first:rounded-l-lg last:rounded-r-lg"
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

const Table = {
  Root: TableRoot,
  Header: TableHeader,
  Head: TableHead,
  Body: TableBody,
  Row: TableRow,
  Footer: TableFooter,
  Cell: TableCell,
  Caption: TableCaption
};

export default Table;
