import { afterEach, describe, expect, it, vi } from 'vitest';
import { compileJsonSchema } from './json-schema-validator';

const compile = (schema: unknown) => compileJsonSchema(JSON.stringify(schema));

const draft07 = 'http://json-schema.org/draft-07/schema#';

describe('compileJsonSchema: without string-to-code evaluation', () => {
  // Browsers refuse eval() and new Function() under the dashboard's
  // Content-Security-Policy; make them fail the same way here. A function
  // declaration, so that `new Function(...)` reaches it as a constructor.
  function refuse(): never {
    throw new EvalError(
      "Refused to evaluate a string as JavaScript because 'unsafe-eval' is not an allowed source of script"
    );
  }

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('compiles and validates while eval and new Function are refused', () => {
    vi.stubGlobal('eval', refuse);
    vi.stubGlobal('Function', refuse);
    expect(() => new Function('return 1')).toThrow(EvalError);

    const validate = compile({
      type: 'object',
      required: ['name'],
      properties: { name: { type: 'string' } }
    })!;
    expect(validate({ name: 'a' })).toEqual({ valid: true });
    expect(validate({})).toEqual({
      valid: false,
      detail: '/ Instance does not have required property "name".'
    });
  });
});

describe('compileJsonSchema: schema definition', () => {
  it('rejects text that is not JSON or not an object or boolean schema', () => {
    expect(compileJsonSchema('{ not json')).toBeNull();
    for (const schema of [[], 'str', null, 42]) {
      expect(compile(schema)).toBeNull();
    }
  });

  it('rejects keywords with a shape no draft accepts', () => {
    const invalid = [
      { type: 'nonsense' },
      { type: ['string', 'nonsense'] },
      { properties: 'oops' },
      { required: 'name' },
      { pattern: '(' },
      { minLength: -1 },
      { minLength: 1.5 },
      { items: 'x' },
      { allOf: {} },
      { properties: { a: { type: 'nonsense' } } },
      { $defs: { a: { properties: { b: { required: 'x' } } } } }
    ];
    for (const schema of invalid) {
      expect(compile(schema), JSON.stringify(schema)).toBeNull();
    }
  });

  it('accepts boolean schemas, unknown keywords and older draft forms', () => {
    const valid = [
      true,
      false,
      { type: 'object' },
      { 'x-custom': 1, type: 'string' },
      { $schema: draft07, type: 'string' },
      { minimum: 1, exclusiveMinimum: true },
      { items: [{ type: 'string' }, { type: 'number' }] }
    ];
    for (const schema of valid) {
      expect(compile(schema), JSON.stringify(schema)).not.toBeNull();
    }
  });

  it('resolves $ref to local definitions', () => {
    const validate = compile({
      $ref: '#/$defs/str',
      $defs: { str: { type: 'string' } }
    })!;
    expect(validate('a')).toEqual({ valid: true });
    expect(validate(1)).toEqual({
      valid: false,
      detail: '/ Instance type "number" is invalid. Expected "string".'
    });
  });

  it('reports a $ref that cannot be resolved when it is evaluated', () => {
    const validate = compile({ $ref: '#/$defs/missing' })!;
    const result = validate({});
    expect(result.valid).toBe(false);
    expect(result.detail).toContain('Unresolved $ref "#/$defs/missing".');
  });
});

describe('compileJsonSchema: format', () => {
  it('treats format as an annotation by default, like the backend', () => {
    const validate = compile({
      type: 'object',
      properties: { email: { type: 'string', format: 'email' } }
    })!;
    expect(validate({ email: 'not-an-email' })).toEqual({ valid: true });
  });

  it('asserts format for schemas that declare draft-07 or older', () => {
    const validate = compile({
      $schema: draft07,
      type: 'string',
      format: 'email'
    })!;
    expect(validate('a@b.co')).toEqual({ valid: true });
    expect(validate('not-an-email')).toEqual({
      valid: false,
      detail: '/ String does not match format "email".'
    });
  });

  it('keeps a property that is named format', () => {
    const validate = compile({
      type: 'object',
      properties: { format: { type: 'string' } }
    })!;
    expect(validate({ format: 'x' })).toEqual({ valid: true });
    expect(validate({ format: 1 })).toEqual({
      valid: false,
      detail: '/format Instance type "number" is invalid. Expected "string".'
    });
  });
});

describe('compileJsonSchema: detail', () => {
  const validate = compile({
    type: 'object',
    required: ['name'],
    additionalProperties: false,
    properties: {
      name: { type: 'string' },
      theme: { enum: ['light', 'dark'] },
      tags: { type: 'array', items: { type: 'string' } },
      'a/b': { type: 'string' },
      'my key': { type: 'string' }
    }
  })!;

  it('is omitted for a conforming value', () => {
    expect(validate({ name: 'a', theme: 'dark', tags: ['x'] })).toEqual({
      valid: true
    });
  });

  it('reports the specific error, not the container keyword around it', () => {
    expect(validate({ name: 'a', theme: 'blue' })).toEqual({
      valid: false,
      detail: '/theme Instance does not match any of ["light", "dark"].'
    });
    expect(validate({ name: 'a', tags: ['x', 1] })).toEqual({
      valid: false,
      detail: '/tags/1 Instance type "number" is invalid. Expected "string".'
    });
  });

  it('reports the keyword that rejected an additional property', () => {
    expect(validate({ name: 'a', extraField: 'x' })).toEqual({
      valid: false,
      detail:
        '/ Property "extraField" does not match additional properties schema.'
    });
  });

  it('reports a false root schema as such', () => {
    expect(compile(false)!('anything')).toEqual({
      valid: false,
      detail: '/ False boolean schema.'
    });
  });

  it('writes paths as JSON pointers', () => {
    expect(validate({ name: 'a', 'a/b': 1, 'my key': 2 })).toEqual({
      valid: false,
      detail:
        '/a~1b Instance type "number" is invalid. Expected "string".; ' +
        '/my key Instance type "number" is invalid. Expected "string".'
    });
  });

  it('caps the listed enum values', () => {
    const sizes = compile({
      type: 'object',
      properties: { size: { enum: ['a', 'b', 'c', 'd', 'e', 'f', 'g'] } }
    })!;
    expect(sizes({ size: 'z' })).toEqual({
      valid: false,
      detail:
        '/size Instance does not match any of ["a", "b", "c", "d", "e", +2 more].'
    });
  });

  it('caps the number of errors with a remainder count', () => {
    const required = compile({
      type: 'object',
      required: ['a', 'b', 'c', 'd']
    })!;
    expect(required({})).toEqual({
      valid: false,
      detail:
        '/ Instance does not have required property "a".; ' +
        '/ Instance does not have required property "b".; ' +
        '/ Instance does not have required property "c". (+1 more)'
    });
  });
});
