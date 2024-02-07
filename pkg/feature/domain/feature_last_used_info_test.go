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

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestFeatureLastUsedInfoID(t *testing.T) {
	patterns := []struct {
		featureID, expect string
		version           int32
	}{
		{
			featureID: "feature-id",
			version:   10,
			expect:    "feature-id:10",
		},
	}

	for _, p := range patterns {
		assert.Equal(t, p.expect, FeatureLastUsedInfoID(p.featureID, p.version))
	}
}

func TestID(t *testing.T) {
	patterns := []struct {
		featureID, expect string
		version           int32
	}{
		{
			featureID: "feature-id",
			version:   10,
			expect:    "feature-id:10",
		},
	}

	for _, p := range patterns {
		f := FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
			FeatureId: p.featureID,
			Version:   p.version,
		}}
		assert.Equal(t, p.expect, f.ID())
	}
}

func TestNewFeatureLastUsedInfo(t *testing.T) {
	patterns := []struct {
		featureID     string
		version       int32
		lastUsedAt    int64
		createdAt     int64
		clientVersion string
		expect        *FeatureLastUsedInfo
	}{
		{
			featureID:     "feature-id",
			version:       10,
			createdAt:     123445566,
			lastUsedAt:    123445566,
			clientVersion: "1.0.0",
		},
	}

	for _, p := range patterns {
		expect := &FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
			FeatureId:  p.featureID,
			Version:    p.version,
			LastUsedAt: p.lastUsedAt,
			CreatedAt:  p.createdAt,
		}}
		expect.SetClientVersion(p.clientVersion)
		assert.Equal(t, expect, NewFeatureLastUsedInfo(p.featureID, p.version, p.lastUsedAt, p.clientVersion))
	}
}

func TestUsedAt(t *testing.T) {
	patterns := []struct {
		v, expect int64
		desc      string
	}{
		{
			v:      0,
			expect: 1,
		},
		{
			v:      1,
			expect: 1,
		},
		{
			v:      2,
			expect: 2,
		},
	}
	for _, p := range patterns {
		featureLastUsed := NewFeatureLastUsedInfo("id", 10, 1, "1.0.0")
		featureLastUsed.UsedAt(p.v)
		assert.Equal(t, p.expect, featureLastUsed.LastUsedAt, p.desc)
	}
}

func TestSetOldestClientVersion(t *testing.T) {
	patterns := []struct {
		currentVersion, clientVersion, expect, desc string
		expectError                                 bool
	}{
		{
			currentVersion: "",
			clientVersion:  "1.0.0",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "empty_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "empty_client_version",
		},
		{
			currentVersion: "10",
			clientVersion:  "1.0.0",
			expect:         "10",
			expectError:    true,
			desc:           "invalid_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "10",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "invalid_client_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "1.0.1",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "client_version_greater_than_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "1.0.0",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "client_version_equal_current_version",
		},
		{
			currentVersion: "1.0.9",
			clientVersion:  "1.0.1",
			expect:         "1.0.1",
			expectError:    false,
			desc:           "client_version_less_than_current_version",
		},
	}
	for _, p := range patterns {
		info := &FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
			FeatureId:           "id",
			Version:             10,
			LastUsedAt:          1,
			CreatedAt:           1,
			ClientOldestVersion: p.currentVersion,
		}}
		err := info.setOldestClientVersion(p.clientVersion)
		assert.Equal(t, p.expectError, err != nil, p.desc)
		assert.Equal(t, p.expect, info.ClientOldestVersion, p.desc)
	}
}

func TestSetLatestClientVersion(t *testing.T) {
	patterns := []struct {
		currentVersion, clientVersion, expect, desc string
		expectError                                 bool
	}{
		{
			currentVersion: "",
			clientVersion:  "1.0.0",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "empty_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "empty_client_version",
		},
		{
			currentVersion: "10",
			clientVersion:  "1.0.0",
			expect:         "10",
			expectError:    true,
			desc:           "invalid_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "10",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "invalid_client_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "1.0.1",
			expect:         "1.0.1",
			expectError:    false,
			desc:           "client_version_greater_than_current_version",
		},
		{
			currentVersion: "1.0.0",
			clientVersion:  "1.0.0",
			expect:         "1.0.0",
			expectError:    false,
			desc:           "client_version_equal_current_version",
		},
		{
			currentVersion: "1.0.9",
			clientVersion:  "1.0.1",
			expect:         "1.0.9",
			expectError:    false,
			desc:           "client_version_less_than_current_version",
		},
	}
	for _, p := range patterns {
		info := &FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
			FeatureId:           "id",
			Version:             10,
			LastUsedAt:          1,
			CreatedAt:           1,
			ClientLatestVersion: p.currentVersion,
		}}
		err := info.setLatestClientVersion(p.clientVersion)
		assert.Equal(t, p.expectError, err != nil, p.desc)
		assert.Equal(t, p.expect, info.ClientLatestVersion, p.desc)
	}
}
