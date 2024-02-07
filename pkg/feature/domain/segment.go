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
	"time"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

type Segment struct {
	*featureproto.Segment
}

func NewSegment(name string, description string) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Segment{
		Segment: &featureproto.Segment{
			Id:          id.String(),
			Name:        name,
			Description: description,
			Version:     1,
			CreatedAt:   time.Now().Unix(),
		},
	}, nil
}

func (s *Segment) SetDeleted() error {
	s.Segment.Deleted = true
	s.Segment.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) ChangeName(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) ChangeDescription(description string) error {
	s.Description = description
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) AddRule(rule *featureproto.Rule) error {
	if _, err := s.findRuleIndex(rule.Id); err == nil {
		return errRuleAlreadyExists
	}
	s.Rules = append(s.Rules, rule)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) DeleteRule(rule string) error {
	idx, err := s.findRuleIndex(rule)
	if err != nil {
		return err
	}
	s.Rules = append(s.Rules[:idx], s.Rules[idx+1:]...)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) AddClause(ruleID string, clause *featureproto.Clause) error {
	idx, err := s.findRuleIndex(ruleID)
	if err != nil {
		return err
	}
	rule := s.Rules[idx]
	if _, err := s.findClauseIndex(clause.Id, rule.Clauses); err == nil {
		return errClauseAlreadyExists
	}
	rule.Clauses = append(rule.Clauses, clause)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) DeleteClause(ruleID string, clauseID string) error {
	ruleIdx, err := s.findRuleIndex(ruleID)
	if err != nil {
		return err
	}
	rule := s.Rules[ruleIdx]
	if len(rule.Clauses) <= 1 {
		return errRuleMustHaveAtLeastOneClause
	}
	clauseIdx, err := s.findClauseIndex(clauseID, rule.Clauses)
	if err != nil {
		return err
	}
	rule.Clauses = append(rule.Clauses[:clauseIdx], rule.Clauses[clauseIdx+1:]...)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) ChangeClauseAttribute(ruleID string, clauseID string, attribute string) error {
	clause, err := s.findClause(ruleID, clauseID)
	if err != nil {
		return err
	}
	clause.Attribute = attribute
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) ChangeClauseOperator(ruleID string, clauseID string, operator featureproto.Clause_Operator) error {
	clause, err := s.findClause(ruleID, clauseID)
	if err != nil {
		return err
	}
	clause.Operator = operator
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) AddClauseValue(ruleID string, clauseID string, value string) error {
	clause, err := s.findClause(ruleID, clauseID)
	if err != nil {
		return err
	}
	clause.Values = append(clause.Values, value)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) RemoveClauseValue(ruleID string, clauseID string, value string) error {
	clause, err := s.findClause(ruleID, clauseID)
	if err != nil {
		return err
	}
	if len(clause.Values) <= 1 {
		return errClauseMustHaveAtLeastOneValue
	}
	idx, err := index(value, clause.Values)
	if err != nil {
		return err
	}
	clause.Values = append(clause.Values[:idx], clause.Values[idx+1:]...)
	s.UpdatedAt = time.Now().Unix()
	return nil
}

func (s *Segment) AddIncludedUserCount(count int64) {
	s.IncludedUserCount += count
	s.UpdatedAt = time.Now().Unix()
}

func (s *Segment) RemoveIncludedUserCount(count int64) {
	s.IncludedUserCount -= count
	s.UpdatedAt = time.Now().Unix()
}

func (s *Segment) SetIncludedUserCount(count int64) {
	s.IncludedUserCount = count
	s.UpdatedAt = time.Now().Unix()
}

func (s *Segment) findRuleIndex(id string) (int, error) {
	for i, r := range s.Rules {
		if r.Id == id {
			return i, nil
		}
	}
	return -1, errRuleNotFound
}

func (s *Segment) findClauseIndex(clauseID string, clauses []*featureproto.Clause) (int, error) {
	for i, c := range clauses {
		if c.Id == clauseID {
			return i, nil
		}
	}
	return -1, errClauseNotFound
}

func (s *Segment) findClause(ruleID string, clauseID string) (*featureproto.Clause, error) {
	ruleIdx, err := s.findRuleIndex(ruleID)
	if err != nil {
		return nil, err
	}
	rule := s.Rules[ruleIdx]
	clauseIdx, err := s.findClauseIndex(clauseID, rule.Clauses)
	if err != nil {
		return nil, errClauseNotFound
	}
	return rule.Clauses[clauseIdx], nil
}

func (s *Segment) SetStatus(status featureproto.Segment_Status) {
	s.Status = status
	s.UpdatedAt = time.Now().Unix()
}
