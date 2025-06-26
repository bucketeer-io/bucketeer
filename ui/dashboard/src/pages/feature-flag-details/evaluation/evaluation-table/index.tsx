import { useCallback } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { EvaluationCounter, Feature } from '@types';
import { getVariationColor } from 'utils/style';
import { DatasetReduceType } from 'pages/experiment-details/collection-loader/results/goal-results/timeseries-area-line-chart';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Checkbox from 'components/checkbox';

const EvaluationTable = ({
  feature,
  timeRangeLabel,
  countData,
  dataSets,
  onToggleShowData
}: {
  feature: Feature;
  timeRangeLabel: string;
  countData: EvaluationCounter[];
  dataSets: DatasetReduceType[];
  onToggleShowData: (variationId: string) => void;
}) => {
  const { t } = useTranslation(['table']);
  const getVariation = useCallback(
    (item: EvaluationCounter) =>
      feature.variations.find(v => v.id === item.variationId),
    [feature]
  );
  const getVariationLabel = useCallback(
    (item: EvaluationCounter) => {
      if (item.variationId === 'default') return 'default value';
      const variation = getVariation(item);
      return variation?.name || variation?.value || '';
    },
    [feature]
  );

  return (
    <div className="flex flex-col gap-y-6 w-full min-w-[650px]">
      <div className="flex items-center w-full">
        <div className="w-[40%] typo-para-medium text-gray-700">
          {t('evaluation.variation-counts')}
        </div>
        <div className="w-[60%] typo-para-medium text-gray-700">
          <Trans
            i18nKey="table:evaluation.total-evaluations"
            components={{
              desc: <span className="text-gray-500" />
            }}
            values={{
              value: `(${timeRangeLabel})`
            }}
          />
        </div>
      </div>
      {countData?.map((item, index) => {
        const isHidden = dataSets?.find(
          dataset => dataset.label === item?.variationId
        )?.hidden;
        return (
          <div
            key={index}
            className="flex items-center w-full px-4 py-5 bg-white rounded-lg shadow-card"
          >
            <div className="flex items-center w-[40%] gap-x-2">
              <Checkbox
                checked={!isHidden}
                onCheckedChange={() =>
                  onToggleShowData(String(item.variationId))
                }
              />
              <Polygon
                className="border-none size-3"
                style={{
                  background: getVariationColor(index),
                  zIndex: index
                }}
              />
              <p className="typo-para-small text-gray-700">
                {getVariationLabel(item)}
              </p>
            </div>
            <div className="w-[60%] typo-para-medium text-gray-700">
              {Number(item.timeseries?.totalCounts || 0)?.toLocaleString()}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default EvaluationTable;
