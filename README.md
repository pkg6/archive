## archive
A simple Go archiving library.

## Example usage

~~~
package main

import (
	"github.com/pkg6/archive"
	"log"
	"os"
)

func main() {
	file, err := os.Create("file.zip")
	if err != nil {
		// deal with the error
	}
	zip, err := archive.New(file, "zip")
	if err != nil {
		log.Fatal(err)
	}
	defer zip.Close()
	zip.Add(archive.File{
		Destination: "file.txt",
		Source:      "/path/to/file.txt",
	})
}
~~~

## Support compression format

- tar.gz
- tgz
- tar
- gz
- tar.xz
- txz
- tar.zst
- tzst
- zip
