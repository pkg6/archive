package archive

import (
	"fmt"
	"io"
	"os"
)

// Archive represents a compression archive files from disk can be written to.
type Archive interface {
	Close() error
	Add(f File) error
}

// New archive.
func New(w io.Writer, format string) (Archive, error) {
	switch format {
	case "tar.gz", "tgz":
		return NewTarGz(w), nil
	case "tar":
		return NewTar(w), nil
	case "gz":
		return NewGzip(w), nil
	case "tar.xz", "txz":
		return NewTarXz(w), nil
	case "tar.zst", "tzst":
		return NewTarZst(w), nil
	case "zip":
		return NewZip(w), nil
	}
	return nil, fmt.Errorf("invalid archive format: %s", format)
}

// Copy copies the source archive into a new one, which can be appended at.
// Source needs to be in the specified format.
func Copy(r *os.File, w io.Writer, format string) (Archive, error) {
	switch format {
	case "tar.gz", "tgz":
		return TarGzCopy(r, w)
	case "tar":
		return TarCopy(r, w)
	case "zip":
		return ZipCopy(r, w)
	}
	return nil, fmt.Errorf("invalid archive format: %s", format)
}
