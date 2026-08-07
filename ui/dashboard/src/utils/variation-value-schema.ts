import Ajv2020 from 'ajv/dist/2020';
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

// Mirrors the backend's strconv.ParseFloat, which rejects surrounding
// whitespace and non-finite values.
const isStrictFiniteNumber = (value: string): boolean =>
  value !== '' && value === value.trim() && Number.isFinite(Number(value));

export type SchemaDefinitionError =
  | 'enum-empty'
  | 'enum-not-number'
  | 'regex-empty'
  | 'json-schema-empty'
  | 'json-schema-invalid'
  | 'type-unsupported';

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
      // The backend compiles patterns with Go RE2, whose syntax differs from
      // JS RegExp (e.g. inline flags like (?U) are valid in Go but throw in
      // JS). Never block on JS compilation; the backend is the source of
      // truth for pattern validity.
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

// Returns a value validator, or null when client-side validation is not
// possible (e.g. the pattern only compiles with Go RE2). The backend remains
// the source of truth in that case.
export const createValueValidator = (
  schema: VariationValueSchema,
  variationType: FeatureVariationType
): ((value: string) => boolean) | null => {
  switch (schema.type) {
    case 'ENUM': {
      const values = schema.enumValidator?.values ?? [];
      if (variationType === 'NUMBER') {
        const allowed = values.map(Number);
        return value =>
          isStrictFiniteNumber(value) && allowed.includes(Number(value));
      }
      return value => values.includes(value);
    }
    case 'REGEX': {
      try {
        const pattern = new RegExp(schema.regexValidator?.pattern ?? '');
        return value => pattern.test(value);
      } catch {
        return null;
      }
    }
    case 'JSON_SCHEMA': {
      try {
        const validate = compileJsonSchema(
          schema.jsonSchemaValidator?.schema ?? ''
        );
        return value => {
          try {
            return validate(JSON.parse(value)) === true;
          } catch {
            return false;
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
