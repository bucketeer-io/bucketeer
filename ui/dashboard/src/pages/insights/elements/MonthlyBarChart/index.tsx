import { useMemo } from 'react';
import { Bar } from 'react-chartjs-2';
import { useTranslation } from 'i18n';
import { InsightsMonthlySummaryResponse } from '@types';
import Spinner from 'components/spinner';
import EmptyState from 'elements/empty-state';
import { ChartCard, formatLargeNumber, getColor } from '../chart-utils';

interface MonthlyBarChartProps {
  title: string;
  summary?: InsightsMonthlySummaryResponse;
  isLoading: boolean;
  field: 'mau' | 'requests';
  months: string[]; // yearmonth strings e.g. ['202401', ...]
  labels: string[]; // formatted labels e.g. ['Jan 24', ...]
}

const MonthlyBarChart = ({
  title,
  summary,
  isLoading,
  field,
  months,
  labels
}: MonthlyBarChartProps) => {
  const { t } = useTranslation(['common']);
  const datasets = useMemo(() => {
    if (!summary?.series?.length) return [];
    const totals = months.map(m =>
      summary.series.reduce((sum, series) => {
        const dp = series.data.find(d => d.yearmonth === m);
        return sum + (dp ? Number(dp[field]) : 0);
      }, 0)
    );
    const baseColor = getColor(0);
    const backgroundColor = totals.map((_, i) =>
      i === totals.length - 1 ? baseColor : `${baseColor}66`
    );
    return [
      {
        label: field === 'mau' ? t('insights.mau') : t('insights.requests'),
        data: totals,
        backgroundColor,
        borderRadius: 4
      }
    ];
  }, [summary, field, months, t]);

  const totals = datasets[0]?.data ?? [];
  const currentMonth =
    totals.length > 0 ? totals[totals.length - 1] : undefined;
  const lastMonth = totals.length > 1 ? totals[totals.length - 2] : undefined;

  return (
    <ChartCard title={title} currentMonth={currentMonth} lastMonth={lastMonth}>
      {isLoading ? (
        <div className="h-[300px] flex items-center justify-center">
          <Spinner />
        </div>
      ) : !datasets.length ? (
        <div className="h-[300px] flex items-center justify-center">
          <EmptyState.Root variant="no-data" size="sm">
            <EmptyState.Illustration />
            <EmptyState.Body>
              <EmptyState.Title>{t('no-data')}</EmptyState.Title>
            </EmptyState.Body>
          </EmptyState.Root>
        </div>
      ) : (
        <div className="h-[300px]">
          <Bar
            data={{ labels, datasets }}
            options={{
              responsive: true,
              maintainAspectRatio: false,
              plugins: { legend: { display: false } },
              scales: {
                x: {
                  grid: { display: false },
                  ticks: { color: '#94A3B8', font: { size: 12 } }
                },
                y: {
                  grid: { color: '#E2E8F0' },
                  ticks: {
                    color: '#94A3B8',
                    font: { size: 12 },
                    callback: (value: number | string) =>
                      formatLargeNumber(Number(value))
                  }
                }
              }
            }}
          />
        </div>
      )}
    </ChartCard>
  );
};

export default MonthlyBarChart;
