import { describe, expect, it } from 'vitest';
import { API_KEY_LABEL_ELLIPSIS, formatAPIKeyLabel } from './utils';

const MASKED_KEY = '0123....cdef';
const FULL_KEY =
  '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef';

describe('formatAPIKeyLabel', () => {
  it.each([
    {
      desc: 'shows the stored prefix followed by an ellipsis',
      input: { apiKey: MASKED_KEY, apiKeyPrefix: '01234567' },
      expected: `01234567${API_KEY_LABEL_ELLIPSIS}`
    },
    {
      desc: 'uses the prefix even when the masked key is empty',
      input: { apiKey: '', apiKeyPrefix: '01234567' },
      expected: `01234567${API_KEY_LABEL_ELLIPSIS}`
    },
    {
      desc: 'trims whitespace around the prefix',
      input: { apiKey: MASKED_KEY, apiKeyPrefix: ' 01234567 ' },
      expected: `01234567${API_KEY_LABEL_ELLIPSIS}`
    },
    {
      desc: 'falls back to the masked key when there is no prefix yet',
      input: { apiKey: MASKED_KEY, apiKeyPrefix: '' },
      expected: MASKED_KEY
    },
    {
      desc: 'falls back to the masked key when the prefix is missing',
      input: { apiKey: MASKED_KEY } as { apiKey: string; apiKeyPrefix: string },
      expected: MASKED_KEY
    },
    {
      desc: 'falls back to the masked key when the prefix is blank',
      input: { apiKey: MASKED_KEY, apiKeyPrefix: '   ' },
      expected: MASKED_KEY
    },
    {
      desc: 'shortens a long fallback value around the middle',
      input: { apiKey: FULL_KEY, apiKeyPrefix: '' },
      expected: '0123456...9abcdef'
    },
    {
      desc: 'returns an empty label when nothing is available',
      input: { apiKey: '', apiKeyPrefix: '' },
      expected: ''
    }
  ])('$desc', ({ input, expected }) => {
    expect(formatAPIKeyLabel(input)).toBe(expected);
  });

  it('never exposes more than the prefix when one is stored', () => {
    const label = formatAPIKeyLabel({
      apiKey: FULL_KEY,
      apiKeyPrefix: FULL_KEY.slice(0, 8)
    });
    expect(label).toBe(`${FULL_KEY.slice(0, 8)}${API_KEY_LABEL_ELLIPSIS}`);
    expect(label).not.toContain(FULL_KEY.slice(8));
  });
});
