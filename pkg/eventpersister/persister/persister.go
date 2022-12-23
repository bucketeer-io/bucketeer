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

package persister

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	aoclient "github.com/bucketeer-io/bucketeer/pkg/autoops/client"
	aodomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	"github.com/bucketeer-io/bucketeer/pkg/cache"
	"github.com/bucketeer-io/bucketeer/pkg/errgroup"
	"github.com/bucketeer-io/bucketeer/pkg/eventpersister/datastore"
	storage "github.com/bucketeer-io/bucketeer/pkg/eventpersister/storage/v2"
	ec "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	featuredomain "github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	featurestorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	"github.com/bucketeer-io/bucketeer/pkg/health"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	aoproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	exproto "github.com/bucketeer-io/bucketeer/proto/experiment"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var (
	ErrUnexpectedMessageType     = errors.New("eventpersister: unexpected message type")
	ErrAutoOpsRulesNotFound      = errors.New("eventpersister: auto ops rules not found")
	ErrExperimentNotFound        = errors.New("eventpersister: experiment not found")
	ErrNoAutoOpsRules            = errors.New("eventpersister: no auto ops rules")
	ErrNoExperiments             = errors.New("eventpersister: no experiments")
	ErrNothingToLink             = errors.New("eventpersister: nothing to link")
	ErrInvalidGoalEventTimestamp = errors.New("eventpersister: invalid goal event timestamp")
)

const (
	listRequestSize        = 500
	furthestEventTimestamp = 24 * time.Hour
	oldestEventTimestamp   = 24 * time.Hour
)

const (
	eventCountKey      = "ec"
	userCountKey       = "uc"
	defaultVariationID = "default"
)

var (
	jpLocation = time.FixedZone("Asia/Tokyo", 9*60*60)
)

type eventMap map[string]proto.Message
type environmentEventMap map[string]eventMap

type options struct {
	maxMPS        int
	numWorkers    int
	flushSize     int
	flushInterval time.Duration
	flushTimeout  time.Duration
	metrics       metrics.Registerer
	logger        *zap.Logger
}

type Option func(*options)

func WithMaxMPS(mps int) Option {
	return func(opts *options) {
		opts.maxMPS = mps
	}
}

func WithNumWorkers(n int) Option {
	return func(opts *options) {
		opts.numWorkers = n
	}
}

func WithFlushSize(s int) Option {
	return func(opts *options) {
		opts.flushSize = s
	}
}

func WithFlushInterval(i time.Duration) Option {
	return func(opts *options) {
		opts.flushInterval = i
	}
}

func WithFlushTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.flushTimeout = timeout
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

type Persister struct {
	experimentClient      ec.Client
	featureClient         featureclient.Client
	autoOpsClient         aoclient.Client
	puller                puller.RateLimitedPuller
	datastore             datastore.Writer
	userEvaluationStorage featurestorage.UserEvaluationsStorage
	group                 errgroup.Group
	opts                  *options
	logger                *zap.Logger
	ctx                   context.Context
	cancel                func()
	doneCh                chan struct{}
	mysqlClient           mysql.Client
	evaluationCountCacher cache.MultiGetDeleteCountCache
	flightgroup           singleflight.Group
}

