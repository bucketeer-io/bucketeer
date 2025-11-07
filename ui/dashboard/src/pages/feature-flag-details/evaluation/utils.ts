export const getLogBase = (maxValue: number): number => {
  if (maxValue <= 5e2) return 2;
  if (maxValue <= 1e4) return 5;
  return 10;
};

/**
 * Applies a symmetric logarithmic (symlog) transformation to the input value.
 * This transformation behaves linearly near zero and logarithmically for large magnitudes,
 * preserving the sign of the input. It is useful for visualizing data that spans several orders of magnitude,
 * including both positive and negative values.
 *
 * @param {number} value - The input value to transform.
 * @param {number} [c=1] - The linear region threshold; controls the transition between linear and logarithmic behavior.
 * @param {number} [base=10] - The logarithm base to use for the transformation.
 * @returns {number} The symmetrically log-transformed value.
 */

export const symlog = (value: number, c = 1, base = 10): number => {
  return (
    (Math.sign(value) * Math.log(1 + Math.abs(value / c))) / Math.log(base)
  );
};

/**
 * Reverses the symlog transformation, converting a symlog-scaled value back to its original scale.
 *
 * @param value The symlog-transformed value to invert.
 * @param c The linear scaling constant used in the original symlog transformation (default is 1).
 * @param base The logarithm base used in the original symlog transformation (default is 10).
 * @returns The original value before the symlog transformation.
 */

export const symlogInverse = (value: number, c = 1, base = 10): number => {
  return Math.sign(value) * ((Math.pow(base, Math.abs(value)) - 1) * c);
};

export const formatLabel = (value: number): string => {
  const abs = Math.round(Math.abs(value));
  if (abs >= 1e9)
    return `${Math.sign(value) * Math.round(Math.abs(value) / 1e9)}B`;
  if (abs >= 1e6)
    return `${Math.sign(value) * Math.round(Math.abs(value) / 1e6)}M`;
  if (abs >= 1e3)
    return `${Math.sign(value) * Math.round(Math.abs(value) / 1e3)}K`;
  return value.toFixed(2).replace(/\.00$/, '');
};
