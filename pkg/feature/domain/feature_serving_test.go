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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestFilterForServingKeepsRequiredDependencies(t *testing.T) {
	t.Parallel()
	disabled := &ftproto.Feature{
		Id:      "flag",
		Enabled: false,
		Prerequisites: []*ftproto.Prerequisite{
			{FeatureId: "nested", VariationId: "on"},
		},
	}
	nested := &ftproto.Feature{Id: "nested", Enabled: false}
	dependent := &ftproto.Feature{
		Id:           "dependent",
		Enabled:      true,
		OffVariation: "off",
		Prerequisites: []*ftproto.Prerequisite{
			{FeatureId: "flag", VariationId: "on"},
		},
	}
	rule := &ftproto.Feature{
		Id:           "rule-dependent",
		Enabled:      true,
		OffVariation: "off",
		Rules: []*ftproto.Rule{{
			Clauses: []*ftproto.Clause{{
				Operator:  ftproto.Clause_FEATURE_FLAG,
				Attribute: "flag",
				Values:    []string{"on"},
			}},
		}},
	}
	orphan := &ftproto.Feature{Id: "orphan", Enabled: false}
	onlyDropped := &ftproto.Feature{
		Id:      "only-dropped",
		Enabled: false,
		Prerequisites: []*ftproto.Prerequisite{
			{FeatureId: "orphan", VariationId: "on"},
		},
	}
	archived := &ftproto.Feature{
		Id:        "old",
		Enabled:   false,
		Archived:  true,
		UpdatedAt: time.Now().AddDate(0, 0, -40).Unix(),
	}
	needsArchived := &ftproto.Feature{
		Id:           "needs-old",
		Enabled:      true,
		OffVariation: "off",
		Prerequisites: []*ftproto.Prerequisite{
			{FeatureId: "old", VariationId: "on"},
		},
	}
	targeted := &ftproto.Feature{Id: "targeted", Enabled: false}

	got := FilterForServing([]*ftproto.Feature{
		nil,
		dependent,
		rule,
		disabled,
		nested,
		orphan,
		onlyDropped,
		archived,
		needsArchived,
		targeted,
	}, map[string]struct{}{"targeted": {}})
	ids := make([]string, 0, len(got))
	for _, feature := range got {
		ids = append(ids, feature.Id)
	}
	assert.Equal(t, []string{
		"dependent",
		"rule-dependent",
		"flag",
		"nested",
		"needs-old",
		"targeted",
	}, ids)
}
