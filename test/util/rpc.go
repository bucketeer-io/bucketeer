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

package util

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	gstatus "google.golang.org/grpc/status"
)

// compareErrorDetails compares the content of error details without considering their order
func CompareErrorDetails(expected, actual error) bool {
	expectedStatus, _ := gstatus.FromError(expected)
	actualStatus, _ := gstatus.FromError(actual)

	// Compare basic error information (code, message)
	if expectedStatus.Code() != actualStatus.Code() {
		return false
	}

	// Get error details
	expectedDetails := expectedStatus.Details()
	actualDetails := actualStatus.Details()

	// Compare number of details
	if len(expectedDetails) != len(actualDetails) {
		return false
	}

	// Compare content of each detail (order doesn't matter)
	expectedMap := make(map[string]string)
	actualMap := make(map[string]string)

	for _, detail := range expectedDetails {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			for k, v := range info.Metadata {
				expectedMap[k] = v
			}
		}
	}

	for _, detail := range actualDetails {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			for k, v := range info.Metadata {
				actualMap[k] = v
			}
		}
	}

	// Compare map contents
	for k, v := range expectedMap {
		if actualMap[k] != v {
			return false
		}
	}

	return true
}
