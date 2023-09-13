package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("non-existent directory", func(t *testing.T) {
		_, err := ReadDir("test")
		require.Error(t, err, "actual err - %v", err)
	})

	t.Run("key with '='", func(t *testing.T) {
		e, err := ReadDir("testdata/equal")
		require.NoError(t, err, "actual err - %v", err)
		require.Equalf(t, 0, len(e), "Invalid length of map")
	})

	t.Run("lowcase", func(t *testing.T) {
		e, err := ReadDir("testdata/lowcase")
		require.NoError(t, err, "actual err - %v", err)
		require.Equalf(t, "bar", e["bar"].Value, "No such env")
	})

	t.Run("tabs and spaces", func(t *testing.T) {
		e, err := ReadDir("testdata/tab-space")
		require.NoError(t, err, "actual err - %v", err)
		require.Equalf(t, "bar", e["SPACE"].Value, "Invalid value")
		require.Equalf(t, "bar", e["TAB"].Value, "Invalid value")
	})

	t.Run("from env dir", func(t *testing.T) {
		e, err := ReadDir("testdata/env")
		require.NoError(t, err, "actual err - %v", err)
		require.Equalf(t, "bar", e["BAR"].Value, "Invalid value")
		require.Equalf(t, "", e["EMPTY"].Value, "Invalid value")
		require.Equalf(t, "   foo\nwith new line", e["FOO"].Value, "Invalid value")
		require.Equalf(t, "\"hello\"", e["HELLO"].Value, "Invalid value")
	})

	t.Run("need remove", func(t *testing.T) {
		os.Setenv("bar", "b")
		e, err := ReadDir("testdata/lowcase")
		require.NoError(t, err, "actual err - %v", err)
		require.Equalf(t, true, e["bar"].NeedRemove, "Invalid value")
	})
}
