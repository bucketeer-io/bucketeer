import { useTranslation } from 'i18n';
import Table from 'components/table';
import { getColor } from '../chart-utils';

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
                  className="typo-para-medium text-gray-700 pr-4"
                >
                  {fmt(min)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700 pr-4"
                >
                  {fmt(max)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700 pr-4"
                >
                  {fmt(avg)}
                </Table.Cell>
                <Table.Cell
                  align="right"
                  className="typo-para-medium text-gray-700 pr-4"
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
