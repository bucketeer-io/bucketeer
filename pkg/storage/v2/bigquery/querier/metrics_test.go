package querier

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/googleapi"
)

func TestGetCodeFromError(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		input    error
		expected string
	}{
		{
			desc:     "codeOK",
			input:    nil,
			expected: codeOK,
		},
		{
			desc:     "codeUnknown: error is not googleapi.Error",
			input:    errors.New("this is not googleapi.Error"),
			expected: codeUnknown,
		},
		{
			desc:     "codeUnknown: error code is not unexpected",
			input:    &googleapi.Error{Code: http.StatusVariantAlsoNegotiates},
			expected: codeUnknown,
		},
		{
			desc:     "codeBadRequest",
			input:    &googleapi.Error{Code: http.StatusBadRequest},
			expected: codeBadRequest,
		},
		{
			desc:     "codeForbidden",
			input:    &googleapi.Error{Code: http.StatusForbidden},
			expected: codeForbidden,
		},
		{
			desc:     "codeNotFound",
			input:    &googleapi.Error{Code: http.StatusNotFound},
			expected: codeNotFound,
		},
		{
			desc:     "codeConflict",
			input:    &googleapi.Error{Code: http.StatusConflict},
			expected: codeConflict,
		},
		{
			desc:     "codeInternalServerError",
			input:    &googleapi.Error{Code: http.StatusInternalServerError},
			expected: codeInternalServerError,
		},
		{
			desc:     "codeNotImplemented",
			input:    &googleapi.Error{Code: http.StatusNotImplemented},
			expected: codeNotImplemented,
		},
		{
			desc:     "codeServiceUnavailable",
			input:    &googleapi.Error{Code: http.StatusServiceUnavailable},
			expected: codeServiceUnavailable,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getCodeFromError(p.input)
			assert.Equal(t, p.expected, actual, "%s", p.desc)
		})
	}
}
