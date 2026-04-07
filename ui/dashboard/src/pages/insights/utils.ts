import { DateTime } from 'luxon';
import { InsightSourceId, InsightApiId } from '@types';

export type TimeRangePreset = '1h' | '6h' | '24h' | '7d' | '30d' | 'this_month';

export interface InsightsFilters {
  projectId: string;
  environmentId: string;
  sourceId: InsightSourceId | '';
  apiId: InsightApiId | '';
  timeRange: TimeRangePreset;
  customStartAt?: string;
  customEndAt?: string;
}

const presetMap: Record<TimeRangePreset, (now: DateTime) => DateTime> = {
  '1h': now => now.minus({ hours: 1 }),
  '6h': now => now.minus({ hours: 6 }),
  '24h': now => now.minus({ hours: 24 }),
  '7d': now => now.minus({ days: 7 }),
  '30d': now => now.minus({ days: 30 }),
  this_month: now => now.startOf('month')
};

export const computeTimeRange = (
  preset: TimeRangePreset,
  customStartAt?: string,
  customEndAt?: string
): { startAt: string; endAt: string } => {
  if (customStartAt && customEndAt) {
    return { startAt: customStartAt, endAt: customEndAt };
  }
  const now = DateTime.now();
  const startAt = presetMap[preset](now);
  return {
    startAt: String(Math.floor(startAt.toSeconds())),
    endAt: String(Math.floor(now.toSeconds()))
  };
};

export const normalizeEnvId = (id: string) => (id === '' ? 'production' : id);

export const formatYAxis = (value: number): string => {
  const abs = Math.abs(value);
  if (abs === 0) return '0';
  if (abs >= 1e10) return `${(value / 1e9).toFixed(0)}B`;
  if (abs >= 1e9) return `${(value / 1e9).toFixed(1)}B`;
  if (abs >= 1e8) return `${(value / 1e6).toFixed(0)}M`;
  if (abs >= 1e7) return `${(value / 1e6).toFixed(0)}M`;
  if (abs >= 1e6) return `${(value / 1e6).toFixed(1)}M`;
  if (abs >= 1e3) return `${(value / 1e3).toFixed(0)}K`;
  if (abs >= 1) return value.toFixed(0);
  if (abs >= 0.1) return value.toFixed(1);
  if (abs >= 0.01) return value.toFixed(2);
  if (abs >= 0.001) return value.toFixed(3);
  if (abs >= 0.0001) return value.toFixed(4);
  if (abs >= 0.00001) return value.toFixed(5);
  return value.toFixed(6);
};
