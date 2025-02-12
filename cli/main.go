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
	Force    bool
}

func init() {
	flag.StringVar(&flags.FileName, "f", "archive.zip", "<filename>  Location of cli")
	flag.StringVar(&flags.Dir, "C", ".", "Change to <dir> before processing remaining files")
	flag.BoolVar(&flags.Force, "force", false, "Do you want to perform forced compression")
}
func main() {
	flag.Parse()
	tempPath := path.Join(os.TempDir(), "archive")
	tempFileName := path.Join(tempPath, flags.FileName)
	defer func() {
		_ = os.RemoveAll(tempPath)
	}()
	if err := os.MkdirAll(tempPath, 0o755|os.ModeDir); err != nil {
		log.Println(err)
		return
	}
	format := strings.Replace(path.Ext(flags.FileName), ".", "", 1)
	file, err := os.Create(tempFileName)
	if err != nil {
		log.Println(err)
		return
	}
	myArchive, err := archive.New(file, format)
	if err != nil {
		log.Println(err)
		return
	}
	defer myArchive.Close()
	if err := filepath.Walk(flags.Dir, func(path string, info os.FileInfo, err error) (rErr error) {
		destination := strings.TrimPrefix(path, flags.Dir)
		rErr = myArchive.Add(archive.File{
			Destination: destination,
			Source:      path,
		})
		if rErr == nil {
			fmt.Println(fmt.Sprintf("Adding: %s ", destination))
			return nil
		}
		if flags.Force {
			return nil
		}
		return rErr
	}); err != nil {
		log.Println(err)
		return
	}
	if err := os.Rename(tempFileName, flags.FileName); err != nil {
		log.Println(err)
		return
	}
}
