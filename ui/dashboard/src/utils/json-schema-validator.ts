import {
  Validator,
  schemaArrayKeyword,
  schemaKeyword,
  schemaMapKeyword
} from '@cfworker/json-schema';
import type { OutputUnit, Schema, SchemaDraft } from '@cfworker/json-schema';

// Validates JSON values against user-provided JSON Schemas by interpreting
// the schema. A code-generating validator builds the validation function
// with `new Function`, which browsers refuse under the dashboard's
// Content-Security-Policy (script-src without 'unsafe-eval').

export interface JsonSchemaValidationResult {
  valid: boolean;
  // Technical detail from the validator: the failing JSON pointer and rule.
  // Not localized; shown as-is next to the localized message.
  detail?: string;
}

export type JsonSchemaValidator = (
  instance: unknown
) => JsonSchemaValidationResult;

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

const isSchema = (value: unknown): value is Schema | boolean =>
  typeof value === 'boolean' || isRecord(value);

const simpleTypes = [
  'array',
  'boolean',
  'integer',
  'null',
  'number',
  'object',
  'string'
];
const subschema: Schema = { $ref: '#' };
const subschemaMap: Schema = {
  type: 'object',
  additionalProperties: subschema
};
const subschemaList: Schema = { type: 'array', items: subschema };
const nonNegativeInteger: Schema = { type: 'integer', minimum: 0 };

// A schema of schemas: the subset of the official meta-schema whose keyword
// shapes are the same in every draft the backend accepts, so a schema that
// fails here is rejected by the backend as well. Unknown keywords are
// allowed, as in the official meta-schema.
const schemaShape: Schema = {
  type: ['object', 'boolean'],
  properties: {
    $ref: { type: 'string' },
    $defs: subschemaMap,
    definitions: subschemaMap,
    title: { type: 'string' },
    description: { type: 'string' },
    type: {
      anyOf: [
        { enum: simpleTypes },
        { type: 'array', items: { enum: simpleTypes } }
      ]
    },
    enum: { type: 'array' },
    properties: subschemaMap,
    patternProperties: subschemaMap,
    additionalProperties: subschema,
    propertyNames: subschema,
    required: { type: 'array', items: { type: 'string' } },
    minProperties: nonNegativeInteger,
    maxProperties: nonNegativeInteger,
    // Older drafts also accept a list of schemas.
    items: { anyOf: [subschema, subschemaList] },
    prefixItems: subschemaList,
    contains: subschema,
    minItems: nonNegativeInteger,
    maxItems: nonNegativeInteger,
    uniqueItems: { type: 'boolean' },
    minimum: { type: 'number' },
    maximum: { type: 'number' },
    // Draft 4 uses booleans that qualify minimum/maximum.
    exclusiveMinimum: { type: ['number', 'boolean'] },
    exclusiveMaximum: { type: ['number', 'boolean'] },
    multipleOf: { type: 'number' },
    minLength: nonNegativeInteger,
    maxLength: nonNegativeInteger,
    pattern: { type: 'string', format: 'regex' },
    format: { type: 'string' },
    allOf: subschemaList,
    anyOf: subschemaList,
    oneOf: subschemaList,
    not: subschema,
    if: subschema,
    then: subschema,
    else: subschema
  }
};
const schemaShapeValidator = new Validator(schemaShape, '2020-12');

// The backend defaults to draft 2020-12 and follows $schema when present.
const draftOf = (schema: Schema | boolean): SchemaDraft => {
  const id = typeof schema === 'boolean' ? '' : (schema.$schema ?? '');
  if (id.includes('draft-04')) return '4';
  if (id.includes('draft-06') || id.includes('draft-07')) return '7';
  if (id.includes('2019-09')) return '2019-09';
  return '2020-12';
};

// From draft 2019-09 on, "format" is an annotation unless the validator opts
// in to asserting it. The backend does not opt in, so a value it accepts
// must not be rejected here. Returns a copy of the schema without the
// keyword; property names, enum values and other data are left untouched.
const withoutFormat = (schema: Schema | boolean): Schema | boolean => {
  if (typeof schema === 'boolean') return schema;
  const copy: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(schema)) {
    if (key === 'format') continue;
    if (key in schemaMapKeyword && isRecord(value)) {
      copy[key] = Object.fromEntries(
        Object.entries(value).map(([name, sub]) => [
          name,
          isSchema(sub) ? withoutFormat(sub) : sub
        ])
      );
    } else if (key in schemaArrayKeyword && Array.isArray(value)) {
      copy[key] = value.map(sub => (isSchema(sub) ? withoutFormat(sub) : sub));
    } else if (key in schemaKeyword && isSchema(value)) {
      copy[key] = withoutFormat(value);
    } else {
      copy[key] = value;
    }
  }
  return copy as Schema;
};

