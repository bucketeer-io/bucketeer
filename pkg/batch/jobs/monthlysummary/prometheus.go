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
	"time"

	"github.com/prometheus/common/model"
)

const metricRequestTotal = "bucketeer_gateway_api_request_total"

// queryRequestCounts retrieves request counts from Prometheus for all environment/sourceID combinations.
// Returns map[environmentID]map[sourceID]count.
func (m *monthlySummarizer) queryRequestCounts(
	ctx context.Context,
	targetDate time.Time,
) (map[string]map[string]int64, error) {
	duration, evalTime := calculateTimeParams(targetDate)
	query := requestCountIncreaseQuery(duration)

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
	evaluationTime = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day()+1, 0, 0, 0, 0, targetDate.Location())
	duration = evaluationTime.Sub(monthStart)
	return
}

// requestCountIncreaseQuery builds an instant query for total request count over a duration.
// Use with Instant Query at the end of the duration.
func requestCountIncreaseQuery(duration time.Duration) string {
	seconds := int(duration.Seconds())
	return fmt.Sprintf(
		`sum by (environment_id,source_id) (increase(%s[%ds]))`,
		metricRequestTotal,
		seconds,
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
		if result[envID] == nil {
			result[envID] = make(map[string]int64)
		}
		result[envID][sourceID] = int64(sample.Value)
	}
	return result
}
