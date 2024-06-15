import dayjs from 'dayjs';

export function addDays(date: Date, days: number) {
  const result = new Date(date);
  result.setDate(result.getDate() + days);
  return result;
}

interface FormatDate {
  date: Date;
  showDateOnly?: boolean;
}

export function formatDate({ date, showDateOnly }: FormatDate): string {
  if (showDateOnly) {
    return dayjs(date).format('YYYY/MM/D');
  }
  return dayjs(date).format('YYYY/MM/D H:mm');
}
