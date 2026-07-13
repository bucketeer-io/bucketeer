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
	"math"
	"regexp"
	"strconv"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"

	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

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

func (f *Feature) newVariationValueValidator() (func(string) error, error) {
	schema := f.VariationValueSchema
	if schema == nil {
		return func(string) error { return nil }, nil
	}
	switch schema.Type {
	case featureproto.VariationValueSchema_ENUM:
		validator := schema.GetEnumValidator()
		if err := validateEnumSchemaDefinition(f.VariationType, validator); err != nil {
			return nil, err
		}
		return func(value string) error {
			return f.validateEnumVariationValue(validator, value)
		}, nil
	case featureproto.VariationValueSchema_REGEX:
		validator := schema.GetRegexValidator()
		if err := validateRegexSchemaDefinition(f.VariationType, validator); err != nil {
			return nil, err
		}
		pattern, err := compileRegexVariationValueValidator(validator)
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
		validator := schema.GetJsonSchemaValidator()
		if err := validateJSONSchemaDefinition(f.VariationType, validator); err != nil {
			return nil, err
		}
		compiled, err := compileJSONSchema(validator.Schema)
		if err != nil {
			return nil, errVariationValueSchemaInvalid
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

func validateEnumSchemaDefinition(
	variationType featureproto.Feature_VariationType,
	validator *featureproto.VariationValueSchema_EnumValidator,
) error {
	if variationType != featureproto.Feature_STRING && variationType != featureproto.Feature_NUMBER {
		return errVariationValueSchemaTypeUnmatched
	}
	if validator == nil || len(validator.Values) == 0 {
		return errVariationValueSchemaInvalid
	}
	if variationType == featureproto.Feature_NUMBER {
		for _, value := range validator.Values {
			if _, err := parseFiniteFloat(value); err != nil {
				return errVariationValueSchemaInvalid
			}
		}
	}
	return nil
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
		target, err := parseFiniteFloat(value)
		if err != nil {
			return errVariationTypeUnmatched
		}
		for _, enumValue := range validator.Values {
			allowed, err := parseFiniteFloat(enumValue)
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

func validateRegexSchemaDefinition(
	variationType featureproto.Feature_VariationType,
	validator *featureproto.VariationValueSchema_RegexValidator,
) error {
	if variationType != featureproto.Feature_STRING {
		return errVariationValueSchemaTypeUnmatched
	}
	if validator == nil || validator.Pattern == "" {
		return errVariationValueSchemaInvalid
	}
	return nil
}

func parseFiniteFloat(value string) (float64, error) {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return 0, fmt.Errorf("feature: number must be finite")
	}
	return parsed, nil
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

func validateJSONSchemaDefinition(
	variationType featureproto.Feature_VariationType,
	validator *featureproto.VariationValueSchema_JsonSchemaValidator,
) error {
	if variationType != featureproto.Feature_JSON {
		return errVariationValueSchemaTypeUnmatched
	}
	if validator == nil {
		return errVariationValueSchemaInvalid
	}
	if validator.Schema == "" {
		return errVariationValueSchemaInvalid
	}
	return nil
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
