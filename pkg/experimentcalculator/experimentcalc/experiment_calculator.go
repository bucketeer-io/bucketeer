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
//

package experimentcalc

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-gota/gota/dataframe"
	"go.uber.org/zap"

	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/domain"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	v2es "github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/proto/eventcounter"
	"github.com/bucketeer-io/bucketeer/proto/experiment"
)

var (
	errFailedToSample                = errors.New("calculator: failed to get all the samples")
	errFailedToGetEvalVariationCount = errors.New("calculator: failed to get eval variation count")
	errFailedToGetGoalEventCount     = errors.New("calculator: failed to get goal event count")
)

const (
	day         = 24 * 60 * 60
	numOfChains = 5
)

type ExperimentCalculator struct {
	httpStan    *stan.Stan
	stanModelID string

	environmentClient  envclient.Client
	eventCounterClient ecclient.Client
	experimentClient   experimentclient.Client
	mysqlClient        mysql.Client
	metrics            metrics.Registerer

	location *time.Location
	logger   *zap.Logger
}

func NewExperimentCalculator(
	httpStan *stan.Stan,
	stanModelID string,
	environmentClient envclient.Client,
	eventCounterClient ecclient.Client,
	experimentClient experimentclient.Client,
	mysqlClient mysql.Client,
	metrics metrics.Registerer,
	loc *time.Location,
	logger *zap.Logger,
) *ExperimentCalculator {
	registerMetrics(metrics)
	return &ExperimentCalculator{
		httpStan:           httpStan,
		stanModelID:        stanModelID,
		environmentClient:  environmentClient,
		eventCounterClient: eventCounterClient,
		experimentClient:   experimentClient,
		mysqlClient:        mysqlClient,
		metrics:            metrics,
		location:           loc,

		logger: logger.Named("experiment-calculator"),
	}
}

func (e ExperimentCalculator) Run(ctx context.Context, request *domain.ExperimentCalculatorReq) error {
	startTime := time.Now()
	experimentResult, calculationErr := e.createExperimentResult(ctx, request.EnvironmentId, request.Experiment)
	if calculationErr != nil {
		e.logger.Error("ExperimentCalculator failed to calculate experiment result",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environmentId", request.EnvironmentId),
				zap.Any("experiment", request.Experiment),
				zap.Error(calculationErr),
			)...,
		)
		return calculationErr
	}
	if err := v2es.NewExperimentResultStorage(e.mysqlClient).
		UpdateExperimentResult(ctx, request.EnvironmentId, &domain.ExperimentResult{
			ExperimentResult: experimentResult,
		}); err != nil {
		e.logger.Error("ExperimentCalculator failed to update experiment result",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("environmentId", request.EnvironmentId),
				zap.Any("experiment", request.Experiment),
				zap.Any("experimentResult", experimentResult),
				zap.Error(err),
			)...,
		)
		return err
	}
	e.logger.Info("ExperimentCalculator calculated successfully",
		log.FieldsFromImcomingContext(ctx).AddFields(
			zap.String("environmentId", request.EnvironmentId),
			zap.Any("experiment", request.Experiment),
			zap.Duration("elapsedTime", time.Since(startTime)),
		)...,
	)
	return nil
}

