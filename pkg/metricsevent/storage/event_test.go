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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/log"
	metricsmock "github.com/bucketeer-io/bucketeer/pkg/metrics/mock"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return().Times(1)
	logger, err := log.NewLogger()
	require.NoError(t, err)
	s := NewStorage(logger, registerer)
	assert.IsType(t, &storage{}, s)
}
