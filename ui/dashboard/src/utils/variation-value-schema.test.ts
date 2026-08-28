import { describe, expect, it } from 'vitest';
import { VariationValueSchema } from '@types';
import {
  createValueValidator,
  getSupportedSchemaTypes,
  isSchemaSupported,
  validateSchemaDefinition
} from './variation-value-schema';

const enumSchema = (values: string[]): VariationValueSchema => ({
  type: 'ENUM',
  enumValidator: { values }
});

const regexSchema = (pattern: string): VariationValueSchema => ({
  type: 'REGEX',
  regexValidator: { pattern }
});

const jsonSchema = (schema: object): VariationValueSchema => ({
  type: 'JSON_SCHEMA',
  jsonSchemaValidator: { schema: JSON.stringify(schema) }
});

describe('getSupportedSchemaTypes / isSchemaSupported', () => {
  it('matches the backend v1 type matrix', () => {
    expect(getSupportedSchemaTypes('STRING')).toEqual(['ENUM', 'REGEX']);
    expect(getSupportedSchemaTypes('NUMBER')).toEqual(['ENUM']);
    expect(getSupportedSchemaTypes('JSON')).toEqual(['JSON_SCHEMA']);
    expect(getSupportedSchemaTypes('BOOLEAN')).toEqual([]);
    expect(getSupportedSchemaTypes('YAML')).toEqual([]);
  });

  it('reports support only for types with at least one validator', () => {
    expect(isSchemaSupported('STRING')).toBe(true);
    expect(isSchemaSupported('BOOLEAN')).toBe(false);
  });
});

describe('validateSchemaDefinition', () => {
  it('rejects schema types unsupported for the variation type', () => {
    expect(validateSchemaDefinition(regexSchema('a+'), 'NUMBER')).toBe(
      'type-unsupported'
    );
  });

  it('rejects empty enums', () => {
    expect(validateSchemaDefinition(enumSchema([]), 'STRING')).toBe(
      'enum-empty'
    );
  });

  it('accepts plain decimal enum values for NUMBER flags', () => {
    expect(
      validateSchemaDefinition(enumSchema(['1', '-1.5', '1e3', '.5']), 'NUMBER')
    ).toBeNull();
  });

  it('rejects non-decimal syntaxes the backend parser rejects', () => {
    for (const value of ['0x10', '0b10', 'Infinity', ' 1', '']) {
      expect(validateSchemaDefinition(enumSchema([value]), 'NUMBER')).toBe(
        'enum-not-number'
      );
    }
  });

  it('rejects empty regex patterns', () => {
    expect(validateSchemaDefinition(regexSchema(''), 'STRING')).toBe(
      'regex-empty'
    );
  });

  it('accepts Go-only RE2 constructs like inline flags', () => {
    expect(
      validateSchemaDefinition(regexSchema('(?i)^a+$'), 'STRING')
    ).toBeNull();
  });

  it('rejects Perl-only constructs like lookaheads, matching Go RE2', () => {
    expect(validateSchemaDefinition(regexSchema('(?=a)'), 'STRING')).toBe(
      'regex-invalid'
    );
  });

  it('rejects empty and malformed JSON Schemas', () => {
    expect(
      validateSchemaDefinition(
        { type: 'JSON_SCHEMA', jsonSchemaValidator: { schema: '  ' } },
        'JSON'
      )
    ).toBe('json-schema-empty');
    expect(
      validateSchemaDefinition(
        { type: 'JSON_SCHEMA', jsonSchemaValidator: { schema: '{ not json' } },
        'JSON'
      )
    ).toBe('json-schema-invalid');
  });

  it('accepts a valid JSON Schema', () => {
    expect(
      validateSchemaDefinition(jsonSchema({ type: 'object' }), 'JSON')
    ).toBeNull();
  });
});

