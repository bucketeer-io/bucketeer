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

package monthlysummary

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/prometheus/common/model"
)

const (
	recordingRuleRequestTotal = "environment_id:source_id:method:bucketeer_gateway_api_request_total:rate5m"
	// Step size must match the recording rule's rate window (5m).
	recordingRuleStep = 5 * time.Minute
)

// queryRequestCounts retrieves request counts from Prometheus for all environment/sourceID combinations.
// Returns map[environmentID]map[sourceID]count.
func (m *monthlySummarizer) queryRequestCounts(
	ctx context.Context,
	targetDate time.Time,
) (map[string]map[string]int64, error) {
	duration, evalTime := calculateTimeParams(targetDate)
	query := requestCountQuery(duration)

	vector, err := m.promClient.QueryInstant(ctx, query, evalTime)
	if err != nil {
		return nil, err
	}

	return parseRequestCountVector(vector), nil
}

func calculateTimeParams(targetDate time.Time) (duration time.Duration, evaluationTime time.Time) {
	// Month start: first day of targetDate's month at 00:00
	monthStart := time.Date(targetDate.Year(), targetDate.Month(), 1, 0, 0, 0, 0, targetDate.Location())
	// Evaluation time: start of the day after targetDate
	nextDay := targetDate.Day() + 1
	loc := targetDate.Location()
	evaluationTime = time.Date(targetDate.Year(), targetDate.Month(), nextDay, 0, 0, 0, 0, loc)
	duration = evaluationTime.Sub(monthStart)
	return
}

// requestCountQuery builds an instant query using the recording rule.
// sum_over_time aggregates the rate5m values at 5m intervals,
// then * 300 converts the sum of rates to total request count.
// Subtract one step from the range because the subquery includes both ends.
func requestCountQuery(duration time.Duration) string {
	rangeDuration := duration - recordingRuleStep
	if rangeDuration < time.Second {
		rangeDuration = time.Second
	}
	seconds := int(rangeDuration.Seconds())
	stepSeconds := int(recordingRuleStep / time.Second)
	return fmt.Sprintf(
		`sum by (environment_id,source_id) (sum_over_time(%s[%ds:%dm])) * %d`,
		recordingRuleRequestTotal,
		seconds,
		int(recordingRuleStep.Minutes()),
		stepSeconds,
	)
}

func parseRequestCountVector(vector model.Vector) map[string]map[string]int64 {
	result := make(map[string]map[string]int64)
	for _, sample := range vector {
		envID := string(sample.Metric[model.LabelName("environment_id")])
		sourceID := string(sample.Metric[model.LabelName("source_id")])
		if envID == "" || sourceID == "" {
			continue
		}
		v := float64(sample.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
			continue
		}
		if result[envID] == nil {
			result[envID] = make(map[string]int64)
		}
		result[envID][sourceID] = int64(math.Round(v))
	}
	return result
}
