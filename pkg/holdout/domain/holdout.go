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
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pkgErr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

const (
	maxNameLength = 511
	// Holdout and comparison arcs are adjacent on [0, 1) and must not overlap.
	maxPercentage = 50.0
	// Bucket starts are stored as DECIMAL(6,3) on [0, 100).
	bucketMilliUnits = 100000
)

var (
	ErrNameRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "name is required", "name",
	).WithMessageKey("HoldoutNameRequired")
	ErrNameTooLong = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "name is too long", "name",
	).WithMessageKey("HoldoutNameTooLong")
	ErrPercentageOutOfRange = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "percentage must be in (0, 50]", "percentage",
	).WithMessageKey("HoldoutPercentageOutOfRange")
	ErrSegmentRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "segment_id is required", "segment_id",
	).WithMessageKey("HoldoutSegmentRequired")
	ErrPercentageNotAllowed = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "percentage is only valid for hash assignment", "percentage",
	).WithMessageKey("HoldoutPercentageNotAllowed")
	ErrSegmentNotAllowed = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "segment_id is only valid for segment assignment", "segment_id",
	).WithMessageKey("HoldoutSegmentNotAllowed")
	ErrClientAssignmentRejected = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName,
		"bucket_start and sampling_seed are assigned by the server",
		"bucket_start",
	).WithMessageKey("HoldoutClientAssignmentRejected")
	ErrPeriodInvalid = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "start_at must be before stop_at", "start_at",
	).WithMessageKey("HoldoutPeriodInvalid")
	ErrStopAtBeforeNow = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName, "stop_at must be in the future", "stop_at",
	).WithMessageKey("HoldoutStopAtBeforeNow")
	ErrAnalysisPeriodInvalid = pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName,
		"start_at < start_analysis_at < stop_at is required when an analysis period is scheduled",
		"start_analysis_at",
	).WithMessageKey("HoldoutAnalysisPeriodInvalid")
	ErrInvalidStatus = pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName, "holdout group status does not allow this operation",
	).WithMessageKey("HoldoutInvalidStatus")
	ErrBeforeStop = pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName, "stop_at has not passed",
	).WithMessageKey("HoldoutBeforeStop")
	ErrAnalysisFrozen = pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName, "target flags cannot be added during the analysis period",
	).WithMessageKey("HoldoutAnalysisFrozen")
	ErrGoalsCannotBeRemoved = pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName, "goals cannot be removed while the group is serving",
	).WithMessageKey("HoldoutGoalsCannotBeRemoved")
	ErrStartAtImmutable = pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName, "start_at cannot be changed while the group is serving",
	).WithMessageKey("HoldoutStartAtImmutable")
	ErrNothingToUpdate = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "no fields to update", "holdout_group",
	).WithMessageKey("HoldoutNothingToUpdate")
	ErrFeatureRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "at least one feature is required", "feature_id",
	).WithMessageKey("HoldoutFeatureRequired")
	ErrFeatureIDRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "feature_id is required", "feature_id",
	).WithMessageKey("HoldoutFeatureIDRequired")
	ErrVariationRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "holdout_variation_id is required", "holdout_variation_id",
	).WithMessageKey("HoldoutVariationRequired")
	ErrDuplicateFeature = pkgErr.NewErrorInvalidArgDuplicated(
		pkgErr.HoldoutPackageName, "feature_id is duplicated", "feature_id",
	).WithMessageKey("HoldoutDuplicateFeature")
	ErrGoalIDRequired = pkgErr.NewErrorInvalidArgEmpty(
		pkgErr.HoldoutPackageName, "goal_id is required", "goal_id",
	).WithMessageKey("HoldoutGoalIDRequired")

	randomBucketStart = defaultRandomBucketStart
	newSamplingSeed   = defaultNewSamplingSeed
)

// GroupSnapshot is the slice of a group needed for period and serving guards.
type GroupSnapshot struct {
	ID        string
	Name      string
	Status    featureproto.HoldoutGroup_Status
	StartAt   int64
	StopAt    int64
	SegmentID string
}

// FeatureTarget is a current target row of an active group.
type FeatureTarget struct {
	GroupID     string
	GroupName   string
	FeatureID   string
	VariationID string
	Status      featureproto.HoldoutGroup_Status
}

type HoldoutGroup struct {
	*featureproto.HoldoutGroup
}

type HoldoutGroupFeature struct {
	*featureproto.HoldoutGroupFeature
}

