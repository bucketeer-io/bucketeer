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

package rpc

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"

	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gatewayproto "github.com/bucketeer-io/bucketeer/proto/gateway"
)

func TestJSONPMarshaller(t *testing.T) {
	var marshalingTests = []struct {
		desc string
		pb   proto.Message
		json string
	}{
		{
			"default values",
			&gatewayproto.GetEvaluationsResponse{State: featureproto.UserEvaluations_QUEUED},
			`{"state":"QUEUED","evaluations":null,"userEvaluationsId":""}`,
		},
		{
			"non-default values",
			&gatewayproto.GetEvaluationsResponse{State: featureproto.UserEvaluations_FULL},
			`{"state":"FULL","evaluations":null,"userEvaluationsId":""}`,
		},
	}
	for _, tt := range marshalingTests {
		json, err := json.Marshal(makeMarshallable(tt.pb))
		if err != nil {
			t.Errorf("%s: marshaling error: %v", tt.desc, err)
		} else if tt.json != string(json) {
			t.Errorf("%s: got [%v] want [%v]", tt.desc, string(json), tt.json)
		}
	}
}
