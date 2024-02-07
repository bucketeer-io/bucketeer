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
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestNewExperiment(t *testing.T) {
	featureID := "id"
	featureVersion := int32(1)
	variations := []*featureproto.Variation{
		{
			Value:       "A",
			Name:        "Variation A",
			Description: "Thing does A",
		},
		{
			Value:       "B",
			Name:        "Variation B",
			Description: "Thing does B",
		},
		{
			Value:       "C",
			Name:        "Variation C",
			Description: "Thing does C",
		},
	}
	goalIDs := []string{"id-1", "id-2"}
	startAt := int64(10)
	stopAt := int64(20)
	name := "name"
	description := "description"
	baseVariationId := "baseVariationId"
	maintainer := "bucketeer@example.com"

	e, err := NewExperiment(
		featureID,
		featureVersion,
		variations,
		goalIDs,
		startAt,
		stopAt,
		name,
		description,
		baseVariationId,
		maintainer,
	)

	assert.NoError(t, err)
	assert.Equal(t, featureID, e.FeatureId)
	assert.Equal(t, featureVersion, e.FeatureVersion)
	if !reflect.DeepEqual(variations, e.Variations) {
		t.Fatal("Variations not equal")
	}
	if !reflect.DeepEqual(goalIDs, e.GoalIds) {
		t.Fatal("GoalIDs not equal")
	}
	assert.Equal(t, startAt, e.StartAt)
	assert.Equal(t, stopAt, e.StopAt)
	assert.Equal(t, name, e.Name)
	assert.Equal(t, description, e.Description)
	assert.Equal(t, baseVariationId, e.BaseVariationId)
	assert.Equal(t, maintainer, e.Maintainer)
}

func TestRenameExperiment(t *testing.T) {
	t.Parallel()
	e := newExperiment(t)
	newName := "newGName"
	err := e.ChangeName(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, e.Name)
}

func TestChangeDescriptionExperiment(t *testing.T) {
	t.Parallel()
	e := newExperiment(t)
	newDesc := "newGDesc"
	err := e.ChangeDescription(newDesc)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, e.Description)
}

func TestSetArchivedExperiment(t *testing.T) {
	t.Parallel()
	e := newExperiment(t)
	err := e.SetArchived()
	assert.NoError(t, err)
	assert.True(t, e.Archived)
}

func TestSetDeletedExperiment(t *testing.T) {
	t.Parallel()
	e := newExperiment(t)
	err := e.SetDeleted()
	assert.NoError(t, err)
	assert.True(t, e.Deleted)
}

func TestChangeStartAt(t *testing.T) {
	t.Parallel()
	stopAt := time.Now().Unix()
	e := &Experiment{
		&experimentproto.Experiment{StopAt: stopAt},
	}
	patterns := []*struct {
		startAt  int64
		expected error
	}{
		{
			startAt:  stopAt + 1,
			expected: ErrExperimentStartIsAfterStop,
		},
		{
			startAt:  stopAt - 1,
			expected: nil,
		},
	}
	for i, p := range patterns {
		actual := e.ChangeStartAt(p.startAt)
		assert.Equal(t, p.expected, actual, "i=%s", i)
	}
}

func TestChangeStopAt(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	startAt := now - 10
	stopAt := now + 10
	e := &Experiment{
		&experimentproto.Experiment{StartAt: startAt},
	}
	patterns := []*struct {
		stopAt   int64
		expected error
	}{
		{
			stopAt:   startAt - 10,
			expected: ErrExperimentStopIsBeforeStart,
		},
		{
			stopAt:   now - 5,
			expected: ErrExperimentStopIsBeforeNow,
		},
		{
			stopAt:   stopAt,
			expected: nil,
		},
	}
	for i, p := range patterns {
		actual := e.ChangeStopAt(p.stopAt)
		assert.Equal(t, p.expected, actual, "i=%s", i)
	}
}

