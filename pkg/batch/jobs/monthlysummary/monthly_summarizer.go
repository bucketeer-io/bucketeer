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

package monthlysummary

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/v2/pkg/batch/jobs"
	cachev3 "github.com/bucketeer-io/bucketeer/v2/pkg/cache/v3"
	envclient "github.com/bucketeer-io/bucketeer/v2/pkg/environment/client"
	insightsstorage "github.com/bucketeer-io/bucketeer/v2/pkg/insights/storage/v2"
	envproto "github.com/bucketeer-io/bucketeer/v2/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/client"
)

type monthlySummarizer struct {
	envClient             envclient.Client
	mauCache              cachev3.MAUCache
	monthlySummaryStorage insightsstorage.MonthlySummaryStorage
	location              *time.Location
	opts                  *jobs.Options
	logger                *zap.Logger
}

func NewMonthlySummarizer(
	envClient envclient.Client,
	mauCache cachev3.MAUCache,
	monthlySummaryStorage insightsstorage.MonthlySummaryStorage,
	location *time.Location,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 30 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &monthlySummarizer{
		envClient:             envClient,
		mauCache:              mauCache,
		monthlySummaryStorage: monthlySummaryStorage,
		location:              location,
		opts:                  dopts,
		logger:                dopts.Logger.Named("monthly-summarizer"),
	}
}

func (m *monthlySummarizer) Run(ctx context.Context) (lastErr error) {
	startTime := time.Now()
	defer func() {
		jobs.RecordJob(jobs.JobMonthlySummarizer, lastErr, time.Since(startTime))
	}()

	ctx, cancel := context.WithTimeout(ctx, m.opts.Timeout)
	defer cancel()

	m.logger.Info("MonthlySummarizer start running")

	yesterday := time.Now().In(m.location).AddDate(0, 0, -1)
	yearmonth := yesterday.Format("200601")

	envs, err := m.listEnvironments(ctx)
	if err != nil {
		m.logger.Error("Failed to list environments", zap.Error(err))
		return err
	}

	sourceIDs := listSourceIDs()

	records := make([]insightsstorage.MonthlySummaryRecord, 0, len(envs)*len(sourceIDs))

	for _, env := range envs {
		mauCounts, err := m.mauCache.MergeIntoMAUBatch(env.Id, sourceIDs, yesterday)
		if err != nil {
			m.logger.Warn("Failed to merge MAU batch, skipping environment",
				zap.String("environmentId", env.Id),
				zap.Error(err),
			)
			continue
		}
		for _, sourceID := range sourceIDs {
			records = append(records, insightsstorage.MonthlySummaryRecord{
				Yearmonth:     yearmonth,
				EnvironmentID: env.Id,
				SourceID:      sourceID,
				MAU:           mauCounts[sourceID],
				Requests:      0, // TODO: Get from Prometheus
			})
		}
	}

	if len(records) == 0 {
		m.logger.Info("No MAU records to upsert")
		return nil
	}

	if err := m.monthlySummaryStorage.UpsertMonthlySummaryBatch(ctx, records); err != nil {
		m.logger.Error("Failed to upsert MAU batch", zap.Error(err))
		return err
	}

	m.logger.Info("MonthlySummarizer completed successfully",
		zap.Duration("elapsed", time.Since(startTime)),
		zap.Int("recordCount", len(records)),
	)
	return nil
}

func (m *monthlySummarizer) listEnvironments(ctx context.Context) ([]*envproto.EnvironmentV2, error) {
	resp, err := m.envClient.ListEnvironmentsV2(ctx, &envproto.ListEnvironmentsV2Request{
		PageSize: 0,
		Archived: &wrapperspb.BoolValue{Value: false},
	})
	if err != nil {
		return nil, err
	}
	return resp.Environments, nil
}

func listSourceIDs() []string {
	ids := make([]string, 0, len(eventproto.SourceId_value)-1)
	for _, v := range eventproto.SourceId_value {
		id := eventproto.SourceId(v)
		if id == eventproto.SourceId_UNKNOWN {
			continue
		}
		ids = append(ids, id.String())
	}
	return ids
}
