// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func enumSchema(values ...string) *featureproto.VariationValueSchema {
	return &featureproto.VariationValueSchema{
		Type: featureproto.VariationValueSchema_ENUM,
		Validator: &featureproto.VariationValueSchema_EnumValidator_{
			EnumValidator: &featureproto.VariationValueSchema_EnumValidator{
				Values: values,
			},
		},
	}
}

func regexSchema(pattern string) *featureproto.VariationValueSchema {
	return &featureproto.VariationValueSchema{
		Type: featureproto.VariationValueSchema_REGEX,
		Validator: &featureproto.VariationValueSchema_RegexValidator_{
			RegexValidator: &featureproto.VariationValueSchema_RegexValidator{
				Pattern: pattern,
			},
		},
	}
}

func jsonSchema(schema string) *featureproto.VariationValueSchema {
	return &featureproto.VariationValueSchema{
		Type: featureproto.VariationValueSchema_JSON_SCHEMA,
		Validator: &featureproto.VariationValueSchema_JsonSchemaValidator_{
			JsonSchemaValidator: &featureproto.VariationValueSchema_JsonSchemaValidator{
				Schema: schema,
			},
		},
	}
}

func TestValidateVariationValueSchemaDefinition(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		variationType featureproto.Feature_VariationType
		schema        *featureproto.VariationValueSchema
		expected      error
	}{
		{
			desc:          "string enum is valid",
			variationType: featureproto.Feature_STRING,
			schema:        enumSchema("small", "large"),
		},
		{
			desc:          "number enum is valid",
			variationType: featureproto.Feature_NUMBER,
			schema:        enumSchema("1", "2.5"),
		},
		{
			desc:          "number enum rejects non-number option",
			variationType: featureproto.Feature_NUMBER,
			schema:        enumSchema("one"),
			expected:      errVariationValueSchemaInvalid,
		},
		{
			desc:          "boolean enum is not supported",
			variationType: featureproto.Feature_BOOLEAN,
			schema:        enumSchema("true", "false"),
			expected:      errVariationValueSchemaTypeUnmatched,
		},
		{
			desc:          "string regex is valid",
			variationType: featureproto.Feature_STRING,
			schema:        regexSchema(`^[a-z]+$`),
		},
		{
			desc:          "number regex is not supported",
			variationType: featureproto.Feature_NUMBER,
			schema:        regexSchema(`^[0-9]+$`),
			expected:      errVariationValueSchemaTypeUnmatched,
		},
		{
			desc:          "invalid regex is rejected",
			variationType: featureproto.Feature_STRING,
			schema:        regexSchema(`[`),
			expected:      errVariationValueSchemaInvalid,
		},
		{
			desc:          "json schema is valid",
			variationType: featureproto.Feature_JSON,
			schema: jsonSchema(`{
				"type": "object",
				"properties": {"theme": {"type": "string"}},
				"required": ["theme"]
			}`),
		},
		{
			desc:          "string json schema is not supported",
			variationType: featureproto.Feature_STRING,
			schema:        jsonSchema(`{"type":"object"}`),
			expected:      errVariationValueSchemaTypeUnmatched,
		},
		{
			desc:          "invalid json schema is rejected",
			variationType: featureproto.Feature_JSON,
			schema:        jsonSchema(`{"type":`),
			expected:      errVariationValueSchemaInvalid,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			err := validateVariationValueSchemaDefinition(p.variationType, p.schema)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestValidateVariationValueAgainstSchema(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc          string
		variationType featureproto.Feature_VariationType
		schema        *featureproto.VariationValueSchema
		value         string
		expected      error
	}{
		{
			desc:          "string enum accepts allowed value",
			variationType: featureproto.Feature_STRING,
			schema:        enumSchema("small", "large"),
			value:         "small",
		},
		{
			desc:          "string enum rejects disallowed value",
			variationType: featureproto.Feature_STRING,
			schema:        enumSchema("small", "large"),
			value:         "medium",
			expected:      errVariationValueSchemaViolation,
		},
		{
			desc:          "number enum compares normalized numbers",
			variationType: featureproto.Feature_NUMBER,
			schema:        enumSchema("1", "2"),
			value:         "1.0",
		},
		{
			desc:          "regex accepts matching string",
			variationType: featureproto.Feature_STRING,
			schema:        regexSchema(`^[a-z]+-[0-9]+$`),
			value:         "item-1",
		},
		{
			desc:          "regex rejects non-matching string",
			variationType: featureproto.Feature_STRING,
			schema:        regexSchema(`^[a-z]+-[0-9]+$`),
			value:         "item",
			expected:      errVariationValueSchemaViolation,
		},
		{
			desc:          "json schema accepts valid value",
			variationType: featureproto.Feature_JSON,
			schema: jsonSchema(`{
				"type": "object",
				"properties": {
					"theme": {"type": "string", "enum": ["light", "dark"]},
					"maxItems": {"type": "integer", "minimum": 1}
				},
				"required": ["theme", "maxItems"],
				"additionalProperties": false
			}`),
			value: `{"theme":"dark","maxItems":20}`,
		},
		{
			desc:          "json schema rejects invalid value",
			variationType: featureproto.Feature_JSON,
			schema: jsonSchema(`{
				"type": "object",
				"properties": {"theme": {"type": "string", "enum": ["light", "dark"]}},
				"required": ["theme"],
				"additionalProperties": false
			}`),
			value:    `{"theme":"blue","extraField":"x"}`,
			expected: errVariationValueSchemaViolation,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			f := &Feature{Feature: &featureproto.Feature{
				VariationType:        p.variationType,
				VariationValueSchema: p.schema,
			}}
			err := f.validateVariationValueAgainstSchema(p.value)
			assert.Equal(t, p.expected, err)
		})
	}
}

