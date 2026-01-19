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
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
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

func TestClientVersionVPrefixNormalization(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc                  string
		clientVersion         string
		expectedOldestVersion string
		expectedLatestVersion string
	}{
		{
			desc:                  "v_prefix_is_stripped",
			clientVersion:         "v10.154.2",
			expectedOldestVersion: "10.154.2",
			expectedLatestVersion: "10.154.2",
		},
		{
			desc:                  "no_v_prefix_works",
			clientVersion:         "10.154.2",
			expectedOldestVersion: "10.154.2",
			expectedLatestVersion: "10.154.2",
		},
		{
			desc:                  "v_prefix_with_prerelease",
			clientVersion:         "v1.0.0-alpha",
			expectedOldestVersion: "1.0.0-alpha",
			expectedLatestVersion: "1.0.0-alpha",
		},
		{
			desc:                  "v_prefix_with_build_metadata",
			clientVersion:         "v1.0.0+build123",
			expectedOldestVersion: "1.0.0+build123",
			expectedLatestVersion: "1.0.0+build123",
		},
	}
	for _, p := range patterns {
		p := p
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			info := NewFeatureLastUsedInfo("feature-id", 1, 123456789, p.clientVersion)
			assert.Equal(t, p.expectedOldestVersion, info.ClientOldestVersion)
			assert.Equal(t, p.expectedLatestVersion, info.ClientLatestVersion)
		})
	}
}

func TestClientVersionVPrefixComparison(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc                  string
		currentOldest         string
		currentLatest         string
		newVersion            string
		expectedOldestVersion string
		expectedLatestVersion string
	}{
		{
			desc:                  "v_prefix_new_version_updates_oldest",
			currentOldest:         "2.0.0",
			currentLatest:         "2.0.0",
			newVersion:            "v1.0.0",
			expectedOldestVersion: "1.0.0",
			expectedLatestVersion: "2.0.0",
		},
		{
			desc:                  "v_prefix_new_version_updates_latest",
			currentOldest:         "1.0.0",
			currentLatest:         "1.0.0",
			newVersion:            "v2.0.0",
			expectedOldestVersion: "1.0.0",
			expectedLatestVersion: "2.0.0",
		},
		{
			desc:                  "mixed_v_prefix_comparison",
			currentOldest:         "1.5.0",
			currentLatest:         "1.5.0",
			newVersion:            "v1.5.0",
			expectedOldestVersion: "1.5.0",
			expectedLatestVersion: "1.5.0",
		},
	}
	for _, p := range patterns {
		p := p
		t.Run(p.desc, func(t *testing.T) {
			t.Parallel()
			info := &FeatureLastUsedInfo{FeatureLastUsedInfo: &proto.FeatureLastUsedInfo{
				FeatureId:           "id",
				Version:             10,
				LastUsedAt:          1,
				CreatedAt:           1,
				ClientOldestVersion: p.currentOldest,
				ClientLatestVersion: p.currentLatest,
			}}
			err := info.SetClientVersion(p.newVersion)
			assert.NoError(t, err)
			assert.Equal(t, p.expectedOldestVersion, info.ClientOldestVersion)
			assert.Equal(t, p.expectedLatestVersion, info.ClientLatestVersion)
		})
	}
}
