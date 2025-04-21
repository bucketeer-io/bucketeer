import { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { IconInfoFilled } from '@icons';
import { EvaluationFeatureAccount } from 'pages/debugger/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import ResultName from './result-name';

export const useColumns = ({
  isFlag,
  handleGetMaintainerInfo
}: {
  isFlag: boolean;
  handleGetMaintainerInfo: (email: string) => string;
}): ColumnDef<EvaluationFeatureAccount>[] => {
  const { t } = useTranslation(['common', 'table']);

  return [
    {
      accessorKey: isFlag ? 'feature_id' : 'user_id',
      header: `${t(isFlag ? 'name' : 'user-id')}`,
      size: 500,
      cell: ({ row }) => {
        const evaluationFeature = row.original;
        const { feature, feature_id, user_id } = evaluationFeature;
        return (
          <ResultName
            feature={feature}
            id={isFlag ? feature_id : user_id}
            isFlag={isFlag}
            name={isFlag ? feature.name : user_id}
            variationType={feature.variationType}
            maintainer={handleGetMaintainerInfo(feature.maintainer)}
            onTable
          />
        );
      }
    },
    {
      accessorKey: 'variation_name',
      header: `${t('table:feature-flags.variation')}`,
      size: 150,
      cell: ({ row }) => {
        const evaluationFeature = row.original;
        const { variation_name, variation_id } = evaluationFeature;
        return (
          <div className="flex w-full col-span-5 gap-x-1.5">
            <FlagVariationPolygon index={row.index} />
            <div className="flex flex-col flex-1 w-full gap-y-1.5">
              <p className="typo-para-small text-gray-700 break-all">
                {variation_name}
              </p>
              <p className="typo-para-small text-gray-500 break-all">
                {variation_id}
              </p>
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'reason',
      header: `${t('table:reason')}`,
      size: 150,
      cell: ({ row }) => {
        const evaluationFeature = row.original;
        return (
          <Tooltip
            content={evaluationFeature.reason}
            trigger={
              <div className="flex items-center gap-x-2">
                <p className="typo-para-medium text-gray-700">
                  {evaluationFeature.reason}
                </p>
                <Icon icon={IconInfoFilled} color="gray-500" />
              </div>
            }
          />
        );
      }
    }
  ];
};
