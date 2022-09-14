package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOutput(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		output, err := parseOutput(strings.NewReader("test name\nsupplementary message\n1.0 1.0\n\n"))
		require.NoError(t, err)
		assert.NotEmpty(t, output)
		assert.Len(t, output, 1)
		assert.Equal(t, output[0].Score, float64(1))
		assert.True(t, output[0].Pass)
	})
	t.Run("empty", func(t *testing.T) {
		output, err := parseOutput(strings.NewReader("name\n1.0 1.0\n\n"))
		require.NoError(t, err)
		assert.NotEmpty(t, output)
		assert.Len(t, output, 1)
		assert.Equal(t, output[0].Score, float64(1))
	})
}
