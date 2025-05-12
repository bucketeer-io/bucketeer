import { useMemo } from 'react';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import { GroupByType } from 'pages/debugger/page-content';
import { EvaluationFeature } from 'pages/debugger/types';
import Icon from 'components/icon';
import { DataTable } from 'elements/data-table';
import { useColumns } from './data-collection';
import ResultName from './result-name';

interface Props {
  group: EvaluationFeature[];
  isExpand: boolean;
  groupBy: GroupByType;
  handleGetMaintainerInfo: (email: string) => string;
  onToggleExpandItem?: () => void;
}

const ResultItem = ({
  group,
  isExpand,
  groupBy,
  handleGetMaintainerInfo,
  onToggleExpandItem
}: Props) => {
  const isFlag = useMemo(() => groupBy === 'FLAG', [groupBy]);

  const columns = useColumns({
    isFlag,
    handleGetMaintainerInfo
  });

  return (
    <div
      className={cn(
        'flex flex-col w-full p-5 pb-0 rounded-lg shadow-card min-h-[86px] transition-all duration-200',
        {
          'pb-5 min-h-fit': isExpand
        }
      )}
    >
      <div
        className={cn('flex items-center w-full justify-between gap-x-4 pb-5', {
          'border-b border-gray-200': isExpand
        })}
      >
        <ResultName
          feature={group[0]?.feature}
          id={isFlag ? group[0]?.featureId : group[0]?.userId}
          isFlag={isFlag}
          name={isFlag ? group[0]?.feature.name : group[0]?.userId}
          variationType={group[0]?.feature.variationType}
          maintainer={handleGetMaintainerInfo(group[0]?.feature.maintainer)}
        />
        {onToggleExpandItem && (
          <button
            className={cn(
              'flex-center size-5 rotate-0 transition-all duration-200',
              {
                'rotate-180': isExpand
              }
            )}
            onClick={onToggleExpandItem}
          >
            <Icon icon={IconChevronDown} />
          </button>
        )}
      </div>
      <div
        className={cn(
          '[&>table]:m-0 [&>table]:border-collapse [&>table>tbody]:divide-y [&>table>tbody]:divide-gray-200 [&>table>tbody>tr]:rounded-none [&>table>tbody>tr]:shadow-none [&>table>tbody>tr>td]:rounded-none [&>table>tbody>tr:last-child]:rounded-b-lg [&>table>tbody>tr>td:first-child]:pl-0 [&>table>tbody>tr>td:last-child]:pr-0 [&>table>thead>tr>th:first-child]:pl-0 [&>table>thead>tr>th:last-child]:pr-0 h-0 opacity-0 transition-all duration-200 z-[-1] [&>table>tbody>tr>td]:py-4',
          {
            'h-fit opacity-100 z-0': isExpand
          }
        )}
      >
        <DataTable columns={columns} data={group} manualSorting={false} />
      </div>
    </div>
  );
};

export default ResultItem;
