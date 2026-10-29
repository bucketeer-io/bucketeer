import Cookies from 'js-cookie';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { removeCookieState, setCookieState } from './index';

vi.mock('js-cookie', () => ({
  default: { get: vi.fn(), set: vi.fn(), remove: vi.fn() }
}));

const stubLocationProtocol = (protocol: string) =>
  vi.stubGlobal('window', { location: { protocol } });

afterEach(() => {
  vi.unstubAllGlobals();
  vi.clearAllMocks();
});

describe('setCookieState', () => {
  it('stores the state as a secure, lax cookie when served over https', () => {
    stubLocationProtocol('https:');
    setCookieState('abc');
    expect(Cookies.set).toHaveBeenCalledWith(
      'state',
      'abc',
      expect.objectContaining({ sameSite: 'lax', secure: true })
    );
  });

  it('omits the secure attribute when served over plain http', () => {
    stubLocationProtocol('http:');
    setCookieState('abc');
    expect(Cookies.set).toHaveBeenCalledWith(
      'state',
      'abc',
      expect.objectContaining({ sameSite: 'lax', secure: false })
    );
  });

  it('keeps the one hour expiry', () => {
    stubLocationProtocol('https:');
    setCookieState('abc');
    expect(Cookies.set).toHaveBeenCalledWith(
      'state',
      'abc',
      expect.objectContaining({ expires: 1 / 24 })
    );
  });
});

describe('removeCookieState', () => {
  it('removes the stored state', () => {
    removeCookieState();
    expect(Cookies.remove).toHaveBeenCalledWith('state');
  });
});
