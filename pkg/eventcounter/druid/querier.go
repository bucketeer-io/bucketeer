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
package druid

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ca-dp/godruid"
	"go.uber.org/zap"

	storagedruid "github.com/bucketeer-io/bucketeer/pkg/storage/druid"
	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

const (
	DataTypeEvaluationEvents = "evaluation_events"
	DataTypeGoalEvents       = "goal_events"
	DataTypeUserEvents       = "user_events"
	ColumnVariation          = "Variation"
	ColumnFeatureVersion     = "Feature version"
	ColumnCVR                = "Conversion rate"
	ColumnGoalUser           = "Goal user"
	ColumnGoalTotal          = "Goal total"
	ColumnGoalValueMean      = "Goal value mean"
	ColumnGoalValueTotal     = "Goal value total"
	ColumnGoalValueVariance  = "Goal value variance"
	ColumnEvaluationUser     = "Evaluation user"
	ColumnEvaluationTotal    = "Evaluation total"
	ColumnUser               = "User"
	ColumnUserCount          = "User count"
	ColumnUserTotal          = "User total"
)

var (
	variationRegex = regexp.MustCompile(`^.*:.*:(.*)$`)
)

type Querier interface {
	QuerySegmentMetadata(ctx context.Context, environmentNamespace, dataType string) ([]string, error)
	QueryGoalCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		goalID, featureID string,
		featureVersion int32,
		reason string,
		segmnets []string,
		filters []*ecproto.Filter,
	) (*ecproto.Row, []*ecproto.Row, error)
	QueryEvaluationCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		featureID string,
		featureVersion int32,
		reason string,
		segmnets []string,
		filters []*ecproto.Filter,
	) (*ecproto.Row, []*ecproto.Row, error)
	QueryEvaluationTimeseriesCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		featureID string,
		featureVersion int32,
		variationID string,
	) (map[string]*ecproto.VariationTimeseries, error)
	QueryUserCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
	) (*ecproto.Row, []*ecproto.Row, error)
	QueryCount(
		ctx context.Context,
		environmentNamespace string,
		startAt, endAt time.Time,
		goalID, featureID string,
		featureVersion int32,
		reason string,
		segmnets []string,
		filters []*ecproto.Filter,
	) (*ecproto.Row, []*ecproto.Row, error)
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type druidQuerier struct {
	brokerClient     *storagedruid.BrokerClient
	datasourcePrefix string
	opts             *options
	logger           *zap.Logger
}

func NewDruidQuerier(
	brokerClient *storagedruid.BrokerClient,
	datasourcePrefix string,
	opts ...Option,
) Querier {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &druidQuerier{
		brokerClient:     brokerClient,
		datasourcePrefix: datasourcePrefix,
		opts:             dopts,
		logger:           dopts.logger.Named("druid-querier"),
	}
}

func (q *druidQuerier) QuerySegmentMetadata(
	ctx context.Context,
	environmentNamespace, dataType string,
) ([]string, error) {
	datasource := storagedruid.Datasource(q.datasourcePrefix, dataType)
	endAt := time.Now()
	// Fetch metadata from last 7 days data.
	startAt := endAt.Add(-7 * 24 * time.Hour)
	query := querySegmentMetadata(datasource, startAt, endAt)
	if err := q.brokerClient.Query(query, ""); err != nil {
		b, _ := json.Marshal(query)
		q.logger.Error("Failed to query segment metadata", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("datastore", datasource),
			zap.String("query", string(b)))
		return nil, err
	}
	data := []string{}
	if len(query.QueryResult) == 0 {
		return data, nil
	}
	userDataRegex := regexp.MustCompile(userDataPattern(environmentNamespace))
	for k := range query.QueryResult[0].Columns {
		if userDataRegex.MatchString(k) {
			data = append(data, removeEnvFromUserData(k, userDataRegex))
		}
	}
	return data, nil
}

func (q *druidQuerier) QueryGoalCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
	goalID, featureID string,
	featureVersion int32,
	reason string,
	segments []string,
	filters []*ecproto.Filter,
) (*ecproto.Row, []*ecproto.Row, error) {
	datasource := storagedruid.Datasource(q.datasourcePrefix, DataTypeGoalEvents)
	envSegments := convToEnvSegments(environmentNamespace, segments)
	envFilters := convToEnvFilters(environmentNamespace, filters)
	query := queryGoalGroupBy(
		datasource,
		startAt,
		endAt,
		environmentNamespace,
		goalID,
		featureID,
		featureVersion,
		reason,
		envSegments,
		envFilters,
	)
	if err := q.brokerClient.Query(query, ""); err != nil {
		b, _ := json.Marshal(query)
		q.logger.Error("Failed to query goal counts", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("datastore", datasource),
			zap.String("query", string(b)))
		return nil, nil, err
	}
	columns := generateColumns(
		featureID,
		envSegments,
		[]string{
			ColumnGoalUser,
			ColumnGoalTotal,
			ColumnGoalValueTotal,
			ColumnGoalValueMean,
			ColumnGoalValueVariance,
		},
	)
	headers := q.convHeaders(environmentNamespace, columns)
	rows, errs := q.convToTable(query.QueryResult, columns)
	if len(errs) > 0 {
		q.logger.Error("Failed to convert query result to table",
			zap.Errors("errs", errs),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("goalID", goalID),
			zap.String("featureID", featureID),
			zap.Strings("segments", segments),
			zap.Time("startAt", startAt),
			zap.Time("endAt", endAt))
	}
	return headers, rows, nil
}

