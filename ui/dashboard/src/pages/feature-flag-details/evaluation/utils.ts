export const getLogBase = (maxValue: number): number => {
  if (maxValue <= 5e2) return 2;
  if (maxValue <= 1e4) return 5;
  return 10;
};

export const symlog = (value: number, c = 1, base = 10): number => {
  return (
    (Math.sign(value) * Math.log(1 + Math.abs(value / c))) / Math.log(base)
  );
};

export const symlogInverse = (value: number, c = 1, base = 10): number => {
  return Math.sign(value) * ((Math.pow(base, Math.abs(value)) - 1) * c);
};

export const fomatLabel = (value: number): string => {
  const abs = Math.abs(value);
  if (abs >= 1e9) return `${Math.round(value / 1e9)}B`;
  if (abs >= 1e6) return `${Math.round(value / 1e6)}M`;
  if (abs >= 1e3) return `${Math.round(value / 1e3)}K`;
  return value.toFixed(2).replace(/\.00$/, '');
};
