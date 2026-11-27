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

package processor

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	evaluation "github.com/bucketeer-io/bucketeer/v2/evaluation/go"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	holdoutdomain "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/domain"
	v2hs "github.com/bucketeer-io/bucketeer/v2/pkg/holdout/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/lock"
	redisv3 "github.com/bucketeer-io/bucketeer/v2/pkg/redis/v3"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
	epproto "github.com/bucketeer-io/bucketeer/v2/proto/eventpersisterdwh"
	exproto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
	ftproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	holdoutExposureLockKind = "holdout_exposure"
	// holdoutExposureLockTTL is how long the key of a written row lives. A
	// key is per interval of one day (holdoutExposureIntervals), so one day
	// from the first row covers the rest of the interval; the extra hour
	// covers events of its last minutes that are delivered after it ended.
	// An event of the interval arriving after the key expired writes one
	// more row for the user. The calculator counts distinct users, so that
	// changes the evaluation total of the variation and not its user count.
	holdoutExposureLockTTL = 25 * time.Hour
	// secondsPerDay is the distance between two cutoffs of the calculator.
	secondsPerDay = int64(24 * 60 * 60)
)

// holdoutExposure is a row derived from an evaluation event on a target flag
// of a measured holdout group, written against the synthetic identity of the
// group. The internal experiments of the group count their users from these
// rows: no SDK sends an evaluation event for the group itself.
type holdoutExposure struct {
	row *epproto.EvaluationEvent
	// eventID is the ID of the evaluation event the row was derived from.
	eventID string
	// lockKey and lockValue identify the key that makes the row unique.
	lockKey   string
	lockValue string
}

// holdoutExposureMatch is what an evaluation event counts for.
type holdoutExposureMatch struct {
	groupID     string
	variationID string
	// intervals names the interval of one day the event is in, for every
	// internal experiment of the group whose window contains the event
	// (holdoutExposureIntervals). With the group and the user it makes the
	// row unique.
	intervals string
	userID    string
}

func (m *holdoutExposureMatch) lockKey(environmentID string) string {
	return fmt.Sprintf(
		"%s:%s:%s:%s:%s",
		environmentID,
		holdoutExposureLockKind,
		m.groupID,
		m.intervals,
		m.userID,
	)
}

// measuredHoldoutGroup is a measured group together with its internal
// experiments that the environment's experiments list contains.
type measuredHoldoutGroup struct {
	group       *ftproto.HoldoutGroup
	experiments []*exproto.Experiment
}

// holdoutExposureDeriver derives the holdoutExposure rows.
type holdoutExposureDeriver struct {
	storage     v2hs.HoldoutStorage
	cache       cachev3.HoldoutMeasuredGroupsCache
	lock        *lock.DistributedLock
	flightgroup singleflight.Group
	now         func() time.Time
	logger      *zap.Logger
}

func newHoldoutExposureDeriver(
	storage v2hs.HoldoutStorage,
	cache cachev3.HoldoutMeasuredGroupsCache,
	lockClient redisv3.Client,
	logger *zap.Logger,
) *holdoutExposureDeriver {
	return &holdoutExposureDeriver{
		storage: storage,
		cache:   cache,
		lock:    lock.NewDistributedLock(lockClient, holdoutExposureLockTTL),
		now:     time.Now,
		logger:  logger,
	}
}

// measuredGroups returns the measured groups that own an internal experiment
// in experiments. When experiments has no internal experiment it returns nil
// and reads nothing, so an environment whose groups have no goals costs
// nothing here.
func (d *holdoutExposureDeriver) measuredGroups(
	ctx context.Context,
	environmentID string,
	experiments []*exproto.Experiment,
) ([]*measuredHoldoutGroup, error) {
	internal := make(map[string][]*exproto.Experiment)
	for _, e := range experiments {
		if e.Kind == exproto.Experiment_HOLDOUT {
			internal[e.FeatureId] = append(internal[e.FeatureId], e)
		}
	}
	if len(internal) == 0 {
		return nil, nil
	}
	groups, err := d.listGroups(ctx, environmentID, internal)
	if err != nil {
		return nil, err
	}
	measured := make([]*measuredHoldoutGroup, 0, len(internal))
	for _, group := range groups.GetHoldoutGroups() {
		if exps, ok := internal[group.Id]; ok {
			measured = append(measured, &measuredHoldoutGroup{group: group, experiments: exps})
		}
	}
	return measured, nil
}

