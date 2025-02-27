// Copyright 2025 The Bucketeer Authors.
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

package evaluation

import (
	"fmt"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

type strategyEvaluator struct {
}

func (e *strategyEvaluator) Evaluate(
	strategy *feature.Strategy,
	userID string,
	variations []*feature.Variation,
	featureID string,
	samplingSeed string,
) (*feature.Variation, error) {
	switch strategy.Type {
	case feature.Strategy_FIXED:
		return findVariation(strategy.FixedStrategy.Variation, variations)
	case feature.Strategy_ROLLOUT:
		variationID, err := e.rollout(strategy.RolloutStrategy, featureID, userID, samplingSeed)
		if err != nil {
			return nil, err
		}
		return findVariation(variationID, variations)
	}
	return nil, ErrUnsupportedStrategy
}

func (e *strategyEvaluator) rollout(
	strategy *feature.RolloutStrategy,
	featureID, userID, samplingSeed string,
) (string, error) {
	b := bucketeer{}
	bucket := b.bucket(fmt.Sprintf("%s-%s-%s", featureID, userID, samplingSeed))
	// Iterate through the variant and increment the threshold by the percentage of each variant.
	// return the first variant where the bucket is smaller than the threshold.
	rangeEnd := 0.0
	for i := range strategy.Variations {
		rangeEnd += float64(strategy.Variations[i].Weight) / 100000.0
		if bucket < rangeEnd {
			return strategy.Variations[i].Variation, nil
		}
	}
	return "", ErrVariationNotFound
}
