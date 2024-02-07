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
//

package stan

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/log"
)

var (
	//go:embed experiment.stan
	stanProgramCode string
)

const (
	HmcNUTSFunction = "stan::services::sample::hmc_nuts_diag_e_adapt"
)

type Stan struct {
	client *resty.Client

	logger *zap.Logger
}

type ModelCompileReq struct {
	ProgramCode string `json:"program_code"`
}

type ModelCompileResp struct {
	Name           string `json:"name"`
	CompilerOutput string `json:"compiler_output"`
	StancWarnings  string `json:"stanc_warnings"`
}

type CreateFitReq struct {
	Chain          int                    `json:"chain,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	Delta          int                    `json:"delta,omitempty"`
	Function       string                 `json:"function,omitempty"`
	Gamma          int                    `json:"gamma,omitempty"`
	Init           map[string]interface{} `json:"init,omitempty"`
	InitBuffer     int                    `json:"init_buffer,omitempty"`
	InitRadius     int                    `json:"init_radius,omitempty"`
	Kappa          int                    `json:"kappa,omitempty"`
	MaxDepth       int                    `json:"max_depth,omitempty"`
	NumSamples     int                    `json:"num_samples,omitempty"`
	NumThin        int                    `json:"num_thin,omitempty"`
	NumWarmup      int                    `json:"num_warmup,omitempty"`
	RandomSeed     int                    `json:"random_seed,omitempty"`
	Refresh        int                    `json:"refresh,omitempty"`
	SaveWarmup     bool                   `json:"save_warmup,omitempty"`
	Stepsize       int                    `json:"stepsize,omitempty"`
	StepsizeJitter int                    `json:"stepsize_jitter,omitempty"`
	T0             int                    `json:"t0,omitempty"`
	TermBuffer     int                    `json:"term_buffer,omitempty"`
	Window         int                    `json:"window,omitempty"`
}

type CreateFitResp struct {
	Name string `json:"name"`
}

type CreateFitErrResp struct {
	Code    int `json:"code"`
	Details []struct {
	} `json:"details"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// GetOperationResp example json:
//
//	{
//	 "name": "operations/cpyhg3rz",
//	 "done": true,
//	 "metadata": {
//	   "fit": {
//	     "name": "models/nld3pk7n/fits/cpyhg3rz"
//	   },
//	   "progress": "Iteration: 22000 / 22000 [100%]  (Sampling)"
//	 },
//	 "result": {
//	   "name": "models/nld3pk7n/fits/cpyhg3rz"
//	 }
//	}
type GetOperationResp struct {
	Name     string `json:"name"`
	Done     bool   `json:"done"`
	Metadata struct {
		Fit struct {
			Name string `json:"name"`
		} `json:"fit"`
		Progress string `json:"progress"`
	} `json:"metadata"`
	Result struct {
		Name string `json:"name"`
	} `json:"result"`
}

type GetParamsResp struct {
	Id     string `json:"id"`
	Params []struct {
		ConstrainedNames []string `json:"constrained_names"`
		Dims             []int    `json:"dims"`
		Name             string   `json:"name"`
	} `json:"params"`
}

func NewStan(host, port string) *Stan {
	client := resty.New()
	client.SetBaseURL(fmt.Sprintf("http://%s:%s", host, port))
	return &Stan{
		client: client,
		logger: zap.NewNop(),
	}
}

// CompileModel compiles the model
// Status Codes:
//
//	201 Created – Identifier for compiled Stan model and compiler output.
//	400 Bad Request – Error associated with compile request.
func (s Stan) CompileModel(ctx context.Context, programCode string) (ModelCompileResp, error) {
	compileResp := ModelCompileResp{}
	compileReq := ModelCompileReq{ProgramCode: programCode}

	resp, err := s.client.R().
		SetBody(compileReq).
		SetResult(&compileResp).
		Post("/v1/models")
	if err != nil {
		s.logger.Error("HttpStan failed to compile model",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return compileResp, err
	}
	if resp.StatusCode() != http.StatusCreated {
		s.logger.Error("HttpStan failed to compile model",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("status", resp.Status()),
				zap.Int("status_code", resp.StatusCode()),
			)...,
		)
		return compileResp, fmt.Errorf("HttpStan failed to compile model: %s", resp.Status())
	}

	s.logger.Info("HttpStan compiled model",
		zap.String("status", resp.Status()),
		zap.Int64("time_cost", resp.Time().Milliseconds()),
	)

	return compileResp, nil
}

// CreateFit Post request /v1/models/{model_id}/fits to start a long-running fit operation
//
// Status Codes:
//
//	201 Created – Identifier for completed Stan fit
//	400 Bad Request – Error associated with request.
//	404 Not Found – Fit not found.
func (s Stan) CreateFit(ctx context.Context, modelID string, req CreateFitReq) (string, error) {
	createFitResp := CreateFitResp{}
	crateFitErrResp := CreateFitErrResp{}

	resp, err := s.client.R().
		SetPathParam("model_id", modelID).
		SetBody(req).
		SetResult(&createFitResp).
		SetError(&crateFitErrResp).
		Post("/v1/models/{model_id}/fits")
	if err != nil {
		s.logger.Error("HttpStan failed to create fit",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("model_id", modelID),
			)...,
		)
		return "", err
	}
	if resp.StatusCode() != http.StatusCreated {
		s.logger.Error("HttpStan failed to create fit",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("model_id", modelID),
				zap.String("status", resp.Status()),
				zap.Int64("time_cost", resp.Time().Milliseconds()),
				zap.String("error_message", crateFitErrResp.Message),
			)...,
		)
		return "", fmt.Errorf("HttpStan failed to create fit: %s", resp.Status())
	}
	return createFitResp.Name, nil
}

