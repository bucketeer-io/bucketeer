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
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// sha1("") and sha1("abc"), as redis.sha1hex would return them.
	sha1Empty = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	sha1ABC   = "a9993e364706816aba3e25717850c26c9cd0d89d"
)

func TestFormatGeneration(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		input       int64
		expected    string
		expectedErr error
	}{
		{desc: "zero", input: 0, expected: "00000000000000000000"},
		{desc: "small", input: 42, expected: "00000000000000000042"},
		{
			desc:     "unix nano magnitude",
			input:    1_791_600_000_123_456_789,
			expected: "01791600000123456789",
		},
		{desc: "max int64", input: math.MaxInt64, expected: "09223372036854775807"},
		{desc: "negative", input: -1, expectedErr: ErrInvalidGeneration},
		{desc: "min int64", input: math.MinInt64, expectedErr: ErrInvalidGeneration},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual, err := FormatGeneration(p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
			if err == nil {
				assert.Len(t, actual, GenerationDigits)
			}
		})
	}
}

// Fixed-width encoding must order exactly like the integers it encodes; this
// is what the Lua script and IsStaleGeneration rely on.
func TestFormatGenerationPreservesOrdering(t *testing.T) {
	t.Parallel()
	values := []int64{
		0, 1, 9, 10, 99, 100,
		1_791_600_000_000_000_000,
		1_791_600_000_000_000_001,
		1_791_600_000_000_000_255, // within double-precision rounding of the previous
		math.MaxInt64 - 1,
		math.MaxInt64,
	}
	for i := 1; i < len(values); i++ {
		prev, err := FormatGeneration(values[i-1])
		require.NoError(t, err)
		cur, err := FormatGeneration(values[i])
		require.NoError(t, err)
		assert.Equal(t, -1, strings.Compare(prev, cur), "%d must sort before %d", values[i-1], values[i])
	}
}

func TestValueDigest(t *testing.T) {
	t.Parallel()
	assert.Equal(t, sha1Empty, ValueDigest(nil))
	assert.Equal(t, sha1Empty, ValueDigest([]byte{}))
	assert.Equal(t, sha1ABC, ValueDigest([]byte("abc")))
	assert.NotEqual(t, ValueDigest([]byte("abc")), ValueDigest([]byte("abd")))
	assert.Len(t, ValueDigest([]byte("x")), generationDigestHexLen)
}

func TestEncodeDecodeGenerationMarker(t *testing.T) {
	t.Parallel()
	gen, err := FormatGeneration(123)
	require.NoError(t, err)

	marker := EncodeGenerationMarker(gen, []byte("abc"))
	assert.Equal(t, "00000000000000000123:3:"+sha1ABC, marker)

	decodedGen, decodedLen, decodedDigest, ok := DecodeGenerationMarker(marker)
	assert.True(t, ok)
	assert.Equal(t, gen, decodedGen)
	assert.Equal(t, 3, decodedLen)
	assert.Equal(t, sha1ABC, decodedDigest)

	empty := EncodeGenerationMarker(gen, nil)
	assert.Equal(t, "00000000000000000123:0:"+sha1Empty, empty)
	decodedGen, decodedLen, decodedDigest, ok = DecodeGenerationMarker(empty)
	assert.True(t, ok)
	assert.Equal(t, gen, decodedGen)
	assert.Equal(t, 0, decodedLen)
	assert.Equal(t, sha1Empty, decodedDigest)
}

func TestDecodeGenerationMarkerRejectsMalformed(t *testing.T) {
	t.Parallel()
	const gen = "00000000000000000123"
	patterns := []struct {
		desc  string
		input string
	}{
		{desc: "empty", input: ""},
		{desc: "too short", input: "123:4:" + sha1ABC},
		{desc: "generation only", input: gen},
		{desc: "generation and length only (previous format)", input: gen + ":3"},
		{desc: "missing digest", input: gen + ":3:"},
		{desc: "extra field", input: gen + ":3:" + sha1ABC + ":x"},
		{desc: "wrong separator", input: gen + ";3;" + sha1ABC},
		{desc: "generation too long", input: "0" + gen + ":3:" + sha1ABC},
		{desc: "generation too short", input: gen[1:] + ":3:" + sha1ABC},
		{desc: "non-digit in generation", input: "0000000000000000012x:3:" + sha1ABC},
		{desc: "non-digit in length", input: gen + ":4a:" + sha1ABC},
		{desc: "negative length", input: gen + ":-1:" + sha1ABC},
		{desc: "empty length", input: gen + "::" + sha1ABC},
		{desc: "length overflow", input: gen + ":99999999999999999999999:" + sha1ABC},
		{desc: "digest too short", input: gen + ":3:" + sha1ABC[:39]},
		{desc: "digest too long", input: gen + ":3:" + sha1ABC + "0"},
		{desc: "digest upper-case hex", input: gen + ":3:" + strings.ToUpper(sha1ABC)},
		{desc: "digest non-hex", input: gen + ":3:" + "z" + sha1ABC[1:]},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			g, n, d, ok := DecodeGenerationMarker(p.input)
			assert.False(t, ok)
			assert.Equal(t, "", g)
			assert.Equal(t, 0, n)
			assert.Equal(t, "", d)
		})
	}
}

