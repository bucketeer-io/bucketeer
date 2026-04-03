import { DateTime } from 'luxon';
import { InsightsMonthlySummaryResponse } from '@types';

type CSVValue = string | number | null | undefined;

export const exportCSV = (
  header: CSVValue[],
  rows: CSVValue[][],
  filename: string
): void => {
  const escape = (value: CSVValue): string => {
    const s = String(value ?? '');
    return s.includes(',') || s.includes('"') || s.includes('\n')
      ? `"${s.replace(/"/g, '""')}"`
      : s;
  };

  const toRow = (fields: CSVValue[]) => fields.map(escape).join(',');

  const csvContent = [toRow(header), ...rows.map(toRow)].join('\n');

  const blob = new Blob(['\uFEFF' + csvContent], {
    type: 'text/csv;charset=utf-8;'
  });

  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');

  a.href = url;
  a.download = `${filename}-${DateTime.now().toFormat('yyyyMMdd')}.csv`;

  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);

  URL.revokeObjectURL(url);
};

export const exportMonthlySummaryCSV = (
  summary: InsightsMonthlySummaryResponse,
  months: string[]
): void => {
  if (!summary.series?.length) return;

  const rows: (string | number)[][] = [];

  for (const series of summary.series) {
    const dataMap = new Map(series.data.map(d => [d.yearmonth, d]));

    for (const m of months) {
      const dp = dataMap.get(m);

      rows.push([
        series.projectName,
        series.environmentName,
        series.sourceId,
        DateTime.fromFormat(m, 'yyyyMM').toFormat('yyyy-MM'),
        dp?.mau ?? 0,
        dp?.requests ?? 0
      ]);
    }
  }

  exportCSV(
    ['Project', 'Environment', 'SDK', 'Month', 'MAU', 'Requests'],
    rows,
    'insights-monthly'
  );
};
