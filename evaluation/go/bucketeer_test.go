package evaluation

import (
	"testing"
)

func TestMurmur128(t *testing.T) {
	b := bucketeer{}
	input := "fid-uid-sampling-seed"

	// Use reflection to access the private murmur128 method
	high, low := b.murmur128(input)

	expectedHigh := uint64(2548757552806388169)
	expectedLow := uint64(9787172855444729749)

	if high != expectedHigh {
		t.Fatalf("Expected high %d, but got %d", expectedHigh, high)
	}
	if low != expectedLow {
		t.Fatalf("Expected low %d, but got %d", expectedLow, low)
	}
}

func TestToFloat64(t *testing.T) {
	b := bucketeer{}
	high := uint64(2548757552806388169)
	low := uint64(9787172855444729749)

	// Calculate the normalized value
	normalized := b.toFloat64(high, low)

	expectedNormalized := 0.1381684237945762

	if normalized != expectedNormalized {
		t.Fatalf("Expected normalized value %f, but got %f", expectedNormalized, normalized)
	}
}