func TestIsStaleGeneration(t *testing.T) {
	t.Parallel()
	gen1, err := FormatGeneration(1)
	require.NoError(t, err)
	gen2, err := FormatGeneration(2)
	require.NoError(t, err)
	gen3, err := FormatGeneration(3)
	require.NoError(t, err)
	abc := []byte("abc")

	patterns := []struct {
		desc          string
		valueExists   bool
		currentValue  []byte
		currentMarker string
		newGeneration string
		expected      bool
	}{
		{
			desc:          "no current value: always accept even if marker is newer",
			valueExists:   false,
			currentValue:  nil,
			currentMarker: EncodeGenerationMarker(gen3, nil),
			newGeneration: gen1,
			expected:      false,
		},
		{
			desc:          "value exists, no marker: accept",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: "",
			newGeneration: gen1,
			expected:      false,
		},
		{
			desc:          "value exists, malformed marker: accept",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: "garbage",
			newGeneration: gen1,
			expected:      false,
		},
		{
			desc:          "marker newer and describes current value: reject",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: EncodeGenerationMarker(gen2, abc),
			newGeneration: gen1,
			expected:      true,
		},
		{
			desc:          "marker newer but value length differs (legacy overwrite): accept",
			valueExists:   true,
			currentValue:  []byte("abcd"),
			currentMarker: EncodeGenerationMarker(gen2, abc),
			newGeneration: gen1,
			expected:      false,
		},
		{
			desc:          "marker newer, same length but different content (legacy overwrite): accept",
			valueExists:   true,
			currentValue:  []byte("abd"),
			currentMarker: EncodeGenerationMarker(gen2, abc),
			newGeneration: gen1,
			expected:      false,
		},
		{
			desc:          "equal generation: accept (idempotent rewrite)",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: EncodeGenerationMarker(gen2, abc),
			newGeneration: gen2,
			expected:      false,
		},
		{
			desc:          "newer generation: accept",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: EncodeGenerationMarker(gen2, abc),
			newGeneration: gen3,
			expected:      false,
		},
		{
			desc:          "empty current value described by newer marker: reject",
			valueExists:   true,
			currentValue:  []byte{},
			currentMarker: EncodeGenerationMarker(gen2, nil),
			newGeneration: gen1,
			expected:      true,
		},
		{
			desc:          "newer marker for empty value but current value non-empty: accept",
			valueExists:   true,
			currentValue:  abc,
			currentMarker: EncodeGenerationMarker(gen2, nil),
			newGeneration: gen1,
			expected:      false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := IsStaleGeneration(p.valueExists, p.currentValue, p.currentMarker, p.newGeneration)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestIsAllDigits(t *testing.T) {
	t.Parallel()
	assert.False(t, isAllDigits(""))
	assert.True(t, isAllDigits("0"))
	assert.True(t, isAllDigits("0123456789"))
	assert.False(t, isAllDigits("12a"))
	assert.False(t, isAllDigits("-1"))
	assert.False(t, isAllDigits(" 1"))
	assert.False(t, isAllDigits("1."))
}

func TestIsAllLowerHex(t *testing.T) {
	t.Parallel()
	assert.False(t, isAllLowerHex(""))
	assert.True(t, isAllLowerHex("0123456789abcdef"))
	assert.False(t, isAllLowerHex("ABCDEF"))
	assert.False(t, isAllLowerHex("0g"))
	assert.False(t, isAllLowerHex("a b"))
}
