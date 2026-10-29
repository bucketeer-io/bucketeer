import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { MarkdownContent } from './markdown-content';

const render = (source: string) =>
  renderToStaticMarkup(<MarkdownContent source={source} />);

describe('MarkdownContent raw HTML handling', () => {
  it('removes inline frames together with their srcdoc attribute', () => {
    const html = render(
      'Hello <iframe srcdoc="<script>alert(1)</script>"></iframe> world'
    );
    expect(html).not.toContain('<iframe');
    expect(html).not.toContain('srcdoc');
    expect(html).toContain('Hello');
    expect(html).toContain('world');
  });

  it('removes meta refresh, style and form elements', () => {
    const html = render(
      [
        '<meta http-equiv="refresh" content="0;url=https://example.com/">',
        '',
        '<style>body { display: none }</style>',
        '',
        '<form action="https://example.com/collect" method="post"><input name="q"></form>'
      ].join('\n')
    );
    expect(html).not.toContain('<meta');
    expect(html).not.toContain('<style');
    expect(html).not.toContain('<form');
    expect(html).not.toContain('http-equiv');
    expect(html).not.toContain('action=');
  });

  it('removes script elements including their content', () => {
    const html = render('before\n\n<script>alert(1)</script>\n\nafter');
    expect(html).not.toContain('<script');
    expect(html).not.toContain('alert(1)');
    expect(html).toContain('before');
    expect(html).toContain('after');
  });

  it('drops event handler attributes but keeps the element', () => {
    const html = render(
      '<img src="https://example.com/a.png" onerror="alert(1)" alt="logo">'
    );
    expect(html).toContain('<img');
    expect(html).toContain('src="https://example.com/a.png"');
    expect(html).toContain('alt="logo"');
    expect(html).not.toContain('onerror');
    expect(html).not.toContain('alert(1)');
  });

  it('drops links and images that use a javascript: or data: URL', () => {
    const html = render(
      [
        '[click me](javascript:alert(1))',
        '',
        '<a href="javascript:alert(2)">anchor</a>',
        '',
        '![pixel](data:image/svg+xml;base64,PHN2Zz48L3N2Zz4=)'
      ].join('\n')
    );
    expect(html).not.toContain('javascript:');
    expect(html).not.toContain('data:image');
    expect(html).toContain('click me');
    expect(html).toContain('anchor');
  });

  it('removes svg, object, embed, base and link elements', () => {
    const html = render(
      [
        '<svg onload="alert(1)"><circle r="1"></circle></svg>',
        '',
        '<object data="https://example.com/x.swf"></object>',
        '',
        '<embed src="https://example.com/x.swf">',
        '',
        '<base href="https://example.com/">',
        '',
        '<link rel="stylesheet" href="https://example.com/x.css">'
      ].join('\n')
    );
    for (const tag of ['<svg', '<object', '<embed', '<base', '<link']) {
      expect(html).not.toContain(tag);
    }
    expect(html).not.toContain('onload');
  });

  it('removes style attributes and unknown class names', () => {
    const html = render(
      '<div style="position:fixed;inset:0" class="overlay">covered</div>'
    );
    expect(html).not.toContain('style=');
    expect(html).not.toContain('overlay');
    expect(html).toContain('covered');
  });
});

describe('MarkdownContent formatting', () => {
  it('keeps headings, emphasis, lists and links', () => {
    const html = render(
      [
        '# Release notes',
        '',
        'Some **bold** and _italic_ text.',
        '',
        '- first',
        '- second',
        '',
        '1. one',
        '2. two',
        '',
        '[Docs](https://docs.bucketeer.io/) and <mailto@example.com>'
      ].join('\n')
    );
    expect(html).toContain('<h1');
    expect(html).toContain('Release notes');
    expect(html).toContain('<strong>bold</strong>');
    expect(html).toContain('<em>italic</em>');
    expect(html).toContain('<ul');
    expect(html).toContain('<ol');
    expect(html).toContain('href="https://docs.bucketeer.io/"');
    expect(html).toContain('href="mailto:mailto@example.com"');
  });

  it('keeps GFM tables, task lists and strikethrough', () => {
    const html = render(
      [
        '| Flag | Status |',
        '| ---- | ------ |',
        '| a    | on     |',
        '',
        '- [x] done',
        '- [ ] todo',
        '',
        '~~removed~~'
      ].join('\n')
    );
    expect(html).toContain('<table');
    expect(html).toContain('<th');
    expect(html).toContain('<td');
    expect(html).toContain('type="checkbox"');
    expect(html).toContain('checked=""');
    expect(html).toContain('<del>removed</del>');
  });

  it('keeps fenced code blocks and inline code', () => {
    const html = render(
      ['```ts', 'const enabled = true;', '```', '', 'Run `yarn build`.'].join(
        '\n'
      )
    );
    expect(html).toContain('<pre');
    expect(html).toContain('language-ts');
    expect(html).toContain('enabled');
    expect(html).toContain('<code>yarn build</code>');
  });

  it('keeps images with an http(s) source and simple inline HTML', () => {
    const html = render(
      [
        '![logo](https://example.com/logo.png)',
        '',
        'Press <kbd>Ctrl</kbd>+<kbd>C</kbd><br>then H<sub>2</sub>O'
      ].join('\n')
    );
    expect(html).toContain('<img src="https://example.com/logo.png"');
    expect(html).toContain('<kbd>Ctrl</kbd>');
    expect(html).toContain('<br');
    expect(html).toContain('<sub>2</sub>');
  });
});
