package archive

import (
	"archive/tar"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg6/archive/testlib"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTarZstFile(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.Create(filepath.Join(tmp, "test.tar.zst"))
	require.NoError(t, err)
	defer f.Close()
	archive := NewTarZst(f)
	defer archive.Close()

	require.Error(t, archive.Add(File{
		Source:      "testdata/nope.txt",
		Destination: "nope.txt",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/foo.txt",
		Destination: "foo.txt",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/sub1",
		Destination: "sub1",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/sub1/bar.txt",
		Destination: "sub1/bar.txt",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/sub1/executable",
		Destination: "sub1/executable",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/sub1/sub2",
		Destination: "sub1/sub2",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/sub1/sub2/subfoo.txt",
		Destination: "sub1/sub2/subfoo.txt",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/regular.txt",
		Destination: "regular.txt",
	}))
	require.NoError(t, archive.Add(File{
		Source:      "testdata/link.txt",
		Destination: "link.txt",
	}))

	require.NoError(t, archive.Close())
	require.Error(t, archive.Add(File{
		Source:      "tar.go",
		Destination: "tar.go",
	}))
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	require.Lessf(t, info.Size(), int64(500), "archived file should be smaller than %d", info.Size())

	zstf, err := zstd.NewReader(f)
	require.NoError(t, err)

	var paths []string
	r := tar.NewReader(zstf)
	for {
		next, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		paths = append(paths, next.Name)
		if testlib.IsWindows() {
			// both of the following checks don't work on windows.
			continue
		}
		if next.Name == "sub1/executable" {
			require.NotEqualf(
				t,
				0,
				next.FileInfo().Mode()&0o111,
				"expected executable perms, got %s",
				next.FileInfo().Mode().String(),
			)
		}
		if next.Name == "link.txt" {
			require.Equal(t, "regular.txt", next.Linkname)
		}
	}
	require.Equal(t, []string{
		"foo.txt",
		"sub1",
		"sub1/bar.txt",
		"sub1/executable",
		"sub1/sub2",
		"sub1/sub2/subfoo.txt",
		"regular.txt",
		"link.txt",
	}, paths)
}

func TestTarZstFileInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	f, err := os.Create(filepath.Join(t.TempDir(), "test.tar.gz"))
	require.NoError(t, err)
	defer f.Close()
	archive := NewTarZst(f)
	defer archive.Close()

	require.NoError(t, archive.Add(File{
		Source:      "testdata/foo.txt",
		Destination: "nope.txt",
		Info: FileInfo{
			Mode:        0o755,
			Owner:       "carlos",
			Group:       "root",
			ParsedMTime: now,
		},
	}))

	require.NoError(t, archive.Close())
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	zstf, err := zstd.NewReader(f)
	require.NoError(t, err)

	var found int
	r := tar.NewReader(zstf)
	for {
		next, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		found++
		require.Equal(t, "nope.txt", next.Name)
		require.Equal(t, now, next.ModTime)
		require.Equal(t, fs.FileMode(0o755), next.FileInfo().Mode())
		require.Equal(t, "carlos", next.Uname)
		require.Equal(t, 0, next.Uid)
		require.Equal(t, "root", next.Gname)
		require.Equal(t, 0, next.Gid)
	}
	require.Equal(t, 1, found)
}