func TestChangePeriod(t *testing.T) {
	t.Parallel()
	now := time.Now().Unix()
	startAt := now - 10
	stopAt := now + 10
	e := &Experiment{
		&experimentproto.Experiment{StartAt: startAt, StopAt: stopAt},
	}
	patterns := []*struct {
		startAt  int64
		stopAt   int64
		expected error
	}{
		{
			startAt:  startAt - 10,
			stopAt:   startAt - 10,
			expected: ErrExperimentStartIsAfterStop,
		},
		{
			startAt:  startAt - 10,
			stopAt:   startAt - 11,
			expected: ErrExperimentStartIsAfterStop,
		},
		{
			startAt:  startAt - 10,
			stopAt:   startAt - 9,
			expected: ErrExperimentStopIsBeforeNow,
		},
		{
			startAt:  startAt + 10,
			stopAt:   stopAt + 10,
			expected: nil,
		},
	}
	for i, p := range patterns {
		actual := e.ChangePeriod(p.startAt, p.stopAt)
		assert.Equal(t, p.expected, actual, "i=%s", i)
	}
}

func TestRemoveDuplicated(t *testing.T) {
	t.Parallel()
	patterns := []*struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"gid"},
			expected: []string{"gid"},
		},
		{
			input:    []string{"gid", "gid"},
			expected: []string{"gid"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
	}
	for i, p := range patterns {
		actual := removeDuplicated(p.input)
		assert.Equal(t, p.expected, actual, "i=%s", i)
	}
}

func TestStartExperiment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       *Experiment
		expectedErr error
	}{
		{
			desc: "error not waiting",
			input: &Experiment{&experimentproto.Experiment{
				Id:     "eID",
				Status: experimentproto.Experiment_RUNNING,
			}},
			expectedErr: ErrExperimentStatusInvalid,
		},
		{
			desc: "error before start at",
			input: &Experiment{&experimentproto.Experiment{
				Id:      "eID",
				Status:  experimentproto.Experiment_WAITING,
				StartAt: time.Now().AddDate(0, 0, 1).Unix(),
			}},
			expectedErr: ErrExperimentBeforeStart,
		},
		{
			desc: "success",
			input: &Experiment{&experimentproto.Experiment{
				Id:      "eID",
				Status:  experimentproto.Experiment_WAITING,
				StartAt: time.Now().Unix(),
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.input.Start()
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, p.input.Experiment.Status, experimentproto.Experiment_RUNNING)
			}
		})
	}
}

func TestFinishExperiment(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       *Experiment
		expectedErr error
	}{
		{
			desc: "error invalid status",
			input: &Experiment{&experimentproto.Experiment{
				Id:     "eID",
				Status: experimentproto.Experiment_STOPPED,
			}},
			expectedErr: ErrExperimentStatusInvalid,
		},
		{
			desc: "error before stop at",
			input: &Experiment{&experimentproto.Experiment{
				Id:     "eID",
				Status: experimentproto.Experiment_RUNNING,
				StopAt: time.Now().AddDate(0, 0, 1).Unix(),
			}},
			expectedErr: ErrExperimentBeforeStop,
		},
		{
			desc: "success: waiting",
			input: &Experiment{&experimentproto.Experiment{
				Id:     "eID",
				Status: experimentproto.Experiment_WAITING,
				StopAt: time.Now().Unix(),
			}},
			expectedErr: nil,
		},
		{
			desc: "success: running",
			input: &Experiment{&experimentproto.Experiment{
				Id:     "eID",
				Status: experimentproto.Experiment_RUNNING,
				StopAt: time.Now().Unix(),
			}},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.input.Finish()
			if p.expectedErr != nil {
				assert.Equal(t, p.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, p.input.Experiment.Status, experimentproto.Experiment_STOPPED)
			}
		})
	}
}

func TestStopExperiment(t *testing.T) {
	t.Parallel()
	patterns := []*struct {
		input       *Experiment
		expectedErr error
	}{
		{
			input:       &Experiment{&experimentproto.Experiment{Id: "eID", Stopped: true, StopAt: time.Now().Unix()}},
			expectedErr: ErrExperimentAlreadyStopped,
		},
		{
			input:       &Experiment{&experimentproto.Experiment{Id: "eID"}},
			expectedErr: nil,
		},
	}
	for i, p := range patterns {
		oldStoppedAt := p.input.StoppedAt
		err := p.input.Stop()
		if p.expectedErr != nil {
			assert.Equal(t, p.expectedErr, err, "i=%s", i)
		} else {
			assert.NoError(t, err, "i=%s", i)
			assert.True(t, p.input.Stopped, "i=%s", i)
			assert.NotEqual(t, oldStoppedAt, p.input.StoppedAt, "i=%s", i)
			assert.Equal(t, p.input.Experiment.Status, experimentproto.Experiment_FORCE_STOPPED)
		}
	}
}

