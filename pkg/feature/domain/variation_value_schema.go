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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func (f *Feature) validateVariationValueSchema() error {
	return validateVariationValueSchemaDefinition(f.VariationType, f.VariationValueSchema)
}

func (f *Feature) validateAllVariationValuesAgainstSchema() error {
	validateValue, err := f.newVariationValueValidator()
	if err != nil {
		return err
	}
	for _, variation := range f.Variations {
		if variation == nil {
			return errVariationRequired
		}
		if err := validateValue(variation.Value); err != nil {
			return err
		}
	}
	return nil
}

func validateVariationValueSchemaDefinition(
	variationType featureproto.Feature_VariationType,
	schema *featureproto.VariationValueSchema,
) error {
	if schema == nil {
		return nil
	}
	switch schema.Type {
	case featureproto.VariationValueSchema_ENUM:
		if variationType != featureproto.Feature_STRING && variationType != featureproto.Feature_NUMBER {
			return errVariationValueSchemaTypeUnmatched
		}
		validator := schema.GetEnumValidator()
		if validator == nil || len(validator.Values) == 0 {
			return errVariationValueSchemaInvalid
		}
		if variationType == featureproto.Feature_NUMBER {
			for _, value := range validator.Values {
				if _, err := strconv.ParseFloat(value, 64); err != nil {
					return errVariationValueSchemaInvalid
				}
			}
		}
		return nil
	case featureproto.VariationValueSchema_REGEX:
		if variationType != featureproto.Feature_STRING {
			return errVariationValueSchemaTypeUnmatched
		}
		validator := schema.GetRegexValidator()
		if validator == nil || validator.Pattern == "" {
			return errVariationValueSchemaInvalid
		}
		if _, err := regexp.Compile(validator.Pattern); err != nil {
			return errVariationValueSchemaInvalid
		}
		return nil
	case featureproto.VariationValueSchema_JSON_SCHEMA:
		if variationType != featureproto.Feature_JSON {
			return errVariationValueSchemaTypeUnmatched
		}
		validator := schema.GetJsonSchemaValidator()
		if validator == nil || validator.Schema == "" {
			return errVariationValueSchemaInvalid
		}
		if _, err := compileJSONSchema(validator.Schema); err != nil {
			return errVariationValueSchemaInvalid
		}
		return nil
	default:
		return errVariationValueSchemaInvalid
	}
}

func (f *Feature) validateVariationValueAgainstSchema(value string) error {
	validateValue, err := f.newVariationValueValidator()
	if err != nil {
		return err
	}
	return validateValue(value)
}

func (f *Feature) newVariationValueValidator() (func(string) error, error) {
	schema := f.VariationValueSchema
	if schema == nil {
		return func(string) error { return nil }, nil
	}
	if err := f.validateVariationValueSchema(); err != nil {
		return nil, err
	}
	switch schema.Type {
	case featureproto.VariationValueSchema_ENUM:
		validator := schema.GetEnumValidator()
		return func(value string) error {
			return f.validateEnumVariationValue(validator, value)
		}, nil
	case featureproto.VariationValueSchema_REGEX:
		pattern, err := compileRegexVariationValueValidator(schema.GetRegexValidator())
		if err != nil {
			return nil, err
		}
		return func(value string) error {
			if !pattern.MatchString(value) {
				return errVariationValueSchemaViolation
			}
			return nil
		}, nil
	case featureproto.VariationValueSchema_JSON_SCHEMA:
		compiled, err := compileJSONSchemaVariationValueValidator(schema.GetJsonSchemaValidator())
		if err != nil {
			return nil, err
		}
		return func(value string) error {
			jsonValue, err := jsonschema.UnmarshalJSON(strings.NewReader(value))
			if err != nil {
				return errVariationTypeUnmatched
			}
			if err := compiled.Validate(jsonValue); err != nil {
				return errVariationValueSchemaViolation
			}
			return nil
		}, nil
	default:
		return nil, errVariationValueSchemaInvalid
	}
}

func (f *Feature) validateEnumVariationValue(
	validator *featureproto.VariationValueSchema_EnumValidator,
	value string,
) error {
	if validator == nil {
		return errVariationValueSchemaInvalid
	}
	switch f.VariationType {
	case featureproto.Feature_NUMBER:
		target, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errVariationTypeUnmatched
		}
		for _, enumValue := range validator.Values {
			allowed, err := strconv.ParseFloat(enumValue, 64)
			if err != nil {
				return errVariationValueSchemaInvalid
			}
			if target == allowed {
				return nil
			}
		}
	case featureproto.Feature_STRING:
		for _, enumValue := range validator.Values {
			if value == enumValue {
				return nil
			}
		}
	}
	return errVariationValueSchemaViolation
}

func compileRegexVariationValueValidator(
	validator *featureproto.VariationValueSchema_RegexValidator,
) (*regexp.Regexp, error) {
	if validator == nil {
		return nil, errVariationValueSchemaInvalid
	}
	pattern, err := regexp.Compile(validator.Pattern)
	if err != nil {
		return nil, errVariationValueSchemaInvalid
	}
	return pattern, nil
}

func compileJSONSchemaVariationValueValidator(
	validator *featureproto.VariationValueSchema_JsonSchemaValidator,
) (*jsonschema.Schema, error) {
	if validator == nil {
		return nil, errVariationValueSchemaInvalid
	}
	schema, err := compileJSONSchema(validator.Schema)
	if err != nil {
		return nil, errVariationValueSchemaInvalid
	}
	return schema, nil
}

func compileJSONSchema(schema string) (*jsonschema.Schema, error) {
	document, err := jsonschema.UnmarshalJSON(strings.NewReader(schema))
	if err != nil {
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	if err := compiler.AddResource("schema.json", document); err != nil {
		return nil, err
	}
	compiled, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}
	if compiled == nil {
		return nil, fmt.Errorf("feature: compiled JSON schema is nil")
	}
	return compiled, nil
}