type CreateParams struct {
	EnvironmentID   string
	Name            string
	Description     string
	AssignmentType  featureproto.HoldoutGroup_AssignmentType
	Percentage      float64
	SegmentID       string
	StartAt         int64
	StopAt          int64
	StartAnalysisAt int64
	GoalIDs         []string
	CreatedBy       string
	Now             int64
	// ExistingSeed is the environment's canonical HASH seed, empty when this is the first group.
	ExistingSeed string
}

type UpdateParams struct {
	Name            *wrapperspb.StringValue
	Description     *wrapperspb.StringValue
	StartAt         *wrapperspb.Int64Value
	StopAt          *wrapperspb.Int64Value
	StartAnalysisAt *wrapperspb.Int64Value
	// GoalIDs is nil when the caller did not send the field.
	GoalIDsSet bool
	GoalIDs    []string
	UpdatedBy  string
	Now        int64
}

func IsServing(status featureproto.HoldoutGroup_Status) bool {
	return status == featureproto.HoldoutGroup_RUNNING || status == featureproto.HoldoutGroup_ANALYSIS
}

func IsActive(status featureproto.HoldoutGroup_Status) bool {
	return status == featureproto.HoldoutGroup_WAITING || IsServing(status)
}

// PeriodsOverlap reports whether half-open windows [aStart, aStop) and [bStart, bStop) intersect.
func PeriodsOverlap(aStart, aStop, bStart, bStop int64) bool {
	return aStart < bStop && bStart < aStop
}

// RejectClientAssignment rejects a create request that supplies bucket_start or sampling_seed.
func RejectClientAssignment(bucketStart *wrapperspb.DoubleValue, samplingSeed *wrapperspb.StringValue) error {
	if bucketStart != nil || samplingSeed != nil {
		return ErrClientAssignmentRejected
	}
	return nil
}

func NormalizePercentage(percentage float64) (float64, error) {
	rounded := math.Round(percentage*1000) / 1000
	if rounded <= 0 || rounded > maxPercentage {
		return 0, ErrPercentageOutOfRange
	}
	return rounded, nil
}

func ValidateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return ErrNameTooLong
	}
	return nil
}

func ValidatePeriod(startAt, stopAt, startAnalysisAt, now int64, requireFutureStop bool) error {
	if startAt >= stopAt {
		return ErrPeriodInvalid
	}
	if requireFutureStop && stopAt <= now {
		return ErrStopAtBeforeNow
	}
	if startAnalysisAt != 0 && !(startAt < startAnalysisAt && startAnalysisAt < stopAt) {
		return ErrAnalysisPeriodInvalid
	}
	return nil
}

