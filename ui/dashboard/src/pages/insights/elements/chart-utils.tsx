import { memo, ReactNode } from 'react';
import { COLORS } from 'constants/styles';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { formatYAxis } from '../utils';

export { formatYAxis };
export const formatLargeNumber = formatYAxis;

export const getColor = (i: number) => COLORS[i % COLORS.length];

interface ChartCardProps {
  title: string;
  description?: ReactNode;
  currentMonth?: number;
  lastMonth?: number;
  children: ReactNode;
  legendDatasets?: { label?: string; backgroundColor: string | string[] }[];
}

export const ChartCard = memo(function ChartCard({
  title,
  description,
  currentMonth,
  lastMonth,
  children,
  legendDatasets
}: ChartCardProps) {
  const { t } = useTranslation(['common']);

  const pctChange =
    currentMonth != null && lastMonth != null && lastMonth !== 0
      ? ((currentMonth - lastMonth) / lastMonth) * 100
      : null;

  return (
    <div className="bg-white">
      <div className="flex items-center gap-x-1 mb-4">
        <h3 className="typo-para-small text-gray-600">{title}</h3>

        {!!description && (
          <Tooltip
            content={description}
            trigger={
              <button type="button" className="flex-center size-fit ">
                <Icon icon={IconInfo} size="xs" color="gray-500" />
              </button>
            }
          />
        )}
      </div>

      <div className="flex items-center gap-x-5 min-w-0">
        {currentMonth != null && (
          <div className="flex-shrink-0">
            <p className="typo-para-small text-gray-500">
              {t('insights.current-month')}
            </p>
            <h1 className="text-4xl font-bold py-3">
              {formatLargeNumber(currentMonth)}
            </h1>
            {pctChange != null && (
              <p
                className={`typo-para-small ${pctChange >= 0 ? 'text-green-600' : 'text-red-500'}`}
              >
                {pctChange >= 0 ? '+' : ''}
                {pctChange.toFixed(1)}% {t('insights.vs-last-month')}
              </p>
            )}
          </div>
        )}
        <div className="flex-1 min-w-0">
          {children}
          {legendDatasets && legendDatasets.length > 0 && (
            <div className="flex flex-wrap justify-center gap-x-4 gap-y-1 mt-3">
              {legendDatasets.map((ds, i) => {
                const color = Array.isArray(ds.backgroundColor)
                  ? ds.backgroundColor[0]
                  : ds.backgroundColor;
                return (
                  <div key={i} className="flex items-center gap-1.5">
                    <span
                      className="inline-block w-10 h-2.5 rounded-sm flex-shrink-0"
                      style={{ backgroundColor: color as string }}
                    />
                    <span className="typo-para-small text-gray-500 truncate max-w-[120px]">
                      {ds.label}
                    </span>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
});
