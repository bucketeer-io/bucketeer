package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc     string
		input    string
		expected *time.Location
		invalid  bool
	}{
		{
			desc:     "tokyo",
			input:    "Asia/Tokyo",
			expected: time.FixedZone("Asia/Tokyo", 9*60*60),
			invalid:  false,
		},
		{
			desc:     "UTC",
			input:    "UTC",
			expected: time.FixedZone("UTC", 0),
			invalid:  false,
		},
		{
			desc:     "invalid",
			input:    "invalid",
			expected: nil,
			invalid:  true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := &server{
				timezone: &p.input,
			}
			actual, err := s.getLocation()
			assert.Equal(t, actual, p.expected)
			if p.invalid {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