func (e ExperimentCalculator) createExperimentResult(
	ctx context.Context,
	envNamespace string,
	experiment *experiment.Experiment,
) (*eventcounter.ExperimentResult, error) {
	experimentResult := &eventcounter.ExperimentResult{
		Id:           experiment.Id,
		ExperimentId: experiment.Id,
		UpdatedAt:    time.Now().In(e.location).Unix(),
	}
	var variationIDs []string
	for _, variation := range experiment.Variations {
		variationIDs = append(variationIDs, variation.Id)
	}
	endAts := listEndAt(experiment.StartAt, experiment.StopAt, time.Now().In(e.location).Unix())
	for _, goalID := range experiment.GoalIds {
		goalResult := &eventcounter.GoalResult{
			GoalId:  goalID,
			Summary: &eventcounter.Summary{},
		}
		for _, v := range experiment.Variations {
			goalResult.VariationResults = append(goalResult.VariationResults, &eventcounter.VariationResult{
				VariationId: v.Id,
			})
		}
		for _, timestamp := range endAts {
			evalVc, evalErr := e.getEvaluationCount(ctx, &eventcounter.GetExperimentEvaluationCountRequest{
				EnvironmentId:  envNamespace,
				StartAt:        experiment.StartAt,
				EndAt:          timestamp,
				FeatureId:      experiment.FeatureId,
				FeatureVersion: experiment.FeatureVersion,
				VariationIds:   variationIDs,
			})
			if evalErr != nil {
				e.logger.Error("ExperimentCalculator failed to get evaluation count",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.String("namespace", envNamespace),
						zap.Any("experiment", experiment),
						zap.Error(evalErr),
					)...,
				)
				return nil, errFailedToGetEvalVariationCount
			}
			goalVc, goalErr := e.getGoalCount(ctx, &eventcounter.GetExperimentGoalCountRequest{
				EnvironmentId:  envNamespace,
				StartAt:        experiment.StartAt,
				EndAt:          timestamp,
				GoalId:         goalID,
				FeatureId:      experiment.FeatureId,
				FeatureVersion: experiment.FeatureVersion,
				VariationIds:   variationIDs,
			})
			if goalErr != nil {
				e.logger.Error("ExperimentCalculator failed to get goal count",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.String("namespace", envNamespace),
						zap.Any("experiment", experiment),
						zap.Error(goalErr),
					)...,
				)
				return nil, errFailedToGetGoalEventCount
			}
			gr := e.calcGoalResult(ctx, evalVc, goalVc, experiment)
			gr.GoalId = goalID
			e.appendVariationResult(ctx, timestamp, goalResult, gr.VariationResults)
		}

		// Calculate conversion rate for each variation after all timestamps have been processed
		for _, vr := range goalResult.VariationResults {
			if vr.EvaluationCount != nil && vr.EvaluationCount.UserCount > 0 {
				vr.ConversionRate = float64(vr.ExperimentCount.UserCount) / float64(vr.EvaluationCount.UserCount) * 100
			}
		}

		// Calculate expected loss for each variation
		e.calculateExpectedLoss(goalResult.VariationResults)

		// Calculate Summary for this goal result
		e.calculateSummary(ctx, goalResult)

		experimentResult.GoalResults = append(experimentResult.GoalResults, goalResult)
	}

	return experimentResult, nil
}

func (e ExperimentCalculator) getEvaluationCount(
	ctx context.Context,
	req *eventcounter.GetExperimentEvaluationCountRequest,
) (map[string]*eventcounter.VariationCount, error) {
	variationCounts := make(map[string]*eventcounter.VariationCount)
	evaluationCount, err := e.eventCounterClient.GetExperimentEvaluationCount(ctx, req)
	if err != nil {
		return variationCounts, err
	}
	for _, variationCount := range evaluationCount.VariationCounts {
		variationCounts[variationCount.VariationId] = variationCount
	}
	return variationCounts, nil
}

func (e ExperimentCalculator) getGoalCount(
	ctx context.Context,
	req *eventcounter.GetExperimentGoalCountRequest,
) (map[string]*eventcounter.VariationCount, error) {
	variationCounts := make(map[string]*eventcounter.VariationCount)
	goalCount, err := e.eventCounterClient.GetExperimentGoalCount(ctx, req)
	if err != nil {
		return variationCounts, err
	}
	for _, variationCount := range goalCount.VariationCounts {
		variationCounts[variationCount.VariationId] = variationCount
	}
	return variationCounts, nil
}