func NewPersister(
	experimentClient ec.Client,
	featureClient featureclient.Client,
	autoOpsClient aoclient.Client,
	p puller.Puller,
	ds datastore.Writer,
	bt bigtable.Client,
	mysqlClient mysql.Client,
	v3Cache cache.MultiGetDeleteCountCache,
	opts ...Option,
) *Persister {
	dopts := &options{
		maxMPS:        1000,
		numWorkers:    1,
		flushSize:     50,
		flushInterval: 5 * time.Second,
		flushTimeout:  20 * time.Second,
		logger:        zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	if dopts.metrics != nil {
		registerMetrics(dopts.metrics)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		experimentClient:      experimentClient,
		featureClient:         featureClient,
		autoOpsClient:         autoOpsClient,
		puller:                puller.NewRateLimitedPuller(p, dopts.maxMPS),
		datastore:             ds,
		userEvaluationStorage: featurestorage.NewUserEvaluationsStorage(bt),
		opts:                  dopts,
		logger:                dopts.logger.Named("persister"),
		ctx:                   ctx,
		cancel:                cancel,
		doneCh:                make(chan struct{}),
		mysqlClient:           mysqlClient,
		evaluationCountCacher: v3Cache,
	}
}

func (p *Persister) Run() error {
	defer close(p.doneCh)
	p.group.Go(func() error {
		return p.puller.Run(p.ctx)
	})
	for i := 0; i < p.opts.numWorkers; i++ {
		p.group.Go(p.batch)
	}
	return p.group.Wait()
}

func (p *Persister) Stop() {
	p.cancel()
	<-p.doneCh
}

func (p *Persister) Check(ctx context.Context) health.Status {
	select {
	case <-p.ctx.Done():
		p.logger.Error("Unhealthy due to context Done is closed", zap.Error(p.ctx.Err()))
		return health.Unhealthy
	default:
		if p.group.FinishedCount() > 0 {
			p.logger.Error("Unhealthy", zap.Int32("FinishedCount", p.group.FinishedCount()))
			return health.Unhealthy
		}
		return health.Healthy
	}
}

func (p *Persister) batch() error {
	batch := make(map[string]*puller.Message)
	timer := time.NewTimer(p.opts.flushInterval)
	defer timer.Stop()
	for {
		select {
		case msg, ok := <-p.puller.MessageCh():
			if !ok {
				return nil
			}
			receivedCounter.Inc()
			id := msg.Attributes["id"]
			if id == "" {
				msg.Ack()
				// TODO: better log format for msg data
				handledCounter.WithLabelValues(codes.MissingID.String()).Inc()
				continue
			}
			if previous, ok := batch[id]; ok {
				previous.Ack()
				p.logger.Warn("Message with duplicate id", zap.String("id", id))
				handledCounter.WithLabelValues(codes.DuplicateID.String()).Inc()
			}
			batch[id] = msg
			if len(batch) < p.opts.flushSize {
				continue
			}
			p.send(batch)
			batch = make(map[string]*puller.Message)
			timer.Reset(p.opts.flushInterval)
		case <-timer.C:
			if len(batch) > 0 {
				p.send(batch)
				batch = make(map[string]*puller.Message)
			}
			timer.Reset(p.opts.flushInterval)
		case <-p.ctx.Done():
			batchSize := len(batch)
			p.logger.Info("Context is done", zap.Int("batchSize", batchSize))
			if len(batch) > 0 {
				p.send(batch)
				p.logger.Info("All the left messages are processed successfully", zap.Int("batchSize", batchSize))
			}
			return nil
		}
	}
}

func (p *Persister) send(messages map[string]*puller.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), p.opts.flushTimeout)
	defer cancel()
	envEvents := p.extractEvents(messages)
	if len(envEvents) == 0 {
		p.logger.Error("all messages were bad")
		return
	}
	fails := make(map[string]bool, len(messages))
	for environmentNamespace, events := range envEvents {
		evs := make(map[string]string, len(events))
		for id, event := range events {
			if err := p.upsertMAU(ctx, event, environmentNamespace); err != nil {
				p.logger.Error(
					"failed to store a mau",
					zap.Error(err),
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = true
				continue
			}
			eventJSON, repeatable, err := p.marshalEvent(ctx, event, environmentNamespace)
			if err != nil {
				if !repeatable {
					p.logger.Error(
						"failed to marshal an unrepeatable event",
						zap.Error(err),
						zap.String("id", id),
						zap.String("environmentNamespace", environmentNamespace),
					)
				}
				fails[id] = repeatable
				continue
			}
			if err := p.upsertEvaluationCount(event, environmentNamespace); err != nil {
				p.logger.Error(
					"failed to upsert an evaluation event in redis",
					zap.Error(err),
					zap.String("id", id),
					zap.String("environmentNamespace", environmentNamespace),
				)
				fails[id] = true
				continue
			}
			evs[id] = eventJSON
		}
		if len(evs) > 0 {
			fs, err := p.datastore.Write(ctx, evs, environmentNamespace)
			if err != nil {
				p.logger.Error(
					"could not write to datastore",
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)
			}
			for id, f := range fs {
				fails[id] = f
			}
		}
	}
	for id, m := range messages {
		if repeatable, ok := fails[id]; ok {
			if repeatable {
				m.Nack()
				handledCounter.WithLabelValues(codes.RepeatableError.String()).Inc()
			} else {
				m.Ack()
				handledCounter.WithLabelValues(codes.NonRepeatableError.String()).Inc()
			}
			continue
		}
		m.Ack()
		handledCounter.WithLabelValues(codes.OK.String()).Inc()
	}
}

