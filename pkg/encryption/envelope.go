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

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// envelopePrefix marks every encrypted value, whatever its version.
	envelopePrefix = "enc:v"
	// envelopeV1Prefix is the header of the only format written today:
	// "enc:v1:<key-id>:<base64url payload>". The payload is a serialized
	// go-kms-wrapping BlobInfo.
	envelopeV1Prefix = envelopePrefix + "1:"
	envelopeSep      = ":"
)

// IsEncrypted reports whether value is an encryption envelope rather than a
// legacy plaintext value. The stored secrets (a JSON document or an https URL)
// can never start with the envelope prefix.
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, envelopePrefix)
}

func encodeEnvelope(keyID string, payload []byte) string {
	return envelopeV1Prefix + keyID + envelopeSep + base64.RawURLEncoding.EncodeToString(payload)
}

func decodeEnvelope(value string) (keyID string, payload []byte, err error) {
	rest, ok := strings.CutPrefix(value, envelopeV1Prefix)
	if !ok {
		return "", nil, fmt.Errorf("%w: unsupported envelope version", ErrMalformedValue)
	}
	keyID, encoded, ok := strings.Cut(rest, envelopeSep)
	if !ok || keyID == "" || encoded == "" {
		return "", nil, fmt.Errorf("%w: missing key id or payload", ErrMalformedValue)
	}
	payload, err = base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrMalformedValue, err)
	}
	return keyID, payload, nil
}
