// Copyright 2022 The Bucketeer Authors.
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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package targetstore

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	autoopsservice "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentservice "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

const (
	listRequestSize = 500
)

type EnvironmentLister interface {
	GetEnvironments(ctx context.Context) []*environmentdomain.Environment
}

type AutoOpsRuleLister interface {
	GetAutoOpsRules(ctx context.Context, environmentNamespace string) []*autoopsdomain.AutoOpsRule
}

type ProgressiveRolloutLister interface {
	GetProgressiveRollouts(ctx context.Context, environmentNamespace string) []*autoopsdomain.ProgressiveRollout
}

type TargetStore interface {
	Run()
	Stop()
	EnvironmentLister
	AutoOpsRuleLister
	ProgressiveRolloutLister
}

type options struct {
	refreshInterval time.Duration
	metrics         metrics.Registerer
	logger          *zap.Logger
}

type Option func(*options)

func WithRefreshInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.refreshInterval = interval
	}
}

func WithMetrics(r metrics.Registerer) Option {
	return func(opts *options) {
		opts.metrics = r
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type targetStore struct {
	timeNow                func() time.Time
	environmentClient      environmentservice.Client
	autoOpsClient          autoopsservice.Client
	autoOpsRules           map[string][]*autoopsdomain.AutoOpsRule
	progressiveRollouts    map[string][]*autoopsdomain.ProgressiveRollout
	autoOpsRulesMtx        sync.Mutex
	progressiveRolloutsMtx sync.Mutex
	environments           atomic.Value
	opts                   *options
	logger                 *zap.Logger
	ctx                    context.Context
	cancel                 func()
	doneCh                 chan struct{}
}

func NewTargetStore(
	environmentClient environmentservice.Client,
	autoOpsClient autoopsservice.Client,
	opts ...Option,
) TargetStore {
	dopts := &options{
		refreshInterval: 2 * time.Minute,
		logger:          zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	store := &targetStore{
		timeNow:           time.Now,
		environmentClient: environmentClient,
		autoOpsClient:     autoOpsClient,
		autoOpsRules:      make(map[string][]*autoopsdomain.AutoOpsRule),
		opts:              dopts,
		logger:            dopts.logger.Named("targetstore"),
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
	store.environments.Store(make([]*environmentdomain.Environment, 0))
	return store
}

func (s *targetStore) Run() {
	s.logger.Info("Run started")
	defer close(s.doneCh)
	s.refresh()
	ticker := time.NewTicker(s.opts.refreshInterval)
	defer func() {
		ticker.Stop()
		s.logger.Info("Run finished")
	}()
	for {
		select {
		case <-ticker.C:
			s.refresh()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *targetStore) Stop() {
	s.logger.Info("ops-event: transformer: targetstore: stop started")
	s.cancel()
	<-s.doneCh
	s.logger.Info("ops-event: transformer: targetstore: stop finished")
}

func (s *targetStore) refresh() {
	ctx, cancel := context.WithTimeout(s.ctx, time.Minute)
	defer cancel()
	err := s.refreshEnvironments(ctx)
	if err != nil {
		s.logger.Error("Failed to refresh environments", zap.Error(err))
	}
	err = s.refreshAutoOpsRules(ctx)
	if err != nil {
		s.logger.Error("Failed to refresh auto ops rules", zap.Error(err))
	}
	err = s.refreshProgressiveRollout(ctx)
	if err != nil {
		s.logger.Error("Failed to refresh progressive rollout", zap.Error(err))
	}
}

func (s *targetStore) refreshEnvironments(ctx context.Context) error {
	pbEnvironments, err := s.listEnvironments(ctx)
	if err != nil {
		s.logger.Error("Failed to list environments", zap.Error(err))
		return err
	}
	domainEnvironments := []*environmentdomain.Environment{}
	for _, e := range pbEnvironments {
		domainEnvironments = append(domainEnvironments, &environmentdomain.Environment{Environment: e})
	}
	s.environments.Store(domainEnvironments)
	itemsGauge.WithLabelValues(typeEnvironment).Set(float64(len(domainEnvironments)))
	return nil
}

func (s *targetStore) listEnvironments(ctx context.Context) ([]*environmentproto.Environment, error) {
	environments := []*environmentproto.Environment{}
	cursor := ""
	for {
		resp, err := s.environmentClient.ListEnvironments(ctx, &environmentproto.ListEnvironmentsRequest{
			PageSize: listRequestSize,
			Cursor:   cursor,
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		environmentSize := len(resp.Environments)
		if environmentSize == 0 || environmentSize < listRequestSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}

func (s *targetStore) refreshAutoOpsRules(ctx context.Context) error {
	autoOpsRulesMap := make(map[string][]*autoopsdomain.AutoOpsRule)
	environments := s.GetEnvironments(ctx)
	for _, e := range environments {
		autoOpsRules, err := s.listTargetAutoOpsRules(ctx, e.Namespace)
		if err != nil {
			s.logger.Error("Failed to list auto ops rules", zap.Error(err), zap.String("environmentNamespace", e.Namespace))
			continue
		}
		s.logger.Debug("Succeeded to list auto ops rules", zap.String("environmentNamespace", e.Namespace))
		autoOpsRulesMap[e.Namespace] = autoOpsRules
	}
	s.autoOpsRulesMtx.Lock()
	s.autoOpsRules = autoOpsRulesMap
	s.autoOpsRulesMtx.Unlock()
	itemsGauge.WithLabelValues(typeAutoOpsRule).Set(float64(len(autoOpsRulesMap)))
	return nil
}

func (s *targetStore) listTargetAutoOpsRules(
	ctx context.Context,
	environmentNamespace string,
) ([]*autoopsdomain.AutoOpsRule, error) {
	pbAutoOpsRules, err := s.listAutoOpsRules(ctx, environmentNamespace)
	if err != nil {
		return nil, err
	}
	targetAutoOpsRules := []*autoopsdomain.AutoOpsRule{}
	for _, a := range pbAutoOpsRules {
		da := &autoopsdomain.AutoOpsRule{AutoOpsRule: a}
		if da.AlreadyTriggered() {
			continue
		}
		targetAutoOpsRules = append(targetAutoOpsRules, da)
	}
	return targetAutoOpsRules, nil
}

func (s *targetStore) listAutoOpsRules(
	ctx context.Context,
	environmentNamespace string,
) ([]*autoopsproto.AutoOpsRule, error) {
	autoOpsRules := []*autoopsproto.AutoOpsRule{}
	cursor := ""
	for {
		resp, err := s.autoOpsClient.ListAutoOpsRules(ctx, &autoopsproto.ListAutoOpsRulesRequest{
			EnvironmentNamespace: environmentNamespace,
			PageSize:             listRequestSize,
			Cursor:               cursor,
		})
		if err != nil {
			return nil, err
		}
		autoOpsRules = append(autoOpsRules, resp.AutoOpsRules...)
		autoOpsRulesSize := len(resp.AutoOpsRules)
		if autoOpsRulesSize == 0 || autoOpsRulesSize < listRequestSize {
			return autoOpsRules, nil
		}
		cursor = resp.Cursor
	}
}

func (s *targetStore) refreshProgressiveRollout(
	ctx context.Context,
) error {
	progressiveRolloutMap := make(map[string][]*autoopsdomain.ProgressiveRollout)
	environments := s.GetEnvironments(ctx)
	for _, e := range environments {
		progressiveRollouts, err := s.listTargetProgressiveRollouts(ctx, e.Namespace)
		if err != nil {
			s.logger.Error(
				"Failed to list progressive rollouts",
				zap.Error(err),
				zap.String("environmentNamespace", e.Namespace),
			)
			continue
		}
		s.logger.Debug("Succeeded to list progressive rollouts", zap.String("environmentNamespace", e.Namespace))
		progressiveRolloutMap[e.Namespace] = progressiveRollouts
	}
	s.progressiveRolloutsMtx.Lock()
	s.progressiveRollouts = progressiveRolloutMap
	s.progressiveRolloutsMtx.Unlock()
	return nil
}

func (s *targetStore) listTargetProgressiveRollouts(
	ctx context.Context,
	environmentNamespace string,
) ([]*autoopsdomain.ProgressiveRollout, error) {
	pbProgressiveRollouts, err := s.listProgressiveRollouts(ctx, environmentNamespace)
	if err != nil {
		return nil, err
	}
	targetProgressiveRollouts := make([]*autoopsdomain.ProgressiveRollout, 0, len(pbProgressiveRollouts))
	for _, p := range pbProgressiveRollouts {
		dp := &autoopsdomain.ProgressiveRollout{ProgressiveRollout: p}
		if dp.IsFinished() {
			continue
		}
		targetProgressiveRollouts = append(targetProgressiveRollouts, dp)
	}
	return targetProgressiveRollouts, nil
}

func (s *targetStore) listProgressiveRollouts(
	ctx context.Context,
	environmentNamespace string,
) ([]*autoopsproto.ProgressiveRollout, error) {
	progressiveRollouts := make([]*autoopsproto.ProgressiveRollout, 0)
	cursor := ""
	for {
		resp, err := s.autoOpsClient.ListProgressiveRollouts(
			ctx,
			&autoopsproto.ListProgressiveRolloutsRequest{
				EnvironmentNamespace: environmentNamespace,
				PageSize:             listRequestSize,
				Cursor:               cursor,
			},
		)
		if err != nil {
			return nil, err
		}
		progressiveRollouts = append(progressiveRollouts, resp.ProgressiveRollouts...)
		size := len(progressiveRollouts)
		if size == 0 || size < listRequestSize {
			return progressiveRollouts, nil
		}
		cursor = resp.Cursor
	}
}

func (s *targetStore) GetEnvironments(ctx context.Context) []*environmentdomain.Environment {
	return s.environments.Load().([]*environmentdomain.Environment)
}

func (s *targetStore) GetAutoOpsRules(ctx context.Context, environmentNamespace string) []*autoopsdomain.AutoOpsRule {
	s.autoOpsRulesMtx.Lock()
	autoOpsRules, ok := s.autoOpsRules[environmentNamespace]
	s.autoOpsRulesMtx.Unlock()
	if !ok {
		return nil
	}
	return autoOpsRules
}

func (s *targetStore) GetProgressiveRollouts(
	ctx context.Context,
	environmentNamespace string,
) []*autoopsdomain.ProgressiveRollout {
	s.progressiveRolloutsMtx.Lock()
	progressiveRollouts, ok := s.progressiveRollouts[environmentNamespace]
	s.progressiveRolloutsMtx.Unlock()
	if !ok {
		return nil
	}
	return progressiveRollouts
}
