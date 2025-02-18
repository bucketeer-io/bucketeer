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
	"crypto/sha256"
	"encoding/binary"
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
		variationID, err := e.rollout(strategy.RolloutStrategy, userID, featureID, samplingSeed)
		if err != nil {
			return nil, err
		}
		return findVariation(variationID, variations)
	}
	return nil, ErrUnsupportedStrategy
}

func (e *strategyEvaluator) rollout(
	strategy *feature.RolloutStrategy,
	userID, featureID, samplingSeed string,
) (string, error) {
	bucket := e.bucketSHA256(featureID, userID, samplingSeed)
	sum := 0.0
	for i := range strategy.Variations {
		sum += float64(strategy.Variations[i].Weight) / 100000.0
		if bucket < sum {
			return strategy.Variations[i].Variation, nil
		}
	}
	return "", ErrVariationNotFound
}

// SHA-256 Bucketing Function
func (e *strategyEvaluator) bucketSHA256(featureID, userID, samplingSeed string) float64 {
	// Format input string correctly
	input := fmt.Sprintf("%s-%s-%s", featureID, userID, samplingSeed)

	// Compute SHA-256 hash
	hash := sha256.Sum256([]byte(input))

	// Convert the first 8 bytes to a uint64 using BigEndian
	num := binary.BigEndian.Uint64(hash[:8])

	// Normalize to [0,1) range
	return float64(num) / float64(^uint64(0)) // 2^64 - 1
}
