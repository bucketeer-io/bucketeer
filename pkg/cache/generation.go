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

package cache

import (
	"crypto/sha1" // content fingerprint, not a security primitive; Redis Lua only offers sha1hex
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GenerationDigits is the fixed width of an encoded generation. int64 needs at
// most 19 decimal digits; 20 keeps a leading zero so that two encoded
// generations always compare correctly as plain strings. Fixed-width string
// comparison is what the Redis Lua script relies on: Lua numbers are doubles
// and lose precision above 2^53, which Unix nanosecond timestamps exceed.
const GenerationDigits = 20

// generationDigestHexLen is the length of a hex-encoded SHA-1 digest.
const generationDigestHexLen = sha1.Size * 2

// generationSeparator separates the fields of a stored generation marker.
const generationSeparator = ":"

// ErrInvalidGeneration is returned when a negative generation is supplied.
var ErrInvalidGeneration = errors.New("cache: generation must not be negative")

// FormatGeneration encodes generation as a zero-padded, fixed-width decimal
// string suitable for lexicographic comparison.
func FormatGeneration(generation int64) (string, error) {
	if generation < 0 {
		return "", ErrInvalidGeneration
	}
	return fmt.Sprintf("%0*d", GenerationDigits, generation), nil
}

// ValueDigest returns the hex SHA-1 fingerprint of value, as Redis's
// redis.sha1hex would compute it.
func ValueDigest(value []byte) string {
	sum := sha1.Sum(value)
	return hex.EncodeToString(sum[:])
}

// EncodeGenerationMarker builds the value stored under a generation key:
//
//	"<fixed-width generation>:<length of the described value>:<sha1 hex of the described value>"
//
// Recording what the generation describes lets a writer detect that the value
// was since overwritten by a path that does not maintain the marker (for
// example a pod running a build that predates conditional puts doing a plain
// SET during a rolling deploy). In that case the marker no longer describes
// the stored value and must not be used to reject writes. The length is a
// cheap first check; the digest catches same-length overwrites.
func EncodeGenerationMarker(formattedGeneration string, value []byte) string {
	return formattedGeneration +
		generationSeparator + strconv.Itoa(len(value)) +
		generationSeparator + ValueDigest(value)
}

// DecodeGenerationMarker parses a marker produced by EncodeGenerationMarker.
// ok is false for anything that does not match the exact format.
func DecodeGenerationMarker(marker string) (formattedGeneration string, valueLen int, digest string, ok bool) {
	parts := strings.Split(marker, generationSeparator)
	if len(parts) != 3 {
		return "", 0, "", false
	}
	formattedGeneration = parts[0]
	if len(formattedGeneration) != GenerationDigits || !isAllDigits(formattedGeneration) {
		return "", 0, "", false
	}
	if !isAllDigits(parts[1]) {
		return "", 0, "", false
	}
	valueLen, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, "", false
	}
	digest = parts[2]
	if len(digest) != generationDigestHexLen || !isAllLowerHex(digest) {
		return "", 0, "", false
	}
	return formattedGeneration, valueLen, digest, true
}

// IsStaleGeneration reports whether a write carrying newFormattedGeneration
// must be rejected given the current generation marker and the value
// currently stored. currentValue is nil when no value exists.
//
// A write is rejected only when all of the following hold:
//   - a value currently exists (valueExists),
//   - the marker is well-formed,
//   - the marker describes the current value (length and digest match), and
//   - the marker's generation is strictly greater than the new one.
//
// Equal generations are accepted (idempotent rewrite). The digest is only
// computed when every cheaper check already points at rejection, so the
// common accept path never hashes the current value.
func IsStaleGeneration(
	valueExists bool,
	currentValue []byte,
	currentMarker string,
	newFormattedGeneration string,
) bool {
	if !valueExists {
		return false
	}
	currentGeneration, describedLen, describedDigest, ok := DecodeGenerationMarker(currentMarker)
	if !ok {
		return false
	}
	if strings.Compare(newFormattedGeneration, currentGeneration) >= 0 {
		return false
	}
	if describedLen != len(currentValue) {
		return false
	}
	return describedDigest == ValueDigest(currentValue)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func isAllLowerHex(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
