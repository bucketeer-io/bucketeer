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

package domain

import (
	"encoding/json"
	"fmt"

	domaineventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
)

// ExtractAPIKeySecrets returns distinct API key secret strings from a domain
// event's entity data snapshots. Both previous and current snapshots are
// checked so that key rotations evict both the old and new secrets.
// Returns an error if the JSON is malformed or the api_key field is missing
// from all non-empty snapshots.
func ExtractAPIKeySecrets(e *domaineventproto.Event) ([]string, error) {
	seen := make(map[string]struct{})
	var out []string
	var lastErr error
	for _, raw := range []string{e.GetPreviousEntityData(), e.GetEntityData()} {
		if raw == "" {
			continue
		}
		sec, err := parseAPIKeySecret(raw)
		if err != nil {
			lastErr = err
			continue
		}
		if sec == "" {
			continue
		}
		if _, ok := seen[sec]; ok {
			continue
		}
		seen[sec] = struct{}{}
		out = append(out, sec)
	}
	if len(out) == 0 && lastErr != nil {
		return nil, fmt.Errorf("failed to extract api_key from entity data: %w", lastErr)
	}
	return out, nil
}

func parseAPIKeySecret(raw string) (string, error) {
	var v struct {
		APIKey string `json:"api_key"`
	}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("invalid JSON in entity data: %w", err)
	}
	return v.APIKey, nil
}
