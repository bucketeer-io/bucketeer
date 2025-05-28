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

package demo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	v2es "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

var (
	targetEntities = []string{
		"subscription",
		"experiment_result",
		"push",
		"ops_count",
		"auto_ops_rule",
		"segment_user",
		"segment",
		"goal",
		"experiment",
		"tag",
		"ops_progressive_rollout",
		"flag_trigger",
		"code_reference",
		"feature",
		"api_key",
		"audit_log",
	}
	targetEntitiesInOrganization = []string{
		"account_v2",
	}
)

type demoDataDeleter struct {
	mysqlClient         mysql.Client
	organizationStorage v2es.OrganizationStorage
	environmentStorage  v2es.EnvironmentStorage
	opts                *jobs.Options
	logger              *zap.Logger
}

func NewDemoDataDeleter(
	mysqlClient mysql.Client,
	opts ...jobs.Option,
) jobs.Job {
	dopts := &jobs.Options{
		Timeout: 1 * time.Minute,
		Logger:  zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}

	return &demoDataDeleter{
		mysqlClient:         mysqlClient,
		organizationStorage: v2es.NewOrganizationStorage(mysqlClient),
		environmentStorage:  v2es.NewEnvironmentStorage(mysqlClient),
		opts:                dopts,
		logger:              dopts.Logger.Named("demo-environment-deleter"),
	}
}

func (d *demoDataDeleter) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, d.opts.Timeout)
	defer cancel()
	demoEnabled, err := strconv.ParseBool(os.Getenv("BUCKETEER_BATCH_DEMO_ENABLED"))
	if err != nil {
		d.logger.Error("Could not parse BUCKETEER_BATCH_DEMO_ENABLED", zap.Error(err))
		return err
	}
	if !demoEnabled {
		d.logger.Info("Batch demo data deleter is disabled")
		return nil
	}

	d.logger.Info("Starting to delete old data from demo site")

	outdatedOrganizations, err := d.getOutdatedOrganizations(ctx)
	if err != nil {
		return err
	}
	if len(outdatedOrganizations) == 0 {
		d.logger.Info("No outdated organizations found")
		return nil
	}

	organizationIDs := make([]string, len(outdatedOrganizations))
	for i, org := range outdatedOrganizations {
		organizationIDs[i] = org.Id
	}
	outdatedEnvironments, err := d.getOutdatedEnvironments(ctx, organizationIDs)
	if err != nil {
		return err
	}
	for _, env := range outdatedEnvironments {
		err = d.deleteDataFromEnvironment(ctx, env.Id)
		if err != nil {
			return err
		}
	}

	err = d.deleteEnvironments(ctx, organizationIDs)
	if err != nil {
		return err
	}
	err = d.deleteProjects(ctx, organizationIDs)
	if err != nil {
		return err
	}
	err = d.deleteOrganizations(ctx, organizationIDs)
	if err != nil {
		return err
	}

	d.logger.Info("Finished deleting old data from demo site")
	return nil
}

func (d *demoDataDeleter) getOutdatedOrganizations(ctx context.Context) ([]*envproto.Organization, error) {
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

func (d *demoDataDeleter) getOutdatedEnvironments(
	ctx context.Context,
	organizationIDs []string,
) ([]*envproto.EnvironmentV2, error) {
	inFilters := []*mysql.InFilter{
		{
			Column: "environment_v2.organization_id",
			Values: convToInterfaceSlice(organizationIDs),
		},
	}
	options := &mysql.ListOptions{
		Limit:       10000,
		Offset:      0,
		Filters:     nil,
		InFilters:   inFilters,
		NullFilters: nil,
		JSONFilters: nil,
		SearchQuery: nil,
		Orders:      nil,
	}
	environments, _, _, err := d.environmentStorage.ListEnvironmentsV2(ctx, options)
	if err != nil {
		d.logger.Error("Could not list environments", zap.Error(err))
		return nil, err
	}
	return environments, nil
}

func (d *demoDataDeleter) deleteDataFromEnvironment(
	ctx context.Context,
	environmentID string,
) error {
	args := []interface{}{
		environmentID,
	}

	return d.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		for _, target := range targetEntities {
			query := fmt.Sprintf(`
				DELETE FROM
					%s
				WHERE
					environment_id = ?`,
				target)
			_, err := d.mysqlClient.ExecContext(
				ctxWithTx,
				query,
				args...,
			)
			if err != nil {
				d.logger.Error("Failed to delete data from environment",
					zap.String("table", target),
					zap.String("environmentId", environmentID),
					zap.Error(err),
				)
				return err
			}
		}
		return nil
	})
}

func (d *demoDataDeleter) deleteEnvironments(ctx context.Context, organizationIDs []string) error {
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("organization_id", convToInterfaceSlice(organizationIDs)),
	}
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf("DELETE FROM environment_v2 %s", whereSQL)
	_, err := d.mysqlClient.ExecContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		d.logger.Error("Failed to delete environments", zap.Error(err))
		return err
	}
	return nil
}

func (d *demoDataDeleter) deleteProjects(
	ctx context.Context,
	organizationIDs []string,
) error {
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("organization_id", convToInterfaceSlice(organizationIDs)),
	}
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	query := fmt.Sprintf("DELETE FROM project %s", whereSQL)
	_, err := d.mysqlClient.ExecContext(
		ctx,
		query,
		whereArgs...,
	)
	if err != nil {
		d.logger.Error("Failed to delete projects", zap.Error(err))
		return err
	}
	return nil
}

func (d *demoDataDeleter) deleteOrganizations(ctx context.Context, organizationIDs []string) error {
	whereParts := []mysql.WherePart{
		mysql.NewInFilter("organization_id", convToInterfaceSlice(organizationIDs)),
	}
	whereSQL, whereArgs := mysql.ConstructWhereSQLString(whereParts)
	return d.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, _ mysql.Transaction) error {
		for _, target := range targetEntitiesInOrganization {
			query := fmt.Sprintf("DELETE FROM %s %s", target, whereSQL)
			_, err := d.mysqlClient.ExecContext(
				ctxWithTx,
				query,
				whereArgs...,
			)
			if err != nil {
				d.logger.Error("Failed to delete organization entity",
					zap.Error(err),
					zap.String("table", target),
				)
				return err
			}
		}
		whereParts = []mysql.WherePart{
			mysql.NewInFilter("id", convToInterfaceSlice(organizationIDs)),
		}
		whereSQL, whereArgs = mysql.ConstructWhereSQLString(whereParts)
		query := fmt.Sprintf("DELETE FROM organization %s", whereSQL)
		_, err := d.mysqlClient.ExecContext(
			ctxWithTx,
			query,
			whereArgs...,
		)
		if err != nil {
			d.logger.Error("Failed to delete organizations", zap.Error(err))
			return err
		}
		return nil
	})
}

func convToInterfaceSlice(
	slice []string,
) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, element := range slice {
		result = append(result, element)
	}
	return result
}
