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
	"time"

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	proto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

var (
	errClauseNotFound = errors.New("autoOpsRule: clause not found")
	errClauseEmpty    = errors.New("autoOpsRule: clause cannot be empty")

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
	for _, c := range opsEventRateClauses {
		if _, err := autoOpsRule.AddOpsEventRateClause(c); err != nil {
			return nil, err
		}
	}
	for _, c := range datetimeClauses {
		if _, err := autoOpsRule.AddDatetimeClause(c); err != nil {
			return nil, err
		}
	}
	if len(autoOpsRule.Clauses) == 0 {
		return nil, errClauseEmpty
	}
	return autoOpsRule, nil

}

func (a *AutoOpsRule) SetDeleted() {
	a.AutoOpsRule.Deleted = true
	a.AutoOpsRule.UpdatedAt = time.Now().Unix()
}

func (a *AutoOpsRule) SetTriggeredAt() {
	now := time.Now().Unix()
	a.AutoOpsRule.TriggeredAt = now
	a.AutoOpsRule.UpdatedAt = now
}

func (a *AutoOpsRule) AlreadyTriggered() bool {
	return a.TriggeredAt > 0
}

func (a *AutoOpsRule) SetOpsType(opsType proto.OpsType) {
	a.AutoOpsRule.OpsType = opsType
	a.AutoOpsRule.UpdatedAt = time.Now().Unix()
	a.AutoOpsRule.TriggeredAt = 0
}

func (a *AutoOpsRule) AddOpsEventRateClause(oerc *proto.OpsEventRateClause) (*proto.Clause, error) {
	ac, err := ptypes.MarshalAny(oerc)
	if err != nil {
		return nil, err
	}
	return a.addClause(ac)
}

func (a *AutoOpsRule) AddDatetimeClause(dc *proto.DatetimeClause) (*proto.Clause, error) {
	ac, err := ptypes.MarshalAny(dc)
	if err != nil {
		return nil, err
	}
	return a.addClause(ac)
}

func (a *AutoOpsRule) addClause(ac *any.Any) (*proto.Clause, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	clause := &proto.Clause{
		Id:     id.String(),
		Clause: ac,
	}
	a.AutoOpsRule.Clauses = append(a.AutoOpsRule.Clauses, clause)
	a.AutoOpsRule.UpdatedAt = time.Now().Unix()
	a.AutoOpsRule.TriggeredAt = 0
	return clause, nil
}

func (a *AutoOpsRule) ChangeOpsEventRateClause(id string, oerc *proto.OpsEventRateClause) error {
	return a.changeClause(id, oerc)
}

func (a *AutoOpsRule) ChangeDatetimeClause(id string, dc *proto.DatetimeClause) error {
	return a.changeClause(id, dc)
}

func (a *AutoOpsRule) changeClause(id string, mc pb.Message) error {
	a.AutoOpsRule.UpdatedAt = time.Now().Unix()
	a.AutoOpsRule.TriggeredAt = 0
	for _, c := range a.Clauses {
		if c.Id == id {
			clause, err := ptypes.MarshalAny(mc)
			if err != nil {
				return err
			}
			c.Clause = clause
			return nil
		}
	}
	return errClauseNotFound
}

func (a *AutoOpsRule) DeleteClause(id string) error {
	if len(a.Clauses) <= 1 {
		return errClauseEmpty
	}
	a.AutoOpsRule.UpdatedAt = time.Now().Unix()
	a.AutoOpsRule.TriggeredAt = 0
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

func (a *AutoOpsRule) ExtractDatetimeClauses() ([]*proto.DatetimeClause, error) {
	datetimeClauses := []*proto.DatetimeClause{}
	for _, c := range a.Clauses {
		datetimeClause, err := a.unmarshalDatetimeClause(c)
		if err != nil {
			return nil, err
		}
		if datetimeClause == nil {
			continue
		}
		datetimeClauses = append(datetimeClauses, datetimeClause)
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
