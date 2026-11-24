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

package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The experiment filter of the feature list joins the experiments on a feature
// flag only. An experiment whose kind is HOLDOUT stores a holdout group ID in
// feature_id, so it must not match a feature flag that has the same ID.
func TestFeaturesByExperimentQueriesJoinOnlyExperimentsOnFeatureFlags(t *testing.T) {
	t.Parallel()
	for name, query := range map[string]string{
		"select_features_by_experiment.sql": selectFeaturesByExperimentSQLQuery,
		"count_features_by_experiment.sql":  countFeaturesByExperimentSQLQuery,
	} {
		assert.Contains(t, query, "experiment.kind = 0", name)
	}
}