func (e ExperimentCalculator) calcGoalResult(
	ctx context.Context,
	evalVariationCounts,
	goalVariationCounts map[string]*eventcounter.VariationCount,
	experiment *experiment.Experiment,
) *eventcounter.GoalResult {
	goalResult := &eventcounter.GoalResult{}
	length := len(goalVariationCounts)
	vids := make([]string, 0, length)
	goalUc, evalUc := make([]int64, 0, length), make([]int64, 0, length)
	vrs := make(map[string]*eventcounter.VariationResult, length)
	valueMeans, valueVars := make([]float64, 0, length), make([]float64, 0, length)
	baseLineIdx, loopIdx := 0, 0

	for vid, goalVariationCount := range goalVariationCounts {
		if _, ok := evalVariationCounts[vid]; !ok {
			calculationExceptionCounter.WithLabelValues(evalVariationCountNotFound).Inc()
			e.logger.Error("Variation not found in evaluation count",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.String("variation_id", vid),
					zap.Any("experiment", experiment),
				)...,
			)
			return goalResult
		}
		vids = append(vids, vid)
		evalVariationCount := evalVariationCounts[vid]
		vr := &eventcounter.VariationResult{
			VariationId:     vid,
			ExperimentCount: copyVariationCount(goalVariationCount),
			EvaluationCount: copyVariationCount(evalVariationCount),
		}
		evalUc = append(evalUc, evalVariationCount.UserCount)
		goalUc = append(goalUc, goalVariationCount.UserCount)
		valueMeans = append(valueMeans, goalVariationCount.ValueSumPerUserMean)
		valueVars = append(valueVars, goalVariationCount.ValueSumPerUserVariance)

		goalResult.VariationResults = append(goalResult.VariationResults, vr)
		vrs[vid] = vr
		if vid == experiment.BaseVariationId {
			baseLineIdx = loopIdx
		}
		loopIdx++
	}
	// Skip the calculation if evaluation count is less than goal count.
	for i := 0; i < len(evalUc); i++ {
		if evalUc[i] < goalUc[i] {
			calculationExceptionCounter.WithLabelValues(evaluationCountLessThanGoalEvent).Inc()
			e.logger.Error("Evaluation count is less than goal count",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.String("variation_id", vids[i]),
					zap.Int64("evaluation_count", evalUc[i]),
					zap.Int64("goal_count", goalUc[i]),
					zap.Any("experiment", experiment),
				)...,
			)
			return goalResult
		}
	}

	cvrResult, sampleErr := e.binomialModelSample(ctx, vids, goalUc, evalUc, baseLineIdx, experiment)
	if sampleErr != nil {
		calculationCounter.WithLabelValues(calculationFail).Inc()
		e.logger.Error("BinomialModelSample error",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(sampleErr),
				zap.Any("experiment", experiment),
				zap.Strings("variation_ids", vids),
				zap.Int64s("goal_user_counts", goalUc),
				zap.Int64s("eval_user_counts", evalUc),
			)...,
		)
		return goalResult
	}
	for vid, vr := range cvrResult {
		vrs[vid].CvrProb = copyDistributionSummary(vr.CvrProb)
		vrs[vid].CvrProbBest = copyDistributionSummary(vr.CvrProbBest)
		vrs[vid].CvrProbBeatBaseline = copyDistributionSummary(vr.CvrProbBeatBaseline)
	}
	// Skip the calculation if values are zero.
	for i := 0; i < len(vids); i++ {
		if goalUc[i] == 0 || valueMeans[i] == 0 || valueVars[i] == 0 {
			calculationExceptionCounter.WithLabelValues(valuesAreZero).Inc()
			// skip it if the value is zero
			return goalResult
		}
	}
	valueResult := normalInverseGamma(ctx, vids, valueMeans, valueVars, goalUc, baseLineIdx, 25000)
	for vid, vr := range valueResult {
		vrs[vid].GoalValueSumPerUserProb = copyDistributionSummary(vr.GoalValueSumPerUserProb)
		vrs[vid].GoalValueSumPerUserProbBest = copyDistributionSummary(vr.GoalValueSumPerUserProbBest)
		vrs[vid].GoalValueSumPerUserProbBeatBaseline = copyDistributionSummary(vr.GoalValueSumPerUserProbBeatBaseline)
	}
	calculationCounter.WithLabelValues(calculationSuccess).Inc()
	return goalResult
}

