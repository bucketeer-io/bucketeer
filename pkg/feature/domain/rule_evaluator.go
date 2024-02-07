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
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

type ruleEvaluator struct {
	clauseEvaluator
}

func (e *ruleEvaluator) Evaluate(
	rules []*featureproto.Rule,
	user *userproto.User,
	segmentUsers []*featureproto.SegmentUser,
) *featureproto.Rule {
	for _, rule := range rules {
		if e.evaluateRule(rule, user, segmentUsers) {
			return rule
		}
	}
	return nil
}

func (e *ruleEvaluator) evaluateRule(
	rule *featureproto.Rule,
	user *userproto.User,
	segmentUsers []*featureproto.SegmentUser,
) bool {
	for _, clause := range rule.Clauses {
		if !e.evaluateClause(clause, user, segmentUsers) {
			return false
		}
	}
	return true
}

func (e *ruleEvaluator) evaluateClause(
	clause *featureproto.Clause,
	user *userproto.User,
	segmentUsers []*featureproto.SegmentUser,
) bool {
	var targetAttr string
	if clause.Attribute == "id" {
		targetAttr = user.Id
	} else {
		targetAttr = user.Data[clause.Attribute]
	}
	return e.clauseEvaluator.Evaluate(targetAttr, clause, user.Id, segmentUsers)
}
