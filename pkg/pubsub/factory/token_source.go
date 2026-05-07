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

package factory

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

var errTokenNil = errors.New("resilientTokenSource: base returned nil token")

// resilientTokenSource wraps an oauth2.TokenSource and returns the last
// successfully fetched token when the underlying source cannot refresh
// (e.g. GCE metadata server unreachable during spot VM preemption).
// This gives in-flight RPCs a chance to succeed with the still-valid cached
// token instead of failing immediately with an authentication error.
type resilientTokenSource struct {
	base     oauth2.TokenSource
	mu       sync.Mutex
	lastGood *oauth2.Token
}

func (s *resilientTokenSource) Token() (*oauth2.Token, error) {
	tok, err := s.base.Token()
	if err == nil && tok != nil {
		s.mu.Lock()
		s.lastGood = tok
		s.mu.Unlock()
		return tok, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastGood != nil &&
		(s.lastGood.Expiry.IsZero() || time.Now().Before(s.lastGood.Expiry)) {
		return s.lastGood, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, errTokenNil
}