func (e ExperimentCalculator) appendVariationResult(
	ctx context.Context,
	timestamp int64,
	goalResult *eventcounter.GoalResult,
	srcVrs []*eventcounter.VariationResult,
) {
	for i := 0; i < len(goalResult.VariationResults); i++ {
		variationResult := getVariationResult(srcVrs, goalResult.VariationResults[i].VariationId)
		if variationResult == nil {
			e.logger.Error("Variation result not found",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.String("variation_id", goalResult.VariationResults[i].VariationId),
				)...,
			)
			continue
		}
		goalResult.VariationResults[i].ExperimentCount = copyVariationCount(variationResult.ExperimentCount)
		goalResult.VariationResults[i].EvaluationCount = copyVariationCount(variationResult.EvaluationCount)

		goalResult.VariationResults[i].CvrProb = copyDistributionSummary(variationResult.CvrProb)
		goalResult.VariationResults[i].CvrProbBest = copyDistributionSummary(variationResult.CvrProbBest)
		goalResult.VariationResults[i].CvrProbBeatBaseline = copyDistributionSummary(variationResult.CvrProbBeatBaseline)

		goalResult.VariationResults[i].GoalValueSumPerUserProb =
			copyDistributionSummary(variationResult.GoalValueSumPerUserProb)
		goalResult.VariationResults[i].GoalValueSumPerUserProbBest =
			copyDistributionSummary(variationResult.GoalValueSumPerUserProbBest)
		goalResult.VariationResults[i].GoalValueSumPerUserProbBeatBaseline =
			copyDistributionSummary(variationResult.GoalValueSumPerUserProbBeatBaseline)

		// Append to EvaluationUserCountTimeseries
		goalResult.VariationResults[i].EvaluationUserCountTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].EvaluationUserCountTimeseries,
			timestamp,
			float64(goalResult.VariationResults[i].EvaluationCount.UserCount),
		)

		// Append to EvaluationEventCountTimeseries
		goalResult.VariationResults[i].EvaluationEventCountTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].EvaluationEventCountTimeseries,
			timestamp,
			float64(goalResult.VariationResults[i].EvaluationCount.EventCount),
		)

		// Append to GoalUserCountTimeseries
		goalResult.VariationResults[i].GoalUserCountTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalUserCountTimeseries,
			timestamp,
			float64(goalResult.VariationResults[i].ExperimentCount.UserCount),
		)

		// Append to GoalEventCountTimeseries
		goalResult.VariationResults[i].GoalEventCountTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalEventCountTimeseries,
			timestamp,
			float64(goalResult.VariationResults[i].ExperimentCount.EventCount),
		)

		// Append to GoalValueSumTimeseries
		goalResult.VariationResults[i].GoalValueSumTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalValueSumTimeseries,
			timestamp,
			goalResult.VariationResults[i].ExperimentCount.ValueSum,
		)

		// Append to CvrMedianTimeseries
		goalResult.VariationResults[i].CvrMedianTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].CvrMedianTimeseries,
			timestamp,
			goalResult.VariationResults[i].CvrProb.Median,
		)

		// Append to CvrPercentile025Timeseries
		goalResult.VariationResults[i].CvrPercentile025Timeseries = appendToTimeseries(
			goalResult.VariationResults[i].CvrPercentile025Timeseries,
			timestamp,
			goalResult.VariationResults[i].CvrProb.Percentile025,
		)

		// Append to CvrPercentile975Timeseries
		goalResult.VariationResults[i].CvrPercentile975Timeseries = appendToTimeseries(
			goalResult.VariationResults[i].CvrPercentile975Timeseries,
			timestamp,
			goalResult.VariationResults[i].CvrProb.Percentile975,
		)

		// Calculate CVR
		cvr := 0.0
		if goalResult.VariationResults[i].EvaluationCount.UserCount != 0 {
			cvr = float64(goalResult.VariationResults[i].ExperimentCount.UserCount) /
				float64(goalResult.VariationResults[i].EvaluationCount.UserCount)
		}

		// Append to CvrTimeseries
		goalResult.VariationResults[i].CvrTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].CvrTimeseries,
			timestamp,
			cvr,
		)

		// Calculate valuePerUser
		valuePerUser := 0.0
		if goalResult.VariationResults[i].ExperimentCount.UserCount != 0 {
			valuePerUser = goalResult.VariationResults[i].ExperimentCount.ValueSum /
				float64(goalResult.VariationResults[i].ExperimentCount.UserCount)
		}

		// Append to GoalValueSumPerUserTimeseries
		goalResult.VariationResults[i].GoalValueSumPerUserTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalValueSumPerUserTimeseries,
			timestamp,
			valuePerUser,
		)

		// Append to GoalValueSumPerUserMedianTimeseries
		goalResult.VariationResults[i].GoalValueSumPerUserMedianTimeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalValueSumPerUserMedianTimeseries,
			timestamp,
			goalResult.VariationResults[i].GoalValueSumPerUserProb.Median,
		)

		// Append to GoalValueSumPerUserPercentile025Timeseries
		goalResult.VariationResults[i].GoalValueSumPerUserPercentile025Timeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalValueSumPerUserPercentile025Timeseries,
			timestamp,
			goalResult.VariationResults[i].GoalValueSumPerUserProb.Percentile025,
		)

		// Append to GoalValueSumPerUserPercentile975Timeseries
		goalResult.VariationResults[i].GoalValueSumPerUserPercentile975Timeseries = appendToTimeseries(
			goalResult.VariationResults[i].GoalValueSumPerUserPercentile975Timeseries,
			timestamp,
			goalResult.VariationResults[i].GoalValueSumPerUserProb.Percentile975,
		)
	}
}

