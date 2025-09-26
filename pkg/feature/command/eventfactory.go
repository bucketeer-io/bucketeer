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

package command

import (
	"github.com/golang/protobuf/proto" // nolint:staticcheck

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type FeatureEventFactory struct {
	editor          *eventproto.Editor
	feature         *domain.Feature
	previousFeature *domain.Feature
	environmentId   string
	comment         string
}

func (s *FeatureEventFactory) CreateEvent(
	eventType eventproto.Event_Type,
	event proto.Message,
) (*eventproto.Event, error) {
	var prev *featureproto.Feature
	if s.previousFeature != nil && s.previousFeature.Feature != nil {
		prev = s.previousFeature.Feature
	}
	return domainevent.NewEvent(
		s.editor,
		eventproto.Event_FEATURE,
		s.feature.Id,
		eventType,
		event,
		s.environmentId,
		s.feature.Feature,
		prev,
		domainevent.WithComment(s.comment),
		domainevent.WithNewVersion(s.feature.Version),
	)
}