func DedupeStrings(ids []string) ([]string, error) {
	if len(ids) == 0 {
		return []string{}, nil
	}
	out := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if id == "" {
			return nil, ErrGoalIDRequired
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

// NewHoldoutGroup validates a create and assigns the server-owned hash fields.
func NewHoldoutGroup(p CreateParams) (*HoldoutGroup, error) {
	if err := ValidateName(p.Name); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(p.StartAt, p.StopAt, p.StartAnalysisAt, p.Now, true); err != nil {
		return nil, err
	}
	goalIDs, err := DedupeStrings(p.GoalIDs)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	group := &HoldoutGroup{
		HoldoutGroup: &featureproto.HoldoutGroup{
			Id:                   id.String(),
			EnvironmentId:        p.EnvironmentID,
			Name:                 p.Name,
			Description:          p.Description,
			AssignmentType:       p.AssignmentType,
			Status:               featureproto.HoldoutGroup_WAITING,
			StartAt:              p.StartAt,
			StopAt:               p.StopAt,
			StartAnalysisAt:      p.StartAnalysisAt,
			GoalIds:              goalIDs,
			ExperimentId:         "",
			AnalysisExperimentId: "",
			CreatedBy:            p.CreatedBy,
			UpdatedBy:            p.CreatedBy,
			CreatedAt:            p.Now,
			UpdatedAt:            p.Now,
		},
	}
	switch p.AssignmentType {
	case featureproto.HoldoutGroup_HASH:
		if p.SegmentID != "" {
			return nil, ErrSegmentNotAllowed
		}
		percentage, err := NormalizePercentage(p.Percentage)
		if err != nil {
			return nil, err
		}
		seed, bucket, err := AssignHash(p.ExistingSeed)
		if err != nil {
			return nil, err
		}
		group.Percentage = percentage
		group.SamplingSeed = seed
		group.BucketStart = bucket
	case featureproto.HoldoutGroup_SEGMENT:
		if p.Percentage != 0 {
			return nil, ErrPercentageNotAllowed
		}
		if p.SegmentID == "" {
			return nil, ErrSegmentRequired
		}
		group.SegmentId = p.SegmentID
	default:
		return nil, pkgErr.NewErrorInvalidArgUnknown(
			pkgErr.HoldoutPackageName, "assignment_type is invalid", "assignment_type",
		).WithMessageKey("HoldoutAssignmentTypeInvalid")
	}
	return group, nil
}

// AssignHash copies the environment seed or issues one, and draws bucket_start.
func AssignHash(existingSeed string) (seed string, bucketStart float64, err error) {
	if existingSeed != "" {
		seed = existingSeed
	} else {
		seed, err = newSamplingSeed()
		if err != nil {
			return "", 0, err
		}
	}
	bucketStart, err = randomBucketStart()
	if err != nil {
		return "", 0, err
	}
	return seed, bucketStart, nil
}

func defaultRandomBucketStart() (float64, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint64(buf[:]) % bucketMilliUnits
	return float64(n) / 1000, nil
}

func defaultNewSamplingSeed() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// OverlappingGroup returns another active group whose period intersects [start, stop).
func OverlappingGroup(start, stop int64, selfID string, others []GroupSnapshot) (GroupSnapshot, bool) {
	for _, other := range others {
		if other.ID == selfID || !IsActive(other.Status) {
			continue
		}
		if PeriodsOverlap(start, stop, other.StartAt, other.StopAt) {
			return other, true
		}
	}
	return GroupSnapshot{}, false
}

// OtherServingGroup returns another group that is already RUNNING or ANALYSIS.
func OtherServingGroup(selfID string, others []GroupSnapshot) (GroupSnapshot, bool) {
	for _, other := range others {
		if other.ID == selfID {
			continue
		}
		if IsServing(other.Status) {
			return other, true
		}
	}
	return GroupSnapshot{}, false
}

// ServingGroupForSegment returns a serving group that references segmentID.
func ServingGroupForSegment(segmentID string, others []GroupSnapshot) (GroupSnapshot, bool) {
	for _, other := range others {
		if other.SegmentID == segmentID && IsServing(other.Status) {
			return other, true
		}
	}
	return GroupSnapshot{}, false
}

// ConflictingFeatureTarget returns a current target of another active group.
func ConflictingFeatureTarget(featureID, selfGroupID string, targets []FeatureTarget) (FeatureTarget, bool) {
	for _, target := range targets {
		if target.FeatureID != featureID || target.GroupID == selfGroupID {
			continue
		}
		if IsActive(target.Status) {
			return target, true
		}
	}
	return FeatureTarget{}, false
}

func NewPeriodOverlapError(other GroupSnapshot) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf("period overlaps holdout group %s (%s)", other.Name, other.ID),
	).WithMessageKey("HoldoutPeriodOverlap")
}

func NewServingGroupExistsError(other GroupSnapshot) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf("holdout group %s (%s) is already serving", other.Name, other.ID),
	).WithMessageKey("HoldoutServingGroupExists")
}

func NewFeatureInUseError(featureID string, target FeatureTarget) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf(
			"feature %s is a target of holdout group %s (%s)",
			featureID, target.GroupName, target.GroupID,
		),
	).WithMessageKey("HoldoutFeatureInUse")
}

func NewVariationInUseError(variationID string, target FeatureTarget) error {
	return pkgErr.NewErrorVariationInUse(
		pkgErr.HoldoutPackageName,
		pkgErr.ErrorTypeVariationInUseByHoldout,
		fmt.Sprintf(
			"variation %s is the holdout variation of group %s (%s)",
			variationID, target.GroupName, target.GroupID,
		),
		map[string]string{
			"variationId":      variationID,
			"holdoutGroupId":   target.GroupID,
			"holdoutGroupName": target.GroupName,
		},
	)
}

func NewFeatureTargetedError(featureID string, target FeatureTarget) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf(
			"feature %s is a target of holdout group %s (%s)",
			featureID, target.GroupName, target.GroupID,
		),
	).WithMessageKey("HoldoutFeatureTargeted")
}

func NewSegmentLockedError(segmentID string, group GroupSnapshot) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf(
			"segment %s is referenced by serving holdout group %s (%s)",
			segmentID, group.Name, group.ID,
		),
	).WithMessageKey("HoldoutSegmentLocked")
}

