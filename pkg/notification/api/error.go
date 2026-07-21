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

package api

import (
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	bkterr "github.com/bucketeer-io/bucketeer/v2/pkg/error"
)

var statusNotImplemented = gstatus.Error(codes.Unimplemented, "notification: not implemented")

var (
	statusUnauthenticated = api.NewGRPCStatus(
		bkterr.NewErrorUnauthenticated(bkterr.NotificationPackageName, "unauthenticated"))
	statusPermissionDenied = api.NewGRPCStatus(
		bkterr.NewErrorPermissionDenied(bkterr.NotificationPackageName, "permission denied"))
	statusLocalizationRequired = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.NotificationPackageName,
			"at least one localization must be specified",
			"Localizations"))
	statusLanguageRequired = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.NotificationPackageName,
			"language must be specified",
			"Language"))
	statusDuplicatedLanguage = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgDuplicated(
			bkterr.NotificationPackageName,
			"language is duplicated",
			"Language"))
	statusTitleRequired = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.NotificationPackageName,
			"title must be specified",
			"Title"))
	statusContentRequired = api.NewGRPCStatus(
		bkterr.NewErrorInvalidArgEmpty(
			bkterr.NotificationPackageName,
			"content must be specified",
			"Content"))
	statusNotificationAlreadyExists = api.NewGRPCStatus(
		bkterr.NewErrorAlreadyExists(bkterr.NotificationPackageName, "already exists"))
)
