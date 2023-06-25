package helpers

import (
	"errors"
	"testing"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/stretchr/testify/assert"
)

func TestParseClientError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		action   string
		name     string
		err      error
		expected string
	}{
		"openapi": {
			action:   "create",
			name:     "readarr_tag",
			err:      &readarr.GenericOpenAPIError{},
			expected: "Unable to create readarr_tag, got error: \nDetails:\n",
		},
		"generic": {
			action:   "create",
			name:     "readarr_tag",
			err:      errors.New("other error"),
			expected: "Unable to create readarr_tag, got error: other error",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseClientError(test.action, test.name, test.err))
		})
	}
}

func TestParseNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind     string
		field    string
		search   string
		expected string
	}{
		"generic": {
			kind:     "readarr_tag",
			field:    "label",
			search:   "test",
			expected: "Unable to find readarr_tag, got error: data source not found: no readarr_tag with label 'test'",
		},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ParseNotFoundError(test.kind, test.field, test.search))
		})
	}
}
