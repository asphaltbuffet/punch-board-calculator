package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvelopeCommand(t *testing.T) {
	got := NewEnvelopeCommand()

	assert.Equal(t, "envelope", got.Name())
	assert.True(t, got.Runnable())
}
