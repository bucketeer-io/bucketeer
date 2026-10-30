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

export interface CompiledJsonSchema {
  // Null when the schema is well-formed but relies on keywords the
  // interpreter does not evaluate; the backend still validates values
  // against it on save.
  validate: JsonSchemaValidator | null;
}

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
    $schema: { type: 'string' },
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

interface Dialect {
  // Part of the $schema URI that identifies the draft.
  id: string;
  draft: SchemaDraft;
  // Evaluated keywords this draft introduced; the backend ignores them in
  // earlier drafts.
  introduced: readonly string[];
  // Keywords the backend stops evaluating from this draft on.
  dropped: readonly string[];
}

// The drafts the backend accepts, oldest first. The interpreter evaluates
// every keyword it knows in every draft, so the keywords a schema's draft
// does not evaluate are removed before the schema is interpreted.
const dialects: readonly Dialect[] = [
  { id: 'draft-04', draft: '4', introduced: [], dropped: [] },
  {
    id: 'draft-06',
    // The interpreter has no draft-06 mode; draft-07 is the closest.
    draft: '7',
    introduced: ['const', 'contains', 'propertyNames'],
    dropped: []
  },
  {
    id: 'draft-07',
    draft: '7',
    introduced: ['if', 'then', 'else'],
    dropped: []
  },
  {
    id: '2019-09',
    draft: '2019-09',
    introduced: [
      '$anchor',
      '$recursiveRef',
      '$recursiveAnchor',
      'dependentRequired',
      'dependentSchemas',
      'unevaluatedItems',
      'unevaluatedProperties',
      'minContains',
      'maxContains'
    ],
    // From draft 2019-09 on, "format" is an annotation unless the validator
    // opts in to asserting it. The backend does not opt in, so a value it
    // accepts must not be rejected here.
    dropped: ['format']
  },
  {
    id: '2020-12',
    draft: '2020-12',
    introduced: ['prefixItems', '$dynamicRef', '$dynamicAnchor'],
    dropped: ['additionalItems']
  }
];

// Keywords the interpreter knows but does not evaluate: a schema that relies
// on them cannot be tested client-side.
const unsupportedKeywords = ['$dynamicRef', '$dynamicAnchor'];

interface DraftHandling {
  draft: SchemaDraft;
  ignoredKeywords: readonly string[];
}

// The backend defaults to draft 2020-12 and follows $schema when present.
const draftHandlingOf = (schema: Schema | boolean): DraftHandling => {
  const id =
    typeof schema !== 'boolean' && typeof schema.$schema === 'string'
      ? schema.$schema
      : '';
  const found = dialects.findIndex(dialect => id.includes(dialect.id));
  const index = found === -1 ? dialects.length - 1 : found;
  return {
    draft: dialects[index].draft,
    ignoredKeywords: [
      ...dialects.slice(index + 1).flatMap(dialect => dialect.introduced),
      ...dialects.slice(0, index + 1).flatMap(dialect => dialect.dropped)
    ]
  };
};

// "dependencies" maps property names to either a list of property names or
// a subschema; the interpreter's keyword tables leave it out.
const isSchemaMapKeyword = (key: string): boolean =>
  key in schemaMapKeyword || key === 'dependencies';

// Returns a copy of the schema without the given keywords in any subschema;
// property names, enum values and other data are left untouched.
const withoutKeywords = (
  schema: Schema | boolean,
  keywords: readonly string[]
): Schema | boolean => {
  if (typeof schema === 'boolean') return schema;
  const strip = (sub: unknown): unknown =>
    isSchema(sub) ? withoutKeywords(sub, keywords) : sub;
  const copy: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(schema)) {
    if (keywords.includes(key)) continue;
    if (isSchemaMapKeyword(key) && isRecord(value)) {
      copy[key] = Object.fromEntries(
        Object.entries(value).map(([name, sub]) => [name, strip(sub)])
      );
    } else if (key in schemaArrayKeyword && Array.isArray(value)) {
      copy[key] = value.map(strip);
    } else if (key in schemaKeyword) {
      copy[key] = strip(value);
    } else {
      copy[key] = value;
    }
  }
  return copy as Schema;
};

// Whether the schema or any of its subschemas uses one of the keywords.
const usesKeyword = (
  schema: Schema | boolean,
  keywords: readonly string[]
): boolean => {
  if (typeof schema === 'boolean') return false;
  return Object.entries(schema).some(([key, value]) => {
    if (keywords.includes(key)) return true;
    let subschemas: unknown[] = [];
    if (isSchemaMapKeyword(key) && isRecord(value)) {
      subschemas = Object.values(value);
    } else if (key in schemaArrayKeyword && Array.isArray(value)) {
      subschemas = value;
    } else if (key in schemaKeyword) {
      subschemas = [value];
    }
    return subschemas.some(sub => isSchema(sub) && usesKeyword(sub, keywords));
  });
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

// Resolves the schema node a keyword location points at, through objects
// (properties, $defs, ...) and arrays (allOf, anyOf, prefixItems, ...).
// Locations that go through a $ref cannot be followed in the root schema
// and yield undefined.
const schemaNodeAt = (
  schema: Schema | boolean,
  keywordLocation: string
): unknown =>
  decodePointer(keywordLocation)
    .split('/')
    .slice(1)
    .reduce<unknown>((node, token) => {
      const key = token.replace(/~1/g, '/').replace(/~0/g, '~');
      if (Array.isArray(node)) return node[Number(key)];
      return isRecord(node) ? node[key] : undefined;
    }, schema);

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

// Compiles the serialized schema, or returns null when it is not valid JSON
// or not a well-formed JSON Schema. The backend remains the source of truth
// either way.
export const compileJsonSchema = (
  schema: string
): CompiledJsonSchema | null => {
  let parsed: unknown;
  try {
    parsed = JSON.parse(schema);
  } catch {
    return null;
  }
  if (!isSchema(parsed)) return null;
  const { draft, ignoredKeywords } = draftHandlingOf(parsed);
  // Keywords the schema's draft does not evaluate are removed first, so the
  // shape check only constrains keywords that count for that draft.
  const evaluated = withoutKeywords(parsed, ignoredKeywords);
  if (!schemaShapeValidator.validate(evaluated).valid) return null;
  if (usesKeyword(evaluated, unsupportedKeywords)) return { validate: null };
  let validator: Validator;
  try {
    validator = new Validator(evaluated, draft, false);
  } catch {
    return null;
  }
  return {
    validate: instance => {
      try {
        const result = validator.validate(instance);
        if (result.valid) return { valid: true };
        return {
          valid: false,
          detail: formatErrors(result.errors, evaluated)
        };
      } catch (error) {
        // Schema problems only found while evaluating a branch, such as an
        // unresolved $ref or a pattern that is not a valid regular
        // expression.
        return { valid: false, detail: firstLine(error) };
      }
    }
  };
};