// listGroups reads the measured groups from the cache. It reads them from
// the storage instead, and writes the result to the cache, when the cache
// cannot be read or when the cached list lacks the group of an internal
// experiment, which means the list is older than the experiments list.
func (d *holdoutExposureDeriver) listGroups(
	ctx context.Context,
	environmentID string,
	internal map[string][]*exproto.Experiment,
) (*ftproto.HoldoutGroups, error) {
	if groups, err := d.cache.Get(environmentID); err == nil && hasEveryGroup(groups, internal) {
		return groups, nil
	}
	groups, err, _ := d.flightgroup.Do(environmentID, func() (interface{}, error) {
		// Captured before the read so this write cannot replace a newer list.
		generation := d.now().UnixNano()
		groups, err := v2hs.MeasuredGroups(ctx, d.storage, environmentID, d.now())
		if err != nil {
			return nil, err
		}
		if _, err := d.cache.PutIfNewer(groups, environmentID, generation); err != nil {
			d.logger.Warn("Failed to cache the measured holdout groups",
				zap.Error(err),
				zap.String("environmentId", environmentID),
			)
		}
		return groups, nil
	})
	if err != nil {
		return nil, err
	}
	return groups.(*ftproto.HoldoutGroups), nil
}

func hasEveryGroup(groups *ftproto.HoldoutGroups, internal map[string][]*exproto.Experiment) bool {
	found := 0
	for _, group := range groups.GetHoldoutGroups() {
		if _, ok := internal[group.Id]; ok {
			found++
		}
	}
	return found >= len(internal)
}

// match returns what evt counts for: the measured group that targeted the
// flag at the time of the event, before the stop of the group and inside the
// window of one of its internal experiments, and the variation of the user.
// It returns nil when the event counts for no group. It reads nothing: the
// variation is decided here, so an event of a user in neither variation
// never reaches Redis.
func (d *holdoutExposureDeriver) match(
	evt *eventproto.EvaluationEvent,
	measured []*measuredHoldoutGroup,
) *holdoutExposureMatch {
	if len(measured) == 0 {
		return nil
	}
	userID := evt.UserId
	if userID == "" {
		userID = evt.GetUser().GetId()
	}
	if userID == "" {
		return nil
	}
	for _, m := range measured {
		// The experiments list and the measured groups are cached apart, so
		// a group that has stopped can meet an internal experiment that
		// still has its window before the stop. The row of an event after
		// the stop would be outside the window the calculator reads and
		// would take the key of the user's row of that interval.
		if evt.Timestamp >= holdoutdomain.InternalExperimentStopAt(m.group) {
			continue
		}
		intervals, ok := holdoutExposureIntervals(m.experiments, evt.Timestamp)
		if !ok || !isHoldoutTargetAt(m.group, evt.FeatureId, evt.Timestamp) {
			continue
		}
		variationID := holdoutExposureVariation(m.group, evt, userID)
		if variationID == "" {
			continue
		}
		return &holdoutExposureMatch{
			groupID:     m.group.Id,
			variationID: variationID,
			intervals:   intervals,
			userID:      userID,
		}
	}
	return nil
}

// acquire returns the row to write for a match, or nil when the row of the
// user, group and intervals is already written. base is the row of the
// event as it is written for a flag experiment; the derived row copies the
// user, the reason, the tag, the source and the timestamp from it.
//
// seen holds the keys this batch has already asked for, so the events of one
// user in one batch cost one Redis call.
func (d *holdoutExposureDeriver) acquire(
	ctx context.Context,
	match *holdoutExposureMatch,
	eventID string,
	base *epproto.EvaluationEvent,
	seen map[string]struct{},
) (*holdoutExposure, error) {
	key := match.lockKey(base.EnvironmentId)
	if _, ok := seen[key]; ok {
		return nil, nil
	}
	acquired, value, err := d.lock.Lock(ctx, key)
	if err != nil {
		return nil, err
	}
	seen[key] = struct{}{}
	if !acquired {
		return nil, nil
	}
	return &holdoutExposure{
		row: &epproto.EvaluationEvent{
			Id:             holdoutExposureRowID(eventID, match.groupID),
			FeatureId:      match.groupID,
			FeatureVersion: holdoutdomain.InternalExperimentFeatureVersion,
			UserData:       base.UserData,
			UserId:         match.userID,
			VariationId:    match.variationID,
			Reason:         base.Reason,
			Tag:            base.Tag,
			SourceId:       base.SourceId,
			EnvironmentId:  base.EnvironmentId,
			Timestamp:      base.Timestamp,
		},
		eventID:   eventID,
		lockKey:   key,
		lockValue: value,
	}, nil
}