func TestNewFeatureValidatesVariationValueSchema(t *testing.T) {
	t.Parallel()
	v1, err := uuid.NewUUID()
	require.NoError(t, err)
	v2, err := uuid.NewUUID()
	require.NoError(t, err)
	_, err = NewFeature(
		"feature-id",
		"feature name",
		"",
		featureproto.Feature_STRING,
		enumSchema("A"),
		[]*featureproto.Variation{
			{Id: v1.String(), Value: "A", Name: "A"},
			{Id: v2.String(), Value: "B", Name: "B"},
		},
		[]string{"tag"},
		0,
		1,
		"test@example.com",
	)
	assert.Equal(t, errVariationValueSchemaViolation, err)
}

func TestUpdateVariationValueSchemaValidatesFinalValues(t *testing.T) {
	t.Parallel()
	t.Run("rejects schema that existing variations violate", func(t *testing.T) {
		t.Parallel()
		original := makeFeature("test-feature")
		_, err := original.Update(
			nil, nil, nil, nil, nil, nil, nil, false,
			nil, nil, nil, nil, nil, nil, nil,
			&VariationValueSchemaUpdate{Schema: enumSchema("A", "B")},
		)
		assert.Equal(t, errVariationValueSchemaViolation, err)
	})
	t.Run("validates values against schema changed in same update", func(t *testing.T) {
		t.Parallel()
		original := makeFeature("test-feature")
		updated, err := original.Update(
			nil, nil, nil, nil, nil, nil, nil, false,
			nil, nil, nil,
			[]*featureproto.VariationChange{
				{
					ChangeType: featureproto.ChangeType_UPDATE,
					Variation: &featureproto.Variation{
						Id:          "variation-C",
						Value:       "D",
						Name:        "Variation C",
						Description: "Thing does C",
					},
				},
			},
			nil, nil, nil,
			&VariationValueSchemaUpdate{Schema: enumSchema("A", "B", "D")},
		)
		require.NoError(t, err)
		assert.Equal(t, "D", updated.Variations[2].Value)
	})
}
