import { DateTime } from 'luxon';
import { InsightsMonthlySummaryResponse } from '@types';

export type CSVValue = string | number | null | undefined;

// Spreadsheet applications evaluate a cell as a formula when its text starts
// with one of these characters, even when the file is a plain CSV.
const FORMULA_LEADING_CHARACTER = /^[=+\-@\t\r]/;

// Converts a single value into a CSV field. A string that a spreadsheet would
// evaluate as a formula is prefixed with a single quote so it is displayed as
// text; numbers are written as-is so they stay numeric. Fields containing a
// comma, a double quote or a line break are quoted.
export const escapeCSVValue = (value: CSVValue): string => {
  let s = String(value ?? '');
  if (typeof value === 'string' && FORMULA_LEADING_CHARACTER.test(s)) {
    s = `'${s}`;
  }
  return s.includes(',') || s.includes('"') || s.includes('\n')
    ? `"${s.replace(/"/g, '""')}"`
    : s;
};

export const toCSVContent = (
  header: CSVValue[],
  rows: CSVValue[][]
): string => {
  const toRow = (fields: CSVValue[]) => fields.map(escapeCSVValue).join(',');
  return [toRow(header), ...rows.map(toRow)].join('\n');
};

export const exportCSV = (
  header: CSVValue[],
  rows: CSVValue[][],
  filename: string
): void => {
  const csvContent = toCSVContent(header, rows);

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