func (q *druidQuerier) QueryEvaluationCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
	reason string,
	segments []string,
	filters []*ecproto.Filter,
) (*ecproto.Row, []*ecproto.Row, error) {
	datasource := storagedruid.Datasource(q.datasourcePrefix, DataTypeEvaluationEvents)
	envSegments := convToEnvSegments(environmentNamespace, segments)
	envFilters := convToEnvFilters(environmentNamespace, filters)
	query := queryEvaluationGroupBy(
		datasource,
		startAt,
		endAt,
		environmentNamespace,
		featureID,
		featureVersion,
		reason,
		envSegments,
		envFilters,
	)
	if err := q.brokerClient.Query(query, ""); err != nil {
		b, _ := json.Marshal(query)
		q.logger.Error("Failed to query evaluation counts", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("datastore", datasource),
			zap.String("query", string(b)))
		return nil, nil, err
	}
	columns := generateColumns(featureID, envSegments, []string{ColumnEvaluationUser, ColumnEvaluationTotal})
	headers := q.convHeaders(environmentNamespace, columns)
	rows, errs := q.convToTable(query.QueryResult, columns)
	if len(errs) > 0 {
		q.logger.Error("Failed to convert query result to table",
			zap.Errors("errs", errs),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureID", featureID),
			zap.Strings("segments", segments),
			zap.Time("startAt", startAt),
			zap.Time("endAt", endAt))
	}
	return headers, rows, nil
}

func (q *druidQuerier) QueryEvaluationTimeseriesCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
	featureID string,
	featureVersion int32,
	variationID string,
) (map[string]*ecproto.VariationTimeseries, error) {
	datasource := storagedruid.Datasource(q.datasourcePrefix, DataTypeEvaluationEvents)
	query := queryEvaluationTimeseries(
		datasource,
		startAt,
		endAt,
		environmentNamespace,
		featureID,
		featureVersion,
		variationID,
	)
	if err := q.brokerClient.Query(query, ""); err != nil {
		b, _ := json.Marshal(query)
		q.logger.Error("Failed to query evaluation counts", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("datastore", datasource),
			zap.String("query", string(b)))
		return nil, err
	}
	ts, errs := q.convToTimeseries(
		variationID,
		query.QueryResult,
		[]string{ColumnEvaluationTotal, ColumnEvaluationUser},
	)
	if len(errs) > 0 {
		q.logger.Error("Failed to convert query result to table",
			zap.Errors("errs", errs),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("featureID", featureID),
			zap.Int32("featureVersion", featureVersion),
			zap.String("variationID", variationID),
			zap.Time("startAt", startAt),
			zap.Time("endAt", endAt))
	}
	return ts, nil
}

func (q *druidQuerier) QueryUserCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
) (*ecproto.Row, []*ecproto.Row, error) {
	datasource := storagedruid.Datasource(q.datasourcePrefix, DataTypeUserEvents)
	query := queryUserGroupBy(datasource, environmentNamespace, startAt, endAt)
	if err := q.brokerClient.Query(query, ""); err != nil {
		b, _ := json.Marshal(query)
		q.logger.Error("Failed to query user count", zap.Error(err),
			zap.String("environmentNamespace", environmentNamespace),
			zap.String("datastore", datasource),
			zap.String("query", string(b)))
		return nil, nil, err
	}
	columns := []string{ColumnUserTotal, ColumnUserCount}
	headers := q.convHeaders(environmentNamespace, columns)
	rows, errs := q.convToTable(query.QueryResult, columns)
	if len(errs) > 0 {
		q.logger.Error("Failed to convert query result to table for user count",
			zap.Errors("errs", errs),
			zap.String("environmentNamespace", environmentNamespace),
			zap.Time("startAt", startAt),
			zap.Time("endAt", endAt))
	}
	return headers, rows, nil
}

