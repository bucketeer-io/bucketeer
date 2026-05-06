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

package database

// QueryNoLimit is a list limit meaning "no row cap" for semantic list parameters
// and storage that forwards into dialect list builders (same numeric value as
// pkg/storage/v2/mysql and pkg/storage/v2/postgres).
const QueryNoLimit = 0

// QueryNoOffset is a pagination offset meaning "start at the first row".
const QueryNoOffset = 0
