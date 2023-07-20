package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	cmp := equalfile.New(nil, equalfile.Options{})

	t.Run("negative argument", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "tmp_file.txt")
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err := Copy("testdata/input.txt", tmpFile.Name(), -1, 0)
		require.Truef(t, errors.Is(err, ErrNegativeArgument), "actual error %q", err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "tmp_file.txt")
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err := Copy("testdata", tmpFile.Name(), 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "tmp_file.txt")
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		err := Copy("testdata/input.txt", tmpFile.Name(), 10000, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("files in testdata", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "tmp_file.txt")
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 0, 0), "failed copy file")
		equal, _ := cmp.CompareFile(tmpFile.Name(), "testdata/out_offset0_limit0.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 0, 10), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset0_limit10.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 0, 1000), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset0_limit1000.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 0, 10000), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset0_limit10000.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 10, 0), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset10_limit0.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 100, 1000), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset100_limit1000.txt")
		require.Truef(t, equal, "files not equal")

		require.NoError(t, Copy("testdata/input.txt", tmpFile.Name(), 6000, 1000), "failed copy file")
		equal, _ = cmp.CompareFile(tmpFile.Name(), "testdata/out_offset6000_limit1000.txt")
		require.Truef(t, equal, "files not equal")
	})
}