// Caps keep the detail usable in a single-line inline form error; a badly
// broken document can produce dozens of errors.
const MAX_DETAIL_ERRORS = 3;
const MAX_ALLOWED_VALUES = 5;

// Locations are JSON pointers in URI fragment form ("#/properties/a").
const decodePointer = (location: string): string =>
  decodeURI(location.slice(1));

// The validator reports a container keyword (properties, allOf, $ref, ...)
// followed by the specific errors found inside it. Only the specific ones
// say what to fix, so a unit whose next unit is nested in it is dropped.
// A `false` subschema is reported at the instance location instead, so the
// keyword that applied it (e.g. additionalProperties) is kept and the bare
// "False boolean schema." unit is dropped.
const specificErrors = (errors: OutputUnit[]): OutputUnit[] =>
  errors.filter((unit, i) => {
    if (unit.keyword === 'false') return errors.length === 1;
    const next = errors[i + 1];
    return (
      !next || !next.keywordLocation.startsWith(`${unit.keywordLocation}/`)
    );
  });

// Resolves the schema node a keyword location points at. Locations that go
// through a $ref cannot be followed in the root schema and yield undefined.
const schemaNodeAt = (
  schema: Schema | boolean,
  keywordLocation: string
): unknown =>
  decodePointer(keywordLocation)
    .split('/')
    .slice(1)
    .reduce<unknown>(
      (node, token) =>
        isRecord(node)
          ? node[token.replace(/~1/g, '/').replace(/~0/g, '~')]
          : undefined,
      schema
    );

// The validator's enum message lists every allowed value; cap the list like
// the error count so a long enum does not flood the detail.
const formatError = (unit: OutputUnit, schema: Schema | boolean): string => {
  let message = unit.error;
  if (unit.keyword === 'enum') {
    const allowed = schemaNodeAt(schema, unit.keywordLocation);
    if (Array.isArray(allowed)) {
      const shown = allowed
        .slice(0, MAX_ALLOWED_VALUES)
        .map(value => JSON.stringify(value))
        .join(', ');
      const more =
        allowed.length > MAX_ALLOWED_VALUES
          ? `, +${allowed.length - MAX_ALLOWED_VALUES} more`
          : '';
      message = `Instance does not match any of [${shown}${more}].`;
    }
  }
  return `${decodePointer(unit.instanceLocation) || '/'} ${message}`;
};

const formatErrors = (
  errors: OutputUnit[],
  schema: Schema | boolean
): string => {
  const specific = specificErrors(errors);
  const shown = specific
    .slice(0, MAX_DETAIL_ERRORS)
    .map(unit => formatError(unit, schema));
  const more =
    specific.length > MAX_DETAIL_ERRORS
      ? ` (+${specific.length - MAX_DETAIL_ERRORS} more)`
      : '';
  return shown.join('; ') + more;
};

const firstLine = (error: unknown): string =>
  error instanceof Error ? error.message.split('\n')[0] : String(error);

// Returns a validator for the serialized schema, or null when the schema is
// not valid JSON or not a well-formed JSON Schema. The backend remains the
// source of truth either way.
export const compileJsonSchema = (
  schema: string
): JsonSchemaValidator | null => {
  let parsed: unknown;
  try {
    parsed = JSON.parse(schema);
  } catch {
    return null;
  }
  if (!isSchema(parsed) || !schemaShapeValidator.validate(parsed).valid) {
    return null;
  }
  const draft = draftOf(parsed);
  const asserted =
    draft === '4' || draft === '7' ? parsed : withoutFormat(parsed);
  let validator: Validator;
  try {
    validator = new Validator(asserted, draft, false);
  } catch {
    return null;
  }
  return instance => {
    try {
      const result = validator.validate(instance);
      if (result.valid) return { valid: true };
      return { valid: false, detail: formatErrors(result.errors, parsed) };
    } catch (error) {
      // Schema problems only found while evaluating a branch, such as an
      // unresolved $ref or a pattern that is not a valid regular expression.
      return { valid: false, detail: firstLine(error) };
    }
  };
};
