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

package api

import (
	"context"

	prdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/autoops/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func executeStopProgressiveRolloutOperation(
	ctx context.Context,
	storage v2as.ProgressiveRolloutStorage,
	featureIDs []interface{},
	environmentNamespace string,
	operation autoopsproto.ProgressiveRollout_StoppedBy,
) error {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_namespace", "=", environmentNamespace),
		mysql.NewInFilter("feature_id", featureIDs),
	}
	list, _, _, err := storage.ListProgressiveRollouts(ctx, whereParts, nil, 0, 0)
	if err != nil {
		return err
	}
	for _, rollout := range list {
		r := &prdomain.ProgressiveRollout{ProgressiveRollout: rollout}
		if r.IsWaiting() || r.IsRunning() {
			if err := r.Stop(operation); err != nil {
				return err
			}
			if err := storage.UpdateProgressiveRollout(ctx, r, environmentNamespace); err != nil {
				return err
			}
		}
	}
	return nil
}