func (q *druidQuerier) QueryCount(
	ctx context.Context,
	environmentNamespace string,
	startAt, endAt time.Time,
	goalID string,
	featureID string,
	featureVersion int32,
	reason string,
	segments []string,
	filters []*ecproto.Filter,
) (*ecproto.Row, []*ecproto.Row, error) {
	goalHeader, goalRows, err := q.QueryGoalCount(
		ctx,
		environmentNamespace,
		startAt,
		endAt,
		goalID,
		featureID,
		featureVersion,
		reason,
		segments,
		filters,
	)
	if err != nil {
		return nil, nil, err
	}
	// If featureID is not provided, return goal count.
	if featureID == "" {
		return goalHeader, goalRows, nil
	}
	evalHeader, evalRows, err := q.QueryEvaluationCount(
		ctx,
		environmentNamespace,
		startAt,
		endAt,
		featureID,
		featureVersion,
		reason,
		segments,
		filters,
	)
	if err != nil {
		return nil, nil, err
	}
	headers, rows := convToResult(ctx, evalHeader, goalHeader, evalRows, goalRows, segments)
	return headers, rows, nil
}

func (q *druidQuerier) convToTable(queryResult []godruid.GroupbyItem, columns []string) ([]*ecproto.Row, []error) {
	rows := []*ecproto.Row{}
	errs := []error{}
	for _, item := range queryResult {
		cells := []*ecproto.Cell{}
		for _, column := range columns {
			value, ok := item.Event[column]
			if !ok {
				errs = append(errs, fmt.Errorf("Column %s does not exist", column))
				cells = append(cells, &ecproto.Cell{Type: ecproto.Cell_STRING, Value: ""})
				continue
			}
			var cell *ecproto.Cell
			switch value := value.(type) {
			case float64:
				cell = &ecproto.Cell{Type: ecproto.Cell_DOUBLE, ValueDouble: value}
			case string:
				if column == ColumnVariation {
					groups := variationRegex.FindStringSubmatch(value)
					if groups != nil {
						cell = &ecproto.Cell{Type: ecproto.Cell_STRING, Value: groups[1]}
					} else {
						cell = &ecproto.Cell{Type: ecproto.Cell_STRING, Value: value}
					}
				} else {
					cell = &ecproto.Cell{Type: ecproto.Cell_STRING, Value: value}
				}
			default:
				cell = &ecproto.Cell{Type: ecproto.Cell_STRING, Value: ""}
				errs = append(errs, fmt.Errorf("Value %v type of column %s is unknown", value, column))
			}
			cells = append(cells, cell)
		}
		rows = append(rows, &ecproto.Row{Cells: cells})
	}
	return rows, errs
}

func (q *druidQuerier) convHeaders(environmentNamespace string, columns []string) *ecproto.Row {
	userDataRegex := regexp.MustCompile(userDataPattern(environmentNamespace))
	headers := &ecproto.Row{Cells: []*ecproto.Cell{}}
	for _, column := range columns {
		if userDataRegex.MatchString(column) {
			headers.Cells = append(
				headers.Cells,
				&ecproto.Cell{Type: ecproto.Cell_STRING, Value: removeEnvFromUserData(column, userDataRegex)},
			)
			continue
		}
		headers.Cells = append(headers.Cells, &ecproto.Cell{Type: ecproto.Cell_STRING, Value: column})
	}
	return headers
}

func generateColumns(featureID string, segments, valueColumns []string) []string {
	columns := make([]string, len(valueColumns))
	copy(columns, valueColumns)
	if featureID != "" {
		columns = append([]string{ColumnVariation}, valueColumns...)
	}
	return append(segments, columns...)
}

