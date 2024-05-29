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
	"context"
	_ "embed"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var stan *Stan

func TestMain(m *testing.M) {
	// setup stan
	stan = NewStan("localhost", "8080")
	// sleep for 5 seconds to wait for stan to be ready
	time.Sleep(5 * time.Second)
	m.Run()
}

func TestStanCompileModel(t *testing.T) {
	compileModel(t)
}

func TestStanSample(t *testing.T) {
	// compile model
	compileModel := compileModel(t)
	// create fit
	modelID := compileModel.Name[len("models/"):]
	req := CreateFitReq{
		Chain: 1,
		Data: map[string]interface{}{
			"g": 3,
			"x": []int{1, 2, 3},
			"n": []int{10, 10, 10},
		},
		Function:   HmcNUTSFunction,
		NumSamples: 1000,
		NumWarmup:  1000,
		RandomSeed: 1234,
	}
	fit, err := stan.CreateFit(context.TODO(), modelID, req)
	assert.NoError(t, err, "Failed to create fit")
	if assert.NotEmpty(t, fit) {
		fmt.Printf("HttpStan created fit name: %s\n", fit)
	}
	// get operation details
	fitId := fit[len("operations/"):]
	checkOperationUntilDone(t, fitId)
	// get fit
	res, err := stan.GetFitResult(context.TODO(), modelID, fitId)
	assert.NoError(t, err, "Failed to get fit result")
	stanDataframe := stan.ExtractFromFitResult(context.TODO(), res)
	assert.Equal(t, req.NumSamples, stanDataframe.Nrow(), "Failed to get fit result")
}

func TestStanStanParams(t *testing.T) {
	// compile model
	compileModel := compileModel(t)
	constrainedNames, paramNames := stan.StanParams(context.TODO(),
		compileModel.Name[len("models/"):],
		map[string]interface{}{
			"g": 3,
			"x": []int{1, 2, 3},
			"n": []int{10, 10, 10},
		})
	if assert.NotEmpty(t, constrainedNames, "Failed to get constrained names") {
		fmt.Printf("HttpStan constrained names: %v\n", constrainedNames)
	}
	if assert.NotEmpty(t, paramNames, "Failed to get param names") {
		fmt.Printf("HttpStan param names: %v\n", paramNames)
	}
}

func compileModel(t *testing.T) ModelCompileResp {
	// compile model
	compileModel, compileErr := stan.CompileModel(context.TODO(), ModelCode())
	assert.NoError(t, compileErr, "Failed to compile model")
	if assert.NotEmpty(t, compileModel) {
		if assert.NotEmpty(t, compileModel.Name) {
			fmt.Printf("HttpStan compiled model name: %s\n", compileModel.Name)
		}
		assert.NotEmpty(t, compileModel.CompilerOutput)
		assert.NotEmpty(t, compileModel.StancWarnings)
	}
	return compileModel
}

func checkOperationUntilDone(t *testing.T, fitId string) {
	for {
		details, err := stan.GetOperationDetails(context.TODO(), fitId)
		assert.NoError(t, err, "Failed to get operation details")
		if assert.NotEmpty(t, details) {
			if details.Done {
				fmt.Printf("HttpStan operation details name: %s\n", details.Name)
				fmt.Printf("HttpStan operation details done: %t\n", details.Done)
				fmt.Printf("HttpStan operation details progress: %s\n", details.Metadata.Progress)
				fmt.Printf("HttpStan operation details result: %s\n", details.Result.Name)
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}