func (p *Persister) extractEvents(messages map[string]*puller.Message) environmentEventMap {
	envEvents := environmentEventMap{}
	handleBadMessage := func(m *puller.Message, err error) {
		m.Ack()
		p.logger.Error("bad message", zap.Error(err), zap.Any("msg", m))
		handledCounter.WithLabelValues(codes.BadMessage.String()).Inc()
	}
	for _, m := range messages {
		event := &eventproto.Event{}
		if err := proto.Unmarshal(m.Data, event); err != nil {
			handleBadMessage(m, err)
			continue
		}
		var innerEvent ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(event.Event, &innerEvent); err != nil {
			handleBadMessage(m, err)
			continue
		}
		if innerEvents, ok := envEvents[event.EnvironmentNamespace]; ok {
			innerEvents[event.Id] = innerEvent.Message
			continue
		}
		envEvents[event.EnvironmentNamespace] = eventMap{event.Id: innerEvent.Message}
	}
	return envEvents
}

func (p *Persister) upsertMAU(ctx context.Context, event proto.Message, environmentNamespace string) error {
	if p.mysqlClient == nil {
		return nil
	}
	if e, ok := event.(*esproto.UserEvent); ok {
		s := storage.NewMysqlMAUStorage(p.mysqlClient)
		return s.UpsertMAU(ctx, e, environmentNamespace)
	}
	return nil
}

func (p *Persister) marshalEvent(
	ctx context.Context,
	event interface{},
	environmentNamespace string,
) (string, bool, error) {
	switch event := event.(type) {
	case *eventproto.EvaluationEvent:
		return p.marshalEvaluationEvent(ctx, event, environmentNamespace)
	case *eventproto.GoalEvent:
		return p.marshalGoalEvent(ctx, event, environmentNamespace)
	case *esproto.UserEvent:
		return p.marshalUserEvent(event, environmentNamespace)
	}
	return "", false, ErrUnexpectedMessageType
}

