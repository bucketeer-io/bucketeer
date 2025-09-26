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

package domain

import (
	"sort"
	"time"

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/anypb"

	err "github.com/bucketeer-io/bucketeer/v2/pkg/error"

	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

var (
	errClauseNotFound   = err.NewErrorNotFound(err.AutoopsPackageName, "clause not found", "clause")
	errClauseEmpty      = err.NewErrorInvalidArgEmpty(err.AutoopsPackageName, "clause cannot be empty", "clause")
	errClauseIDRequired = err.NewErrorInvalidArgEmpty(err.AutoopsPackageName, "clause id is required", "clause_id")

	OpsEventRateClause = &proto.OpsEventRateClause{}
	DatetimeClause     = &proto.DatetimeClause{}
)

type AutoOpsRule struct {
	*proto.AutoOpsRule
}

func NewAutoOpsRule(
	featureID string,
	opsType proto.OpsType,
	opsEventRateClauses []*proto.OpsEventRateClause,
	datetimeClauses []*proto.DatetimeClause,
) (*AutoOpsRule, error) {
	now := time.Now().Unix()
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	autoOpsRule := &AutoOpsRule{&proto.AutoOpsRule{
		Id:        id.String(),
		FeatureId: featureID,
		OpsType:   opsType,
		Clauses:   []*proto.Clause{},
		CreatedAt: now,
		UpdatedAt: now,
	}}
	switch opsType {
	case proto.OpsType_EVENT_RATE:
		for _, c := range opsEventRateClauses {
			if _, err := autoOpsRule.AddOpsEventRateClause(c); err != nil {
				return nil, err
			}
		}
	case proto.OpsType_SCHEDULE:
		for _, c := range datetimeClauses {
			if _, err := autoOpsRule.AddDatetimeClause(c); err != nil {
				return nil, err
			}
		}
	}
	if len(autoOpsRule.Clauses) == 0 {
		return nil, errClauseEmpty
	}
	return autoOpsRule, nil
}

func (a *AutoOpsRule) Update(
	autoOpsStatus *proto.AutoOpsStatus,
	opsEventRateClauses []*proto.OpsEventRateClauseChange,
	datetimeClauses []*proto.DatetimeClauseChange,
) (*AutoOpsRule, error) {
	updated := &AutoOpsRule{}
	if err := copier.Copy(updated, a); err != nil {
		return nil, err
	}

	if autoOpsStatus != nil {
		updated.AutoOpsStatus = *autoOpsStatus
	}

	if err := updated.changeClauses(
		opsEventRateClauses,
		datetimeClauses,
	); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	updated.UpdatedAt = now
	return updated, nil
}

func (a *AutoOpsRule) changeClauses(
	opsEventRateClauses []*proto.OpsEventRateClauseChange,
	datetimeClauses []*proto.DatetimeClauseChange,
) error {
	// We split variation changes into two separate steps to handle dependencies correctly:
	//
	// Step 1: apply clauses creation and update
	// - This ensures newly created or updated clauses exist when validating and applying other changes
	// - Without this step, validations referencing newly created clauses would fail.
	//
	// Step 2: apply clauses deletion
	// - This ensures that we can validate and apply deletions without affecting the existing clauses.
	// - If deletions were processed earlier, we could end up with invalid references to non-existent clauses.

	// Step 1: apply clauses creation and update
	var (
		opsEventRateCreateAndUpdate, opsEventRateDelete []*proto.OpsEventRateClauseChange
		datetimeCreateAndUpdate, datetimeDelete         []*proto.DatetimeClauseChange
	)
	for _, c := range opsEventRateClauses {
		if c.ChangeType == proto.ChangeType_DELETE {
			opsEventRateDelete = append(opsEventRateDelete, c)
		} else {
			opsEventRateCreateAndUpdate = append(opsEventRateCreateAndUpdate, c)
		}
	}
	for _, c := range datetimeClauses {
		if c.ChangeType == proto.ChangeType_DELETE {
			datetimeDelete = append(datetimeDelete, c)
		} else {
			datetimeCreateAndUpdate = append(datetimeCreateAndUpdate, c)
		}
	}
	if err := a.validateGranularChanges(
		opsEventRateCreateAndUpdate,
		datetimeCreateAndUpdate,
	); err != nil {
		return err
	}
	if err := a.applyGranularChanges(
		opsEventRateCreateAndUpdate,
		datetimeCreateAndUpdate,
	); err != nil {
		return err
	}

	// Step 2: apply clauses deletion
	if err := a.validateGranularChanges(
		opsEventRateDelete,
		datetimeDelete,
	); err != nil {
		return err
	}
	if err := a.applyGranularChanges(
		opsEventRateDelete,
		datetimeDelete,
	); err != nil {
		return err
	}
	if len(a.Clauses) == 0 {
		return errClauseEmpty
	}
	return nil
}

func (a *AutoOpsRule) applyGranularChanges(
	opsEventRateClauses []*proto.OpsEventRateClauseChange,
	datetimeClauses []*proto.DatetimeClauseChange,
) error {
	for _, c := range opsEventRateClauses {
		switch c.ChangeType {
		case proto.ChangeType_CREATE:
			ac, err := anypb.New(c.Clause)
			if err != nil {
				return err
			}
			_, err = a.addClause(ac, c.Clause.ActionType)
			if err != nil {
				return err
			}
		case proto.ChangeType_UPDATE:
			err := a.changeClause(c.Id, c.Clause, c.Clause.ActionType)
			if err != nil {
				return err
			}
		case proto.ChangeType_DELETE:
			if err := a.DeleteClause(c.Id); err != nil {
				return err
			}
		}
	}

	for _, c := range datetimeClauses {
		switch c.ChangeType {
		case proto.ChangeType_CREATE:
			ac, err := anypb.New(c.Clause)
			if err != nil {
				return err
			}
			_, err = a.addClause(ac, c.Clause.ActionType)
			if err != nil {
				return err
			}
		case proto.ChangeType_UPDATE:
			err := a.changeClause(c.Id, c.Clause, c.Clause.ActionType)
			if err != nil {
				return err
			}
		case proto.ChangeType_DELETE:
			if err := a.DeleteClause(c.Id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AutoOpsRule) validateGranularChanges(
	opsEventRateClauseChanges []*proto.OpsEventRateClauseChange,
	datetimeClauseChanges []*proto.DatetimeClauseChange,
) error {
	for _, c := range opsEventRateClauseChanges {
		switch c.ChangeType {
		case proto.ChangeType_CREATE:
			if c.Clause == nil {
				return errClauseEmpty
			}
		case proto.ChangeType_UPDATE:
			if c.Id == "" {
				return errClauseIDRequired
			}
			if c.Clause == nil {
				return errClauseEmpty
			}
		case proto.ChangeType_DELETE:
			if c.Id == "" {
				return errClauseIDRequired
			}
		}
	}
	for _, c := range datetimeClauseChanges {
		switch c.ChangeType {
		case proto.ChangeType_CREATE:
			if c.Clause == nil {
				return errClauseEmpty
			}
		case proto.ChangeType_UPDATE:
			if c.Id == "" {
				return errClauseIDRequired
			}
			if c.Clause == nil {
				return errClauseEmpty
			}
		case proto.ChangeType_DELETE:
			if c.Id == "" {
				return errClauseIDRequired
			}
		}
	}
	return nil
}

func (a *AutoOpsRule) SetStopped() {
	a.SetAutoOpsStatus(proto.AutoOpsStatus_STOPPED)
}

func (a *AutoOpsRule) SetDeleted() {
	a.Deleted = true
	a.UpdatedAt = time.Now().Unix()
}

func (a *AutoOpsRule) SetFinished() {
	a.SetAutoOpsStatus(proto.AutoOpsStatus_FINISHED)
}

func (a *AutoOpsRule) IsFinished() bool {
	return a.AutoOpsStatus == proto.AutoOpsStatus_FINISHED
}

func (a *AutoOpsRule) IsStopped() bool {
	return a.AutoOpsStatus == proto.AutoOpsStatus_STOPPED
}

func (a *AutoOpsRule) SetAutoOpsStatus(status proto.AutoOpsStatus) {
	now := time.Now().Unix()
	a.AutoOpsStatus = status
	a.UpdatedAt = now
}

func (a *AutoOpsRule) AddOpsEventRateClause(oerc *proto.OpsEventRateClause) (*proto.Clause, error) {
	ac, err := ptypes.MarshalAny(oerc)
	if err != nil {
		return nil, err
	}
	clause, err := a.addClause(ac, oerc.ActionType)
	if err != nil {
		return nil, err
	}
	a.UpdatedAt = time.Now().Unix()
	return clause, nil
}

func (a *AutoOpsRule) AddDatetimeClause(dc *proto.DatetimeClause) (*proto.Clause, error) {
	ac, err := ptypes.MarshalAny(dc)
	if err != nil {
		return nil, err
	}
	clause, err := a.addClause(ac, dc.ActionType)
	a.sortDatetimeClause()

	if err != nil {
		return nil, err
	}
	a.UpdatedAt = time.Now().Unix()
	return clause, nil
}

func (a *AutoOpsRule) addClause(ac *any.Any, actionType proto.ActionType) (*proto.Clause, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	clause := &proto.Clause{
		Id:         id.String(),
		Clause:     ac,
		ActionType: actionType,
	}
	a.Clauses = append(a.Clauses, clause)
	return clause, nil
}

func (a *AutoOpsRule) sortDatetimeClause() {
	type SortClause struct {
		clause         *proto.Clause
		dataTimeClause *proto.DatetimeClause
	}
	newClauses := []*proto.Clause{}
	sortClauses := []*SortClause{}
	for _, c := range a.Clauses {
		datetimeClause, _ := a.unmarshalDatetimeClause(c)
		if datetimeClause == nil {
			newClauses = append(newClauses, c)
			continue
		}
		s := &SortClause{
			clause:         c,
			dataTimeClause: datetimeClause,
		}
		sortClauses = append(sortClauses, s)
	}

	sort.Slice(sortClauses, func(i, j int) bool {
		return sortClauses[i].dataTimeClause.Time < sortClauses[j].dataTimeClause.Time
	})

	for _, c := range sortClauses {
		newClauses = append(newClauses, c.clause)
	}
	a.Clauses = newClauses
}

func (a *AutoOpsRule) ChangeOpsEventRateClause(id string, oerc *proto.OpsEventRateClause) error {
	err := a.changeClause(id, oerc, oerc.ActionType)
	if err != nil {
		return err
	}
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AutoOpsRule) ChangeDatetimeClause(id string, dc *proto.DatetimeClause) error {
	err := a.changeClause(id, dc, dc.ActionType)
	a.sortDatetimeClause()
	if err != nil {
		return err
	}
	a.UpdatedAt = time.Now().Unix()
	return nil
}

func (a *AutoOpsRule) changeClause(id string, mc pb.Message, actionType proto.ActionType) error {
	for _, c := range a.Clauses {
		if c.Id == id {
			clause, err := ptypes.MarshalAny(mc)
			if err != nil {
				return err
			}
			c.Clause = clause
			c.ActionType = actionType
			return nil
		}
	}
	return errClauseNotFound
}

func (a *AutoOpsRule) DeleteClause(id string) error {
	if len(a.Clauses) <= 1 {
		return errClauseEmpty
	}
	a.UpdatedAt = time.Now().Unix()
	var clauses []*proto.Clause
	for i, c := range a.Clauses {
		if c.Id == id {
			clauses = append(a.Clauses[:i], a.Clauses[i+1:]...)
			continue
		}
	}
	if len(clauses) > 0 {
		a.Clauses = clauses
		return nil
	}
	return errClauseNotFound
}

func (a *AutoOpsRule) HasEventRateOps() (bool, error) {
	clauses, err := a.ExtractOpsEventRateClauses()
	if err != nil {
		return false, err
	}
	return len(clauses) > 0, nil
}

func (a *AutoOpsRule) ExtractOpsEventRateClauses() (map[string]*proto.OpsEventRateClause, error) {
	opsEventRateClauses := map[string]*proto.OpsEventRateClause{}
	for _, c := range a.Clauses {
		opsEventRateClause, err := a.unmarshalOpsEventRateClause(c)
		if err != nil {
			return nil, err
		}
		if opsEventRateClause == nil {
			continue
		}
		opsEventRateClauses[c.Id] = opsEventRateClause
	}
	return opsEventRateClauses, nil
}

func (a *AutoOpsRule) unmarshalOpsEventRateClause(clause *proto.Clause) (*proto.OpsEventRateClause, error) {
	if ptypes.Is(clause.Clause, OpsEventRateClause) {
		c := &proto.OpsEventRateClause{}
		if err := ptypes.UnmarshalAny(clause.Clause, c); err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, nil
}

func (a *AutoOpsRule) HasScheduleOps() (bool, error) {
	clauses, err := a.ExtractDatetimeClauses()
	if err != nil {
		return false, err
	}
	return len(clauses) > 0, nil
}

func (a *AutoOpsRule) ExtractDatetimeClauses() (map[string]*proto.DatetimeClause, error) {
	datetimeClauses := map[string]*proto.DatetimeClause{}
	for _, c := range a.Clauses {
		datetimeClause, err := a.unmarshalDatetimeClause(c)
		if err != nil {
			return nil, err
		}
		if datetimeClause == nil {
			continue
		}
		datetimeClauses[c.Id] = datetimeClause
	}
	return datetimeClauses, nil
}

func (a *AutoOpsRule) unmarshalDatetimeClause(clause *proto.Clause) (*proto.DatetimeClause, error) {
	if ptypes.Is(clause.Clause, DatetimeClause) {
		c := &proto.DatetimeClause{}
		if err := ptypes.UnmarshalAny(clause.Clause, c); err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, nil
}
