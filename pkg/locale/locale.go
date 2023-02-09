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

	"google.golang.org/grpc/metadata"
)

const (
	Ja = "ja"
	En = "en"
)

type locale struct {
	locale string
}

type Locale interface {
	GetLocale() string
}

func NewLocale(l string) Locale {
	return &locale{
		locale: l,
	}
}

func (l *locale) GetLocale() string {
	return l.locale
}

func getAcceptLang(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Ja
	}
	keys, ok := md["accept-language"]
	if !ok || len(keys) == 0 || keys[0] == "" {
		return Ja
	}
	// TODO: change this part after writing en.yaml
	return Ja
}
