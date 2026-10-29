import { describe, expect, it } from 'vitest';
import { escapeCSVValue, toCSVContent } from './csv-export';

describe('escapeCSVValue', () => {
  it('prefixes a string starting with a formula character with a single quote', () => {
    expect(escapeCSVValue('=1+1')).toBe("'=1+1");
    expect(escapeCSVValue('+1+1')).toBe("'+1+1");
    expect(escapeCSVValue('-1+1')).toBe("'-1+1");
    expect(escapeCSVValue('@SUM(A1:A2)')).toBe("'@SUM(A1:A2)");
    expect(escapeCSVValue('\t=1+1')).toBe("'\t=1+1");
    expect(escapeCSVValue('\r=1+1')).toBe("'\r=1+1");
  });

  it('prefixes and quotes a formula that also contains commas and quotes', () => {
    expect(escapeCSVValue('=HYPERLINK("https://example.com","open")')).toBe(
      '"\'=HYPERLINK(""https://example.com"",""open"")"'
    );
  });

  it('prefixes a string that only looks like a negative number', () => {
    expect(escapeCSVValue('-5')).toBe("'-5");
  });

  it('keeps plain text unchanged', () => {
    expect(escapeCSVValue('production')).toBe('production');
    expect(escapeCSVValue('My Project')).toBe('My Project');
    expect(escapeCSVValue('2026-01')).toBe('2026-01');
    expect(escapeCSVValue('123')).toBe('123');
  });

  it('keeps a formula character that is not at the start unchanged', () => {
    expect(escapeCSVValue('a=b')).toBe('a=b');
    expect(escapeCSVValue('foo-bar')).toBe('foo-bar');
    expect(escapeCSVValue('user@example.com')).toBe('user@example.com');
    expect(escapeCSVValue('1+1')).toBe('1+1');
  });

  it('keeps numbers unchanged, including negative numbers', () => {
    expect(escapeCSVValue(0)).toBe('0');
    expect(escapeCSVValue(42)).toBe('42');
    expect(escapeCSVValue(-5)).toBe('-5');
    expect(escapeCSVValue(12.5)).toBe('12.5');
  });

  it('writes null and undefined as an empty field', () => {
    expect(escapeCSVValue(null)).toBe('');
    expect(escapeCSVValue(undefined)).toBe('');
  });

  it('quotes fields containing a comma, a double quote or a line break', () => {
    expect(escapeCSVValue('a,b')).toBe('"a,b"');
    expect(escapeCSVValue('say "hi"')).toBe('"say ""hi"""');
    expect(escapeCSVValue('line1\nline2')).toBe('"line1\nline2"');
  });
});

describe('toCSVContent', () => {
  const header = ['Project', 'Environment', 'SDK', 'Month', 'MAU', 'Requests'];

  it('applies the field escaping to every cell of every row', () => {
    const content = toCSVContent(header, [
      [
        '=HYPERLINK("https://example.com","open")',
        '-production',
        'ANDROID',
        '2026-01',
        '10',
        '20'
      ],
      ['@project', '+staging', 'IOS', '2026-02', '0', '0']
    ]);

    expect(content.split('\n')).toEqual([
      'Project,Environment,SDK,Month,MAU,Requests',
      '"\'=HYPERLINK(""https://example.com"",""open"")",\'-production,ANDROID,2026-01,10,20',
      "'@project,'+staging,IOS,2026-02,0,0"
    ]);
  });

  it('leaves ordinary names and numeric values unchanged', () => {
    const content = toCSVContent(header, [
      ['My Project', 'production', 'ANDROID', '2026-01', 1234, 56789],
      ['Other, Inc.', 'staging', 'GO_SERVER', '2026-02', 0, -1]
    ]);

    expect(content).toBe(
      [
        'Project,Environment,SDK,Month,MAU,Requests',
        'My Project,production,ANDROID,2026-01,1234,56789',
        '"Other, Inc.",staging,GO_SERVER,2026-02,0,-1'
      ].join('\n')
    );
  });
});
