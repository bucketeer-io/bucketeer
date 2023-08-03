// Copyright 2023 The Bucketeer Authors.
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

package api

import "go.uber.org/atomic"

type RunningJobManager struct {
	runningJobs *atomic.Int32
}

func NewRunningJobManager() *RunningJobManager {
	return &RunningJobManager{
		runningJobs: atomic.NewInt32(0),
	}
}

func (m RunningJobManager) AddRunningJob() {
	m.runningJobs.Inc()
}

func (m RunningJobManager) RemoveRunningJob() {
	m.runningJobs.Dec()
}

func (m RunningJobManager) GetCurrentRunningJobs() int32 {
	return m.runningJobs.Load()
}