func (g *HoldoutGroup) Clone() *HoldoutGroup {
	cloned, ok := proto.Clone(g.HoldoutGroup).(*featureproto.HoldoutGroup)
	if !ok {
		return &HoldoutGroup{HoldoutGroup: &featureproto.HoldoutGroup{}}
	}
	return &HoldoutGroup{HoldoutGroup: cloned}
}

func (g *HoldoutGroup) Snapshot() GroupSnapshot {
	return GroupSnapshot{
		ID:        g.Id,
		Name:      g.Name,
		Status:    g.Status,
		StartAt:   g.StartAt,
		StopAt:    g.StopAt,
		SegmentID: g.SegmentId,
	}
}

func (g *HoldoutGroup) Update(p UpdateParams) error {
	if g.Deleted || g.Archived {
		return ErrInvalidStatus
	}
	if !p.dirty() {
		return ErrNothingToUpdate
	}
	serving := IsServing(g.Status)
	stopped := g.Status == featureproto.HoldoutGroup_STOPPED || g.Status == featureproto.HoldoutGroup_FORCE_STOPPED
	if p.Name != nil {
		if err := ValidateName(p.Name.Value); err != nil {
			return err
		}
		g.Name = p.Name.Value
	}
	if p.Description != nil {
		g.Description = p.Description.Value
	}
	if p.StartAt != nil {
		if serving || stopped {
			return ErrStartAtImmutable
		}
		g.StartAt = p.StartAt.Value
	}
	if p.StopAt != nil {
		if stopped {
			return ErrInvalidStatus
		}
		g.StopAt = p.StopAt.Value
	}
	if p.StartAnalysisAt != nil {
		if stopped {
			return ErrInvalidStatus
		}
		if g.Status == featureproto.HoldoutGroup_ANALYSIS && p.StartAnalysisAt.Value == 0 {
			return ErrAnalysisPeriodInvalid
		}
		g.StartAnalysisAt = p.StartAnalysisAt.Value
	}
	if p.GoalIDsSet {
		next, err := DedupeStrings(p.GoalIDs)
		if err != nil {
			return err
		}
		if serving && !isSuperset(next, g.GoalIds) {
			return ErrGoalsCannotBeRemoved
		}
		if stopped {
			return ErrInvalidStatus
		}
		g.GoalIds = next
	}
	requireFutureStop := g.Status == featureproto.HoldoutGroup_WAITING || serving
	if err := ValidatePeriod(g.StartAt, g.StopAt, g.StartAnalysisAt, p.Now, requireFutureStop); err != nil {
		return err
	}
	g.UpdatedBy = p.UpdatedBy
	g.UpdatedAt = p.Now
	return nil
}

func (p UpdateParams) dirty() bool {
	return p.Name != nil || p.Description != nil || p.StartAt != nil || p.StopAt != nil ||
		p.StartAnalysisAt != nil || p.GoalIDsSet
}

func isSuperset(next, current []string) bool {
	have := make(map[string]struct{}, len(next))
	for _, id := range next {
		have[id] = struct{}{}
	}
	for _, id := range current {
		if _, ok := have[id]; !ok {
			return false
		}
	}
	return true
}

// Start moves WAITING to RUNNING. A future start_at is clamped to now.
func (g *HoldoutGroup) Start(now int64) error {
	if g.Deleted || g.Archived || g.Status != featureproto.HoldoutGroup_WAITING {
		return ErrInvalidStatus
	}
	if g.StartAt > now {
		g.StartAt = now
	}
	if err := ValidatePeriod(g.StartAt, g.StopAt, g.StartAnalysisAt, now, true); err != nil {
		return err
	}
	g.Status = featureproto.HoldoutGroup_RUNNING
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) StartAnalysis(now int64) error {
	if g.Deleted || g.Archived || g.Status != featureproto.HoldoutGroup_RUNNING {
		return ErrInvalidStatus
	}
	g.Status = featureproto.HoldoutGroup_ANALYSIS
	g.AnalysisStartedAt = now
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) StopManual(now int64) error {
	if g.Deleted || g.Archived || !IsServing(g.Status) {
		return ErrInvalidStatus
	}
	g.Status = featureproto.HoldoutGroup_FORCE_STOPPED
	g.StoppedAt = now
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) FinishScheduled(now int64) error {
	if g.Deleted || g.Archived || !IsServing(g.Status) {
		return ErrInvalidStatus
	}
	if g.StopAt > now {
		return ErrBeforeStop
	}
	g.Status = featureproto.HoldoutGroup_STOPPED
	g.StoppedAt = now
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) Archive(now int64) error {
	if g.Deleted || g.Archived {
		return ErrInvalidStatus
	}
	if g.Status != featureproto.HoldoutGroup_STOPPED && g.Status != featureproto.HoldoutGroup_FORCE_STOPPED {
		return ErrInvalidStatus
	}
	g.Archived = true
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) Delete(now int64) error {
	if g.Deleted || g.Archived || g.Status != featureproto.HoldoutGroup_WAITING {
		return ErrInvalidStatus
	}
	g.Deleted = true
	g.UpdatedAt = now
	return nil
}

