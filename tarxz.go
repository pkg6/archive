package archive

import (
	"github.com/ulikunitz/xz"
	"io"
)

// TarXz as tar.xz.
type TarXz struct {
	xzw *xz.Writer
	tw  *Tar
}

// NewTarXz tar.xz archive.
func NewTarXz(target io.Writer) TarXz {
	xzw, _ := xz.WriterConfig{DictCap: 16 * 1024 * 1024}.NewWriter(target)
	tw := NewTar(xzw)
	return TarXz{
		xzw: xzw,
		tw:  &tw,
	}
}

// Close all closeables.
func (a TarXz) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.xzw.Close()
}

// Add file to the archive.
func (a TarXz) Add(f File) error {
	return a.tw.Add(f)
}
