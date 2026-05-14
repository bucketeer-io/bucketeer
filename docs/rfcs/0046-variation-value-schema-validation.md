# Summary

This RFC proposes adding schema-based validation for feature flag variation
values. The schema is owned by the feature flag and is used to validate every
variation value before it is created or updated.

The v1 scope supports:

- Enum validation for `STRING` and `NUMBER` variation types
- Regex validation for `STRING` variation types
- JSON Schema validation for `JSON` variation types

Validation will run in the backend before saving direct flag changes and
scheduled flag changes. The UI implementation will be handled separately after
the design direction is confirmed.

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
- Return structured backend errors that the UI can map to variation value fields.
- Keep the implementation incremental so backend work can proceed before the UI
  design is finalized.

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

  Type type = 1;

  // Used when type is ENUM.
  // For NUMBER flags, each value must parse as a number.
  repeated string enum_values = 2;

  // Used when type is REGEX.
  // The pattern must compile using Go regexp syntax.
  string regex_pattern = 3;

  // Used when type is JSON_SCHEMA.
  // The value must be a valid JSON Schema document.
  string json_schema = 4;
}
```

Add `VariationValueSchema variation_value_schema` to:

- `Feature`
- `CreateFeatureRequest`
- `UpdateFeatureRequest`
- Public gateway `CreateFeatureRequest`
- Public gateway `UpdateFeatureRequest`

The update request needs optional semantics so callers can distinguish "do not
change schema" from "clear schema". If plain proto3 message presence is not
sufficient for the generated code and gateway behavior, use an explicit update
field such as:

```proto
VariationValueSchema variation_value_schema = X;
google.protobuf.BoolValue clear_variation_value_schema = Y;
```

The final field shape should be confirmed during RFC review before
implementation.

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

For `STRING` flags, compare the variation value against `enum_values` using exact
string matching.

For `NUMBER` flags:

- Every enum value must parse as a number.
- Every variation value must parse as a number.
- Comparison should use normalized numeric comparison instead of raw string
  comparison, so `1` and `1.0` are treated consistently.

Enum schemas must contain at least one value.

### Regex

Regex validation is supported only for `STRING` flags.

The regex pattern must compile before it can be saved. Variation values are valid
when they match the compiled pattern. The backend should document that Go regexp
syntax is used.

### JSON Schema

JSON Schema validation is supported only for `JSON` flags.

The schema document must be valid JSON and a valid JSON Schema. Variation values
must be valid JSON and must satisfy the schema.

The implementation should use a maintained Go JSON Schema library rather than a
custom validator. The library choice should be made in the backend validation PR.

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

## Error Contract

Schema validation errors should return `InvalidArgument`.

The backend should include structured details so the UI can map errors to
specific form fields. The preferred shape is `google.rpc.ErrorInfo`, matching the
existing backend error design direction.

Example metadata:

```json
{
  "reason": "INVALID",
  "domain": "feature.bucketeer.io",
  "metadata": {
    "messageKey": "VariationValueSchemaValidationError",
    "field": "variations.0.value",
    "variationId": "variation-id",
    "schemaType": "JSON_SCHEMA",
    "path": "/theme",
    "constraint": "enum"
  }
}
```

For JSON Schema validation, the backend may return multiple details, one per
schema violation. If multiple details are too large for a response, the backend
should return the most relevant errors and include a generic summary.

The UI can initially display the backend message as a toast, but the contract
should support field-level rendering in the UI PR.

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

- Add schema configuration UI after design direction is confirmed.
- Add client-side schema validation where practical.
- Map structured backend validation errors to variation value fields.
- Add UI tests for create, update, and scheduled-change flows.

## Open Questions

- Should schema changes be allowed while an experiment is waiting or running, or
  should they follow the same restrictions as variation value updates?
- Should scheduled flag changes be able to change the schema itself, or should
  schema updates be immediate-only in v1?
- What is the maximum allowed size for `json_schema`, `regex_pattern`, and enum
  value lists?
- Which Go JSON Schema library should be adopted?
- Should JSON Schema validation support a specific draft version only?
