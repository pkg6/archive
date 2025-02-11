package main

import (
	"flag"
	"fmt"
	"github.com/pkg6/archive"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var flags Flags

type Flags struct {
	FileName string
	Dir      string
}

func init() {
	flag.StringVar(&flags.FileName, "f", "archive.zip", "<filename>  Location of cli")
	flag.StringVar(&flags.Dir, "C", ".", "Change to <dir> before processing remaining files")
}
func main() {
	flag.Parse()
	tempPath := path.Join(os.TempDir(), "archive")
	tempFileName := path.Join(tempPath, flags.FileName)
	defer func() {
		_ = os.RemoveAll(tempPath)
	}()
	if err := os.MkdirAll(tempPath, 0o755|os.ModeDir); err != nil {
		log.Fatal(err)
		return
	}
	format := strings.Replace(path.Ext(flags.FileName), ".", "", 1)
	file, err := os.Create(tempFileName)
	if err != nil {
		_ = os.Remove(flags.FileName)
		log.Fatal(err)
		return
	}
	myArchive, err := archive.New(file, format)
	if err != nil {
		_ = os.Remove(flags.FileName)
		log.Fatal(err)
		return
	}
	defer myArchive.Close()
	if err := filepath.Walk(flags.Dir, func(path string, info os.FileInfo, err error) error {
		destination := strings.TrimPrefix(path, flags.Dir)
		fmt.Println(fmt.Sprintf("adding: %s", destination))
		return myArchive.Add(archive.File{
			Destination: destination,
			Source:      path,
		})
	}); err != nil {
		_ = os.Remove(flags.FileName)
		log.Fatal(err)
		return
	}
	if err := os.Rename(tempFileName, flags.FileName); err != nil {
		_ = os.Remove(flags.FileName)
		log.Fatal(err)
	}
}
