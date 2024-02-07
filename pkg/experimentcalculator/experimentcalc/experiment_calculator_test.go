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
//

package experimentcalc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client/mock"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	metricsmock "github.com/bucketeer-io/bucketeer/pkg/metrics/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func creatExperimentCalculator(mockController *gomock.Controller) *ExperimentCalculator {
	registerer := metricsmock.NewMockRegisterer(mockController)
	registerer.EXPECT().MustRegister(gomock.Any()).Return()
	return NewExperimentCalculator(
		stan.NewStan("localhost", "8080"),
		envclient.NewMockClient(mockController),
		ecclient.NewMockClient(mockController),
		experimentclient.NewMockClient(mockController),
		mysqlmock.NewMockClient(mockController),
		registerer,
		jpLocation,
		zap.NewNop(),
	)
}

func TestExperimentCalculatorBinomialModelSample(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	experimentCalculator := creatExperimentCalculator(mockController)
	assert.NotNil(t, experimentCalculator, "ExperimentCalculator should not be nil")
	ctx := context.TODO()
	vrs, err := experimentCalculator.binomialModelSample(
		ctx,
		[]string{"vid1", "vid2"},
		[]int64{38, 51},
		[]int64{101, 99},
		0,
	)

	assert.NoError(t, err, "BinomialModelSample should not be error")

	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Mean, 0.37)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Mean, 0.38)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Sd, 0.045)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Sd, 0.05)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProb.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid1"].CvrProb.Rhat, 1.1)
	assert.Equal(t, len(vrs["vid1"].CvrProb.Histogram.Hist), 100)
	assert.Equal(t, len(vrs["vid1"].CvrProb.Histogram.Bins), 101)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Mean, 0.023)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Mean, 0.026)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid1"].CvrProbBest.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid1"].CvrProbBest.Rhat, 1.1)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Mean, 0.0)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Sd, 0.0)
	assert.Equal(t, vrs["vid1"].CvrProbBeatBaseline.Rhat, 0.0)

	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Mean, 0.49)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Mean, 0.52)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Sd, 0.045)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Sd, 0.05)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProb.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProb.Rhat, 1.1)
	assert.Equal(t, len(vrs["vid2"].CvrProb.Histogram.Hist), 100)
	assert.Equal(t, len(vrs["vid2"].CvrProb.Histogram.Bins), 101)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Mean, 0.97)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Mean, 0.98)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBest.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBest.Rhat, 1.1)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Mean, 0.97)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Mean, 0.98)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Sd, 0.15)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Sd, 0.16)
	assert.GreaterOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Rhat, 0.9)
	assert.LessOrEqual(t, vrs["vid2"].CvrProbBeatBaseline.Rhat, 1.1)

}

func TestListEndAt(t *testing.T) {
	type args struct {
		startAt int64
		endAt   int64
		now     int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "3 days 18 hours",
			args: args{
				startAt: 1614848400, // 2021-03-04 09:00:00Z
				endAt:   1615086000, // 2021-03-07 03:00:00Z
				now:     32508810000,
			},
			want: []int64{1614934800, 1615021200, 1615086000},
		},
		{
			name: "now is earlier than end_at",
			args: args{
				startAt: 1614848400, // 2021-03-04 09:00:00Z
				endAt:   1615086000, // 2021-03-07 03:00:00Z
				now:     1614967200, // 2021-03-06 03:00:00Z
			},
			want: []int64{1614934800, 1614967200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, listEndAt(tt.args.startAt, tt.args.endAt, tt.args.now),
				"listEndAt(%v, %v, %v)", tt.args.startAt, tt.args.endAt, tt.args.now)
		})
	}
}
