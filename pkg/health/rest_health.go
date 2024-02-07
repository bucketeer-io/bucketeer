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

package health

import (
	"fmt"
	"net/http"
)

const healthPath = "/health"

type restChecker struct {
	*checker
	version string
	service string
}

func NewRestChecker(version, service string, opts ...option) *restChecker {
	checker := &restChecker{
		checker: newChecker(opts...),
		version: version,
		service: service,
	}
	return checker
}

func (c *restChecker) Register(mux *http.ServeMux) {
	mux.HandleFunc(fmt.Sprintf("%s%s%s", c.version, c.service, healthPath), c.ServeHTTP)
}
