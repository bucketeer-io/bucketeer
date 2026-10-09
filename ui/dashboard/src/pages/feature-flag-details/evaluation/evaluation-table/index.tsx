import { useCallback } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { EvaluationCounter, Feature } from '@types';
import { getVariationColor } from 'utils/style';
import { IconInfo } from '@icons';
import { DatasetReduceType } from 'pages/experiment-details/collection-loader/results/goal-results/timeseries-area-line-chart';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Checkbox from 'components/checkbox';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

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
    <div className="flex flex-col gap-y-6 w-full">
      <div className="flex items-center w-full gap-x-6">
        <div className="w-[60%] sm:w-[40%] typo-para-medium text-gray-700">
          {t('evaluation.variation-counts')}
        </div>
        <div className="w-[40%] sm:w-[60%] typo-para-medium text-gray-700">
          <Trans
            i18nKey="table:evaluation.total-evaluations"
            components={{
              desc: <span className="text-gray-500 hidden sm:inline" />
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
            className="flex items-center w-full px-4 py-5 bg-white rounded-lg shadow-card gap-x-6"
          >
            <div className="flex items-center w-[60%] sm:w-[40%] gap-x-2">
              <Checkbox
                checked={!isHidden}
                onCheckedChange={() =>
                  onToggleShowData(String(item.variationId))
                }
              />
              <Polygon
                className="border-none size-3 min-w-3"
                style={{
                  background: getVariationColor(index),
                  zIndex: index
                }}
              />
              {item.variationId === 'default' ? (
                <Tooltip
                  align="start"
                  side="top"
                  content={t('default-value-tooltip')}
                  trigger={
                    <div className="flex items-center gap-x-2">
                      <p className="line-clamp-2 typo-para-small text-gray-700">
                        {getVariationLabel(item)}
                      </p>
                      <Icon icon={IconInfo} size={'xs'} color="gray-500" />
                    </div>
                  }
                  className="max-w-[300px]"
                />
              ) : (
                <p className="line-clamp-2 typo-para-small text-gray-700">
                  {getVariationLabel(item)}
                </p>
              )}
            </div>
            <div className="w-[40%] sm:w-[60%] typo-para-medium text-gray-700">
              {Number(item.timeseries?.totalCounts || 0)?.toLocaleString()}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default EvaluationTable;
