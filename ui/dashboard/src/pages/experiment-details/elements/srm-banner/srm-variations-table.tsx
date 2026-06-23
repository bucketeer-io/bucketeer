import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { FeatureVariation, SrmVariation } from '@types';
import { cn, getVariationColor } from 'utils/style';
import {
  ResultCell,
  ResultHeaderCell
} from 'pages/experiment-details/collection-loader/results/goal-results/goal-results-table-element';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import NameWithTooltip from 'elements/name-with-tooltip';

const VARIATION_MIN_SIZE = 255;
const NUMERIC_MIN_SIZE = 171.5;

const formatNumber = (value: number, fractionDigits = 1) =>
  value.toLocaleString(undefined, {
    minimumFractionDigits: fractionDigits,
    maximumFractionDigits: fractionDigits
  });

const SrmVariationsTable = ({
  variations,
  experimentVariations,
  className
}: {
  variations: SrmVariation[];
  experimentVariations: FeatureVariation[];
  className?: string;
}) => {
  const { t } = useTranslation(['table']);

  const headerList = useMemo(
    () => [
      { name: 'variation', tooltipKey: '', minSize: VARIATION_MIN_SIZE },
      {
        name: 'srm.table.observed-users',
        tooltipKey: 'srm.table.observed-users-tooltip',
        minSize: NUMERIC_MIN_SIZE
      },
      {
        name: 'srm.table.expected-users',
        tooltipKey: 'srm.table.expected-users-tooltip',
        minSize: NUMERIC_MIN_SIZE
      },
      {
        name: 'srm.table.expected-weight',
        tooltipKey: 'srm.table.expected-weight-tooltip',
        minSize: NUMERIC_MIN_SIZE
      }
    ],
    []
  );

  const variationIndexById = useMemo(() => {
    const map = new Map<string, number>();
    experimentVariations.forEach((v, i) => map.set(v.id, i));
    return map;
  }, [experimentVariations]);

  const variationById = useMemo(() => {
    const map = new Map<string, FeatureVariation>();
    experimentVariations.forEach(v => map.set(v.id, v));
    return map;
  }, [experimentVariations]);

  const resolveName = (variationId: string) => {
    const v = variationById.get(variationId);
    return v?.name || v?.value || variationId;
  };

  return (
    <div className={cn('min-w-fit', className)}>
      <div className="flex w-full">
        {headerList.map((item, index) => (
          <ResultHeaderCell
            key={index}
            text={t(`table:results.${item.name}`)}
            tooltip={
              item.tooltipKey ? t(`table:results.${item.tooltipKey}`) : ''
            }
            isShowIcon={index > 0}
            minSize={item.minSize}
          />
        ))}
      </div>
      <div className="divide-y divide-gray-300">
        {variations.map((v, rowIndex) => {
          const name = resolveName(v.variationId);
          // Prefer the variation's position in experiment.variations so
          // colors stay aligned with the conversion-rate table. Fall back
          // to the row index when the variation isn't in the experiment
          // (demo fixtures + the backend's "leaked traffic" case), so a
          // polygon still renders.
          const experimentIndex = variationIndexById.get(v.variationId);
          const colorIndex =
            typeof experimentIndex === 'number' ? experimentIndex : rowIndex;
          return (
            <div key={v.variationId} className="flex items-center w-full">
              <div
                className="flex items-center size-fit w-full px-4 py-5 gap-x-2 text-gray-500"
                style={{ minWidth: VARIATION_MIN_SIZE }}
              >
                <Polygon
                  className="border-none size-3 shrink-0"
                  style={{
                    background: getVariationColor(colorIndex),
                    zIndex: colorIndex
                  }}
                />
                <NameWithTooltip
                  id={`srm-${v.variationId}`}
                  maxLines={1}
                  content={
                    <NameWithTooltip.Content
                      content={name}
                      id={`srm-${v.variationId}`}
                    />
                  }
                  trigger={
                    <NameWithTooltip.Trigger
                      id={`srm-${v.variationId}`}
                      name={name}
                      maxLines={1}
                      haveAction={false}
                      className="typo-para-medium text-gray-800"
                    />
                  }
                />
              </div>
              <ResultCell
                value={Number(v.observedUserCount).toLocaleString()}
                minSize={NUMERIC_MIN_SIZE}
              />
              <ResultCell
                value={formatNumber(v.expectedUserCount)}
                minSize={NUMERIC_MIN_SIZE}
              />
              <ResultCell
                value={`${formatNumber(v.expectedWeight * 100)}%`}
                minSize={NUMERIC_MIN_SIZE}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default SrmVariationsTable;