func (p *Persister) marshalEvaluationEvent(
	ctx context.Context,
	e *eventproto.EvaluationEvent,
	environmentNamespace string,
) (string, bool, error) {
	evaluation, tag := p.convToEvaluation(ctx, e)
	if err := p.upsertUserEvaluation(ctx, environmentNamespace, tag, evaluation); err != nil {
		handledCounter.WithLabelValues(codeUpsertUserEvaluationFailed).Inc()
		return "", true, err
	}
	m := map[string]string{}
	m["environmentNamespace"] = environmentNamespace
	m["sourceId"] = e.SourceId.String()
	m["tag"] = e.Tag
	m["timestamp"] = time.Unix(e.Timestamp, 0).Format(time.RFC3339)
	m["featureId"] = e.FeatureId
	m["featureVersion"] = strconv.FormatInt(int64(e.FeatureVersion), 10)
	m["userId"] = e.UserId
	m["metric.userId"] = e.UserId
	m["variationId"] = e.VariationId
	if e.Reason != nil {
		m["reason"] = e.Reason.Type.String()
	}
	if e.User != nil {
		for k, v := range e.User.Data {
			m[userMetadataColumn(environmentNamespace, k)] = v
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", false, err
	}
	return string(b), false, nil
}

func (p *Persister) marshalGoalEvent(
	ctx context.Context,
	e *eventproto.GoalEvent,
	environmentNamespace string,
) (string, bool, error) {
	evaluations, retriable, err := p.linkGoalEvent(ctx, e, environmentNamespace)
	if err != nil {
		return "", retriable, err
	}
	m := map[string]interface{}{}
	m["environmentNamespace"] = environmentNamespace
	m["sourceId"] = e.SourceId.String()
	m["tag"] = e.Tag
	m["timestamp"] = time.Unix(e.Timestamp, 0).Format(time.RFC3339)
	m["goalId"] = e.GoalId
	m["userId"] = e.UserId
	m["metric.userId"] = e.UserId
	if e.User != nil {
		for k, v := range e.User.Data {
			m[userMetadataColumn(environmentNamespace, k)] = v
		}
	}
	m["value"] = strconv.FormatFloat(e.Value, 'f', -1, 64)
	m["evaluations"] = evaluations
	b, err := json.Marshal(m)
	if err != nil {
		return "", false, err
	}
	return string(b), false, nil
}

func (p *Persister) convToEvaluation(
	ctx context.Context,
	event *eventproto.EvaluationEvent,
) (*featureproto.Evaluation, string) {
	evaluation := &featureproto.Evaluation{
		Id: featuredomain.EvaluationID(
			event.FeatureId,
			event.FeatureVersion,
			event.UserId,
		),
		FeatureId:      event.FeatureId,
		FeatureVersion: event.FeatureVersion,
		UserId:         event.UserId,
		VariationId:    event.VariationId,
		Reason:         event.Reason,
	}
	// For requests that doesn't have the tag info,
	// it will insert none instead, until all SDK clients are updated
	var tag string
	if event.Tag == "" {
		tag = "none"
	} else {
		tag = event.Tag
	}
	return evaluation, tag
}

func (p *Persister) upsertUserEvaluation(
	ctx context.Context,
	environmentNamespace, tag string,
	evaluation *featureproto.Evaluation,
) error {
	if err := p.userEvaluationStorage.UpsertUserEvaluation(
		ctx,
		evaluation,
		environmentNamespace,
		tag,
	); err != nil {
		return err
	}
	return nil
}

func (p *Persister) marshalUserEvent(e *esproto.UserEvent, environmentNamespace string) (string, bool, error) {
	m := map[string]interface{}{}
	m["environmentNamespace"] = environmentNamespace
	m["sourceId"] = e.SourceId.String()
	m["tag"] = e.Tag
	m["timestamp"] = time.Unix(e.LastSeen, 0).Format(time.RFC3339)
	m["userId"] = e.UserId
	b, err := json.Marshal(m)
	if err != nil {
		return "", false, err
	}
	return string(b), false, nil
}

func userMetadataColumn(environmentNamespace string, key string) string {
	if environmentNamespace == "" {
		return fmt.Sprintf("user.data.%s", key)
	}
	return fmt.Sprintf("%s.user.data.%s", environmentNamespace, key)
}

func getVariationID(reason featureproto.Reason_Type, vID string) string {
	if reason == featureproto.Reason_CLIENT {
		return defaultVariationID
	}
	return vID
}

func (p *Persister) upsertEvaluationCount(event proto.Message, environmentNamespace string) error {
	if e, ok := event.(*eventproto.EvaluationEvent); ok {
		vID := getVariationID(e.Reason.Type, e.VariationId)
		// To avoid duplication when the request fails, we increment the event count in the end
		// because the user count is an unique count, and there is no problem adding the same event more than once
		uck := p.newEvaluationCountkey(userCountKey, e.FeatureId, vID, environmentNamespace, e.Timestamp)
		if err := p.countUser(uck, e.UserId); err != nil {
			return err
		}
		eck := p.newEvaluationCountkey(eventCountKey, e.FeatureId, vID, environmentNamespace, e.Timestamp)
		if err := p.countEvent(eck); err != nil {
			return err
		}
	}
	return nil
}

func (p *Persister) newEvaluationCountkey(
	kind, featureID, variationID, environmentNamespace string,
	timestamp int64,
) string {
	t := time.Unix(timestamp, 0)
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jpLocation)
	return cache.MakeKey(
		kind,
		fmt.Sprintf("%d:%s:%s", date.Unix(), featureID, variationID),
		environmentNamespace,
	)
}

// validateTimestamp limits date range of the given timestamp
func (p *Persister) validateTimestamp(
	timestamp int64,
	oldestTimestampDuration, furthestTimestampDuration time.Duration,
) bool {
	given := time.Unix(timestamp, 0)
	maxPast := time.Now().Add(-oldestTimestampDuration)
	if given.Before(maxPast) {
		return false
	}
	maxFuture := time.Now().Add(furthestTimestampDuration)
	return !given.After(maxFuture)
}

