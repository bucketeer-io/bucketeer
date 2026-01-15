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

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"

	"github.com/bucketeer-io/bucketeer/v2/pkg/account/domain"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

type accountV2CommandHandler struct {
	editor          *eventproto.Editor
	account         *domain.AccountV2
	previousAccount *domain.AccountV2
	publisher       publisher.Publisher
	organizationID  string
}

func NewAccountV2CommandHandler(
	editor *eventproto.Editor,
	account *domain.AccountV2,
	p publisher.Publisher,
	organizationID string,
) (Handler, error) {
	prev := &domain.AccountV2{}
	if err := copier.Copy(prev, account); err != nil {
		return nil, err
	}
	return &accountV2CommandHandler{
		editor:          editor,
		account:         account,
		previousAccount: prev,
		publisher:       p,
		organizationID:  organizationID,
	}, nil
}

func (h *accountV2CommandHandler) Handle(ctx context.Context, cmd Command) error {
	switch c := cmd.(type) {
	case *accountproto.CreateSearchFilterCommand:
		return h.createSearchFilter(ctx, c)
	case *accountproto.ChangeSearchFilterNameCommand:
		return h.changeSearchFilerName(ctx, c)
	case *accountproto.ChangeSearchFilterQueryCommand:
		return h.changeSearchFilterQuery(ctx, c)
	case *accountproto.ChangeDefaultSearchFilterCommand:
		return h.changeDefaultSearchFilter(ctx, c)
	case *accountproto.DeleteSearchFilterCommand:
		return h.deleteSearchFiler(ctx, c)
	default:
		return ErrBadCommand
	}
}

func (h *accountV2CommandHandler) createSearchFilter(
	ctx context.Context,
	cmd *accountproto.CreateSearchFilterCommand) error {
	searchFilter, err := h.account.AddSearchFilter(
		cmd.Name,
		cmd.Query,
		cmd.FilterTargetType,
		cmd.EnvironmentId,
		cmd.DefaultFilter)
	if err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_CREATED_SEARCH_FILTER, &eventproto.SearchFilterCreatedEvent{
		Name:          searchFilter.Name,
		Query:         searchFilter.Query,
		TargetType:    searchFilter.FilterTargetType,
		EnvironmentId: searchFilter.EnvironmentId,
		DefaultFilter: searchFilter.DefaultFilter,
	})
}

func (h *accountV2CommandHandler) changeSearchFilerName(
	ctx context.Context,
	cmd *accountproto.ChangeSearchFilterNameCommand) error {
	if err := h.account.ChangeSearchFilterName(cmd.Id, cmd.Name); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_NANE_CHANGED,
		&eventproto.SearchFilterNameChangedEvent{
			Id:   cmd.Id,
			Name: cmd.Name,
		},
	)
}

func (h *accountV2CommandHandler) changeSearchFilterQuery(
	ctx context.Context,
	cmd *accountproto.ChangeSearchFilterQueryCommand) error {
	if err := h.account.ChangeSearchFilterQuery(cmd.Id, cmd.Query); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_QUERY_CHANGED,
		&eventproto.SearchFilterQueryChangedEvent{
			Id:    cmd.Id,
			Query: cmd.Query,
		},
	)
}

func (h *accountV2CommandHandler) changeDefaultSearchFilter(
	ctx context.Context,
	cmd *accountproto.ChangeDefaultSearchFilterCommand) error {
	if err := h.account.ChangeDefaultSearchFilter(cmd.Id, cmd.DefaultFilter); err != nil {
		return err
	}
	return h.send(
		ctx,
		eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_DEFAULT_CHANGED,
		&eventproto.SearchFilterDefaultChangedEvent{
			Id:            cmd.Id,
			DefaultFilter: cmd.DefaultFilter,
		},
	)
}

func (h *accountV2CommandHandler) deleteSearchFiler(
	ctx context.Context,
	cmd *accountproto.DeleteSearchFilterCommand) error {
	if err := h.account.DeleteSearchFilter(cmd.Id); err != nil {
		return err
	}
	return h.send(ctx, eventproto.Event_ACCOUNT_V2_SEARCH_FILTER_DELETED, &eventproto.SearchFilterDeletedEvent{
		Id: cmd.Id,
	})
}

func (h *accountV2CommandHandler) send(
	ctx context.Context,
	eventType eventproto.Event_Type,
	event proto.Message,
) error {
	var prev *accountproto.AccountV2
	if h.previousAccount != nil && h.previousAccount.AccountV2 != nil {
		prev = h.previousAccount.AccountV2
	}
	e, err := domainevent.NewAdminEvent(
		h.editor,
		eventproto.Event_ACCOUNT,
		h.account.Email,
		eventType,
		event,
		h.account.AccountV2,
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
