import { ReactNode } from 'react';
import { COLORS } from 'constants/styles';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

export const formatLargeNumber = (value: number): string => {
  if (value >= 1e9) {
    const billions = value / 1e9;
    return `${billions >= 10 ? billions.toFixed(0) : billions.toFixed(1)}B`;
  }
  if (value >= 1e6) {
    const millions = value / 1e6;
    return `${millions >= 10 ? millions.toFixed(0) : millions.toFixed(1)}M`;
  }
  if (value >= 1e3) {
    const thousands = value / 1e3;
    return `${thousands >= 10 ? thousands.toFixed(0) : thousands.toFixed(1)}K`;
  }
  return String(value);
};

export const getColor = (i: number) => COLORS[i % COLORS.length];

interface ChartCardProps {
  title: string;
  description?: ReactNode;
  currentMonth?: number;
  lastMonth?: number;
  children: React.ReactNode;
}

export const ChartCard = ({
  title,
  description,
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
        <div className="flex-1 min-w-0">{children}</div>
      </div>
    </div>
  );
};
