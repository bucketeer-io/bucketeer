// Copyright 2024 The Bucketeer Authors.
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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

const max = float64(0xffffffffffffffff)

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
	return nil, errUnsupportedStrategy
}

func (e *strategyEvaluator) rollout(
	strategy *feature.RolloutStrategy,
	userID, featureID, samplingSeed string,
) (string, error) {
	bucket, err := e.bucket(userID, featureID, samplingSeed)
	if err != nil {
		return "", err
	}
	sum := 0.0
	for i := range strategy.Variations {
		sum += float64(strategy.Variations[i].Weight) / 100000.0
		if bucket < sum {
			return strategy.Variations[i].Variation, nil
		}
	}
	return "", errVariationNotFound
}

func (e *strategyEvaluator) bucket(userID string, featureID string, samplingSeed string) (float64, error) {
	hash := e.hash(userID, featureID, samplingSeed)
	// use first 16 characters (hex string) / first 8 bytes (byte array)
	intVal, err := strconv.ParseUint(hex.EncodeToString(hash[:])[:16], 16, 64)
	if err != nil {
		return 0.0, err
	}
	return float64(intVal) / max, nil
}

func (e *strategyEvaluator) hash(userID string, featureID string, samplingSeed string) [16]byte {
	// concat feature test id and user id
	// TODO: explain why this makes sense? Why does it make sense to add 'prerequisit' key here?
	concat := fmt.Sprintf("%s-%s%s", featureID, userID, samplingSeed)
	// returns 16 bytes which, if shown as hex string, has 32 characters
	return md5.Sum([]byte(concat))
}
