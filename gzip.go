// Package archive implements the Gzip interface providing gz archiving
// and compression.
package archive

import (
	"fmt"
	gzip "github.com/klauspost/pgzip"
	"io"
	"os"
)

// Gzip as gz.
type Gzip struct {
	gw *gzip.Writer
}

// NewGzip gz archive.
func NewGzip(target io.Writer) Archive {
	// the error will be nil since the compression level is valid
	gw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	return Gzip{gw: gw}
}

// Close all closeables.
func (a Gzip) Close() error {
	return a.gw.Close()
}

// Add file to the archive.
func (a Gzip) Add(f File) error {
	if a.gw.Header.Name != "" {
		return fmt.Errorf("gzip: failed to add %s, only one file can be archived in gz format", f.Destination)
	}
	file, err := os.Open(f.Source) // #nosec
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	a.gw.Header.Name = f.Destination
	if f.Info.ParsedMTime.IsZero() {
		a.gw.Header.ModTime = info.ModTime()
	} else {
		a.gw.Header.ModTime = f.Info.ParsedMTime
	}
	_, err = io.Copy(a.gw, file)
	return err
}