// appendToTimeseries appends a new timestamp and value to the timeseries and increments total_counts
func appendToTimeseries(ts *eventcounter.Timeseries, timestamp int64, value float64) *eventcounter.Timeseries {
	// If timeseries doesn't exist, create a new one
	if ts == nil {
		return &eventcounter.Timeseries{
			Timestamps:  []int64{timestamp},
			Values:      []float64{value},
			Unit:        eventcounter.Timeseries_DAY,
			TotalCounts: 1,
		}
	}

	// Append the new timestamp and value
	ts.Timestamps = append(ts.Timestamps, timestamp)
	ts.Values = append(ts.Values, value)

	// Increment the total_counts
	ts.TotalCounts++
	return ts
}

func (e ExperimentCalculator) binomialModelSample(
	ctx context.Context,
	vids []string,
	goalUc, evalUc []int64,
	baseLineIdx int,
	experiment *experiment.Experiment,
) (map[string]*eventcounter.VariationResult, error) {
	// The index starts from 1 in PyStan.
	startTime := time.Now()
	baseLineIdx++
	samplesChan := make(chan dataframe.DataFrame)
	wg := sync.WaitGroup{}
	wg.Add(numOfChains)
	for i := 1; i <= numOfChains; i++ {
		go func(chain int) {
			defer wg.Done()
			req := stan.CreateFitReq{
				Chain: chain,
				Data: map[string]interface{}{
					"g": len(goalUc),
					"x": goalUc,
					"n": evalUc,
				},
				Function:   stan.HmcNUTSFunction,
				NumSamples: 21000,
				NumWarmup:  1000,
				RandomSeed: 1234,
			}
			fitResp, err := e.httpStan.CreateFit(ctx, e.stanModelID, req)
			if err != nil {
				e.logger.Error("Failed to create fit",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.String("modelId", e.stanModelID),
						zap.Error(err),
						zap.Any("experiment", experiment),
						zap.Strings("variation_ids", vids),
						zap.Int64s("goal_user_counts", goalUc),
						zap.Int64s("eval_user_counts", evalUc),
					)...,
				)
				return
			}
			fitId := fitResp[len("operations/"):]
			for {
				details, err := e.httpStan.GetOperationDetails(ctx, fitId)
				if err != nil {
					e.logger.Error("Failed to get operation details",
						log.FieldsFromImcomingContext(ctx).AddFields(
							zap.String("fitId", fitId),
							zap.Error(err),
							zap.Any("experiment", experiment),
							zap.Strings("variation_ids", vids),
							zap.Int64s("goal_user_counts", goalUc),
							zap.Int64s("eval_user_counts", evalUc),
						)...,
					)
					return
				}
				if details.Done {
					break
				}
				time.Sleep(50 * time.Millisecond)
			}
			result, err := e.httpStan.GetFitResult(ctx, e.stanModelID, fitId)
			if err != nil {
				e.logger.Error("Failed to get fit result",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.String("fitId", fitId),
						zap.Error(err),
						zap.Any("experiment", experiment),
						zap.Strings("variation_ids", vids),
						zap.Int64s("goal_user_counts", goalUc),
						zap.Int64s("eval_user_counts", evalUc),
					)...,
				)
				return
			}
			fit := e.httpStan.ExtractFromFitResult(ctx, result)
			constrainedNames, _ := e.httpStan.StanParams(ctx, e.stanModelID, req.Data)
			samplesChan <- fit.Select(constrainedNames)
		}(i)
	}
	go func() {
		wg.Wait()
		close(samplesChan)
	}()

	samples := make([]dataframe.DataFrame, 0, numOfChains)
	for sample := range samplesChan {
		samples = append(samples, sample)
	}

	if len(samples) != numOfChains {
		e.logger.Error("Failed to get all samples",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Int("numOfChains", numOfChains),
				zap.Int("numOfSamples", len(samples)),
				zap.Any("experiment", experiment),
				zap.Strings("variation_ids", vids),
				zap.Int64s("goal_user_counts", goalUc),
				zap.Int64s("eval_user_counts", evalUc),
			)...,
		)
		return nil, errFailedToSample
	}

	variationResults := e.convertFitSamples(ctx, samples, vids, baseLineIdx)

	calculationHistogram.WithLabelValues(binomialModelSampleMethod).Observe(time.Since(startTime).Seconds())
	return variationResults, nil
}

