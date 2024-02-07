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

package recorder

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	sqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/proto/event/client"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

type PullerMock struct{}

func (p *PullerMock) Pull(ctx context.Context, f func(context.Context, *puller.Message)) error {
	timer := time.NewTimer(time.Millisecond * 100)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			event := client.EvaluationEvent{
				FeatureId:      "id",
				FeatureVersion: 10,
				Timestamp:      int64(time.Now().Nanosecond()),
			}
			data, _ := proto.Marshal(&event)
			f(ctx, &puller.Message{
				Data: data,
				Ack:  func() {},
			})
		}
	}
}

func TestNewRecorder(t *testing.T) {
	t.Parallel()
	puller := &PullerMock{}
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	assert.IsType(t, &recorder{}, NewRecorder(puller, db))
}

func TestRecorderRun(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	recorder := NewRecorder(&PullerMock{}, db)
	go recorder.Run()
	time.Sleep(time.Second)
	recorder.Stop()
}

func TestUnmarshalMessage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	recorder := NewRecorder(&PullerMock{}, db).(*recorder)
	event := client.Event{Id: "hoge"}
	data, err := proto.Marshal(&event)
	assert.NoError(t, err)
	msg := puller.Message{Data: data}
	e, err := recorder.unmarshalMessage(&msg)
	assert.NoError(t, err)
	assert.Equal(t, event.Id, e.Id)
}

func TestUnmarshalEvent(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	recorder := NewRecorder(&PullerMock{}, db).(*recorder)
	evalEvent := client.EvaluationEvent{FeatureId: "id"}
	any, err := ptypes.MarshalAny(&evalEvent)
	assert.NoError(t, err)
	evale, err := recorder.unmarshalEvent(any)
	assert.NoError(t, err)
	assert.Equal(t, evalEvent.FeatureId, evale.FeatureId)
}

func TestCacheEnvLastUsedInfo(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	recorder := NewRecorder(&PullerMock{}, db, WithFlushInterval(1*time.Second)).(*recorder)
	expectVersion := "1.0.0"
	user := &userproto.User{
		Id:   "id",
		Data: map[string]string{"app_version": expectVersion},
	}
	patterns := []struct {
		input                *client.EvaluationEvent
		environmentNamespace string
		expect               int64
	}{
		{
			input:                &client.EvaluationEvent{FeatureId: "id", FeatureVersion: 10, Timestamp: 0, User: user},
			environmentNamespace: "ns0",
			expect:               1,
		},
		{
			input:                &client.EvaluationEvent{FeatureId: "id", FeatureVersion: 10, Timestamp: 2, User: user},
			environmentNamespace: "ns0",
			expect:               2,
		},
		{
			input:                &client.EvaluationEvent{FeatureId: "id", FeatureVersion: 11, Timestamp: 1, User: user},
			environmentNamespace: "ns0",
			expect:               1,
		},
		{
			input:                &client.EvaluationEvent{FeatureId: "id", FeatureVersion: 10, Timestamp: 1, User: user},
			environmentNamespace: "ns1",
			expect:               1,
		},
	}
	envCache := make(environmentLastUsedInfoCache, 1)
	e := client.EvaluationEvent{FeatureId: "id", FeatureVersion: 10, Timestamp: 1, User: user}
	recorder.cacheEnvLastUsedInfo(&e, envCache, "ns0")
	for i, p := range patterns {
		recorder.cacheEnvLastUsedInfo(p.input, envCache, p.environmentNamespace)
		key := domain.FeatureLastUsedInfoID(p.input.FeatureId, p.input.FeatureVersion)
		assert.Equal(t, p.expect, envCache[p.environmentNamespace][key].LastUsedAt, "i=%d", i)
		assert.Equal(t, expectVersion, envCache[p.environmentNamespace][key].ClientOldestVersion, "i=%d", i)
		assert.Equal(t, expectVersion, envCache[p.environmentNamespace][key].ClientLatestVersion, "i=%d", i)
	}
}

func BenchmarkCacheEnvLastUsedInfo(b *testing.B) {
	mockController := gomock.NewController(b)
	defer mockController.Finish()
	db := sqlmock.NewMockClient(mockController)
	recorder := NewRecorder(&PullerMock{}, db).(*recorder)
	b.ResetTimer()
	envCache := make(environmentLastUsedInfoCache, 1)
	for i := 0; i < b.N; i++ {
		e := client.EvaluationEvent{
			FeatureId:      "id",
			FeatureVersion: 10,
			Timestamp:      int64(time.Now().Nanosecond()),
		}
		recorder.cacheEnvLastUsedInfo(&e, envCache, "ns0")
	}
}
