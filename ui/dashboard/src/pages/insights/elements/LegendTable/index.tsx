import { useMemo, useState } from 'react';
import {
  ColumnDef,
  SortingState,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable
} from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconSorting, IconSortingDown, IconSortingUp } from '@icons';
import Table from 'components/table';
import { Tooltip } from 'components/tooltip';
import EmptyState from 'elements/empty-state';
import { formatYAxis, getColor } from '../chart-utils';

interface LegendRow {
  index: number;
  label: string;
  min: number;
  max: number;
  avg: number;
  last: number;
}

interface LegendTableProps {
  datasets: { label?: string; data: number[] }[];
  formatter?: (value: number) => string;
}

export const LegendTable = ({ datasets, formatter }: LegendTableProps) => {
  const fmt = formatter ?? formatYAxis;
  const { t } = useTranslation(['common']);
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'avg', desc: true }
  ]);

  const rows: LegendRow[] = useMemo(
    () =>
      datasets.map((ds, i) => {
        const nums = ds.data.filter(v => typeof v === 'number' && !isNaN(v));
        const min = nums.length
          ? nums.reduce((a, b) => (b < a ? b : a), nums[0])
          : 0;
        const max = nums.length
          ? nums.reduce((a, b) => (b > a ? b : a), nums[0])
          : 0;
        const avg = nums.length
          ? nums.reduce((a, b) => a + b, 0) / nums.length
          : 0;
        const last = nums.length ? nums[nums.length - 1] : 0;
        return {
          index: i,
          label: ds.label ?? `Series ${i + 1}`,
          min,
          max,
          avg,
          last
        };
      }),
    [datasets]
  );

  const columns = useMemo<ColumnDef<LegendRow>[]>(
    () => [
      {
        accessorKey: 'label',
        header: t('insights.series'),
        enableSorting: true,
        cell: ({ row }) => (
          <div className="flex items-center gap-2 pl-4">
            <span
              className="inline-block w-3 h-3 rounded-full flex-shrink-0"
              style={{ backgroundColor: getColor(row.original.index) }}
            />
            <Tooltip
              content={row.original.label}
              side="top"
              trigger={
                <span className="max-w-[300px] text-gray-700 truncate typo-para-medium">
                  {row.original.label}
                </span>
              }
            />
          </div>
        )
      },
      {
        accessorKey: 'min',
        header: t('insights.min'),
        enableSorting: true,
        cell: ({ getValue }) => fmt(getValue<number>())
      },
      {
        accessorKey: 'max',
        header: t('insights.max'),
        enableSorting: true,
        cell: ({ getValue }) => fmt(getValue<number>())
      },
      {
        accessorKey: 'avg',
        header: t('insights.avg'),
        enableSorting: true,
        cell: ({ getValue }) => fmt(getValue<number>())
      },
      {
        accessorKey: 'last',
        header: t('insights.last'),
        enableSorting: true,
        cell: ({ getValue }) => fmt(getValue<number>())
      }
    ],
    [t, fmt]
  );

  const table = useReactTable({
    data: rows,
    columns,
    state: { sorting },
    onSortingChange: setSorting,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualSorting: false
  });

  return (
    <div className="overflow-x-auto small-scroll h-[250px]">
      <div className="min-w-full">
        <Table.Root>
          <Table.Header className="sticky top-0 z-10 bg-white">
            {table.getHeaderGroups().map(headerGroup => (
              <Table.Row key={headerGroup.id}>
                {headerGroup.headers.map((header, i) => (
                  <Table.Head
                    key={header.id}
                    align={i === 0 ? undefined : 'right'}
                    onClick={header.column.getToggleSortingHandler()}
                    className={cn('select-none', {
                      'cursor-pointer': header.column.getCanSort()
                    })}
                  >
                    <div
                      className={cn('flex items-center gap-1', {
                        'justify-end': i !== 0
                      })}
                    >
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                      {header.column.getCanSort() &&
                        {
                          asc: <IconSortingUp />,
                          desc: <IconSortingDown />,
                          false: <IconSorting />
                        }[header.column.getIsSorted() as string]}
                    </div>
                  </Table.Head>
                ))}
              </Table.Row>
            ))}
          </Table.Header>
          <Table.Body>
            {!datasets.length ? (
              <Table.Row>
                <Table.Cell className="pt-3" colSpan={columns.length}>
                  <div className="w-full h-full flex items-center justify-center">
                    <EmptyState.Root variant="no-data" size="sm">
                      <EmptyState.Illustration />
                    </EmptyState.Root>
                  </div>
                </Table.Cell>
              </Table.Row>
            ) : (
              <>
                {table.getRowModel().rows.map(row => (
                  <Table.Row key={row.id}>
                    {row.getVisibleCells().map((cell, i) => (
                      <Table.Cell
                        key={cell.id}
                        align={i === 0 ? undefined : 'right'}
                        className="typo-para-medium text-gray-700 pr-4"
                      >
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext()
                        )}
                      </Table.Cell>
                    ))}
                  </Table.Row>
                ))}
              </>
            )}
          </Table.Body>
        </Table.Root>
      </div>
    </div>
  );
};
