import test from 'ava';
import { Feature } from '../../proto/feature/feature_pb';
import { Evaluator } from '../../evaluation';
import { SegmentUser } from '../../proto/feature/segment_pb';
import { UserEvaluations } from '../../proto/feature/evaluation_pb';
import { createEvaluation, createPrerequisite, createReason, createUser } from '../../modelFactory';
import { Reason } from '../../proto/feature/reason_pb';
import { NewUserEvaluations } from '../../userEvaluation';
import { newTestFeature } from './evaluate_feature_test';


var allFeaturesForPrerequisiteTest = {
  "featureA": create
	"featureA": {
		Id:   "featureA",
		Name: "featureA",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureE",
			},
			{
				FeatureId: "featureF",
			},
		},
	},
	"featureB": {
		Id:   "featureB",
		Name: "featureB",
	},
	"featureC": {
		Id:   "featureC",
		Name: "featureC",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureL",
			},
		},
	},
	"featureD": {
		Id:   "featureD",
		Name: "featureD",
	},
	"featureE": {
		Id:   "featureE",
		Name: "featureE",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureG",
			},
		},
	},
	"featureF": {
		Id:   "featureF",
		Name: "featureF",
	},
	"featureG": {
		Id:   "featureG",
		Name: "featureG",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureH",
			},
		},
	},
	"featureH": {
		Id:   "featureH",
		Name: "featureH",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureI",
			},
			{
				FeatureId: "featureJ",
			},
		},
	},
	"featureI": {
		Id:   "featureI",
		Name: "featureI",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureK",
			},
		},
	},
	"featureJ": {
		Id:   "featureJ",
		Name: "featureJ",
	},
	"featureK": {
		Id:   "featureK",
		Name: "featureK",
	},
	"featureL": {
		Id:   "featureL",
		Name: "featureL",
		Prerequisites: []*ftproto.Prerequisite{
			{
				FeatureId: "featureM",
			},
			{
				FeatureId: "featureN",
			},
		},
	},
	"featureM": {
		Id:   "featureM",
		Name: "featureM",
	},
	"featureN": {
		Id:   "featureN",
		Name: "featureN",
	},
}