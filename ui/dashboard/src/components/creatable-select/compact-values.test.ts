import { describe, expect, it } from 'vitest';
import { getVisibleValueCount } from './compact-values';

describe('getVisibleValueCount', () => {
  it('shows more values as the field grows', () => {
    expect(getVisibleValueCount(200, 10)).toBe(2);
    expect(getVisibleValueCount(360, 10)).toBe(2);
    expect(getVisibleValueCount(640, 10)).toBe(4);
  });

  it('does not reserve a counter when every value fits', () => {
    expect(getVisibleValueCount(444, 3)).toBe(3);
    expect(getVisibleValueCount(443, 3)).toBe(2);
  });

  it('keeps at least two values without exceeding the selection', () => {
    expect(getVisibleValueCount(640, 0)).toBe(0);
    expect(getVisibleValueCount(0, 0)).toBe(0);
    expect(getVisibleValueCount(100, 1)).toBe(1);
    expect(getVisibleValueCount(100, 2)).toBe(2);
    expect(getVisibleValueCount(0, 10)).toBe(2);
    expect(getVisibleValueCount(100, 10)).toBe(2);
  });

  it('reserves additional space for larger overflow counts', () => {
    expect(getVisibleValueCount(480, 9)).toBe(3);
    expect(getVisibleValueCount(480, 100)).toBe(2);
  });
});
