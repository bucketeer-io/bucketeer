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
	"errors"
	"math"
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	experimentproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrExperimentBeforeStart       = errors.New("experiment: start timestamp is greater than now")
	ErrExperimentBeforeStop        = errors.New("experiment: stop timestamp is greater than now")
	ErrExperimentStatusInvalid     = errors.New("experiment: experiment status is invalid")
	ErrExperimentAlreadyStopped    = errors.New("experiment: experiment is already stopped")
	ErrExperimentStartIsAfterStop  = errors.New("experiment: start is after stop timestamp")
	ErrExperimentStopIsBeforeStart = errors.New("experiment: stop is before start timestamp")
	ErrExperimentStopIsBeforeNow   = errors.New("experiment: stop is same or older than now timestamp")
)

type Experiment struct {
	*experimentproto.Experiment
}

func NewExperiment(
	featureID string,
	featureVersion int32,
	variations []*featureproto.Variation,
	goalIDs []string,
	startAt int64,
	stopAt int64,
	name string,
	description string,
	baseVariationID string,
	maintainer string) (*Experiment, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	goalIDs = removeDuplicated(goalIDs)
	now := time.Now().Unix()
	return &Experiment{
		&experimentproto.Experiment{
			Id:              id.String(),
			FeatureId:       featureID,
			FeatureVersion:  featureVersion,
			Variations:      variations,
			GoalIds:         goalIDs,
			StartAt:         startAt,
			StopAt:          stopAt,
			StoppedAt:       math.MaxInt64,
			CreatedAt:       now,
			UpdatedAt:       now,
			Name:            name,
			Description:     description,
			BaseVariationId: baseVariationID,
			Status:          experimentproto.Experiment_WAITING,
			Maintainer:      maintainer,
		},
	}, nil
}

func removeDuplicated(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for _, v := range args {
		if _, duplicated := encountered[v]; !duplicated {
			results = append(results, v)
			encountered[v] = true
		}
	}
	return results
}

func (e *Experiment) Start() error {
	if e.Status != experimentproto.Experiment_WAITING {
		return ErrExperimentStatusInvalid
	}
	now := time.Now().Unix()
	if e.StartAt > now {
		return ErrExperimentBeforeStart
	}
	e.Experiment.Status = experimentproto.Experiment_RUNNING
	e.Experiment.UpdatedAt = now
	return nil
}

func (e *Experiment) Finish() error {
	if e.Status != experimentproto.Experiment_WAITING && e.Status != experimentproto.Experiment_RUNNING {
		return ErrExperimentStatusInvalid
	}
	now := time.Now().Unix()
	if e.StopAt > now {
		return ErrExperimentBeforeStop
	}
	e.Experiment.Status = experimentproto.Experiment_STOPPED
	e.Experiment.UpdatedAt = now
	return nil
}

func (e *Experiment) Stop() error {
	if e.Stopped {
		return ErrExperimentAlreadyStopped
	}
	now := time.Now().Unix()
	e.Experiment.Stopped = true
	e.Experiment.Status = experimentproto.Experiment_FORCE_STOPPED
	e.Experiment.StoppedAt = now
	e.Experiment.UpdatedAt = now
	return nil
}

func (e *Experiment) ChangePeriod(startAt, stopAt int64) error {
	if err := e.validatePeriod(startAt, stopAt); err != nil {
		return err
	}
	e.Experiment.StartAt = startAt
	e.Experiment.StopAt = stopAt
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) validatePeriod(startAt, stopAt int64) error {
	if startAt >= stopAt {
		return ErrExperimentStartIsAfterStop
	}
	if stopAt <= time.Now().Unix() {
		return ErrExperimentStopIsBeforeNow
	}
	return nil
}

func (e *Experiment) ChangeStartAt(startAt int64) error {
	if err := e.validateStartAt(startAt); err != nil {
		return err
	}
	e.Experiment.StartAt = startAt
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) validateStartAt(startAt int64) error {
	if startAt >= e.Experiment.StopAt {
		return ErrExperimentStartIsAfterStop
	}
	return nil
}

func (e *Experiment) ChangeStopAt(stopAt int64) error {
	if err := e.validateStopAt(stopAt); err != nil {
		return err
	}
	e.Experiment.StopAt = stopAt
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) ChangeName(name string) error {
	e.Experiment.Name = name
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) ChangeDescription(description string) error {
	e.Experiment.Description = description
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) validateStopAt(stopAt int64) error {
	if stopAt <= e.Experiment.StartAt {
		return ErrExperimentStopIsBeforeStart
	}
	if stopAt <= time.Now().Unix() {
		return ErrExperimentStopIsBeforeNow
	}
	return nil
}

func (e *Experiment) SetArchived() error {
	e.Experiment.Archived = true
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

func (e *Experiment) SetDeleted() error {
	e.Experiment.Deleted = true
	e.Experiment.UpdatedAt = time.Now().Unix()
	return nil
}

// SyncGoalIDs syncs goalID and goalIDs.
// FIXME: This function is needed until admin UI implements multiple goals.
func SyncGoalIDs(goalID string, goalIDs []string) (string, []string) {
	if goalID == "" && len(goalIDs) == 0 {
		return "", nil
	}
	if goalID == "" {
		return goalIDs[0], goalIDs
	}
	if len(goalIDs) == 0 {
		return goalID, []string{goalID}
	}
	return goalID, goalIDs
}

// IsNotFinished returns true if the status is either waiting or running.
func (e *Experiment) IsNotFinished(t time.Time) bool {
	if e.Experiment.Deleted {
		return false
	}
	if e.Experiment.StopAt <= t.Unix() {
		return false
	}
	if e.Experiment.StoppedAt <= t.Unix() {
		return false
	}
	return true
}
