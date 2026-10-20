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

package server

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/v2/ui/dashboard"
)

func TestSpaFileSystem_Open_invalidUTF8Path(t *testing.T) {
	t.Parallel()

	sfs := &spaFileSystem{root: http.FS(dashboard.FS)}
	_, err := sfs.Open(".env.local\x85")
	require.ErrorIs(t, err, fs.ErrNotExist)
}

func TestDashboardHandler_invalidUTF8PathReturns404(t *testing.T) {
	t.Parallel()

	handler := dashboardHandler()
	req := httptest.NewRequest(http.MethodGet, "/.env.local%85", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