// GetOperationDetails GET /v1/operations/{operation_id}
//
// Return Operation details. Details about an Operation include
// whether the operation is done and information about the progress of sampling.
//
// Status Codes:
//
//	200 OK – Operation name and metadata.
//	404 Not Found – Operation not found.
func (s Stan) GetOperationDetails(ctx context.Context, operationID string) (GetOperationResp, error) {
	getOperationResp := GetOperationResp{}

	resp, err := s.client.R().
		SetPathParam("operation_id", operationID).
		SetResult(&getOperationResp).
		Get("/v1/operations/{operation_id}")
	if err != nil {
		return getOperationResp, err
	}
	if resp.StatusCode() != http.StatusOK {
		s.logger.Error("HttpStan failed to get operation details",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("operation_id", operationID),
				zap.String("status", resp.Status()),
				zap.Int64("time_cost", resp.Time().Milliseconds()),
			)...,
		)
		return getOperationResp, fmt.Errorf("HttpStan failed to get operation details: %s", resp.Status())
	}
	return getOperationResp, nil
}

// GetFitResult GET /v1/models/{model_id}/fits/{fit_id}
//
// Get results returned by a function.
// Result (draws, logger messages) from calling a function defined in stan::services.
//
// Status Codes:
//
//	200 OK – Newline-delimited JSON-encoded messages from Stan. Includes draws.
//	404 Not Found – Fit not found.
func (s Stan) GetFitResult(ctx context.Context, modelID, fitID string) (io.ReadCloser, error) {
	resp, err := s.client.R().
		SetPathParam("model_id", modelID).
		SetPathParam("fit_id", fitID).
		SetDoNotParseResponse(true).
		Get("/v1/models/{model_id}/fits/{fit_id}")
	if err != nil {
		s.logger.Error("HttpStan failed to get fit result",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("model_id", modelID),
				zap.String("fit_id", fitID),
			)...,
		)
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		s.logger.Error("HttpStan failed to get fit result",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("model_id", modelID),
				zap.String("fit_id", fitID),
				zap.String("status", resp.Status()),
				zap.Int64("time_cost", resp.Time().Milliseconds()),
			)...,
		)
		return nil, fmt.Errorf("HttpStan failed to get fit result: %s", resp.Status())
	}
	return resp.RawBody(), nil
}

// StanParams POST /v1/models/{model_id}/params to get parameter names and dimensions.
//
// Status Codes:
//
//	200 OK – Parameters for Stan Model
//	400 Bad Request – Error associated with request.
//	404 Not Found – Model not found.
func (s Stan) StanParams(ctx context.Context, modelID string, data map[string]interface{}) ([]string, []string) {
	paramsResp := GetParamsResp{}
	resp, err := s.client.R().
		SetPathParam("model_id", modelID).
		SetBody(map[string]interface{}{
			"data": data,
		}).
		SetResult(&paramsResp).
		Post("/v1/models/{model_id}/params")
	if err != nil {
		s.logger.Error("HttpStan failed to get params",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
		return nil, nil
	}
	if resp.StatusCode() != http.StatusOK {
		s.logger.Error("HttpStan failed to get params",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("status", resp.Status()),
				zap.Int64("time_cost", resp.Time().Milliseconds()),
			)...,
		)
		return nil, nil
	}
	var constrainedNames []string
	var paramNames []string
	for _, param := range paramsResp.Params {
		constrainedNames = append(constrainedNames, param.ConstrainedNames...)
		paramNames = append(paramNames, param.Name)
	}
	return constrainedNames, paramNames
}

func (s Stan) ExtractFromFitResult(ctx context.Context, res io.ReadCloser) dataframe.DataFrame {
	defer func(res io.ReadCloser) {
		err := res.Close()
		if err != nil {
			s.logger.Error("HttpStan failed to close response body",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
				)...,
			)
		}
	}(res)
	data := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		lineJSON := make(map[string]interface{})
		unmarshalErr := json.Unmarshal(scanner.Bytes(), &lineJSON)
		if unmarshalErr != nil {
			s.logger.Error("HttpStan failed to unmarshal line",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.String("line", scanner.Text()),
					zap.Error(unmarshalErr),
				)...,
			)
			continue
		}
		if lineJSON["topic"] == "sample" {
			if values, ok := lineJSON["values"].(map[string]interface{}); ok {
				data = append(data, values)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		s.logger.Error("HttpStan failed to scan response body",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
			)...,
		)
	}
	stanDataframe := dataframe.LoadMaps(data)
	return stanDataframe
}

func ModelCode() string {
	return stanProgramCode
}
