package archive

import (
	"github.com/klauspost/compress/zstd"
	"io"
)

// TarZst as tar.zst.
type TarZst struct {
	zstw *zstd.Encoder
	tw   *Tar
}

// NewTarZst tar.zst archive.
func NewTarZst(target io.Writer) TarZst {
	zstw, _ := zstd.NewWriter(target)
	tw := NewTar(zstw)
	return TarZst{
		zstw: zstw,
		tw:   &tw,
	}
}

// Close all closeables.
func (a TarZst) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.zstw.Close()
}

// Add file to the archive.
func (a TarZst) Add(f File) error {
	return a.tw.Add(f)
}
