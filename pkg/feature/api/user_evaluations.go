// Copyright 2022 The Bucketeer Authors.
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

package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (s *FeatureService) GetUserEvaluations(
	ctx context.Context,
	req *featureproto.GetUserEvaluationsRequest,
) (*featureproto.GetUserEvaluationsResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetUserEvaluationsRequest(req); err != nil {
		return nil, err
	}
	evaluations, err := s.userEvaluationStorage.GetUserEvaluations(
		ctx,
		req.UserId,
		req.EnvironmentNamespace,
		req.Tag,
	)
	if err != nil {
		s.logger.Error(
			"Failed to get user evaluations",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("userId", req.UserId),
				zap.String("tag", req.Tag),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.GetUserEvaluationsResponse{
		Evaluations: evaluations,
	}, nil
}

func (s *FeatureService) UpsertUserEvaluation(
	ctx context.Context,
	req *featureproto.UpsertUserEvaluationRequest,
) (*featureproto.UpsertUserEvaluationResponse, error) {
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateUpsertUserEvaluationRequest(req); err != nil {
		return nil, err
	}
	if err := s.userEvaluationStorage.UpsertUserEvaluation(
		ctx,
		req.Evaluation,
		req.EnvironmentNamespace,
		req.Tag,
	); err != nil {
		s.logger.Error(
			"Failed to upsert user evaluation",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
				zap.String("tag", req.Tag),
				zap.Any("evaluation", req.Evaluation),
			)...,
		)
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.UpsertUserEvaluationResponse{}, nil
}