func (p *Persister) linkGoalEvent(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) ([]string, bool, error) {
	if !p.validateTimestamp(event.Timestamp, oldestEventTimestamp, furthestEventTimestamp) {
		handledCounter.WithLabelValues(codeInvalidGoalEventTimestamp).Inc()
		return nil, false, ErrInvalidGoalEventTimestamp
	}
	evaluations := []*featureproto.Evaluation{}
	evalExp, retriable, err := p.linkGoalEventByExperiment(ctx, event, environmentNamespace)
	// If there are no experiments or the goal ID didn't match, it will ignore the error
	// so it can try to link the goal event to auto ops.
	if err != nil && err != ErrNoExperiments && err != ErrExperimentNotFound {
		return nil, retriable, err
	}
	if evalExp != nil {
		evaluations = append(evaluations, evalExp)
	}
	evalAuto, retriable, err := p.linkGoalEventByAutoOps(ctx, event, environmentNamespace)
	// If there are no rules or the goal ID didn't match, it will ignore the error
	// so we can acknowledge the message because the event doesn't belong to no one
	if err != nil && err != ErrNoAutoOpsRules && err != ErrAutoOpsRulesNotFound {
		return nil, retriable, err
	}
	if evalAuto != nil {
		evaluations = append(evaluations, evalAuto...)
	}
	if len(evaluations) == 0 {
		handledCounter.WithLabelValues(codeNothingToLink).Inc()
		return nil, false, ErrNothingToLink
	}
	evalsMap := make(map[string]struct{})
	for _, ev := range evaluations {
		// Check the reason
		reason := ""
		if ev.Reason != nil {
			reason = ev.Reason.Type.String()
		}
		eval := fmt.Sprintf("%s:%d:%s:%s", ev.FeatureId, ev.FeatureVersion, ev.VariationId, reason)
		// Remove duplicates if needed
		evalsMap[eval] = struct{}{}
	}
	// Convert it to slice
	evals := []string{}
	for key := range evalsMap {
		evals = append(evals, key)
	}
	// Sort the slice to avoid errors in the unit tests
	sort.Sort(sort.Reverse(sort.StringSlice(evals)))
	return evals, false, nil
}

// Because the same goal can be used on other feature flags
// we must link the goal event to all flags using the same goal ID
// If it fails even once to link the goal event, it will retry until all events are linked.
func (p *Persister) linkGoalEventByAutoOps(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) ([]*featureproto.Evaluation, bool, error) {
	// List all auto ops rules
	list, err := p.listAutoOpsRules(ctx, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if len(list) == 0 {
		return nil, false, ErrNoAutoOpsRules
	}
	// Find the feature flags by goal ID
	featureIDs := p.findFeatureIDs(event.GoalId, list)
	if len(featureIDs) == 0 {
		return nil, false, ErrAutoOpsRulesNotFound
	}
	// Get the lastest feature version
	resp, err := p.featureClient.GetFeatures(ctx, &featureproto.GetFeaturesRequest{
		EnvironmentNamespace: environmentNamespace,
		Ids:                  featureIDs,
	})
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToGetFeatures).Inc()
		return nil, true, err
	}
	// Get all user evaluations using the feature flag info
	evaluations := make([]*featureproto.Evaluation, 0, len(resp.Features))
	for _, feature := range resp.Features {
		ev, err := p.getUserEvaluation(
			environmentNamespace,
			event.UserId,
			event.Tag,
			feature.Id,
			feature.Version,
		)
		if err != nil {
			return nil, true, err
		}
		evaluations = append(evaluations, ev)
	}
	return evaluations, false, nil
}

func (p *Persister) findFeatureIDs(goalID string, listAutoOpsRules []*aoproto.AutoOpsRule) []string {
	featureIDs := []string{}
	for _, aor := range listAutoOpsRules {
		autoOpsRule := &aodomain.AutoOpsRule{AutoOpsRule: aor}
		// We ignore the rules that are already triggered
		if autoOpsRule.AlreadyTriggered() {
			continue
		}
		clauses, err := autoOpsRule.ExtractOpsEventRateClauses()
		if err != nil {
			handledCounter.WithLabelValues(codeFailedToExtractOpsEventRateClauses).Inc()
			continue
		}
		for _, clause := range clauses {
			if clause.GoalId == goalID {
				featureIDs = append(featureIDs, autoOpsRule.FeatureId)
			}
		}
	}
	return featureIDs
}

func (p *Persister) listAutoOpsRules(
	ctx context.Context,
	environmentNamespace string,
) ([]*aoproto.AutoOpsRule, error) {
	exp, err, _ := p.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listAutoOpsRules"),
		func() (interface{}, error) {
			aor := []*aoproto.AutoOpsRule{}
			cursor := ""
			for {
				resp, err := p.autoOpsClient.ListAutoOpsRules(ctx, &aoproto.ListAutoOpsRulesRequest{
					EnvironmentNamespace: environmentNamespace,
					PageSize:             listRequestSize,
					Cursor:               cursor,
				})
				if err != nil {
					return nil, err
				}
				aor = append(aor, resp.AutoOpsRules...)
				aorSize := len(resp.AutoOpsRules)
				if aorSize == 0 || aorSize < listRequestSize {
					return aor, nil
				}
				cursor = resp.Cursor
			}
		},
	)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToListAutoOpsRules).Inc()
		return nil, err
	}
	return exp.([]*aoproto.AutoOpsRule), nil
}

