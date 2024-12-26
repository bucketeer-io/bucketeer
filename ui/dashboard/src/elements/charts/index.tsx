import { memo, ReactNode } from 'react';
import { DataTable } from 'elements/data-table';
import { TempTableDataType, useColumns } from './data-collection';
import ChartHeader, { HeaderProps } from './header';
import LineChart, { ChartDataType } from './line-chart';

type ChartWrapperProps = HeaderProps & {
  chartData: ChartDataType;
  tableData: TempTableDataType[];
  formatLabel: (index: number) => string;
  renderName: (tempData: TempTableDataType) => ReactNode;
};

const ChartWrapper = memo(
  ({
    chartData,
    tableData,
    formatLabel,
    renderName,
    ...props
  }: ChartWrapperProps) => {
    const columns = useColumns({ renderName });

    return (
      <div className="flex flex-col w-fit xl:w-full border border-gray-200 rounded-2xl">
        <ChartHeader {...props} />
        <div className="flex w-full p-5 pt-0">
          <LineChart chartData={chartData} formatLabel={formatLabel} />
          <div className="flex size-fit xl:w-[40%] min-w-[300px]">
            <DataTable
              data={tableData}
              columns={columns}
              rowClassName="!shadow-none [&>td]:!border-b [&>td]:!rounded-none [&>td]:last:!border-b-0"
              onSortingChange={() => {}}
            />
          </div>
        </div>
      </div>
    );
  }
);

export default ChartWrapper;
