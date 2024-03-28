// Copyright 2024 The Bucketeer Authors.
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

package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	gwapi "github.com/bucketeer-io/bucketeer/pkg/gateway/api"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
	gwproto "github.com/bucketeer-io/bucketeer/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/proto/user"
)

const (
	version          = "/v1"
	service          = "/gateway"
	evaluationsAPI   = "/evaluations"
	evaluationAPI    = "/evaluation"
	eventsAPI        = "/events"
	authorizationKey = "authorization"
)

type successResponse struct {
	Data json.RawMessage `json:"data"`
}

type registerEventsRequest struct {
	Events []Event `json:"events,omitempty"`
}

type registerEventsResponse struct {
	Errors map[string]*gwproto.RegisterEventsResponse_Error `json:"errors,omitempty"`
}

type getEvaluationsRequest struct {
	Tag               string              `json:"tag,omitempty"`
	User              *userproto.User     `json:"user,omitempty"`
	UserEvaluationsID string              `json:"user_evaluations_id,omitempty"`
	SourceID          eventproto.SourceId `json:"source_id,omitempty"`
}

type getEvaluationsResponse struct {
	Evaluations       *featureproto.UserEvaluations `json:"evaluations,omitempty"`
	UserEvaluationsID string                        `json:"user_evaluations_id,omitempty"`
}

type getEvaluationRequest struct {
	Tag       string              `json:"tag,omitempty"`
	User      *userproto.User     `json:"user,omitempty"`
	FeatureID string              `json:"feature_id,omitempty"`
	SourceId  eventproto.SourceId `json:"source_id,omitempty"`
}

type getEvaluationResponse struct {
	Evaluation *featureproto.Evaluation `json:"evaluations,omitempty"`
}

type Event struct {
	ID                   string          `json:"id,omitempty"`
	Event                json.RawMessage `json:"event,omitempty"`
	EnvironmentNamespace string          `json:"environment_namespace,omitempty"`
	Type                 gwapi.EventType `json:"type,omitempty"`
}

func GetEvaluations(t *testing.T, tag, userID, gatewayAddr, apiKeyPath string) *getEvaluationsResponse {
	t.Helper()
	url := fmt.Sprintf("https://%s%s%s%s",
		gatewayAddr,
		version,
		service,
		evaluationsAPI,
	)
	req := &getEvaluationsRequest{
		Tag: tag,
		User: &userproto.User{
			Id: userID,
		},
	}
	resp := SendHTTPRequest(t, url, req, apiKeyPath)
	var ger getEvaluationsResponse
	if err := json.Unmarshal(resp.Data, &ger); err != nil {
		t.Fatal(err)
	}
	return &ger
}

func GetEvaluation(t *testing.T, tag, featureID, userID, gatewayAddr, apiKeyPath string) *getEvaluationResponse {
	t.Helper()
	url := fmt.Sprintf("https://%s%s%s%s",
		gatewayAddr,
		version,
		service,
		evaluationAPI,
	)
	req := &getEvaluationRequest{
		Tag:       tag,
		User:      &userproto.User{Id: userID},
		FeatureID: featureID,
	}
	resp := SendHTTPRequest(t, url, req, apiKeyPath)
	var ger getEvaluationResponse
	if err := json.Unmarshal(resp.Data, &ger); err != nil {
		t.Fatal(err)
	}
	return &ger
}

func RegisterEvents(t *testing.T, events []Event, gatewayAddr, apiKeyPath string) *registerEventsResponse {
	t.Helper()
	url := fmt.Sprintf("https://%s%s%s%s",
		gatewayAddr,
		version,
		service,
		eventsAPI,
	)
	req := &registerEventsRequest{
		Events: events,
	}
	resp := SendHTTPRequest(t, url, req, apiKeyPath)
	var rer registerEventsResponse
	if err := json.Unmarshal(resp.Data, &rer); err != nil {
		t.Fatal(err)
	}
	return &rer
}

func SendHTTPRequest(t *testing.T, url string, body interface{}, apiKeyPath string) *successResponse {
	data, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(encoded))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(authorizationKey, string(data))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Send HTTP request failed: %d", resp.StatusCode)
	}
	var sr successResponse
	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		t.Fatal(err)
	}
	return &sr
}