// release deletes the key of a row that was not written, so that the next
// delivery of the event, or the next event of the user, writes the row.
func (d *holdoutExposureDeriver) release(ctx context.Context, exposure *holdoutExposure) {
	if _, err := d.lock.Unlock(ctx, exposure.lockKey, exposure.lockValue); err != nil {
		d.logger.Error("Failed to release the key of a holdout exposure that was not written",
			zap.Error(err),
			zap.String("eventId", exposure.eventID),
			zap.String("environmentId", exposure.row.EnvironmentId),
			zap.String("holdoutGroupId", exposure.row.FeatureId),
		)
	}
}

// holdoutExposureIntervals reports whether timestamp is inside the window
// [start_at, stop_at) of an internal experiment of the group. For every such
// experiment, oldest first, it names the interval of that experiment the
// timestamp is in, as "<start_at>.<n>" joined by "-".
//
// The calculator counts the users of an experiment from its start_at up to
// each cutoff of that experiment, ends included: start_at plus n days, for
// n = 1, 2 and so on. Interval n is the time after start_at plus n days, up
// to and including start_at plus n+1 days, and interval 0 has start_at too,
// so the events of one interval are on the same side of every cutoff. One
// row per interval is then enough: whichever event of the interval is
// processed first, the calculator counts the user at the cutoffs it would
// count the user at with a row per event.
//
// The intervals of the two experiments of a group line up only when the
// analysis period starts a whole number of days after the first experiment,
// so a user can get two rows within a day once the analysis period started.
// The first event of a user after that start is in an interval of the
// analysis-period experiment that no earlier event is in, so that experiment
// counts a user whose earlier events of the day were before it.
func holdoutExposureIntervals(experiments []*exproto.Experiment, timestamp int64) (string, bool) {
	var startAts []int64
	for _, e := range experiments {
		if timestamp < e.StartAt || timestamp >= e.StopAt {
			continue
		}
		if !slices.Contains(startAts, e.StartAt) {
			startAts = append(startAts, e.StartAt)
		}
	}
	if len(startAts) == 0 {
		return "", false
	}
	slices.Sort(startAts)
	intervals := make([]string, 0, len(startAts))
	for _, startAt := range startAts {
		n := max(timestamp-startAt-1, 0) / secondsPerDay
		intervals = append(intervals, fmt.Sprintf("%d.%d", startAt, n))
	}
	return strings.Join(intervals, "-"), true
}

// lacksInternalExperiment reports whether a measured group names an internal
// experiment that is not among the experiments the group was matched with.
// The holdout service writes a group and its internal experiments together,
// so the experiments list is then older than the group.
func lacksInternalExperiment(measured []*measuredHoldoutGroup) bool {
	for _, m := range measured {
		for _, id := range []string{m.group.ExperimentId, m.group.AnalysisExperimentId} {
			if id == "" {
				continue
			}
			listed := slices.ContainsFunc(m.experiments, func(e *exproto.Experiment) bool {
				return e.Id == id
			})
			if !listed {
				return true
			}
		}
	}
	return false
}

// isHoldoutTargetAt reports whether the flag was a target of the group at
// timestamp: added at or before it, and not removed, or removed after it.
func isHoldoutTargetAt(group *ftproto.HoldoutGroup, featureID string, timestamp int64) bool {
	for _, f := range group.Features {
		if f.FeatureId != featureID || f.AddedAt > timestamp {
			continue
		}
		if f.RemovedAt == 0 || timestamp < f.RemovedAt {
			return true
		}
	}
	return false
}

// holdoutExposureVariation returns the variation of the internal experiments
// the event counts for, or an empty string when it counts for neither.
//
//   - in_holdout: the user was served the holdout of this group, which the
//     event says with the HOLDOUT reason and the ID of the group.
//   - not_in_holdout: the user was served normally and, for a HASH group, is
//     in the comparison range. A SEGMENT group has no hash, so every user
//     served normally counts.
func holdoutExposureVariation(
	group *ftproto.HoldoutGroup,
	evt *eventproto.EvaluationEvent,
	userID string,
) string {
	if evt.GetReason().GetType() == ftproto.Reason_HOLDOUT {
		if evt.GetReason().GetHoldoutGroupId() == group.Id {
			return holdoutdomain.InHoldoutVariationID
		}
		return ""
	}
	switch group.AssignmentType {
	case ftproto.HoldoutGroup_HASH:
		if evaluation.IsInHoldoutComparisonRange(group, userID) {
			return holdoutdomain.NotInHoldoutVariationID
		}
	case ftproto.HoldoutGroup_SEGMENT:
		return holdoutdomain.NotInHoldoutVariationID
	}
	return ""
}

// holdoutExposureRowID is the ID of the row derived from an evaluation
// event. It differs from the ID of the event, which is the ID of the row
// written for a flag experiment, and it is the same on every delivery of the
// event.
func holdoutExposureRowID(eventID, groupID string) string {
	return fmt.Sprintf("%s:%s", eventID, groupID)
}
