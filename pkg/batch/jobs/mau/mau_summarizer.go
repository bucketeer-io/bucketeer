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
