import { describe, expect, it } from 'vitest';
import { convertJSONToRender, formatJSONWithIndent } from './utils';

const KEY_COLOR = '#e439ac';
const STRING_COLOR = '#40BF42';
const NUMBER_COLOR = '#64748B';

const render = (value: unknown) =>
  convertJSONToRender(formatJSONWithIndent(JSON.stringify(value))) ?? '';

// Removes the <span> wrappers that convertJSONToRender adds for highlighting,
// leaving only the text that came from the JSON input.
const stripHighlightSpans = (html: string) =>
  html.replace(/<span style="color: #[0-9A-Fa-f]{6}">|<\/span>/g, '');

const countMatches = (text: string, pattern: RegExp) =>
  (text.match(pattern) ?? []).length;

describe('convertJSONToRender', () => {
  it('returns null for empty input', () => {
    expect(convertJSONToRender(null)).toBeNull();
    expect(convertJSONToRender('')).toBeNull();
  });

  it('highlights keys, string values and numeric values', () => {
    const html = render({ name: 'flag', version: 3 });

    expect(html).toContain(`<span style="color: ${KEY_COLOR}">"name"</span>`);
    expect(html).toContain(
      `<span style="color: ${STRING_COLOR}">"flag"</span>`
    );
    expect(html).toContain(`<span style="color: ${NUMBER_COLOR}">3</span>`);
  });

  it('adds nothing but the highlighting spans to plain JSON', () => {
    const json = formatJSONWithIndent(
      JSON.stringify({
        id: 'feature-1',
        name: 'flag',
        description: 'A plain description, with punctuation: 1, 2, 3.',
        enabled: true,
        tags: ['web', 'ios'],
        version: 12,
        archived: null,
        variations: [{ id: 'a', value: 'true' }]
      })
    );

    expect(json).not.toBeNull();
    expect(stripHighlightSpans(convertJSONToRender(json) ?? '')).toBe(json);
  });

  it('renders an element in a string value as text', () => {
    const html = render({ description: '<img src=x onerror=alert(1)>' });

    expect(html).not.toContain('<img');
    expect(html).toContain(
      `<span style="color: ${STRING_COLOR}">"&lt;img src=x onerror=alert(1)&gt;"</span>`
    );
    expect(stripHighlightSpans(html)).not.toMatch(/[<>]/);
  });

  it('renders an element in a key as text', () => {
    const html = render({ '<b>name</b>': 'flag' });

    expect(html).not.toContain('<b>');
    expect(html).toContain(
      `<span style="color: ${KEY_COLOR}">"&lt;b&gt;name&lt;/b&gt;"</span>`
    );
    expect(stripHighlightSpans(html)).not.toMatch(/[<>]/);
  });

  it('cannot close a highlighting span from inside a value', () => {
    const html = render({
      description: '</span><script>alert(1)</script><span>'
    });

    expect(html).not.toContain('<script');
    expect(html).toContain('&lt;/span&gt;&lt;script&gt;');
    expect(countMatches(html, /<\/span>/g)).toBe(
      countMatches(html, /<span style=/g)
    );
    expect(stripHighlightSpans(html)).not.toMatch(/[<>]/);
  });

  it('shows character references in the data literally', () => {
    const html = render({ description: 'a &lt; b & c &#60;' });

    expect(html).toContain(
      `<span style="color: ${STRING_COLOR}">"a &amp;lt; b &amp; c &amp;#60;"</span>`
    );
  });

  it('keeps highlighting a value that mixes escaped quotes and markup', () => {
    const html = render({ description: 'say "<hi>" & bye', count: 2 });

    expect(html).toContain(
      `<span style="color: ${STRING_COLOR}">"say \\"&lt;hi&gt;\\" &amp; bye"</span>`
    );
    expect(html).toContain(`<span style="color: ${NUMBER_COLOR}">2</span>`);
    expect(stripHighlightSpans(html)).not.toMatch(/[<>]/);
  });

  it('escapes every line of a multi-line document', () => {
    const html = render({
      name: '<i>one</i>',
      nested: { description: '<i>two</i>', values: ['<i>three</i>'] }
    });

    expect(html).not.toContain('<i>');
    expect(countMatches(html, /&lt;i&gt;/g)).toBe(3);
    expect(countMatches(html, /&lt;\/i&gt;/g)).toBe(3);
    expect(stripHighlightSpans(html)).not.toMatch(/[<>]/);
  });
});
