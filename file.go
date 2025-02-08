package archive

import (
	"os"
	"time"
)

type File struct {
	Source      string   `yaml:"src,omitempty" json:"src,omitempty"`
	Destination string   `yaml:"dst,omitempty" json:"dst,omitempty"`
	StripParent bool     `yaml:"strip_parent,omitempty" json:"strip_parent,omitempty"`
	Info        FileInfo `yaml:"info,omitempty" json:"info,omitempty"`
	Default     bool     `yaml:"-" json:"-"`
}

// FileInfo is the file info of a file.
type FileInfo struct {
	Owner       string      `yaml:"owner,omitempty" json:"owner,omitempty"`
	Group       string      `yaml:"group,omitempty" json:"group,omitempty"`
	Mode        os.FileMode `yaml:"mode,omitempty" json:"mode,omitempty"`
	MTime       string      `yaml:"mtime,omitempty" json:"mtime,omitempty"`
	ParsedMTime time.Time   `yaml:"-" json:"-"`
}
