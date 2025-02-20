// Copyright 2025 The Bucketeer Authors.
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

package evaluation

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type bucketeer struct{}

// Calculate the input hash of the target property and map it to a float64 between [0,1]
func (b *bucketeer) bucket(input string) float64 {
	// Compute MurmurHash3 (128-bit) hash.
	// The high and low parts of the hash are returned in big-endian format,
	// which is important for consistency with the Node.js implementation that also uses big-endian.
	high, low := murmur3.Sum128([]byte(input))

	// Convert to float64 using both 64-bit parts (high and low).
	// This maintains the full 128-bit range, ensuring a normalized value between [0, 1].
	return b.toFloat64(high, low)
}

// Convert a 128-bit hash (two uint64 values) to a float64 in the range [0,1]
// Because we bucket millions of users, we use the full 128-bit hash to avoid collisions.
func (*bucketeer) toFloat64(high, low uint64) float64 {
	// Combine the high and low parts into a single floating-point number.
	// This utilizes the full 128-bit range for precision.
	full := (float64(high) * math.Pow(2, 64)) + float64(low)

	// Calculate the maximum value for a 128-bit number.
	maxValue := math.Pow(2, 128) - 1 // Maximum value for 128-bit

	// Normalize the combined value to the range [0,1].
	return full / maxValue
}