// convertFitSamples extracts data from Stan model output and creates VariationResult objects
func (e ExperimentCalculator) convertFitSamples(
	ctx context.Context,
	samples []dataframe.DataFrame,
	vids []string,
	baseLineIdx int,
) map[string]*eventcounter.VariationResult {
	// 1) Merge all chains
	allSample := contactSamples(samples)

	// 2) Build parameter names using the correct Stan column pattern: "p.1", "p.2", etc.
	pNames := make([]string, len(vids))
	for i := range vids {
		pNames[i] = fmt.Sprintf("p.%d", i+1)
	}

	// Get all column names for checking
	colNames := allSample.Names()

	e.logger.Debug("Available columns in Stan samples",
		zap.Strings("columnNames", colNames),
		zap.Strings("expectedColumns", pNames),
	)

	// 3) Extract each p-column into a []float64
	draws := make(map[string][]float64, len(vids))
	for i, name := range pNames {
		// Check if the column exists in the dataframe
		colExists := false
		for _, colName := range colNames {
			if colName == name {
				colExists = true
				break
			}
		}

		if colExists {
			col := allSample.Col(name)
			draws[vids[i]] = col.Float() // Extract samples for this variation
		} else {
			// Log warning if column not found
			e.logger.Warn("Column not found in Stan samples",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.String("columnName", name),
					zap.String("variationId", vids[i]),
				)...,
			)
		}
	}

	// 4) Build your VariationResult map
	variationResults := make(map[string]*eventcounter.VariationResult, len(vids))
	for i, vid := range vids {
		idx := i + 1
		vr := &eventcounter.VariationResult{
			VariationId:         vid,
			CvrProb:             createCvrProb(allSample, samples, idx),
			CvrProbBest:         createCvrProbBest(allSample, samples, idx),
			CvrProbBeatBaseline: createCvrProbBeatBaseline(allSample, samples, baseLineIdx, idx),
		}

		// Add the raw CVR samples if available
		if samples, ok := draws[vid]; ok {
			vr.CvrSamples = samples
		}

		variationResults[vid] = vr
	}

	return variationResults
}

//-------------------------------------utility functions----------------------------------------------

func listEndAt(startAt, endAt, now int64) []int64 {
	var timestamps []int64
	if endAt >= now {
		endAt = now
	}
	for ts := startAt + day; ts < endAt; ts += day {
		timestamps = append(timestamps, ts)
	}
	timestamps = append(timestamps, endAt)
	return timestamps
}

func getVariationResult(vcs []*eventcounter.VariationResult, id string) *eventcounter.VariationResult {
	for _, vc := range vcs {
		if vc.VariationId == id {
			return vc
		}
	}
	return nil
}

func copyVariationCount(from *eventcounter.VariationCount) *eventcounter.VariationCount {
	return &eventcounter.VariationCount{
		VariationId:             from.VariationId,
		UserCount:               from.UserCount,
		EventCount:              from.EventCount,
		ValueSum:                from.ValueSum,
		CreatedAt:               from.CreatedAt,
		VariationValue:          from.VariationValue,
		ValueSumPerUserMean:     from.ValueSumPerUserMean,
		ValueSumPerUserVariance: from.ValueSumPerUserVariance,
	}

}

func copyDistributionSummary(from *eventcounter.DistributionSummary) *eventcounter.DistributionSummary {
	if from == nil {
		return &eventcounter.DistributionSummary{}
	}
	e := &eventcounter.DistributionSummary{
		Mean:          from.Mean,
		Sd:            from.Sd,
		Rhat:          from.Rhat,
		Median:        from.Median,
		Percentile025: from.Percentile025,
		Percentile975: from.Percentile975,
	}
	if from.Histogram != nil {
		hist := make([]int64, 0, len(from.Histogram.Hist))
		bins := make([]float64, 0, len(from.Histogram.Bins))
		copy(hist, from.Histogram.Hist)
		copy(bins, from.Histogram.Bins)
		e.Histogram = &eventcounter.Histogram{
			Hist: hist,
			Bins: bins,
		}
	}
	return e
}

