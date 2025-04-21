import { useMemo } from 'react';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconChevronDown } from '@icons';
import Icon from 'components/icon';
import { DataTable } from 'elements/data-table';
import { GroupByType } from '..';
import { useColumns } from './data-collection';
import ResultName from './result-name';

interface Props {
  feature: Feature;
  featureId: string;
  userId: string;
  maintainer: string;
  isExpand: boolean;
  groupBy: GroupByType;
  handleGetMaintainerInfo: (email: string) => string;
}

const ResultItem = ({
  feature,
  featureId,
  userId,
  isExpand,
  groupBy,
  handleGetMaintainerInfo
}: Props) => {
  const isFlag = useMemo(() => groupBy === 'FLAG', [groupBy]);

  const columns = useColumns({
    isFlag,
    handleGetMaintainerInfo
  });

  return (
    <div className="flex flex-col w-full gap-y-5 p-5">
      <div className="flex items-center w-full justify-between gap-x-4">
        <ResultName
          feature={feature}
          id={isFlag ? featureId : userId}
          isFlag={isFlag}
          name={isFlag ? feature.name : userId}
          variationType={feature.variationType}
          maintainer={handleGetMaintainerInfo(feature.maintainer)}
        />
        <button
          className={cn(
            'flex-center size-5 rotate-0 transition-all duration-200',
            {
              'rotate-180': isExpand
            }
          )}
        >
          <Icon icon={IconChevronDown} />
        </button>
      </div>
      <DataTable columns={columns} data={[]} />
    </div>
  );
};

export default ResultItem;