func convToResult(
	ctx context.Context,
	evalHeader, goalHeader *ecproto.Row,
	evaluationRows, goalRows []*ecproto.Row,
	segments []string,
) (*ecproto.Row, []*ecproto.Row) {
	evalVarIdx := cellIndex(evalHeader, ColumnVariation)
	evalMap := make(map[string]*ecproto.Row, len(evaluationRows))
	for _, evalRow := range evaluationRows {
		key := newKey(evalRow.Cells[evalVarIdx].Value, segmentValuesNew(evalHeader, evalRow, segments))
		evalMap[key] = evalRow
	}

	evalUUIdx := cellIndex(evalHeader, ColumnEvaluationUser)
	goalUUIdx := cellIndex(goalHeader, ColumnGoalUser)
	goalVarIdx := cellIndex(goalHeader, ColumnVariation)
	genRows := []*ecproto.Row{}
	for _, goalRow := range goalRows {
		key := newKey(goalRow.Cells[goalVarIdx].Value, segmentValuesNew(goalHeader, goalRow, segments))
		evalRow, ok := evalMap[key]
		if !ok {
			continue
		}
		cells := []*ecproto.Cell{}

		// Copy evaluation count.
		cells = append(cells, evalRow.Cells...)

		// Copy goal count.
		for i, goalCell := range goalRow.Cells {
			excludes := append(segments, ColumnVariation)
			if contains(excludes, goalHeader.Cells[i].Value) {
				continue
			}
			cells = append(cells, goalCell)
		}

		// Calculate CVR.
		var cvr float64
		if evalRow.Cells[evalUUIdx].ValueDouble > float64(0) {
			euu := evalRow.Cells[evalUUIdx].ValueDouble
			guu := goalRow.Cells[goalUUIdx].ValueDouble
			cvr = guu / euu * 100
		}
		cells = append(cells, &ecproto.Cell{Type: ecproto.Cell_DOUBLE, ValueDouble: cvr})

		genRows = append(genRows, &ecproto.Row{Cells: cells})
	}
	genHeaders := joinHeaders(evalHeader, goalHeader, segments)
	return genHeaders, genRows
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

func segmentValuesNew(header, row *ecproto.Row, segments []string) []string {
	segmentIdx := []int{}
	for _, s := range segments {
		segmentIdx = append(segmentIdx, cellIndex(header, s))
	}
	ss := []string{}
	for _, i := range segmentIdx {
		ss = append(ss, row.Cells[i].Value)
	}
	return ss
}

func newKey(variation string, segments []string) string {
	return strings.Join(append(segments, variation), ",")
}

func cellIndex(header *ecproto.Row, value string) int {
	for i, cell := range header.Cells {
		if cell.Value == value {
			return i
		}
	}
	return -1
}

func joinHeaders(evalHeader, goalHeader *ecproto.Row, segments []string) *ecproto.Row {
	newHeader := []*ecproto.Cell{}
	newHeader = append(newHeader, evalHeader.Cells...)
	excludes := append(segments, ColumnVariation)
	for _, h := range goalHeader.Cells {
		if contains(excludes, h.Value) {
			continue
		}
		newHeader = append(newHeader, h)
	}
	newHeader = append(newHeader, &ecproto.Cell{Value: ColumnCVR})
	return &ecproto.Row{Cells: newHeader}
}

func userDataPattern(environmentNamespace string) string {
	if environmentNamespace == "" {
		return `^user\.data\.(.*)$`
	}
	return fmt.Sprintf(`^%s\.user\.data\.(.*)$`, environmentNamespace)
}

func removeEnvFromUserData(key string, r *regexp.Regexp) string {
	group := r.FindStringSubmatch(key)
	return fmt.Sprintf("user.data.%s", group[1])
}

func convToEnvSegments(environmentNamespace string, segments []string) []string {
	if environmentNamespace == "" {
		return segments
	}
	userDataRegex := regexp.MustCompile(userDataPattern(""))
	sgmt := []string{}
	for _, s := range segments {
		if userDataRegex.MatchString(s) {
			sgmt = append(sgmt, fmt.Sprintf("%s.%s", environmentNamespace, s))
			continue
		}
		sgmt = append(sgmt, s)
	}
	return sgmt
}

func convToEnvFilters(environmentNamespace string, filters []*ecproto.Filter) []*ecproto.Filter {
	fls := []*ecproto.Filter{}
	userDataRegex := regexp.MustCompile(userDataPattern(""))
	for _, f := range filters {
		key := f.Key
		if environmentNamespace != "" && userDataRegex.MatchString(f.Key) {
			key = fmt.Sprintf("%s.%s", environmentNamespace, key)
		}
		switch f.Operator {
		case ecproto.Filter_EQUALS:
			fls = append(fls, &ecproto.Filter{
				Operator: ecproto.Filter_EQUALS,
				Key:      key,
				Values:   f.Values,
			})
		}
	}
	return fls
}

func (q *druidQuerier) convToTimeseries(
	variationID string,
	godruidTS []godruid.Timeseries,
	columns []string,
) (map[string]*ecproto.VariationTimeseries, []error) {
	ts := []int64{}
	values := map[string][]float64{}
	errs := []error{}
	for _, gdts := range godruidTS {
		t, err := time.Parse(time.RFC3339, gdts.Timestamp)
		if err != nil {
			errs = append(errs, fmt.Errorf("time %s cannot be parsed", gdts.Timestamp))
			return nil, errs
		}
		ts = append(ts, t.Unix())
		for _, column := range columns {
			value, ok := gdts.Result[column]
			if !ok {
				errs = append(errs, fmt.Errorf("Column %s does not exist", column))
				values[column] = append(values[column], 0)
				continue
			}
			switch value := value.(type) {
			case float64:
				values[column] = append(values[column], value)
			default:
				values[column] = append(values[column], 0)
				errs = append(errs, fmt.Errorf("Value %s type is unknown", value))
			}
		}
	}
	variationTS := map[string]*ecproto.VariationTimeseries{}
	for column, vals := range values {
		variationTS[column] = &ecproto.VariationTimeseries{
			VariationId: variationID,
			Timeseries: &ecproto.Timeseries{
				Timestamps: ts,
				Values:     vals,
			},
		}
	}
	return variationTS, errs
}