func contactSamples(samples []dataframe.DataFrame) dataframe.DataFrame {
	var df dataframe.DataFrame
	for _, sample := range samples {
		df = df.Concat(sample)
	}
	return df
}

// calculateSummary sets the Summary field of the GoalResult
func (e ExperimentCalculator) calculateSummary(
	ctx context.Context,
	goalResult *eventcounter.GoalResult,
) {
	if len(goalResult.VariationResults) == 0 {
		return
	}

	// 1. Find best variations (cvr_prob_beat_baseline.mean > 95%)
	var bestVariations []*eventcounter.Summary_Variation
	var maxProbability float64
	var maxProbabilityVariationID string

	for _, vr := range goalResult.VariationResults {
		if vr.CvrProbBeatBaseline != nil && vr.CvrProbBeatBaseline.Mean > 0.95 {
			probability := vr.CvrProbBeatBaseline.Mean
			bestVar := &eventcounter.Summary_Variation{
				Id:          vr.VariationId,
				Probability: probability,
				IsBest:      false,
			}

			// Track which variation has the highest probability
			if probability > maxProbability {
				maxProbability = probability
				maxProbabilityVariationID = vr.VariationId
			}

			bestVariations = append(bestVariations, bestVar)
		}
	}

	// Mark the variation with highest probability as outperformed
	for _, bestVar := range bestVariations {
		if bestVar.Id == maxProbabilityVariationID {
			bestVar.IsBest = true
		}
	}

	// 2 & 3. Calculate total counts across all variations
	var totalEvaluationUserCount int64
	var totalGoalUserCount int64

	for _, vr := range goalResult.VariationResults {
		if vr.EvaluationCount != nil {
			totalEvaluationUserCount += vr.EvaluationCount.UserCount
		}
		if vr.ExperimentCount != nil {
			totalGoalUserCount += vr.ExperimentCount.UserCount
		}
	}

	// Set the summary values
	goalResult.Summary.BestVariations = bestVariations
	goalResult.Summary.TotalEvaluationUserCount = totalEvaluationUserCount
	goalResult.Summary.TotalGoalUserCount = totalGoalUserCount
}

// calculateExpectedLoss computes the posterior expected loss (regret) for each variation
// using a Monte Carlo approach over raw CVR samples (vr.CvrSamples).
func (e ExperimentCalculator) calculateExpectedLoss(variationResults []*eventcounter.VariationResult) {
	if len(variationResults) == 0 {
		return
	}

	// 1) Ensure every variation has CVR samples
	for _, vr := range variationResults {
		if len(vr.CvrSamples) == 0 {
			e.logger.Warn("No CVR samples available for expected loss calculation",
				zap.String("variationId", vr.VariationId),
			)
			return
		}
	}

	// 2) Enforce equal sample counts across all variations
	numDraws := len(variationResults[0].CvrSamples)
	for _, vr := range variationResults {
		if len(vr.CvrSamples) != numDraws {
			e.logger.Error("Inconsistent CVR sample lengths",
				zap.String("variationId", vr.VariationId),
				zap.Int("expectedDraws", numDraws),
				zap.Int("actualDraws", len(vr.CvrSamples)),
			)
			return
		}
	}

	// 3) Monte Carlo accumulation of regret
	regretSum := make(map[string]float64, len(variationResults))
	for t := 0; t < numDraws; t++ {
		// find best CVR in this draw
		best := variationResults[0].CvrSamples[t]
		for _, vr := range variationResults[1:] {
			if s := vr.CvrSamples[t]; s > best {
				best = s
			}
		}
		// accumulate (best - this) for each variation
		for _, vr := range variationResults {
			regretSum[vr.VariationId] += best - vr.CvrSamples[t]
		}
	}

	// 4) Average and assign expected loss (×100 to match percent scale)
	for _, vr := range variationResults {
		avgRegret := regretSum[vr.VariationId] / float64(numDraws)
		vr.ExpectedLoss = avgRegret * 100
	}

	e.logger.Debug("Calculated expected loss using Monte Carlo",
		zap.Int("numVariations", len(variationResults)),
		zap.Int("numDraws", numDraws),
	)
}
