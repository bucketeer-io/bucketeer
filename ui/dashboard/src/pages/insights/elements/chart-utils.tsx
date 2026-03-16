import { useTranslation } from 'i18n';
import Table from 'components/table';

export const formatLargeNumber = (value: number): string => {
  if (value >= 10_000_000_000) return `${(value / 1_000_000_000).toFixed(0)}B`;
  if (value >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(1)}B`;
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`;
  if (value >= 1_000) return `${(value / 1_000).toFixed(1)}k`;
  return String(value);
};

const CHART_COLORS = [
  '#7C3AED',
  '#10B981',
  '#F59E0B',
  '#EF4444',
  '#3B82F6',
  '#EC4899',
  '#6366F1',
  '#14B8A6'
];

export const getColor = (i: number) => CHART_COLORS[i % CHART_COLORS.length];

interface LegendTableProps {
  datasets: { label?: string; data: number[] }[];
  formatter?: (value: number) => string;
}

export const LegendTable = ({ datasets, formatter }: LegendTableProps) => {
  const fmt = formatter ?? ((v: number) => v.toFixed(2));
  const { t } = useTranslation(['common']);
  if (!datasets.length) return null;
  return (
    <div className="overflow-y-auto max-h-[256px]">
      <Table.Root>
        <Table.Header className="sticky top-0 z-10 bg-white">
          <Table.Row>
            <Table.Head className="max-w-[200px]">
              {t('insights.series')}
            </Table.Head>
            <Table.Head align="right">{t('insights.min')}</Table.Head>
            <Table.Head align="right">{t('insights.max')}</Table.Head>
            <Table.Head align="right">{t('insights.avg')}</Table.Head>
            <Table.Head align="right">{t('insights.last')}</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {datasets.map((ds, i) => {
            const nums = ds.data.filter(
              v => typeof v === 'number' && !isNaN(v)
            );
            const min = nums.length ? Math.min(...nums) : 0;
            const max = nums.length ? Math.max(...nums) : 0;
            const avg = nums.length
              ? nums.reduce((a, b) => a + b, 0) / nums.length
              : 0;
            const last = nums.length ? nums[nums.length - 1] : 0;
            return (
              <Table.Row key={i}>
                <Table.Cell className="pl-4 max-w-[200px]">
                  <div className="flex items-center gap-2">
                    <span
                      className="inline-block w-3 h-3 rounded-full flex-shrink-0"
                      style={{ backgroundColor: getColor(i) }}
                    />
                    <span className="text-gray-700 truncate typo-para-medium">
                      {ds.label ?? `Series ${i + 1}`}
                    </span>
                  </div>
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700"
                >
                  {fmt(min)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700"
                >
                  {fmt(max)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700"
                >
                  {fmt(avg)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700"
                >
                  {fmt(last)}
                </Table.Cell>
              </Table.Row>
            );
          })}
        </Table.Body>
      </Table.Root>
    </div>
  );
};

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
