package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/bucketeer/proto/event/client"
	"github.com/bucketeer-io/bucketeer/proto/gateway"
)

func TestSplitFullMethodName(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		input    string
		service  string
		method   string
		expected bool
	}{
		"valid": {
			input:   "/bucketeer.gateway.Gateway/GetEvaluations",
			service: "bucketeer.gateway.Gateway",
			method:  "GetEvaluations",
		},
		"invalid": {
			input:   "bucketeer.gateway.Gateway/GetEvaluations",
			service: "unknown",
			method:  "unknown",
		},
		"short": {
			input:   "/GetEvaluations",
			service: "unknown",
			method:  "unknown",
		},
		"empty": {
			input:   "",
			service: "unknown",
			method:  "unknown",
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			service, method := splitFullMethodName(p.input)
			assert.Equal(t, p.service, service)
			assert.Equal(t, p.method, method)
		})
	}
}

func TestExtractRequestLabels(t *testing.T) {
	t.Parallel()
	patterns := map[string]struct {
		methodName    string
		req           interface{}
		sourceID      string
		sdkVersion    string
		tag           string
		shouldCheck   bool
		expectedPanic bool
	}{
		"GetEvaluations": {
			methodName: "GetEvaluations",
			req: &gateway.GetEvaluationsRequest{
				SourceId:   client.SourceId_GO_SERVER,
				SdkVersion: "1.0.0",
				Tag:        "test",
			},
			sourceID:   "GO_SERVER",
			sdkVersion: "1.0.0",
			tag:        "test",
		},
		"GetEvaluation": {
			methodName: "GetEvaluation",
			req: &gateway.GetEvaluationRequest{
				SourceId:   client.SourceId_ANDROID,
				SdkVersion: "2.0.0",
				Tag:        "test-android",
			},
			sourceID:   "ANDROID",
			sdkVersion: "2.0.0",
			tag:        "test-android",
		},
		"GetFeatureFlags": {
			methodName: "GetFeatureFlags",
			req: &gateway.GetFeatureFlagsRequest{
				SourceId:   client.SourceId_IOS,
				SdkVersion: "3.0.0",
			},
			sourceID:   "IOS",
			sdkVersion: "3.0.0",
			tag:        "",
		},
		"GetSegmentUsers": {
			methodName: "GetSegmentUsers",
			req: &gateway.GetSegmentUsersRequest{
				SourceId:   client.SourceId_WEB,
				SdkVersion: "4.0.0",
			},
			sourceID:   "WEB",
			sdkVersion: "4.0.0",
			tag:        "",
		},
		"RegisterEvents": {
			methodName: "RegisterEvents",
			req: &gateway.RegisterEventsRequest{
				SourceId:   client.SourceId_NODE_SERVER,
				SdkVersion: "5.0.0",
			},
			sourceID:   "NODE_SERVER",
			sdkVersion: "5.0.0",
			tag:        "",
		},
		"UnknownMethod": {
			methodName: "Unknown",
			req:        &gateway.GetEvaluationsRequest{},
			sourceID:   "",
			sdkVersion: "",
			tag:        "",
		},
		"WrongRequestType": {
			methodName: "GetEvaluations",
			req:        &gateway.GetEvaluationRequest{},
			sourceID:   "",
			sdkVersion: "",
			tag:        "",
		},
	}
	for msg, p := range patterns {
		t.Run(msg, func(t *testing.T) {
			sourceID, sdkVersion, tag := extractRequestLabels(p.methodName, p.req)
			assert.Equal(t, p.sourceID, sourceID)
			assert.Equal(t, p.sdkVersion, sdkVersion)
			assert.Equal(t, p.tag, tag)
		})
	}
}
