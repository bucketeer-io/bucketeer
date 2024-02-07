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
	"fmt"

	"github.com/blang/semver"

	"github.com/bucketeer-io/bucketeer/proto/feature"
)

type FeatureLastUsedInfo struct {
	*feature.FeatureLastUsedInfo
}

func NewFeatureLastUsedInfo(
	featureID string,
	version int32,
	lastUsedAt int64,
	clientVersion string,
) *FeatureLastUsedInfo {
	info := &FeatureLastUsedInfo{FeatureLastUsedInfo: &feature.FeatureLastUsedInfo{
		FeatureId:  featureID,
		Version:    version,
		LastUsedAt: lastUsedAt,
		CreatedAt:  lastUsedAt,
	}}
	info.SetClientVersion(clientVersion) // nolint:errcheck
	return info
}

func (f *FeatureLastUsedInfo) ID() string {
	return FeatureLastUsedInfoID(f.FeatureId, f.Version)
}

func (f *FeatureLastUsedInfo) UsedAt(v int64) {
	if f.LastUsedAt < v {
		f.LastUsedAt = v
	}
}

func (f *FeatureLastUsedInfo) SetClientVersion(version string) error {
	if err := f.setOldestClientVersion(version); err != nil {
		return err
	}
	if err := f.setLatestClientVersion(version); err != nil {
		return err
	}
	return nil
}

func (f *FeatureLastUsedInfo) setOldestClientVersion(version string) error {
	clientSemVersion, err := f.parseSemver(version)
	if err != nil {
		// Because the client version is optional and
		// it could not be a semantic version, it ignores parse errors
		return nil
	}
	if f.ClientOldestVersion == "" {
		f.ClientOldestVersion = clientSemVersion.String()
		return nil
	}
	currentSemVersion, err := f.parseSemver(f.ClientOldestVersion)
	if err != nil {
		return err
	}
	if currentSemVersion.GT(clientSemVersion) {
		f.ClientOldestVersion = clientSemVersion.String()
	}
	return nil
}

func (f *FeatureLastUsedInfo) setLatestClientVersion(version string) error {
	clientSemVersion, err := f.parseSemver(version)
	if err != nil {
		// Because the client version is optional and
		// it could not be a semantic version, it ignores parse errors
		return nil
	}
	if f.ClientLatestVersion == "" {
		f.ClientLatestVersion = clientSemVersion.String()
		return nil
	}
	currentSemVersion, err := f.parseSemver(f.ClientLatestVersion)
	if err != nil {
		return err
	}
	if currentSemVersion.LT(clientSemVersion) {
		f.ClientLatestVersion = clientSemVersion.String()
	}
	return nil
}

func (f *FeatureLastUsedInfo) parseSemver(value string) (semver.Version, error) {
	version, err := semver.Parse(value)
	if err != nil {
		return semver.Version{}, err
	}
	return version, nil
}

func FeatureLastUsedInfoID(featureID string, version int32) string {
	return fmt.Sprintf("%s:%d", featureID, version)
}
