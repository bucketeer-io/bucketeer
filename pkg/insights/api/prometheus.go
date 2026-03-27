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

package api

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	recordingRuleLatencyAvg       = "environment_id:source_id:method:bucketeer_gateway_api_handling_seconds:avg:rate5m"
	recordingRuleRequestTotal     = "environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m"
	recordingRuleEvaluationsTotal = "environment_id:pod:evaluation_type:source_id:" +
		"bucketeer_api_gateway_evaluations_total:rate5m"
	recordingRuleErrorRate = "environment_id:source_id:method:bucketeer_gateway_api_error_rate:rate5m"
)

func latencyQuery(envIDs, sourceIDs, apiIDs []string) string {
	return recordingRuleLatencyAvg + buildSelector(envIDs, sourceIDs, apiIDs)
}

func requestCountQuery(envIDs, sourceIDs, apiIDs []string) string {
	return recordingRuleRequestTotal + buildSelector(envIDs, sourceIDs, apiIDs)
}

func evaluationsQuery(envIDs, sourceIDs []string) string {
	return fmt.Sprintf(
		"sum by (environment_id, source_id, evaluation_type) (%s%s)",
		recordingRuleEvaluationsTotal,
		buildSelector(envIDs, sourceIDs, nil),
	)
}

func errorRatesQuery(envIDs, sourceIDs, apiIDs []string) string {
	return recordingRuleErrorRate + buildSelector(envIDs, sourceIDs, apiIDs)
}

func buildSelector(envIDs, sourceIDs, apiIDs []string) string {
	var filters []string
	if len(envIDs) > 0 {
		filters = append(filters, fmt.Sprintf(`environment_id=~"%s"`, buildRegex(envIDs)))
	}
	if len(sourceIDs) > 0 {
		filters = append(filters, fmt.Sprintf(`source_id=~"%s"`, buildRegex(sourceIDs)))
	}
	if len(apiIDs) > 0 {
		filters = append(filters, fmt.Sprintf(`method=~"%s"`, buildRegex(apiIDs)))
	}
	if len(filters) == 0 {
		return ""
	}
	return fmt.Sprintf("{%s}", strings.Join(filters, ","))
}

func buildRegex(values []string) string {
	escaped := make([]string, len(values))
	for i, v := range values {
		escaped[i] = regexp.QuoteMeta(v)
	}
	if len(escaped) == 1 {
		return "^" + escaped[0] + "$"
	}
	return "^(" + strings.Join(escaped, "|") + ")$"
}
