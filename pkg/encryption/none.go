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

package encryption

import "context"

// noneEncryptor keeps values as they are. It still refuses to hand back an
// encrypted value as if it were plaintext, so a value written while encryption
// was enabled is reported instead of being used as a credential.
type noneEncryptor struct{}

// NewNone returns the Encryptor used when encryption is disabled.
func NewNone() Encryptor {
	return noneEncryptor{}
}

func (noneEncryptor) Encrypt(_ context.Context, plaintext string) (string, error) {
	return plaintext, nil
}

func (noneEncryptor) Decrypt(_ context.Context, value string) (string, error) {
	if IsEncrypted(value) {
		return "", ErrNotConfigured
	}
	return value, nil
}
