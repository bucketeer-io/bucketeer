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

package command

import (
	"context"

	pb "github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/jinzhu/copier"

	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type projectCommandHandler struct {
	editor          *eventproto.Editor
	project         *domain.Project
	previousProject *domain.Project
	publisher       publisher.Publisher
}

func NewProjectCommandHandler(
	editor *eventproto.Editor,
	project *domain.Project,
	p publisher.Publisher,
) (Handler, error) {
	prev := &domain.Project{}
	if err := copier.Copy(prev, project); err != nil {
		return nil, err
	}
	return &projectCommandHandler{
		editor:          editor,
		project:         project,
		previousProject: prev,
		publisher:       p,
	}, nil
}

func (h *projectCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *proto.CreateTrialProjectCommand:
		return h.createTrial(ctx, c)
	default:
		return errUnknownCommand
	}
}

func (h *projectCommandHandler) createTrial(ctx context.Context, cmd *proto.CreateTrialProjectCommand) error {
	return h.send(ctx, eventproto.Event_PROJECT_TRIAL_CREATED, &eventproto.ProjectTrialCreatedEvent{
		Id:           h.project.Id,
		Name:         h.project.Name,
		UrlCode:      h.project.UrlCode,
		Description:  h.project.Description,
		Disabled:     h.project.Disabled,
		Trial:        h.project.Trial,
		CreatorEmail: h.project.CreatorEmail,
		CreatedAt:    h.project.CreatedAt,
		UpdatedAt:    h.project.UpdatedAt,
	})
}

func (h *projectCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	var prev *proto.Project
	if h.previousProject != nil && h.previousProject.Project != nil {
		prev = h.previousProject.Project
	}
	e, err := domainevent.NewAdminEvent(
		h.editor,
		eventproto.Event_PROJECT,
		h.project.Id,
		eventType,
		event,
		h.project.Project,
		prev,
	)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
