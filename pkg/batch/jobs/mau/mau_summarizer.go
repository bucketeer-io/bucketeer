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

package mau

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

type mauSummarizer struct {
	mysqlClient        mysql.Client
	eventCounterClient ecclient.Client
	location           *time.Location
	opts               *jobs.Options
	logger             *zap.Logger
}

func NewMAUSummarizer(
	mysqlClient mysql.Client,
	eventCounterClient ecclient.Client,
	location *time.Location,
	opts ...jobs.Option) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 5 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &mauSummarizer{
		mysqlClient:        mysqlClient,
		eventCounterClient: eventCounterClient,
		location:           location,
		opts:               dopts,
		logger:             dopts.Logger.Named("mau-count-watcher"),
	}
}

func (s *mauSummarizer) Run(ctx context.Context) error {
	s.logger.Info("MAUSummarizer start running")
	ctx, cancel := context.WithTimeout(ctx, s.opts.Timeout)
	defer cancel()
	now := time.Now().In(s.location)
	yesterday := s.getYesterday(now)
	_, err := s.eventCounterClient.SummarizeMAUCounts(
		ctx,
		&eventcounter.SummarizeMAUCountsRequest{
			YearMonth:  s.newYearMonth(int32(yesterday.Year()), int32(yesterday.Month())),
			IsFinished: s.isEndDateOfMonth(yesterday),
		},
	)
	if err != nil {
		s.logger.Error("Failed to summarize MAU counts",
			zap.Error(err),
			zap.Int32("year", int32(yesterday.Year())),
			zap.Int32("month", int32(yesterday.Month())),
		)
		return err
	}
	s.logger.Info("MAUSummarizer start stopping")
	return nil
}

func (s *mauSummarizer) getYesterday(now time.Time) time.Time {
	return now.AddDate(0, 0, -1)
}

func (s *mauSummarizer) newYearMonth(year, month int32) string {
	return fmt.Sprintf("%d%02d", year, month)
}

func (s *mauSummarizer) isEndDateOfMonth(target time.Time) bool {
	next := target.AddDate(0, 0, 1)
	return next.Day() == 1
}
