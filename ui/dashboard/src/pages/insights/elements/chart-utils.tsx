import { COLORS } from 'constants/styles';
import { useTranslation } from 'i18n';

export const formatLargeNumber = (value: number): string => {
  if (value >= 1e10) return `${(value / 1e10).toFixed(0)}B`;
  if (value >= 1e9) return `${(value / 1e9).toFixed(1)}B`;
  if (value >= 1e6) return `${(value / 1e6).toFixed(1)}M`;
  if (value >= 1e3) return `${(value / 1e3).toFixed(1)}k`;
  return String(value);
};

export const getColor = (i: number) => COLORS[i % COLORS.length];

interface ChartCardProps {
  title: string;
  currentMonth?: number;
  lastMonth?: number;
  children: React.ReactNode;
}

export const ChartCard = ({
  title,
  currentMonth,
  lastMonth,
  children
}: ChartCardProps) => {
  const { t } = useTranslation(['common']);

  const pctChange =
    currentMonth != null && lastMonth != null && lastMonth !== 0
      ? ((currentMonth - lastMonth) / lastMonth) * 100
      : null;

  return (
    <div className="bg-white">
      <h3 className="typo-para-small text-gray-600 mb-4">{title}</h3>
      <div className="flex items-center gap-x-5 min-w-0">
        {!!currentMonth && (
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
        <div className="flex-1 min-w-0">{children}</div>
      </div>
    </div>
  );
};
