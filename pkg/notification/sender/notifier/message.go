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

package notifier

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
)

type msgType int

const (
	_ msgType = iota
	msgTypeFeatureStale
	msgTypeExperimentResult
	msgTypeMAUCount
)

var (
	errUnknownMsgType   = errors.New("notification: unknown message type")
	msgFeatureStaleJaJP = &errdetails.LocalizedMessage{
		Locale:  locale.Ja,
		Message: "%d日以上使用されていないフィーチャーフラグがあります。",
	}
	msgExperimentResultJaJP = &errdetails.LocalizedMessage{
		Locale:  locale.Ja,
		Message: "実行中のエクスペリメントがあります。",
	}
	msgMAUCountJaJP = &errdetails.LocalizedMessage{
		Locale:  locale.Ja,
		Message: "%d月のMAUです。",
	}
)

func localizedMessage(t msgType, loc string) (*errdetails.LocalizedMessage, error) {
	// handle loc if multi-lang is necessary
	switch t {
	case msgTypeFeatureStale:
		return msgFeatureStaleJaJP, nil
	case msgTypeExperimentResult:
		return msgExperimentResultJaJP, nil
	case msgTypeMAUCount:
		return msgMAUCountJaJP, nil
	default:
		return nil, errUnknownMsgType
	}
}
