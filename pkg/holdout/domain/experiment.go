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
	"math"
	"slices"

	"google.golang.org/protobuf/proto"

	experimentdomain "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/domain"
	experimentproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	// InHoldoutVariationID is the base variation of an internal experiment:
	// the users the group holds out.
	InHoldoutVariationID = "in_holdout"
	// NotInHoldoutVariationID is the variation the held-out users are
	// compared with.
	NotInHoldoutVariationID = "not_in_holdout"
	// InternalExperimentFeatureVersion is the feature version of every
	// internal experiment. The feature ID is the ID of the group.
	InternalExperimentFeatureVersion int32 = 1
)

// NewInternalExperiment creates the internal experiment that the group needs
// and does not have, and stores its ID on the group. It returns nil when the
// group needs none. previousStatus is the status of the group before the
// change being applied.
//
// Only a serving group with goals needs one:
//   - RUNNING without experiment_id: the experiment of the whole period. It
//     starts at the group's start_at when the group is starting, and at now
//     when the group was already serving.
//   - ANALYSIS without analysis_experiment_id: the experiment of the analysis
//     period. It starts at analysis_started_at when the group is entering
//     ANALYSIS, and at now when the group was already in ANALYSIS.
//
// A group that gets its first goals while in ANALYSIS has only the second
// one: an experiment of the whole period would have the same window.
func (g *HoldoutGroup) NewInternalExperiment(
	previousStatus featureproto.HoldoutGroup_Status,
	now int64,
) (*experimentdomain.Experiment, error) {
	if len(g.GoalIds) == 0 {
		return nil, nil
	}
	switch {
	case g.Status == featureproto.HoldoutGroup_RUNNING && g.ExperimentId == "":
		startAt := now
		if previousStatus == featureproto.HoldoutGroup_WAITING {
			startAt = g.StartAt
		}
		experiment, err := g.newInternalExperiment(startAt, now)
		if err != nil {
			return nil, err
		}
		g.ExperimentId = experiment.Id
		return experiment, nil
	case g.Status == featureproto.HoldoutGroup_ANALYSIS && g.AnalysisExperimentId == "":
		startAt := now
		if previousStatus == featureproto.HoldoutGroup_RUNNING {
			startAt = g.AnalysisStartedAt
		}
		experiment, err := g.newInternalExperiment(startAt, now)
		if err != nil {
			return nil, err
		}
		g.AnalysisExperimentId = experiment.Id
		return experiment, nil
	}
	return nil, nil
}

func (g *HoldoutGroup) newInternalExperiment(startAt, now int64) (*experimentdomain.Experiment, error) {
	experiment, err := experimentdomain.NewExperiment(
		g.Id,
		InternalExperimentFeatureVersion,
		[]*featureproto.Variation{
			{Id: InHoldoutVariationID, Value: InHoldoutVariationID, Name: InHoldoutVariationID},
			{Id: NotInHoldoutVariationID, Value: NotInHoldoutVariationID, Name: NotInHoldoutVariationID},
		},
		g.GoalIds,
		startAt,
		g.StopAt,
		g.Name,
		g.Description,
		InHoldoutVariationID,
		g.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	experiment.Kind = experimentproto.Experiment_HOLDOUT
	experiment.CreatedAt = now
	g.MirrorToExperiment(experiment, now)
	return experiment, nil
}

// MirrorToExperiment writes the state of the group onto one of its internal
// experiments and reports whether the experiment changed.
//
// The name, the description and the goals are the group's. A serving group
// has a RUNNING experiment that stops at the group's stop_at. A stopped group
// has a STOPPED experiment: it stops at the group's stop_at after a scheduled
// stop, and at the group's stopped_at after a manual stop. It is never
// FORCE_STOPPED, because the readers of the analytics pipeline list only
// RUNNING and STOPPED experiments. start_at, archived and deleted are left as
// they are.
func (g *HoldoutGroup) MirrorToExperiment(experiment *experimentdomain.Experiment, now int64) bool {
	before := proto.Clone(experiment.Experiment)
	experiment.Name = g.Name
	experiment.Description = g.Description
	experiment.GoalIds = slices.Clone(g.GoalIds)
	switch g.Status {
	case featureproto.HoldoutGroup_RUNNING, featureproto.HoldoutGroup_ANALYSIS:
		experiment.Status = experimentproto.Experiment_RUNNING
		experiment.StopAt = g.StopAt
		experiment.StoppedAt = math.MaxInt64
	case featureproto.HoldoutGroup_STOPPED:
		experiment.Status = experimentproto.Experiment_STOPPED
		experiment.StopAt = g.StopAt
		experiment.StoppedAt = g.StoppedAt
	case featureproto.HoldoutGroup_FORCE_STOPPED:
		experiment.Status = experimentproto.Experiment_STOPPED
		experiment.StopAt = g.StoppedAt
		experiment.StoppedAt = g.StoppedAt
	}
	if proto.Equal(before, experiment.Experiment) {
		return false
	}
	experiment.UpdatedAt = now
	return true
}
