package archive

import (
	"compress/gzip"
	"io"
)

type TarGz struct {
	gw *gzip.Writer
	tw *Tar
}

func TarGzCopy(source io.Reader, target io.Writer) (TarGz, error) {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	srcgz, err := gzip.NewReader(source)
	if err != nil {
		return TarGz{}, err
	}
	tw, err := TarCopy(srcgz, gw)
	return TarGz{
		gw: gw,
		tw: &tw,
	}, err
}

// NewTarGz  tar.gz archive.
func NewTarGz(target io.Writer) TarGz {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	tw := NewTar(gw)
	return TarGz{
		gw: gw,
		tw: &tw,
	}
}

// Close all closeables.
func (a TarGz) Close() error {
	if err := a.tw.Close(); err != nil {
		return err
	}
	return a.gw.Close()
}

// Add file to the archive.
func (a TarGz) Add(f File) error {
	return a.tw.Add(f)
}
