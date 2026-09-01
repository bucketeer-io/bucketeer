# Summary

This RFC proposes adding schema-based validation for feature flag variation
values. The schema is owned by the feature flag and is used to validate every
variation value before it is created or updated.

The v1 scope supports:

- Enum validation for `STRING` and `NUMBER` variation types
- Regex validation for `STRING` variation types
- JSON Schema validation for `JSON` variation types

Validation will run in the backend before saving direct flag changes and
scheduled flag changes. The UI will provide schema configuration and client-side
validation.

[Issue](https://github.com/bucketeer-io/bucketeer/issues/2499)

## Background

Variation values are currently stored as strings because a value can represent a
string, number, JSON object, YAML document, or boolean. The feature already has a
`variation_type` field, and the domain layer validates primitive type rules and
uniqueness before saving changes.

The current validation covers only basic type safety:

- `BOOLEAN`: must be `true` or `false`
- `NUMBER`: must parse as a number
- `JSON`: must parse as JSON
- `YAML`: must parse as YAML
- All types: variation values must be unique within a flag

This does not allow teams to constrain values to a product-defined contract, such
as a fixed enum, a string format, or a JSON object shape.

## Goals

- Store a flag-level variation value schema.
- Validate variation values before creating or updating a feature flag.
- Validate scheduled variation changes before scheduling and again before
  execution.
- Return `InvalidArgument` when a schema definition or variation value is
  invalid.
- Keep backend validation authoritative while providing detailed client-side
  validation in the UI.

## Non-goals

- Do not add UI controls in the backend PRs.
- Do not support YAML schema validation in v1.
- Do not support multiple schema validators on the same flag in v1.
- Do not change SDK evaluation behavior. SDKs continue to receive variation
  values as they do today.
- Do not validate user attributes or targeting clause values with this schema.

## Schema Ownership

The schema should be a flag-level field, not a per-variation field.

All variations of the same feature flag share the same `variation_type`, so a
single schema should describe the allowed shape of all values for that flag.
Per-variation schemas would allow a flag to contain incompatible contracts and
would make client-side usage harder to reason about.

In v1, schemas are owned by feature flags because Bucketeer does not currently
have a reusable Variable or Config entity separate from Feature. If reusable
schemas become necessary later, this design can be extended by introducing a
schema registry and storing a schema reference on the feature.

## V1 Type Matrix

| Variation type | Enum | Regex | JSON Schema |
| :-- | :-- | :-- | :-- |
| `STRING` | Supported | Supported | Not supported |
| `NUMBER` | Supported | Not supported | Not supported |
| `JSON` | Not supported | Not supported | Supported |
| `BOOLEAN` | Not supported | Not supported | Not supported |
| `YAML` | Not supported in v1 | Not supported in v1 | Not supported in v1 |

If a schema type is not compatible with the feature's `variation_type`, the
backend rejects the request with `InvalidArgument`.

## Protobuf Changes

Add a new schema message under `proto/feature`.

```proto
message VariationValueSchema {
  enum Type {
    TYPE_UNSPECIFIED = 0;
    ENUM = 1;
    REGEX = 2;
    JSON_SCHEMA = 3;
  }

  message EnumValidator {
    // For NUMBER flags, each value must parse as a number.
    repeated string values = 1;
  }

  message RegexValidator {
    // The pattern must compile using Go regexp syntax.
    string pattern = 1;
  }

  message JsonSchemaValidator {
    // The value must be a valid JSON Schema document.
    string schema = 1;
  }

  Type type = 1;

  oneof validator {
    EnumValidator enum_validator = 2;
    RegexValidator regex_validator = 3;
    JsonSchemaValidator json_schema_validator = 4;
  }

  // Optional human-readable explanation of the schema purpose.
  string description = 5;
}
```

The `oneof` makes the validator payload structurally exclusive, so callers
cannot send multiple validator definitions for a single schema.

`type` and the selected `validator` must describe the same validator. For
example, `type = ENUM` requires `enum_validator`. An unspecified type, a missing
validator, or a mismatch between the two fields is invalid and the backend
rejects it with `InvalidArgument`.

Add `VariationValueSchema variation_value_schema` to:

- `Feature`
- `CreateFeatureRequest`
- `UpdateFeatureRequest`
- Public gateway `CreateFeatureRequest`
- Public gateway `UpdateFeatureRequest`

The update request uses separate fields so callers can distinguish "do not
change schema" from "clear schema":

```proto
VariationValueSchema variation_value_schema = 20;
google.protobuf.BoolValue clear_variation_value_schema = 21;
```

The update semantics are:

- If both fields are absent, the existing schema is unchanged.
- If `variation_value_schema` is present and clear is absent or `false`, the
  schema is set or replaced.
- If `clear_variation_value_schema` is `true`, the schema is cleared. Clear takes
  precedence if both fields are present.
- If clear is `false` and no schema is present, the existing schema is unchanged.

## Persistence

Add a nullable JSON column to the feature table.

For MySQL:

```sql
ALTER TABLE feature
  ADD COLUMN variation_value_schema JSON NULL;
```

For PostgreSQL:

```sql
ALTER TABLE feature
  ADD COLUMN variation_value_schema JSONB NULL;
```

Storage changes are required because feature storage writes individual columns
instead of serializing the entire `Feature` proto.

The generated protobuf `oneof` cannot be restored through the generic
`encoding/json` database adapters. Variation value schemas therefore use a
dedicated scanner and valuer backed by `protojson.Marshal` and
`protojson.Unmarshal`. Unknown fields are discarded while reading so newer
stored schemas remain readable by older binaries.

Update the feature create, update, get, list, and list-by-environment SQL paths
for both MySQL and PostgreSQL.

The column should be omitted or `NULL` for flags without a schema. Existing flags
therefore keep their current behavior.

## Validation Behavior

Validation should run after existing primitive type validation and before
persistence.

The domain layer already centralizes primitive validation in
`Feature.validateVariationValue`. Schema validation should be implemented there
or in a helper called by that method so the following flows share the same rules:

- Create feature
- Update feature variation changes
- Scheduled flag change execution

### Create Feature

When creating a feature with `variation_value_schema`:

1. Validate the schema is compatible with `variation_type`.
2. Validate the schema definition itself.
3. Validate every requested variation value against the schema.
4. Reject the request before saving if any variation is invalid.

### Update Feature

When updating variation values:

1. Load the existing feature and schema.
2. If the request changes the schema, validate the new schema.
3. Validate all existing post-update variation values against the active schema.
4. Reject the request before saving if any variation is invalid.

When adding or changing a schema on an existing flag, all existing variations
must pass the new schema. This prevents a flag from storing a schema that its
current values violate.

### Clearing A Schema

Clearing a schema removes only schema-based validation. Existing primitive type
validation and uniqueness validation continue to apply.

## Validator Semantics

### Enum

Enum validation allows a finite set of values.

For `STRING` flags, compare the variation value against
`enum_validator.values` using exact string matching.

For `NUMBER` flags:

- Every enum and variation value is parsed with `strconv.ParseFloat(value, 64)`.
- Decimal and hexadecimal floating-point forms accepted by `ParseFloat` are
  valid. Surrounding whitespace and other unsupported syntax are invalid.
- `NaN`, positive infinity, and negative infinity are invalid.
- Parsed IEEE 754 `float64` values are compared with numeric equality. Therefore
  values such as `1`, `1.0`, and `1e0` compare equal. Values that round to the
  same `float64` also compare equal; arbitrary-precision comparison is not part
  of v1.

The UI accepts decimal notation as a conservative subset of this backend
grammar.

Enum schemas must contain at least one value.

### Regex

Regex validation is supported only for `STRING` flags.

The regex pattern must compile before it can be saved. Patterns use Go's
`regexp` syntax and `MatchString` semantics, so an unanchored pattern may match a
substring. Callers must use `^` and `$` when the whole variation value must
match.

### JSON Schema

JSON Schema validation is supported only for `JSON` flags.

The schema document must be valid JSON and a valid JSON Schema. Variation values
must be valid JSON and must satisfy the schema.

The backend uses `github.com/santhosh-tekuri/jsonschema/v6` with JSON Schema
Draft 2020-12. The v1 contract permits only self-contained schemas. References
to definitions within the submitted schema are supported, but implementations
must not resolve external file, HTTP, or HTTPS `$ref` resources.

## Scheduled Changes

Scheduled variation changes need validation at two points.

### Schedule Create And Update

When creating or updating a scheduled flag change, the backend should validate
the scheduled payload against the schema that is active at scheduling time.

This gives users immediate feedback and prevents obviously invalid scheduled
changes from being stored.

### Schedule Execution

The backend must validate again at execution time.

The flag schema may change between schedule creation and execution. Execution
must use the schema active at execution time. If the scheduled variation value no
longer satisfies the active schema, execution should fail and the scheduled flag
change should move to the existing failed state with a validation failure reason.

During scheduled-change execution, a validation or other permanent execution
error rolls back the feature transaction. After that transaction returns, the
backend records the scheduled change's failed state in a separate write so that
the failure is not rolled back with the feature update. Transient errors leave
the scheduled change pending for retry.

## Error Contract

Schema definition and value validation errors return `InvalidArgument`. They use
the existing backend error conversion and contain one `google.rpc.ErrorInfo`
detail with the existing generic invalid-argument message key.

V1 does not return one backend detail per JSON Schema violation or guarantee a
variation field path in the error metadata. The UI performs client-side
validation before submission and renders its own validation details. The generic
backend error remains the fallback for API clients and validation races.

## API Compatibility

This change is additive:

- Existing flags have no schema and keep current behavior.
- Existing create/update clients can omit `variation_value_schema`.
- SDK evaluation responses are unchanged.

The only behavior change is that clients providing a schema, or editing a flag
that already has a schema, may receive validation errors before persistence.

## Implementation Plan

### PR 1: RFC

- Add this RFC.
- Confirm schema shape, update semantics, scheduled-change behavior, and error
  contract.

### PR 2: Backend Model And Persistence

- Add `VariationValueSchema` proto.
- Add schema fields to feature and feature create/update requests.
- Add gateway request fields.
- Add MySQL and PostgreSQL migrations.
- Update MySQL and PostgreSQL feature storage SQL.
- Regenerate Go protobuf and OpenAPI output.
- Add storage and model round-trip tests.

### PR 3: Backend Validation

- Add schema compatibility and schema-definition validation.
- Add enum, regex, and JSON Schema value validation.
- Integrate validation with feature create and update flows.
- Validate scheduled flag changes at create/update time.
- Revalidate scheduled flag changes at execution time.
- Add domain, API, and scheduled-change tests.

### PR 4: UI

- Add schema configuration UI.
- Add client-side schema validation where practical.
- Render detailed client-side validation errors and use the generic backend
  error as a fallback.
- Add UI tests for create, update, and scheduled-change flows.

## Open Questions

- Should schema changes be allowed while an experiment is waiting or running, or
  should they follow the same restrictions as variation value updates?
- Should scheduled flag changes be able to change the schema itself, or should
  schema updates be immediate-only in v1?
- What is the maximum allowed size for `json_schema_validator.schema`,
  `regex_validator.pattern`, and enum value lists?
