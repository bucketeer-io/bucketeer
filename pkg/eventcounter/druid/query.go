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

package druid

import (
	"fmt"
	"time"

	"github.com/ca-dp/godruid"

	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

const (
	intervalStr = "2006-01-02T15:04"
)

func querySegmentMetadata(datasource string, startAt, endAt time.Time) *godruid.QuerySegmentMetadata {
	query := &godruid.QuerySegmentMetadata{
		QueryType:  godruid.SEGMENTMETADATA,
		DataSource: godruid.DataSourceTable(datasource),
		Intervals:  toConvInterval(startAt, endAt),
		Merge:      "true",
	}
	return query
}

func toConvInterval(start, end time.Time) string {
	return fmt.Sprintf("%s/%s", start.Format(intervalStr), end.Format(intervalStr))
}

func queryGoalGroupBy(
	datasource string,
	startAt, endAt time.Time,
	environmentNamespace string,
	goalID string,
	featureID string,
	featureVersion int32,
	reason string,
	segments []string,
	fls []*ecproto.Filter,
) *godruid.QueryGroupBy {
	filters := []*godruid.Filter{}
	innerDimensions := []godruid.DimSpec{}
	outerDimensions := []godruid.DimSpec{}
	limitColumns := []godruid.Column{}
	filters = append(filters, godruid.FilterSelector("environmentNamespace", environmentNamespace))
	filters = append(filters, godruid.FilterSelector("goalId", goalID))
	filters = append(filters, convToDruidFilters(fls)...)
	innerDimensions = append(innerDimensions, godruid.DimDefault("userId", ColumnUser))
	for _, segment := range segments {
		limitColumns = append(limitColumns, godruid.Column{Dimension: segment, Direction: godruid.DirectionASC})
		innerDimensions = append(innerDimensions, godruid.DimDefault(segment, segment))
		outerDimensions = append(outerDimensions, godruid.DimDefault(segment, segment))
	}
	if featureID != "" {
		filterEvaluationPattern := ""
		if reason != "" {
			filterEvaluationPattern = fmt.Sprintf("^%s:%d:.*:%s$", featureID, featureVersion, reason)
		} else {
			filterEvaluationPattern = fmt.Sprintf("^%s:%d:.*$", featureID, featureVersion)
		}
		innerEvaluationPattern := fmt.Sprintf("^%s:%d:.*$", featureID, featureVersion)
		filters = append(filters, godruid.FilterRegex("evaluations", filterEvaluationPattern))
		innerDimensions = append(innerDimensions, evaluationsDim(innerEvaluationPattern))
		outerDimensions = append(outerDimensions, godruid.DimDefault(ColumnVariation, ColumnVariation))
		limitColumns = append(limitColumns, godruid.Column{Dimension: ColumnVariation, Direction: godruid.DirectionASC})
	}
	innerQuery := &godruid.QueryGroupBy{
		QueryType:   godruid.GROUPBY,
		DataSource:  godruid.DataSourceTable(datasource),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranAll,
		Filter:      godruid.FilterAnd(filters...),
		Dimensions:  innerDimensions,
		Aggregations: []godruid.Aggregation{
			*godruid.AggLongSum("count", "count"),
			*godruid.AggDoubleSum("valueSum", "valueSum"),
		},
	}
	query := &godruid.QueryGroupBy{
		QueryType:   "groupBy",
		DataSource:  godruid.DataSourceQuery(innerQuery),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranAll,
		LimitSpec:   godruid.LimitDefault(10000, limitColumns),
		Dimensions:  outerDimensions,
		Aggregations: []godruid.Aggregation{
			*godruid.AggLongSum(ColumnGoalTotal, "count"),
			*godruid.AggDoubleSum(ColumnGoalValueTotal, "valueSum"),
			*godruid.AggCount(ColumnGoalUser),
			*godruid.AggRawJson(
				fmt.Sprintf(`{ "type": "doubleMean", "name": "%s", "fieldName": "valueSum" }`, ColumnGoalValueMean),
			),
			*godruid.ExtAggVariance(ColumnGoalValueVariance, "valueSum"),
		},
	}
	return query
}

func queryEvaluationGroupBy(
	datasource string,
	startAt, endAt time.Time,
	environmentNamespace string,
	featureID string,
	featureVersion int32,
	reason string,
	segments []string,
	fls []*ecproto.Filter,
) *godruid.QueryGroupBy {

	filters := []*godruid.Filter{}
	dimensions := []godruid.DimSpec{}
	limitColumns := []godruid.Column{}
	for _, segment := range segments {
		limitColumns = append(limitColumns, godruid.Column{Dimension: segment, Direction: godruid.DirectionASC})
		dimensions = append(dimensions, godruid.DimDefault(segment, segment))
	}
	filters = append(filters, godruid.FilterSelector("environmentNamespace", environmentNamespace))
	filters = append(filters, godruid.FilterSelector("featureId", featureID))
	if featureVersion == 0 {
		dimensions = append(dimensions, godruid.DimDefault("featureVersion", ColumnFeatureVersion))
	} else {
		filters = append(filters, godruid.FilterSelector("featureVersion", featureVersion))
	}
	if reason != "" {
		filters = append(filters, godruid.FilterSelector("reason", reason))
	}
	filters = append(filters, convToDruidFilters(fls)...)
	dimensions = append(dimensions, godruid.DimDefault("variationId", ColumnVariation))
	limitColumns = append(limitColumns, godruid.Column{Dimension: ColumnVariation, Direction: godruid.DirectionASC})
	query := &godruid.QueryGroupBy{
		QueryType:   godruid.GROUPBY,
		DataSource:  godruid.DataSourceTable(datasource),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranAll,
		Filter:      godruid.FilterAnd(filters...),
		LimitSpec:   godruid.LimitDefault(10000, limitColumns),
		Dimensions:  dimensions,
		Aggregations: []godruid.Aggregation{
			*godruid.AggLongSum(ColumnEvaluationTotal, "count"),
			*godruid.AggRawJson(
				fmt.Sprintf(`{ "type": "thetaSketch", "name": "%s", "fieldName": "userIdThetaSketch" }`, ColumnEvaluationUser),
			),
		},
	}
	return query
}

func queryUserGroupBy(datasource, environmentNamespace string, startAt, endAt time.Time) *godruid.QueryGroupBy {
	countFieldName := "count"
	query := &godruid.QueryGroupBy{
		QueryType:   godruid.GROUPBY,
		DataSource:  godruid.DataSourceTable(datasource),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranAll,
		Filter:      godruid.FilterAnd(godruid.FilterSelector("environmentNamespace", environmentNamespace)),
		Dimensions:  []godruid.DimSpec{godruid.DimDefault("userId", ColumnUser)},
		Aggregations: []godruid.Aggregation{
			*godruid.AggCount(countFieldName),
		},
	}
	query = &godruid.QueryGroupBy{
		QueryType:   godruid.GROUPBY,
		DataSource:  godruid.DataSourceQuery(query),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranAll,
		Aggregations: []godruid.Aggregation{
			*godruid.AggLongSum(ColumnUserTotal, countFieldName),
			*godruid.AggCount(ColumnUserCount),
		},
	}
	return query
}

func evaluationsDim(evaluationPattern string) *godruid.DimFiltered {
	variationExFn := godruid.DimExFnRegex("^(.*):.*$")
	variationDelegate := godruid.DimExtraction("evaluations", ColumnVariation, variationExFn)
	return godruid.DimFilteredRegex(variationDelegate, evaluationPattern)
}

func queryEvaluationTimeseries(
	datasource string,
	startAt, endAt time.Time,
	environmentNamespace string,
	featureID string,
	featureVersion int32,
	variationID string,
) *godruid.QueryTimeseries {
	filters := []*godruid.Filter{}
	filters = append(filters, godruid.FilterSelector("environmentNamespace", environmentNamespace))
	filters = append(filters, godruid.FilterSelector("featureId", featureID))
	if featureVersion != 0 {
		filters = append(filters, godruid.FilterSelector("featureVersion", featureVersion))
	}
	filters = append(filters, godruid.FilterSelector("variationId", variationID))
	query := &godruid.QueryTimeseries{
		QueryType:   godruid.TIMESERIES,
		DataSource:  godruid.DataSourceTable(datasource),
		Intervals:   toConvInterval(startAt, endAt),
		Granularity: godruid.GranPeriod("P1D", "Asia/Tokyo", ""),
		Filter:      godruid.FilterAnd(filters...),
		Aggregations: []godruid.Aggregation{
			*godruid.AggLongSum(ColumnEvaluationTotal, "count"),
			*godruid.AggRawJson(
				fmt.Sprintf(`{ "type": "thetaSketch", "name": "%s", "fieldName": "userIdThetaSketch" }`, ColumnEvaluationUser),
			),
		},
	}
	return query
}

func convToDruidFilters(filters []*ecproto.Filter) []*godruid.Filter {
	fls := []*godruid.Filter{}
	for _, f := range filters {
		switch f.Operator {
		case ecproto.Filter_EQUALS:
			for _, v := range f.Values {
				fls = append(fls, godruid.FilterSelector(f.Key, v))
			}
		}
	}
	return fls
}