func (p *Persister) linkGoalEventByExperiment(
	ctx context.Context,
	event *eventproto.GoalEvent,
	environmentNamespace string,
) (*featureproto.Evaluation, bool, error) {
	// List experiments with the following status RUNNING, FORCE_STOPPED, and STOPPED
	experiments, err := p.listExperiments(ctx, environmentNamespace)
	if err != nil {
		return nil, true, err
	}
	if len(experiments) == 0 {
		return nil, false, ErrNoExperiments
	}
	// Find the experiment by goal ID
	// TODO: we must change the console UI not to allow creating
	// multiple experiments running at the same time,
	// using the same feature flag id and goal id
	var experiment *exproto.Experiment
	for _, exp := range experiments {
		if p.findGoalID(event.GoalId, exp.GoalIds) {
			experiment = exp
			break
		}
	}
	if experiment == nil {
		return nil, false, ErrExperimentNotFound
	}
	// Get the user evaluation using the experiment info
	ev, err := p.getUserEvaluation(
		environmentNamespace,
		event.UserId,
		event.Tag,
		experiment.FeatureId,
		experiment.FeatureVersion,
	)
	if err != nil {
		return nil, true, err
	}
	return ev, false, nil
}

func (p *Persister) listExperiments(
	ctx context.Context,
	environmentNamespace string,
) ([]*exproto.Experiment, error) {
	exp, err, _ := p.flightgroup.Do(
		fmt.Sprintf("%s:%s", environmentNamespace, "listExperiments"),
		func() (interface{}, error) {
			experiments := []*exproto.Experiment{}
			cursor := ""
			for {
				resp, err := p.experimentClient.ListExperiments(ctx, &exproto.ListExperimentsRequest{
					PageSize:             listRequestSize,
					Cursor:               cursor,
					EnvironmentNamespace: environmentNamespace,
					Statuses: []exproto.Experiment_Status{
						exproto.Experiment_RUNNING,
						exproto.Experiment_FORCE_STOPPED,
						exproto.Experiment_STOPPED,
					},
					Archived: &wrappers.BoolValue{Value: false},
				})
				if err != nil {
					return nil, err
				}
				experiments = append(experiments, resp.Experiments...)
				experimentSize := len(resp.Experiments)
				if experimentSize == 0 || experimentSize < listRequestSize {
					return experiments, nil
				}
				cursor = resp.Cursor
			}
		},
	)
	if err != nil {
		handledCounter.WithLabelValues(codeFailedToListExperiments).Inc()
		return nil, err
	}
	return exp.([]*exproto.Experiment), nil
}

func (p *Persister) findGoalID(id string, goalIDs []string) bool {
	for _, goalID := range goalIDs {
		if id == goalID {
			return true
		}
	}
	return false
}

func (p *Persister) getUserEvaluation(
	environmentNamespace,
	userID,
	tag,
	featureID string,
	featureVersion int32,
) (*featureproto.Evaluation, error) {
	// For requests with no tag, it will insert "none" instead, until all old SDK clients are updated
	if tag == "" {
		tag = "none"
	}
	evaluation, err := p.userEvaluationStorage.GetUserEvaluation(
		p.ctx,
		userID,
		environmentNamespace,
		tag,
		featureID,
		featureVersion,
	)
	if err != nil {
		if err == bigtable.ErrKeyNotFound {
			handledCounter.WithLabelValues(codeUserEvaluationNotFound).Inc()
		} else {
			handledCounter.WithLabelValues(codeFailedToGetUserEvaluation).Inc()
		}
		return nil, err
	}
	return evaluation, nil
}

func (p *Persister) countEvent(key string) error {
	_, err := p.evaluationCountCacher.Increment(key)
	if err != nil {
		return err
	}
	return nil
}

func (p *Persister) countUser(key, userID string) error {
	_, err := p.evaluationCountCacher.PFAdd(key, userID)
	if err != nil {
		return err
	}
	return nil
}
