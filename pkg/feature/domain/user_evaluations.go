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
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"time"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

type UserEvaluations struct {
	*feature.UserEvaluations
}

func NewUserEvaluations(
	id string,
	evaluations []*feature.Evaluation,
	archivedFeaturesIds []string,
	forceUpdate bool,
) *UserEvaluations {
	now := time.Now().Unix()
	return &UserEvaluations{&feature.UserEvaluations{
		Id:                 id,
		Evaluations:        evaluations,
		CreatedAt:          now,
		ArchivedFeatureIds: archivedFeaturesIds,
		ForceUpdate:        forceUpdate,
	}}
}

func UserEvaluationsID(userID string, userMetadata map[string]string, features []*feature.Feature) string {
	sort.SliceStable(features, func(i, j int) bool {
		return features[i].Id < features[j].Id
	})
	// TODO: consider about a better hash algorithm?
	h := fnv.New64a()
	h.Write([]byte(userID)) // nolint:errcheck
	keys := sortMapKeys(userMetadata)
	for _, key := range keys {
		fmt.Fprintf(h, "%s:%s", key, userMetadata[key])
	}
	for _, feature := range features {
		fmt.Fprintf(h, "%s:%d", feature.Id, feature.Version)
	}
	return strconv.FormatUint(h.Sum64(), 10)
}

func sortMapKeys(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