func (g *HoldoutGroup) CanAddFeatures() error {
	if g.Deleted || g.Archived {
		return ErrInvalidStatus
	}
	if g.Status == featureproto.HoldoutGroup_ANALYSIS {
		return ErrAnalysisFrozen
	}
	if g.Status != featureproto.HoldoutGroup_WAITING && g.Status != featureproto.HoldoutGroup_RUNNING {
		return ErrInvalidStatus
	}
	return nil
}

func (g *HoldoutGroup) CanRemoveFeatures() error {
	if g.Deleted || g.Archived {
		return ErrInvalidStatus
	}
	if g.Status != featureproto.HoldoutGroup_WAITING && !IsServing(g.Status) {
		return ErrInvalidStatus
	}
	return nil
}

func (g *HoldoutGroup) Touch(updatedBy string, now int64) {
	g.UpdatedBy = updatedBy
	g.UpdatedAt = now
}

func ValidateFeatureInputs(inputs []*featureproto.HoldoutGroupFeatureInput) error {
	if len(inputs) == 0 {
		return ErrFeatureRequired
	}
	seen := make(map[string]struct{}, len(inputs))
	for _, input := range inputs {
		if input.GetFeatureId() == "" {
			return ErrFeatureIDRequired
		}
		if input.GetHoldoutVariationId() == "" {
			return ErrVariationRequired
		}
		if _, ok := seen[input.FeatureId]; ok {
			return ErrDuplicateFeature
		}
		seen[input.FeatureId] = struct{}{}
	}
	return nil
}

func NewHoldoutGroupFeature(
	groupID, environmentID, featureID, variationID string,
	addedAt int64,
) (*HoldoutGroupFeature, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &HoldoutGroupFeature{
		HoldoutGroupFeature: &featureproto.HoldoutGroupFeature{
			Id:                 id.String(),
			HoldoutGroupId:     groupID,
			FeatureId:          featureID,
			EnvironmentId:      environmentID,
			HoldoutVariationId: variationID,
			AddedAt:            addedAt,
			RemovedAt:          0,
		},
	}, nil
}

func VariationBelongsToFeature(variationID string, variations []*featureproto.Variation) bool {
	for _, variation := range variations {
		if variation.GetId() == variationID {
			return true
		}
	}
	return false
}

func NewFeatureUnavailableError(featureID string) error {
	return pkgErr.NewErrorFailedPrecondition(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf("feature %s is archived or deleted", featureID),
	).WithMessageKey("HoldoutFeatureUnavailable")
}

func NewVariationNotOnFeatureError(featureID, variationID string) error {
	return pkgErr.NewErrorInvalidArgNotMatchFormat(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf("variation %s is not on feature %s", variationID, featureID),
		"holdout_variation_id",
	).WithMessageKey("HoldoutVariationNotOnFeature")
}

func NewFeatureAlreadyAddedError(featureID string) error {
	return pkgErr.NewErrorAlreadyExists(
		pkgErr.HoldoutPackageName,
		fmt.Sprintf("feature %s is already a target of this holdout group", featureID),
	).WithMessageKey("HoldoutFeatureAlreadyAdded")
}

// DueForScheduledStart reports a WAITING group whose start_at has passed.
func DueForScheduledStart(g *featureproto.HoldoutGroup, now time.Time) bool {
	return g.GetStatus() == featureproto.HoldoutGroup_WAITING && g.GetStartAt() <= now.Unix()
}

// DueForScheduledAnalysis reports a RUNNING group whose analysis time has passed.
func DueForScheduledAnalysis(g *featureproto.HoldoutGroup, now time.Time) bool {
	return g.GetStatus() == featureproto.HoldoutGroup_RUNNING &&
		g.GetStartAnalysisAt() > 0 &&
		g.GetStartAnalysisAt() <= now.Unix() &&
		g.GetStopAt() > now.Unix()
}

// DueForScheduledStop reports a serving group whose stop_at has passed.
func DueForScheduledStop(g *featureproto.HoldoutGroup, now time.Time) bool {
	return IsServing(g.GetStatus()) && g.GetStopAt() <= now.Unix()
}
