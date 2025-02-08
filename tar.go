package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Tar as tar.
type Tar struct {
	tw    *tar.Writer
	files map[string]bool
}

// NewTar tar archive.
func NewTar(target io.Writer) Tar {
	return Tar{
		tw:    tar.NewWriter(target),
		files: map[string]bool{},
	}
}

// TarCopy creates a new tar with the contents of the given tar.
func TarCopy(source io.Reader, target io.Writer) (Tar, error) {
	w := NewTar(target)
	r := tar.NewReader(source)
	for {
		header, err := r.Next()
		if err == io.EOF || header == nil {
			break
		}
		if err != nil {
			return Tar{}, err
		}
		w.files[header.Name] = true
		if err := w.tw.WriteHeader(header); err != nil {
			return w, err
		}
		if _, err := io.Copy(w.tw, r); err != nil {
			return w, err
		}
	}
	return w, nil
}

// Close all closeables.
func (a Tar) Close() error {
	return a.tw.Close()
}

// Add file to the archive.
func (a Tar) Add(f File) error {
	if _, ok := a.files[f.Destination]; ok {
		return &fs.PathError{Err: fs.ErrExist, Path: f.Destination, Op: "add"}
	}
	a.files[f.Destination] = true
	info, err := os.Lstat(f.Source) // #nosec
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	var link string
	if info.Mode()&os.ModeSymlink != 0 {
		link, err = os.Readlink(f.Source) // #nosec
		if err != nil {
			return fmt.Errorf("%s: %w", f.Source, err)
		}
	}
	header, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	header.Name = f.Destination
	if !f.Info.ParsedMTime.IsZero() {
		header.ModTime = f.Info.ParsedMTime
	}
	if f.Info.Mode != 0 {
		header.Mode = int64(f.Info.Mode)
	}
	if f.Info.Owner != "" {
		header.Uid = 0
		header.Uname = f.Info.Owner
	}
	if f.Info.Group != "" {
		header.Gid = 0
		header.Gname = f.Info.Group
	}
	if err = a.tw.WriteHeader(header); err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	file, err := os.Open(f.Source) // #nosec
	if err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	defer file.Close()
	if _, err := io.Copy(a.tw, file); err != nil {
		return fmt.Errorf("%s: %w", f.Source, err)
	}
	return nil
}
