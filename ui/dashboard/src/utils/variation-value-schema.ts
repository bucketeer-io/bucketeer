import { RE2JS } from 're2js';
import {
  FeatureVariationType,
  VariationValueSchema,
  VariationValueSchemaType
} from '@types';
import { compileJsonSchema } from './json-schema-validator';

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
      if (!compileJsonSchema(jsonSchema)) return 'json-schema-invalid';
      return null;
    }
    default:
      return 'type-unsupported';
  }
};

export interface ValueValidationResult {
  valid: boolean;
  // Technical detail from the underlying validator (e.g. the failing JSON
  // path). Not localized; shown as-is next to the localized message.
  detail?: string;
}

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
      const validate = compileJsonSchema(
        schema.jsonSchemaValidator?.schema ?? ''
      )?.validate;
      if (!validate) return null;
      return value => {
        let instance: unknown;
        try {
          instance = JSON.parse(value);
        } catch {
          return { valid: false, detail: 'invalid JSON' };
        }
        return validate(instance);
      };
    }
    default:
      return null;
  }
};
