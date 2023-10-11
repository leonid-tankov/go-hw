package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	logger := New("ERROR", buffer)

	logger.Error("Error")
	require.Contains(t, buffer.String(), "Error")

	logger.Info("Info")
	require.NotContains(t, buffer.String(), "Info")

	logger.Debug("Debug")
	require.NotContains(t, buffer.String(), "Debug")
}
