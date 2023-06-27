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

package notification

import (
	"context"

	"github.com/bucketeer-io/bucketeer/pkg/job"
)

type ExperimentRunningWatcherJob struct {
	ExperimentRunningWatcher job.Job
}

func NewExperimentRunningWatcherJob(experimentRunningWatcher job.Job) *ExperimentRunningWatcherJob {
	return &ExperimentRunningWatcherJob{
		ExperimentRunningWatcher: experimentRunningWatcher,
	}
}

func (j *ExperimentRunningWatcherJob) Run(ctx context.Context) error {
	return j.ExperimentRunningWatcher.Run(ctx)
}
