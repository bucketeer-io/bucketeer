// Copyright 2026 The Bucketeer Authors.
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

package v2

import "errors"

// Shared list-query errors returned by both the SubscriptionStorage and
// AdminSubscriptionStorage implementations.
var (
	ErrInvalidCursor  = errors.New("notification/storage/v2: invalid cursor")
	ErrInvalidOrderBy = errors.New("notification/storage/v2: invalid order by")
)
