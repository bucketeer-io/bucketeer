// Copyright 2025 The Bucketeer Authors.
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

package deleter

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

type demoOrganizationDeleter struct {
	organizationStorage v2es.OrganizationStorage
	environmentClient   environmentclient.Client
	opts                *jobs.Options
	logger              *zap.Logger
}

func NewDemoOrganizationDeleter(
	mysqlClient mysql.Client,
	environmentClient environmentclient.Client,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}

	return &demoOrganizationDeleter{
		organizationStorage: v2es.NewOrganizationStorage(mysqlClient),
		environmentClient:   environmentClient,
		opts:                dopts,
		logger:              dopts.Logger.Named("demo-organization-deleter"),
	}
}

func (d *demoOrganizationDeleter) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, d.opts.Timeout)
	defer cancel()

	outdatedOrganizations, err := d.getOutdatedOrganizations(ctx)
	if err != nil {
		return err
	}
	if len(outdatedOrganizations) == 0 {
		return nil
	}

	outdatedOrganizationIDs := make([]string, 0, len(outdatedOrganizations))
	for _, org := range outdatedOrganizations {
		outdatedOrganizationIDs = append(outdatedOrganizationIDs, org.Id)
	}

	_, err = d.environmentClient.DeleteBucketeerData(ctx, &envproto.DeleteBucketeerDataRequest{
		DeleteOrganizationIds: outdatedOrganizationIDs,
	})
	if err != nil {
		d.logger.Error("Could not delete bucketeer data", zap.Error(err))
		return err
	}
	return nil
}

func (d *demoOrganizationDeleter) getOutdatedOrganizations(ctx context.Context) ([]*envproto.Organization, error) {
	trialPeriod, err := strconv.Atoi(os.Getenv("DEMO_TRIAL_PERIOD_DAY"))
	if err != nil {
		d.logger.Error("Could not parse DEMO_TRIAL_PERIOD_DAY", zap.Error(err))
		return nil, err
	}
	filters := []*mysql.FilterV2{
		{
			Column:   "organization.created_at",
			Operator: mysql.OperatorLessThan,
			Value:    time.Now().AddDate(0, 0, -trialPeriod).Unix(),
		},
		{
			Column:   "organization.system_admin",
			Operator: mysql.OperatorEqual,
			Value:    false,
		},
	}
	options := &mysql.ListOptions{
		Limit:       10000,
		Offset:      0,
		Filters:     filters,
		InFilters:   nil,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}

	organizations, _, _, err := d.organizationStorage.ListOrganizations(ctx, options)
	if err != nil {
		d.logger.Error("Could not list organizations", zap.Error(err))
		return nil, err
	}
	return organizations, nil
}
