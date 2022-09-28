// Copyright 2022 The Bucketeer Authors.
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

package webhookhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/itchyny/gojq"

	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func evaluateClause(
	ctx context.Context,
	clause *autoopsproto.WebhookClause,
	payload interface{},
) (bool, error) {
	if len(clause.Conditions) == 0 {
		return false, fmt.Errorf("WebhookClause has no conditions")
	}
	// All conditions are combined with implicit AND
	for _, condition := range clause.Conditions {
		asmt, err := evaluateCondition(ctx, condition, payload)
		if err != nil {
			return false, err
		}
		if asmt {
			continue
		}
		return false, nil
	}
	return true, nil
}

func evaluateCondition(
	ctx context.Context,
	condition *autoopsproto.WebhookClause_Condition,
	payload interface{},
) (bool, error) {
	filtered, err := filterWebhookPayload(
		ctx,
		condition.Filter,
		payload,
	)
	if err != nil {
		return false, err
	}
	if filtered == nil {
		return false, nil
	}
	var specified interface{}
	if err := json.Unmarshal([]byte(condition.Value), &specified); err != nil {
		return false, err
	}
	cFiltered, err := convertType(reflect.TypeOf(specified), filtered)
	if err != nil {
		return false, err
	}
	switch op := condition.Operator; op {
	case autoopsproto.WebhookClause_Condition_EQUAL:
		return specified == cFiltered, nil
	case autoopsproto.WebhookClause_Condition_NOT_EQUAL:
		return specified != cFiltered, nil
	case autoopsproto.WebhookClause_Condition_MORE_THAN:
		specifiedF, cFilteredF, ok := extractFloat64Values(specified, cFiltered)
		if ok {
			return specifiedF < cFilteredF, nil
		}
		specifiedS, cFilteredS, ok := extractStringValues(specified, cFiltered)
		if ok {
			return specifiedS < cFilteredS, nil
		}
		return false, fmt.Errorf("Failed to evaluate %v < %v", specified, cFiltered)
	case autoopsproto.WebhookClause_Condition_MORE_THAN_OR_EQUAL:
		specifiedF, cFilteredF, ok := extractFloat64Values(specified, cFiltered)
		if ok {
			return specifiedF <= cFilteredF, nil
		}
		specifiedS, cFilteredS, ok := extractStringValues(specified, cFiltered)
		if ok {
			return specifiedS <= cFilteredS, nil
		}
		return false, fmt.Errorf("Failed to evaluate %v <= %v", specified, cFiltered)
	case autoopsproto.WebhookClause_Condition_LESS_THAN:
		specifiedF, cFilteredF, ok := extractFloat64Values(specified, cFiltered)
		if ok {
			return specifiedF > cFilteredF, nil
		}
		specifiedS, cFilteredS, ok := extractStringValues(specified, cFiltered)
		if ok {
			return specifiedS > cFilteredS, nil
		}
		return false, fmt.Errorf("Failed to evaluate %v > %v", specified, cFiltered)
	case autoopsproto.WebhookClause_Condition_LESS_THAN_OR_EQUAL:
		specifiedF, cFilteredF, ok := extractFloat64Values(specified, cFiltered)
		if ok {
			return specifiedF >= cFilteredF, nil
		}
		specifiedS, cFilteredS, ok := extractStringValues(specified, cFiltered)
		if ok {
			return specifiedS >= cFilteredS, nil
		}
		return false, fmt.Errorf("Failed to evaluate %v >= %v", specified, cFiltered)
	default:
		return false, fmt.Errorf("Unknown operation: %s", op)
	}
}

func filterWebhookPayload(
	ctx context.Context,
	filter string,
	payload interface{},
) (interface{}, error) {
	query, err := gojq.Parse(filter)
	if err != nil {
		return nil, err
	}
	iter := query.RunWithContext(ctx, payload)
	var results []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		if v == nil {
			continue
		}
		results = append(results, v)
	}
	if len(results) > 1 {
		return false, fmt.Errorf("Got multiple values: %s", filter)
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return nil, nil
}

func extractFloat64Values(i1, i2 interface{}) (float64, float64, bool) {
	// json.Unmarshal make JSON number into float64
	var (
		f1, f2 float64
		ok     bool
	)
	if f1, ok = i1.(float64); !ok {
		return f1, f2, false
	}
	if f2, ok = i2.(float64); !ok {
		return f1, f2, false
	}
	return f1, f2, true
}

func extractStringValues(i1, i2 interface{}) (string, string, bool) {
	var (
		s1, s2 string
		ok     bool
	)
	if s1, ok = i1.(string); !ok {
		return s1, s2, false
	}
	if s2, ok = i2.(string); !ok {
		return s1, s2, false
	}
	return s1, s2, true
}

func convertType(targetType reflect.Type, original interface{}) (interface{}, error) {
	if targetType == reflect.TypeOf(original) {
		return original, nil
	}
	switch targetType.Kind() {
	case reflect.String:
		if f, ok := original.(float64); ok {
			converted := strconv.FormatFloat(f, 'f', -1, 64)
			return converted, nil
		}
		return nil, fmt.Errorf("Failed to convert %v to %v", original, targetType)
	case reflect.Float64:
		if s, ok := original.(string); ok {
			converted, err := strconv.ParseFloat(s, 64)
			return converted, err
		}
		return nil, fmt.Errorf("Failed to convert %v to %v", original, targetType)
	default:
		return nil, fmt.Errorf("Not supported type %v", targetType)
	}
}
