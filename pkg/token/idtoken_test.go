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

package token

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func TestExtractUserID(t *testing.T) {
	userID := "test-id"
	sub := &authproto.IDTokenSubject{
		UserId: userID,
		ConnId: "test-connector-id",
	}
	data, err := proto.Marshal(sub)
	require.NoError(t, err)
	encodedSub := base64.RawURLEncoding.EncodeToString(data)
	testcases := []struct {
		subject string
		userID  string
		failed  bool
	}{
		{
			subject: "invalid",
			userID:  "",
			failed:  true,
		},
		{
			subject: encodedSub,
			userID:  userID,
			failed:  false,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("index %d", i)
		userID, err := ExtractUserID(tc.subject)
		assert.Equal(t, tc.userID, userID, des)
		assert.Equal(t, tc.failed, err != nil, des)
	}
}
