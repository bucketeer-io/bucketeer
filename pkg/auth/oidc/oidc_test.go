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

package oidc

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestIsBadRequestError(t *testing.T) {
	testcases := []struct {
		err          error
		isBadRequest bool
	}{
		{
			err:          errors.New("test-error"),
			isBadRequest: false,
		},
		{
			err:          &oauth2.RetrieveError{Response: &http.Response{StatusCode: 500}},
			isBadRequest: false,
		},
		{
			err:          &oauth2.RetrieveError{Response: &http.Response{StatusCode: 400}},
			isBadRequest: true,
		},
	}
	for i, tc := range testcases {
		des := fmt.Sprintf("case %d", i)
		result := isBadRequestError(tc.err)
		assert.Equal(t, tc.isBadRequest, result, des)
	}
}
