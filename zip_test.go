package archive

import (
	"archive/zip"
	"bytes"
	"github.com/pkg6/archive/testlib"
	"github.com/stretchr/testify/require"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestZipFile(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.Create(filepath.Join(tmp, "test.zip"))
	require.NoError(t, err)
	defer f.Close()
	archive := NewZip(f)
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

	require.ErrorIs(t, archive.Add(File{
		Source:      "testdata/regular.txt",
		Destination: "link.txt",
	}), fs.ErrExist)

	require.NoError(t, archive.Close())
	require.NoError(t, f.Close())

	f, err = os.Open(f.Name())
	require.NoError(t, err)
	defer f.Close()

	info, err := f.Stat()
	require.NoError(t, err)
	require.Lessf(t, info.Size(), int64(1000), "archived file should be smaller than %d", info.Size())

	r, err := zip.NewReader(f, info.Size())
	require.NoError(t, err)

	paths := make([]string, len(r.File))
	for i, zf := range r.File {
		paths[i] = zf.Name
		if zf.Name == "sub1/executable" && !testlib.IsWindows() {
			require.NotEqualf(
				t,
				0,
				zf.Mode()&0o111,
				"expected executable perms, got %s",
				zf.Mode().String(),
			)
		}
		if zf.Name == "link.txt" {
			require.NotEqual(t, 0, zf.FileInfo().Mode()&os.ModeSymlink)
			rc, err := zf.Open()
			require.NoError(t, err)
			var link bytes.Buffer
			_, err = io.Copy(&link, rc)
			require.NoError(t, err)
			rc.Close()
			require.Equal(t, "regular.txt", link.String())
		}
	}
	require.Equal(t, []string{
		"foo.txt",
		"sub1/bar.txt",
		"sub1/executable",
		"sub1/sub2/subfoo.txt",
		"regular.txt",
		"link.txt",
	}, paths)
}

func TestZipFileInfo(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	f, err := os.Create(filepath.Join(t.TempDir(), "test.zip"))
	require.NoError(t, err)
	defer f.Close()
	archive := NewZip(f)
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

	info, err := f.Stat()
	require.NoError(t, err)

	r, err := zip.NewReader(f, info.Size())
	require.NoError(t, err)

	require.Len(t, r.File, 1)
	for _, next := range r.File {
		require.Equal(t, "nope.txt", next.Name)
		require.Equal(t, now.Unix(), next.Modified.Unix())
		require.Equal(t, fs.FileMode(0o755), next.FileInfo().Mode())
	}
}
