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

package rest

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Data interface{} `json:"data"`
}

// This response is based on https://google.github.io/styleguide/jsoncstyleguide.xml?showone=error#error.
type failureResponse struct {
	Error errorResponse `json:"error"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ReturnFailureResponse(w http.ResponseWriter, err error) {
	status, ok := convertToErrStatus(err)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status.GetStatusCode())
	returnResponse(
		w,
		&failureResponse{
			Error: errorResponse{
				Code:    status.GetStatusCode(),
				Message: status.GetErrMessage(),
			},
		},
	)
}

func ReturnSuccessResponse(w http.ResponseWriter, resp interface{}) {
	returnResponse(w, successResponse{Data: resp})
}

func returnResponse(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoded, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(encoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
