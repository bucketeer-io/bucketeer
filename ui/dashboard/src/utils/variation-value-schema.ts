import Ajv2020 from 'ajv/dist/2020';
import type { ErrorObject } from 'ajv/dist/2020';
import { RE2JS } from 're2js';
import {
  FeatureVariationType,
  VariationValueSchema,
  VariationValueSchemaType
} from '@types';

// Matches the backend v1 type matrix in pkg/feature/domain/variation_value_schema.go
export const getSupportedSchemaTypes = (
  variationType: FeatureVariationType
): VariationValueSchemaType[] => {
  switch (variationType) {
    case 'STRING':
      return ['ENUM', 'REGEX'];
    case 'NUMBER':
      return ['ENUM'];
    case 'JSON':
      return ['JSON_SCHEMA'];
    default:
      return [];
  }
};

export const isSchemaSupported = (
  variationType: FeatureVariationType
): boolean => getSupportedSchemaTypes(variationType).length > 0;

// Reuse a single Ajv instance; constructing one is relatively expensive and
// this runs on every keystroke while editing a JSON Schema.
const ajv2020 = new Ajv2020({ strict: false, allErrors: true });

const compileJsonSchema = (schema: string) =>
  ajv2020.compile(JSON.parse(schema));

// The backend parses enum numbers with strconv.ParseFloat, but Number() is
// more lenient (it also accepts 0x/0b/0o prefixes, 'Infinity', and
// surrounding whitespace). Only accept plain decimal notation so the dialog
// never accepts a value the backend rejects.
const DECIMAL_NUMBER_REGEX = /^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$/;

const isStrictFiniteNumber = (value: string): boolean =>
  DECIMAL_NUMBER_REGEX.test(value) && Number.isFinite(Number(value));

export type SchemaDefinitionError =
  | 'enum-empty'
  | 'enum-not-number'
  | 'regex-empty'
  | 'regex-invalid'
  | 'json-schema-empty'
  | 'json-schema-invalid'
  | 'type-unsupported';

// Compiles the pattern with RE2JS, which implements the same RE2 syntax the
// backend uses via Go's regexp package. This keeps the client and backend in
// agreement: Go-only constructs like inline flags compile, while Perl-only
// constructs like lookarounds are rejected up front.
const compileRe2Pattern = (pattern: string): RE2JS | null => {
  try {
    return RE2JS.compile(pattern);
  } catch {
    return null;
  }
};

// Returns null when the schema definition itself is valid.
export const validateSchemaDefinition = (
  schema: VariationValueSchema,
  variationType: FeatureVariationType
): SchemaDefinitionError | null => {
  if (!getSupportedSchemaTypes(variationType).includes(schema.type)) {
    return 'type-unsupported';
  }
  switch (schema.type) {
    case 'ENUM': {
      const values = schema.enumValidator?.values ?? [];
      if (values.length === 0) return 'enum-empty';
      if (
        variationType === 'NUMBER' &&
        values.some(value => !isStrictFiniteNumber(value))
      ) {
        return 'enum-not-number';
      }
      return null;
    }
    case 'REGEX': {
      const pattern = schema.regexValidator?.pattern ?? '';
      if (pattern === '') return 'regex-empty';
      if (!compileRe2Pattern(pattern)) return 'regex-invalid';
      return null;
    }
    case 'JSON_SCHEMA': {
      const jsonSchema = schema.jsonSchemaValidator?.schema ?? '';
      if (jsonSchema.trim() === '') return 'json-schema-empty';
      try {
        compileJsonSchema(jsonSchema);
        return null;
      } catch {
        return 'json-schema-invalid';
      }
    }
    default:
      return 'type-unsupported';
  }
};

export interface ValueValidationResult {
  valid: boolean;
  // Technical detail from the underlying validator (e.g. the failing JSON
  // path reported by AJV). Not localized; shown as-is next to the
  // localized message.
  detail?: string;
}

// Caps keep the detail usable in a single-line inline form error; a badly
// broken document can produce dozens of AJV errors.
const MAX_DETAIL_ERRORS = 3;
const MAX_ALLOWED_VALUES = 5;

// AJV's default messages omit specifics it collects in error.params (which
// additional property is present, which values an enum allows), so append
// them to make the error actionable without re-reading the schema.
const formatAjvError = (error: ErrorObject): string => {
  let message = error.message ?? '';
  if (
    error.keyword === 'additionalProperties' &&
    typeof error.params.additionalProperty === 'string'
  ) {
    message += ` ('${error.params.additionalProperty}')`;
  } else if (
    error.keyword === 'enum' &&
    Array.isArray(error.params.allowedValues)
  ) {
    const allowed = error.params.allowedValues as unknown[];
    const shown = allowed
      .slice(0, MAX_ALLOWED_VALUES)
      .map(value => JSON.stringify(value))
      .join(', ');
    const more =
      allowed.length > MAX_ALLOWED_VALUES
        ? `, +${allowed.length - MAX_ALLOWED_VALUES} more`
        : '';
    message += `: ${shown}${more}`;
  }
  return `${error.instancePath || '/'} ${message}`.trim();
};

const formatAjvErrors = (
  errors: ErrorObject[] | null | undefined
): string | undefined => {
  if (!errors || errors.length === 0) return undefined;
  const shown = errors.slice(0, MAX_DETAIL_ERRORS).map(formatAjvError);
  const more =
    errors.length > MAX_DETAIL_ERRORS
      ? ` (+${errors.length - MAX_DETAIL_ERRORS} more)`
      : '';
  return shown.join('; ') + more;
};

// Returns a value validator, or null when client-side validation is not
// possible. The backend remains the source of truth either way.
export const createValueValidator = (
  schema: VariationValueSchema,
  variationType: FeatureVariationType
): ((value: string) => ValueValidationResult) | null => {
  switch (schema.type) {
    case 'ENUM': {
      const values = schema.enumValidator?.values ?? [];
      if (variationType === 'NUMBER') {
        const allowed = values.map(Number);
        return value => ({
          valid: isStrictFiniteNumber(value) && allowed.includes(Number(value))
        });
      }
      return value => ({ valid: values.includes(value) });
    }
    case 'REGEX': {
      const pattern = compileRe2Pattern(schema.regexValidator?.pattern ?? '');
      if (!pattern) return null;
      // find() is an unanchored partial match, mirroring the backend's
      // regexp MatchString.
      return value => ({ valid: pattern.matcher(value).find() });
    }
    case 'JSON_SCHEMA': {
      try {
        const validate = compileJsonSchema(
          schema.jsonSchemaValidator?.schema ?? ''
        );
        return value => {
          try {
            if (validate(JSON.parse(value)) === true) return { valid: true };
            return { valid: false, detail: formatAjvErrors(validate.errors) };
          } catch {
            return { valid: false, detail: 'invalid JSON' };
          }
        };
      } catch {
        return null;
      }
    }
    default:
      return null;
  }
};