func TestSyncGoalIDs(t *testing.T) {
	t.Parallel()
	patterns := []*struct {
		goalID    string
		goalIDs   []string
		exGoalID  string
		exGoalIDs []string
	}{
		{
			goalID:    "gid",
			exGoalID:  "gid",
			exGoalIDs: []string{"gid"},
		},
		{
			goalIDs:   []string{"gid"},
			exGoalID:  "gid",
			exGoalIDs: []string{"gid"},
		},
		{
			goalID:    "",
			goalIDs:   []string{"gid-0", "gid-1"},
			exGoalID:  "gid-0",
			exGoalIDs: []string{"gid-0", "gid-1"},
		},
	}
	for i, p := range patterns {
		gid, gids := SyncGoalIDs(p.goalID, p.goalIDs)
		assert.Equal(t, p.exGoalID, gid, "i=%s", i)
		assert.Equal(t, p.exGoalIDs, gids, "i=%s", i)
	}
}

func TestIsNotFinished(t *testing.T) {
	t.Parallel()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t3, err := time.Parse(layout, "2014-01-19 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t4, err := time.Parse(layout, "2014-01-20 23:02:03 +0000 UTC")
	require.NoError(t, err)
	patterns := []struct {
		desc       string
		experiment *Experiment
		input      time.Time
		expected   bool
	}{
		{
			desc: "before StartAt",
			experiment: &Experiment{Experiment: &experimentproto.Experiment{
				Deleted:   false,
				StartAt:   t2.Unix(),
				StopAt:    t3.Unix(),
				StoppedAt: math.MaxInt64,
			}},
			input:    t1,
			expected: true,
		},
		{
			desc: "running",
			experiment: &Experiment{Experiment: &experimentproto.Experiment{
				Deleted:   false,
				StartAt:   t1.Unix(),
				StopAt:    t3.Unix(),
				StoppedAt: math.MaxInt64,
			}},
			input:    t2,
			expected: true,
		},
		{
			desc: "after StoppedAt, before StoppAt",
			experiment: &Experiment{Experiment: &experimentproto.Experiment{
				Deleted:   false,
				StartAt:   t1.Unix(),
				StopAt:    t4.Unix(),
				StoppedAt: t2.Unix(),
			}},
			input:    t3,
			expected: false,
		},
		{
			desc: "after StopAt and StoppedAt",
			experiment: &Experiment{Experiment: &experimentproto.Experiment{
				Deleted:   false,
				StartAt:   t1.Unix(),
				StopAt:    t3.Unix(),
				StoppedAt: t2.Unix(),
			}},
			input:    t4,
			expected: false,
		},
		{
			desc: "Deleted",
			experiment: &Experiment{Experiment: &experimentproto.Experiment{
				Deleted:   true,
				StartAt:   t1.Unix(),
				StopAt:    t3.Unix(),
				StoppedAt: math.MaxInt64,
			}},
			input:    t2,
			expected: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expected, p.experiment.IsNotFinished(p.input))
		})
	}
}

func newExperiment(t *testing.T) *Experiment {
	t.Helper()
	featureID := "id"
	featureVersion := int32(1)
	variations := []*featureproto.Variation{
		{
			Value:       "A",
			Name:        "Variation A",
			Description: "Thing does A",
		},
		{
			Value:       "B",
			Name:        "Variation B",
			Description: "Thing does B",
		},
		{
			Value:       "C",
			Name:        "Variation C",
			Description: "Thing does C",
		},
	}
	goalIDs := []string{"id-1", "id-2"}
	startAt := int64(10)
	stopAt := int64(20)
	name := "name"
	description := "description"
	baseVariationId := "baseVariationId"
	maintainer := "bucketeer@example.com"

	e, err := NewExperiment(
		featureID,
		featureVersion,
		variations,
		goalIDs,
		startAt,
		stopAt,
		name,
		description,
		baseVariationId,
		maintainer,
	)
	assert.NoError(t, err)
	return e
}
