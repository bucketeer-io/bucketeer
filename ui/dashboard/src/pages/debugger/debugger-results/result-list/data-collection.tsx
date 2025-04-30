import { ColumnDef } from '@tanstack/react-table';
import { useTranslation } from 'i18n';
import { EvaluationFeature } from 'pages/debugger/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import ReasonTooltip from './reason-tooltip';
import ResultName from './result-name';

export const useColumns = ({
  isFlag,
  handleGetMaintainerInfo
}: {
  isFlag: boolean;
  handleGetMaintainerInfo: (email: string) => string;
}): ColumnDef<EvaluationFeature>[] => {
  const { t } = useTranslation(['common', 'table', 'form']);

  return [
    {
      accessorKey: !isFlag ? 'featureId' : 'userId',
      header: `${t(!isFlag ? 'name' : 'form:user-id')}`,
      size: 400,
      cell: ({ row }) => {
        const evaluation = row.original;
        const { featureId, userId, feature } = evaluation;
        return (
          <ResultName
            feature={feature}
            id={!isFlag ? featureId : userId}
            isFlag={isFlag}
            name={!isFlag ? feature.name : userId}
            variationType={feature.variationType}
            maintainer={handleGetMaintainerInfo(feature.maintainer)}
            onTable
          />
        );
      }
    },
    {
      accessorKey: 'variationName',
      header: `${t('table:feature-flags.variation')}`,
      size: 500,
      cell: ({ row }) => {
        const evaluationFeature = row.original;
        const { variationName, variationId, variationValue } =
          evaluationFeature;
        return (
          <div className="flex w-full col-span-5 gap-x-1.5">
            <FlagVariationPolygon index={row.index} className="mt-0.5" />
            <div className="flex flex-col flex-1 w-full gap-y-1.5">
              <p className="typo-para-small text-gray-700 break-all">
                {variationName || variationValue}
              </p>
              <p className="typo-para-small text-gray-500 break-all">
                {variationId}
              </p>
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'reason',
      header: `${t('table:reason')}`,
      size: 200,
      cell: ({ row }) => {
        const evaluationFeature = row.original;
        return <ReasonTooltip reason={evaluationFeature.reason} />;
      }
    }
  ];
};