describe('createValueValidator: ENUM', () => {
  it('validates string values by exact match', () => {
    const validate = createValueValidator(
      enumSchema(['ssh', 'email']),
      'STRING'
    )!;
    expect(validate('ssh').valid).toBe(true);
    expect(validate('SSH').valid).toBe(false);
    expect(validate('').valid).toBe(false);
  });

  it('validates number values numerically but rejects non-decimal syntax', () => {
    const validate = createValueValidator(enumSchema(['16', '1.5']), 'NUMBER')!;
    expect(validate('16').valid).toBe(true);
    expect(validate('16.0').valid).toBe(true);
    // Number('0x10') === 16, but the backend's strconv.ParseFloat rejects it.
    expect(validate('0x10').valid).toBe(false);
    expect(validate('2').valid).toBe(false);
  });
});

describe('createValueValidator: REGEX', () => {
  it('returns null for patterns that do not compile as RE2', () => {
    expect(createValueValidator(regexSchema('(?=a)'), 'STRING')).toBeNull();
  });

  it('supports Go-only inline flags', () => {
    const validate = createValueValidator(regexSchema('(?i)^ssh$'), 'STRING')!;
    expect(validate('SSH').valid).toBe(true);
    expect(validate('email').valid).toBe(false);
  });

  it('matches unanchored, mirroring the backend regexp.MatchString', () => {
    const validate = createValueValidator(regexSchema('b+'), 'STRING')!;
    expect(validate('abc').valid).toBe(true);
    expect(validate('ac').valid).toBe(false);
  });
});

describe('createValueValidator: JSON_SCHEMA', () => {
  const schema = jsonSchema({
    type: 'object',
    required: ['name'],
    additionalProperties: false,
    properties: {
      name: { type: 'string' },
      theme: { enum: ['light', 'dark'] }
    }
  });

  it('returns null when the schema itself does not compile', () => {
    const invalid: VariationValueSchema = {
      type: 'JSON_SCHEMA',
      jsonSchemaValidator: { schema: '{ not json' }
    };
    expect(createValueValidator(invalid, 'JSON')).toBeNull();
  });

  it('accepts a conforming value without detail', () => {
    const validate = createValueValidator(schema, 'JSON')!;
    expect(validate('{"name":"a","theme":"dark"}')).toEqual({ valid: true });
  });

  it('reports the failing path for missing required properties', () => {
    const validate = createValueValidator(schema, 'JSON')!;
    const result = validate('{}');
    expect(result.valid).toBe(false);
    expect(result.detail).toBe("/ must have required property 'name'");
  });

  it('names the offending additional property', () => {
    const validate = createValueValidator(schema, 'JSON')!;
    const result = validate('{"name":"a","extraField":"x"}');
    expect(result.valid).toBe(false);
    expect(result.detail).toContain(
      "must NOT have additional properties ('extraField')"
    );
  });

  it('lists the allowed values for enum violations', () => {
    const validate = createValueValidator(schema, 'JSON')!;
    const result = validate('{"name":"a","theme":"blue"}');
    expect(result.valid).toBe(false);
    expect(result.detail).toContain(
      '/theme must be equal to one of the allowed values'
    );
    expect(result.detail).toContain('"light", "dark"');
  });

  it('truncates long enum allowed-value lists', () => {
    const manyValues = jsonSchema({
      type: 'object',
      properties: { size: { enum: ['a', 'b', 'c', 'd', 'e', 'f', 'g'] } }
    });
    const validate = createValueValidator(manyValues, 'JSON')!;
    const result = validate('{"size":"z"}');
    expect(result.valid).toBe(false);
    expect(result.detail).toContain('"a", "b", "c", "d", "e", +2 more');
  });

  it('reports multiple violations, capped with a remainder count', () => {
    const multi = jsonSchema({
      type: 'object',
      required: ['a', 'b', 'c', 'd'],
      properties: {
        a: { type: 'string' },
        b: { type: 'string' },
        c: { type: 'string' },
        d: { type: 'string' }
      }
    });
    const validate = createValueValidator(multi, 'JSON')!;
    const result = validate('{}');
    expect(result.valid).toBe(false);
    expect(result.detail).toBe(
      "/ must have required property 'a'; " +
        "/ must have required property 'b'; " +
        "/ must have required property 'c' (+1 more)"
    );
  });

  it('reports values that are not valid JSON', () => {
    const validate = createValueValidator(schema, 'JSON')!;
    expect(validate('not json')).toEqual({
      valid: false,
      detail: 'invalid JSON'
    });
  });
});
