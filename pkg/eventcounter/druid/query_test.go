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
	"testing"
	"time"

	"github.com/ca-dp/godruid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ecproto "github.com/bucketeer-io/bucketeer/proto/eventcounter"
)

func TestQuerySegmentMetadata(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)

	patterns := map[string]struct {
		inputDatasource string
		inputStartAt    time.Time
		inputEndAt      time.Time
		expected        *godruid.QuerySegmentMetadata
	}{
		"success": {
			inputDatasource: "ds",
			inputStartAt:    t1,
			inputEndAt:      t2,
			expected: &godruid.QuerySegmentMetadata{
				QueryType:  godruid.SEGMENTMETADATA,
				DataSource: godruid.DataSourceTable("ds"),
				Intervals:  "2014-01-17T23:02/2014-01-18T23:02",
				Merge:      "true",
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := querySegmentMetadata(p.inputDatasource, p.inputStartAt, p.inputEndAt)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestQueryGoalGroupBy(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)

	patterns := map[string]struct {
		inputDatasource           string
		inputStartAt              time.Time
		inputEndAt                time.Time
		inputEnvironmentNamespace string
		inputGoalID               string
		inputFeatureID            string
		inputFeatureVersion       int32
		inputReason               string
		inputSegments             []string
		inputFilters              []*ecproto.Filter
		expected                  *godruid.QueryGroupBy
	}{
		"no feature, no segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			expected: &godruid.QueryGroupBy{
				QueryType: "groupBy",
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
							&godruid.Filter{Type: "selector", Dimension: "goalId", Value: "gid"},
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(`{ "type": "longSum", "name": "count", "fieldName": "count" }`),
							*godruid.AggRawJson(`{ "type": "doubleSum", "name": "valueSum", "fieldName": "valueSum" }`),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				LimitSpec:   godruid.LimitDefault(10000, []godruid.Column{}),
				Dimensions:  []godruid.DimSpec{},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Goal total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "doubleSum", "name": "Goal value total", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "count", "name": "Goal user"}`),
					*godruid.AggRawJson(`{ "type": "doubleMean", "name": "Goal value mean", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "variance", "name": "Goal value variance", "fieldName": "valueSum" }`),
				},
			},
		},
		"no feature, segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			inputSegments:             []string{"s1", "s2"},
			expected: &godruid.QueryGroupBy{
				QueryType: "groupBy",
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
							&godruid.Filter{Type: "selector", Dimension: "goalId", Value: "gid"},
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
							godruid.DimDefault("s1", "s1"),
							godruid.DimDefault("s2", "s2"),
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(`{ "type": "longSum", "name": "count", "fieldName": "count" }`),
							*godruid.AggRawJson(`{ "type": "doubleSum", "name": "valueSum", "fieldName": "valueSum" }`),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "s1", Direction: godruid.DirectionASC},
					{Dimension: "s2", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("s1", "s1"),
					godruid.DimDefault("s2", "s2"),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Goal total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "doubleSum", "name": "Goal value total", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "count", "name": "Goal user"}`),
					*godruid.AggRawJson(`{ "type": "doubleMean", "name": "Goal value mean", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "variance", "name": "Goal value variance", "fieldName": "valueSum" }`),
				},
			},
		},
		"feature, no segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(1),
			expected: &godruid.QueryGroupBy{
				QueryType: "groupBy",
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
							&godruid.Filter{Type: "selector", Dimension: "goalId", Value: "gid"},
							&godruid.Filter{Type: "regex", Dimension: "evaluations", Pattern: "^fid:1:.*$"},
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
							&godruid.DimFiltered{
								Type:    "regexFiltered",
								Pattern: "^fid:1:.*$",
								Delegate: &godruid.Dimension{
									Type:       "extraction",
									Dimension:  "evaluations",
									OutputName: "Variation",
									ExtractionFn: &godruid.DimExtractionFn{
										Type: "regex",
										Expr: "^(.*):.*$",
									},
								},
							},
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(`{ "type": "longSum", "name": "count", "fieldName": "count" }`),
							*godruid.AggRawJson(`{ "type": "doubleSum", "name": "valueSum", "fieldName": "valueSum" }`),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				LimitSpec:   godruid.LimitDefault(10000, []godruid.Column{{Dimension: "Variation", Direction: godruid.DirectionASC}}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault(ColumnVariation, ColumnVariation),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Goal total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "doubleSum", "name": "Goal value total", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "count", "name": "Goal user"}`),
					*godruid.AggRawJson(`{ "type": "doubleMean", "name": "Goal value mean", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "variance", "name": "Goal value variance", "fieldName": "valueSum" }`),
				},
			},
		},
		"feature, segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			inputSegments:             []string{"s1", "s2"},
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(1),
			inputFilters: []*ecproto.Filter{
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v1"}},
				{Key: "f1", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
			},
			expected: &godruid.QueryGroupBy{
				QueryType: "groupBy",
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
							&godruid.Filter{Type: "selector", Dimension: "goalId", Value: "gid"},
							&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v0"},
							&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v1"},
							&godruid.Filter{Type: "selector", Dimension: "f1", Value: "v0"},
							&godruid.Filter{Type: "regex", Dimension: "evaluations", Pattern: "^fid:1:.*$"},
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
							godruid.DimDefault("s1", "s1"),
							godruid.DimDefault("s2", "s2"),
							&godruid.DimFiltered{
								Type:    "regexFiltered",
								Pattern: "^fid:1:.*$",
								Delegate: &godruid.Dimension{
									Type:       "extraction",
									Dimension:  "evaluations",
									OutputName: "Variation",
									ExtractionFn: &godruid.DimExtractionFn{
										Type: "regex",
										Expr: "^(.*):.*$",
									},
								},
							},
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(`{ "type": "longSum", "name": "count", "fieldName": "count" }`),
							*godruid.AggRawJson(`{ "type": "doubleSum", "name": "valueSum", "fieldName": "valueSum" }`),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "s1", Direction: godruid.DirectionASC},
					{Dimension: "s2", Direction: godruid.DirectionASC},
					{Dimension: "Variation", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("s1", "s1"),
					godruid.DimDefault("s2", "s2"),
					godruid.DimDefault(ColumnVariation, ColumnVariation),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Goal total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "doubleSum", "name": "Goal value total", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "count", "name": "Goal user"}`),
					*godruid.AggRawJson(`{ "type": "doubleMean", "name": "Goal value mean", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "variance", "name": "Goal value variance", "fieldName": "valueSum" }`),
				},
			},
		},
		"feature, reason, segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			inputSegments:             []string{"s1", "s2"},
			inputFeatureID:            "fid",
			inputReason:               "DEFAULT",
			inputFeatureVersion:       int32(1),
			inputFilters: []*ecproto.Filter{
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v1"}},
				{Key: "f1", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
			},
			expected: &godruid.QueryGroupBy{
				QueryType: "groupBy",
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
							&godruid.Filter{Type: "selector", Dimension: "goalId", Value: "gid"},
							&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v0"},
							&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v1"},
							&godruid.Filter{Type: "selector", Dimension: "f1", Value: "v0"},
							&godruid.Filter{Type: "regex", Dimension: "evaluations", Pattern: "^fid:1:.*:DEFAULT$"},
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
							godruid.DimDefault("s1", "s1"),
							godruid.DimDefault("s2", "s2"),
							&godruid.DimFiltered{
								Type:    "regexFiltered",
								Pattern: "^fid:1:.*$",
								Delegate: &godruid.Dimension{
									Type:       "extraction",
									Dimension:  "evaluations",
									OutputName: "Variation",
									ExtractionFn: &godruid.DimExtractionFn{
										Type: "regex",
										Expr: "^(.*):.*$",
									},
								},
							},
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(`{ "type": "longSum", "name": "count", "fieldName": "count" }`),
							*godruid.AggRawJson(`{ "type": "doubleSum", "name": "valueSum", "fieldName": "valueSum" }`),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "s1", Direction: godruid.DirectionASC},
					{Dimension: "s2", Direction: godruid.DirectionASC},
					{Dimension: "Variation", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("s1", "s1"),
					godruid.DimDefault("s2", "s2"),
					godruid.DimDefault(ColumnVariation, ColumnVariation),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Goal total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "doubleSum", "name": "Goal value total", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "count", "name": "Goal user"}`),
					*godruid.AggRawJson(`{ "type": "doubleMean", "name": "Goal value mean", "fieldName": "valueSum" }`),
					*godruid.AggRawJson(`{ "type": "variance", "name": "Goal value variance", "fieldName": "valueSum" }`),
				},
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := queryGoalGroupBy(
				p.inputDatasource,
				p.inputStartAt,
				p.inputEndAt,
				p.inputEnvironmentNamespace,
				p.inputGoalID,
				p.inputFeatureID,
				p.inputFeatureVersion,
				p.inputReason,
				p.inputSegments,
				p.inputFilters)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestQueryEvaluationGroupBy(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)

	patterns := map[string]struct {
		inputDatasource           string
		inputStartAt              time.Time
		inputEndAt                time.Time
		inputEnvironmentNamespace string
		inputGoalID               string
		inputFeatureID            string
		inputFeatureVersion       int32
		inputReason               string
		inputSegments             []string
		inputFilters              []*ecproto.Filter
		expected                  *godruid.QueryGroupBy
	}{
		"feature, no segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputGoalID:               "gid",
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(1),
			expected: &godruid.QueryGroupBy{
				QueryType:   "groupBy",
				DataSource:  godruid.DataSourceTable("ds"),
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				Filter: godruid.FilterAnd(
					godruid.FilterSelector("environmentNamespace", "ns"),
					godruid.FilterSelector("featureId", "fid"),
					godruid.FilterSelector("featureVersion", int32(1)),
				),
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "Variation", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("variationId", "Variation"),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Evaluation total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "thetaSketch", "name": "Evaluation user", "fieldName": "userIdThetaSketch" }`),
				},
			},
		},
		"feature, segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputSegments:             []string{"s1", "s2"},
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(1),
			inputFilters: []*ecproto.Filter{
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v1"}},
				{Key: "f1", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
			},
			expected: &godruid.QueryGroupBy{
				QueryType:   "groupBy",
				DataSource:  godruid.DataSourceTable("ds"),
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				Filter: godruid.FilterAnd(
					godruid.FilterSelector("environmentNamespace", "ns"),
					godruid.FilterSelector("featureId", "fid"),
					godruid.FilterSelector("featureVersion", int32(1)),
					&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v0"},
					&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v1"},
					&godruid.Filter{Type: "selector", Dimension: "f1", Value: "v0"},
				),
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "s1", Direction: godruid.DirectionASC},
					{Dimension: "s2", Direction: godruid.DirectionASC},
					{Dimension: "Variation", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("s1", "s1"),
					godruid.DimDefault("s2", "s2"),
					godruid.DimDefault("variationId", "Variation"),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Evaluation total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "thetaSketch", "name": "Evaluation user", "fieldName": "userIdThetaSketch" }`),
				},
			},
		},
		"feature, reason, segments": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputSegments:             []string{"s1", "s2"},
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(1),
			inputReason:               "DEFAULT",
			inputFilters: []*ecproto.Filter{
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
				{Key: "f0", Operator: ecproto.Filter_EQUALS, Values: []string{"v1"}},
				{Key: "f1", Operator: ecproto.Filter_EQUALS, Values: []string{"v0"}},
			},
			expected: &godruid.QueryGroupBy{
				QueryType:   "groupBy",
				DataSource:  godruid.DataSourceTable("ds"),
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				Filter: godruid.FilterAnd(
					godruid.FilterSelector("environmentNamespace", "ns"),
					godruid.FilterSelector("featureId", "fid"),
					godruid.FilterSelector("featureVersion", int32(1)),
					godruid.FilterSelector("reason", "DEFAULT"),
					&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v0"},
					&godruid.Filter{Type: "selector", Dimension: "f0", Value: "v1"},
					&godruid.Filter{Type: "selector", Dimension: "f1", Value: "v0"},
				),
				LimitSpec: godruid.LimitDefault(10000, []godruid.Column{
					{Dimension: "s1", Direction: godruid.DirectionASC},
					{Dimension: "s2", Direction: godruid.DirectionASC},
					{Dimension: "Variation", Direction: godruid.DirectionASC},
				}),
				Dimensions: []godruid.DimSpec{
					godruid.DimDefault("s1", "s1"),
					godruid.DimDefault("s2", "s2"),
					godruid.DimDefault("variationId", "Variation"),
				},
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(`{ "type": "longSum", "name": "Evaluation total", "fieldName": "count" }`),
					*godruid.AggRawJson(`{ "type": "thetaSketch", "name": "Evaluation user", "fieldName": "userIdThetaSketch" }`),
				},
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := queryEvaluationGroupBy(
				p.inputDatasource,
				p.inputStartAt,
				p.inputEndAt,
				p.inputEnvironmentNamespace,
				p.inputFeatureID,
				p.inputFeatureVersion,
				p.inputReason,
				p.inputSegments,
				p.inputFilters,
			)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestQueryUserGroupBy(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)
	patterns := map[string]struct {
		inputDatasource           string
		inputStartAt              time.Time
		inputEndAt                time.Time
		inputEnvironmentNamespace string
		inputFilters              []*ecproto.Filter
		expected                  *godruid.QueryGroupBy
	}{
		"success": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			expected: &godruid.QueryGroupBy{
				QueryType: godruid.GROUPBY,
				DataSource: &godruid.DataSource{
					Type: "query",
					Query: &godruid.QueryGroupBy{
						QueryType:   godruid.GROUPBY,
						DataSource:  godruid.DataSourceTable("ds"),
						Intervals:   toConvInterval(t1, t2),
						Granularity: godruid.GranAll,
						Filter: godruid.FilterAnd(
							godruid.FilterSelector("environmentNamespace", "ns"),
						),
						Dimensions: []godruid.DimSpec{
							godruid.DimDefault("userId", ColumnUser),
						},
						Aggregations: []godruid.Aggregation{
							*godruid.AggRawJson(fmt.Sprintf(`{ "type": "count", "name": "count" }`)),
						},
					},
				},
				Intervals:   "2014-01-17T23:02/2014-01-18T23:02",
				Granularity: godruid.GranAll,
				Aggregations: []godruid.Aggregation{
					*godruid.AggRawJson(fmt.Sprintf(`{ "type": "longSum", "name": "%s", "fieldName": "count" }`, ColumnUserTotal)),
					*godruid.AggRawJson(fmt.Sprintf(`{ "type": "count", "name": "%s" }`, ColumnUserCount)),
				},
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := queryUserGroupBy(
				p.inputDatasource,
				p.inputEnvironmentNamespace,
				p.inputStartAt,
				p.inputEndAt,
			)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestQueryEvaluationTimeseries(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	layout := "2006-01-02 15:04:05 -0700 MST"
	t1, err := time.Parse(layout, "2014-01-17 23:02:03 +0000 UTC")
	require.NoError(t, err)
	t2, err := time.Parse(layout, "2014-01-18 23:02:03 +0000 UTC")
	require.NoError(t, err)
	patterns := map[string]struct {
		inputDatasource           string
		inputStartAt              time.Time
		inputEndAt                time.Time
		inputEnvironmentNamespace string
		inputFeatureID            string
		inputFeatureVersion       int32
		inputVariationID          string
		expected                  *godruid.QueryTimeseries
	}{
		"success": {
			inputDatasource:           "ds",
			inputStartAt:              t1,
			inputEndAt:                t2,
			inputEnvironmentNamespace: "ns",
			inputFeatureID:            "fid",
			inputFeatureVersion:       int32(0),
			inputVariationID:          "vid",
			expected: &godruid.QueryTimeseries{
				QueryType:   godruid.TIMESERIES,
				DataSource:  godruid.DataSourceTable("ds"),
				Intervals:   toConvInterval(t1, t2),
				Granularity: godruid.GranPeriod("P1D", "Asia/Tokyo", ""),
				Filter: godruid.FilterAnd(
					godruid.FilterSelector("environmentNamespace", "ns"),
					godruid.FilterSelector("featureId", "fid"),
					godruid.FilterSelector("variationId", "vid"),
				),
				Aggregations: []godruid.Aggregation{
					*godruid.AggLongSum(ColumnEvaluationTotal, "count"),
					*godruid.AggRawJson(
						fmt.Sprintf(`{ "type": "thetaSketch", "name": "%s", "fieldName": "userIdThetaSketch" }`, ColumnEvaluationUser),
					),
				},
			},
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			actual := queryEvaluationTimeseries(
				p.inputDatasource,
				p.inputStartAt,
				p.inputEndAt,
				p.inputEnvironmentNamespace,
				p.inputFeatureID,
				p.inputFeatureVersion,
				p.inputVariationID,
			)
			assert.Equal(t, p.expected, actual)
		})
	}
}
