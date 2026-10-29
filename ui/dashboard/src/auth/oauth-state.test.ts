import { describe, expect, it } from 'vitest';
import { generateOAuthState, isValidOAuthState } from './oauth-state';

describe('isValidOAuthState', () => {
  it('accepts a callback state equal to the stored state', () => {
    expect(isValidOAuthState('abc', 'abc')).toBe(true);
  });

  it('rejects when both the callback state and the stored state are missing', () => {
    expect(isValidOAuthState(undefined, undefined)).toBe(false);
    expect(isValidOAuthState(null, undefined)).toBe(false);
  });

  it('rejects a missing callback state even when a state was stored', () => {
    expect(isValidOAuthState(undefined, 'abc')).toBe(false);
    expect(isValidOAuthState(null, 'abc')).toBe(false);
  });

  it('rejects a callback state when no state was stored', () => {
    expect(isValidOAuthState('abc', undefined)).toBe(false);
  });

  it('rejects empty strings on either side', () => {
    expect(isValidOAuthState('', '')).toBe(false);
    expect(isValidOAuthState('', 'abc')).toBe(false);
    expect(isValidOAuthState('abc', '')).toBe(false);
  });

  it('rejects a callback state different from the stored state', () => {
    expect(isValidOAuthState('abc', 'abd')).toBe(false);
    expect(isValidOAuthState('ABC', 'abc')).toBe(false);
    expect(isValidOAuthState(' abc', 'abc')).toBe(false);
  });

  it('rejects a repeated query parameter even when one value matches', () => {
    expect(isValidOAuthState(['abc', 'xyz'], 'abc')).toBe(false);
    expect(isValidOAuthState(['abc'], 'abc')).toBe(false);
  });
});

describe('generateOAuthState', () => {
  it('returns a version 4 UUID', () => {
    expect(generateOAuthState()).toMatch(
      /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/
    );
  });

  it('returns a different value on every call', () => {
    const states = new Set(
      Array.from({ length: 1000 }, () => generateOAuthState())
    );
    expect(states.size).toBe(1000);
  });

  it('is accepted by isValidOAuthState when echoed back unchanged', () => {
    const state = generateOAuthState();
    expect(isValidOAuthState(state, state)).toBe(true);
  });
});
