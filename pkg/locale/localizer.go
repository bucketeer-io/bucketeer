// Copyright 2022 The Bucketeer Authors.
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

package locale

import (
	"context"
	"embed"
	"fmt"
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var (
	bundle *i18n.Bundle

	//go:embed localizedata
	localizedata embed.FS
)

// status error messages
const (
	// nouns
	FeatureFlagID = "FeatureFlagID"
	// error sentence
	RequiredFieldTemplate = "RequiredField"
	InternalServerError   = "InternalServerError"
	NotFoundError         = "NotFoundError"
	InvalidArgumentError  = "InvalidArgumentError"
	UnauthenticatedError  = "UnauthenticatedError"
	PermissionDenied      = "PermissionDenied"
	AlreadyExistsError    = "AlreadyExistsError"
	AlreadyDeletedError   = "AlreadyDeletedError"
	// event counter
	StartAtIsAfterEndAt = "StartAtIsAfterEndAt"
	// environment
	ProjectDisabled = "ProjectDisabled"
	// feature
	SegmentInUse                  = "SegmentInUse"
	SegmentUsersAlreadyUploading  = "SegmentUsersAlreadyUploading"
	SegmentStatusNotSucceeded     = "SegmentStatusNotSucceeded"
	HasWaitingOrRunningExperiment = "HasWaitingOrRunningExperiment"
	NothingToChange               = "NothingToChange"
	DifferentVariationsSize       = "DifferentVariationsSize"
)

// domain events
const (
	UnknownOperation   = "UnknownOperation"
	CreatedTemplate    = "Created"
	EnabledTemplate    = "Enabled"
	DisabledTemplate   = "Disabled"
	ArchivedTemplate   = "Archived"
	UnarchivedTemplate = "Unarchived"
	DeletedTemplate    = "Deleted"
	AddedTemplate      = "Added"
	ChangedTemplate    = "Changed"
	UpdatedTemplate    = "Updated"
	StartedTemplate    = "Started"
	StoppedTemplate    = "Stopped"
	FinishedTemplate   = "Finished"
	UploadedTemplate   = "Uploaded"
	ResetTemplate      = "Reset"
	ClonedTemplate     = "Cloned"
	TrialConverted     = "TrialConverted"
)

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	files := []string{
		"en.yaml",
		"ja.yaml",
	}
	for _, f := range files {
		data, err := localizedata.ReadFile(fmt.Sprintf("localizedata/%s", f))
		if err != nil {
			panic(fmt.Errorf("Failed to load translation data: %s", f))
		}
		bundle.MustParseMessageFileBytes(data, f)
	}
}

type localizer struct {
	Locale
	*i18n.Localizer
}

type Localizer interface {
	Locale
	MustLocalize(id string) string
	MustLocalizeWithTemplate(id string, fields ...string) string
}

func NewLocalizer(ctx context.Context, fopts ...Option) Localizer {
	locale := NewLocale(getAcceptLang(ctx))
	opts := defaultOptions()
	for _, fo := range fopts {
		fo.apply(&opts)
	}
	return &localizer{
		locale,
		i18n.NewLocalizer(opts.bundle, locale.GetLocale()),
	}
}

func (l *localizer) MustLocalize(id string) string {
	return l.Localizer.MustLocalize(createLocalizeConfig(id))
}

func (l *localizer) MustLocalizeWithTemplate(id string, fields ...string) string {
	return l.Localizer.MustLocalize(createLocalizeConfigWithTemplate(id, fields...))
}

func createLocalizeConfig(id string) *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID: id,
	}
}

func createLocalizeConfigWithTemplate(id string, fields ...string) *i18n.LocalizeConfig {
	td := make(map[string]interface{}, len(fields))
	for i, f := range fields {
		td["Field_"+strconv.Itoa(i+1)] = f
	}
	return &i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: td,
	}
}
