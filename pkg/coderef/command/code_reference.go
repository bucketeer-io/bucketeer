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

package command

import (
	"context"

	pb "github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"

	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	coderefproto "github.com/bucketeer-io/bucketeer/proto/coderef"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

type codeReferenceCommandHandler struct {
	editor                *eventproto.Editor
	codeReference         *domain.CodeReference
	previousCodeReference *domain.CodeReference
	publisher             publisher.Publisher
	environmentID         string
}

func NewCodeReferenceCommandHandler(
	editor *eventproto.Editor,
	codeReference *domain.CodeReference,
	p publisher.Publisher,
	environmentID string,
) (Handler, error) {
	prev := &domain.CodeReference{}
	if err := copier.Copy(prev, codeReference); err != nil {
		return nil, err
	}
	return &codeReferenceCommandHandler{
		editor:                editor,
		codeReference:         codeReference,
		previousCodeReference: prev,
		publisher:             p,
		environmentID:         environmentID,
	}, nil
}

func (h *codeReferenceCommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *coderefproto.CreateCodeReferenceCommand:
		return h.create(ctx, c)
	case *coderefproto.UpdateCodeReferenceCommand:
		return h.update(ctx, c)
	case *coderefproto.DeleteCodeReferenceCommand:
		return h.delete(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *codeReferenceCommandHandler) create(ctx context.Context, cmd *coderefproto.CreateCodeReferenceCommand) error {
	return h.send(ctx, eventproto.Event_CODE_REFERENCE_CREATED, &eventproto.CodeReferenceCreatedEvent{
		Id:               h.codeReference.Id,
		FeatureId:        h.codeReference.FeatureId,
		FilePath:         h.codeReference.FilePath,
		LineNumber:       h.codeReference.LineNumber,
		CodeSnippet:      h.codeReference.CodeSnippet,
		ContentHash:      h.codeReference.ContentHash,
		Aliases:          h.codeReference.Aliases,
		RepositoryName:   h.codeReference.RepositoryName,
		RepositoryOwner:  h.codeReference.RepositoryOwner,
		RepositoryType:   h.codeReference.RepositoryType,
		RepositoryBranch: h.codeReference.RepositoryBranch,
		CommitHash:       h.codeReference.CommitHash,
		EnvironmentId:    h.codeReference.EnvironmentId,
		CreatedAt:        h.codeReference.CreatedAt,
		UpdatedAt:        h.codeReference.UpdatedAt,
	})
}

func (h *codeReferenceCommandHandler) update(ctx context.Context, cmd *coderefproto.UpdateCodeReferenceCommand) error {
	h.codeReference.Update(
		cmd.FilePath,
		cmd.LineNumber,
		cmd.CodeSnippet,
		cmd.ContentHash,
		cmd.Aliases,
		cmd.RepositoryBranch,
		cmd.CommitHash,
	)
	return h.send(ctx, eventproto.Event_CODE_REFERENCE_UPDATED, &eventproto.CodeReferenceUpdatedEvent{
		Id:               h.codeReference.Id,
		FilePath:         h.codeReference.FilePath,
		LineNumber:       h.codeReference.LineNumber,
		CodeSnippet:      h.codeReference.CodeSnippet,
		ContentHash:      h.codeReference.ContentHash,
		Aliases:          h.codeReference.Aliases,
		RepositoryBranch: h.codeReference.RepositoryBranch,
		CommitHash:       h.codeReference.CommitHash,
		EnvironmentId:    h.codeReference.EnvironmentId,
		UpdatedAt:        h.codeReference.UpdatedAt,
	})
}

func (h *codeReferenceCommandHandler) delete(ctx context.Context, cmd *coderefproto.DeleteCodeReferenceCommand) error {
	return h.send(ctx, eventproto.Event_CODE_REFERENCE_DELETED, &eventproto.CodeReferenceDeletedEvent{
		Id:            h.codeReference.Id,
		EnvironmentId: h.codeReference.EnvironmentId,
	})
}

func (h *codeReferenceCommandHandler) send(ctx context.Context, eventType eventproto.Event_Type, event pb.Message) error {
	e, err := domainevent.NewEvent(
		h.editor,
		eventproto.Event_CODEREF,
		h.codeReference.Id,
		eventType,
		event,
		h.environmentID,
		h.codeReference,
		h.previousCodeReference,
	)
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
