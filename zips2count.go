package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"iter"
	"log"
	"os"
)

func zip2count(zrdr *zip.Reader) int {
	return len(zrdr.File)
}

type CountInfo struct {
	ZipFilename   string `json:"zip_filename"`
	NumberOfItems int    `json:"number_of_items"`
	IsValidZip    bool   `json:"is_valid_zip"`
}

type ZipFileLike struct {
	io.ReaderAt
	Size int64
}

func (l ZipFileLike) ToZip() (*zip.Reader, error) {
	return zip.NewReader(l.ReaderAt, l.Size)
}

func (l ZipFileLike) ToCountInfo(zipFilename string) (c CountInfo) {
	c.ZipFilename = zipFilename

	rdr, e := l.ToZip()
	switch e {
	case nil:
		c.NumberOfItems = zip2count(rdr)
		c.IsValidZip = true
	default:
		c.NumberOfItems = 0
		c.IsValidZip = false
	}

	return
}

type OsFile struct{ *os.File }

func (f OsFile) Close() error { return f.File.Close() }

func (f OsFile) Stat() (fs.FileInfo, error) { return f.File.Stat() }

func (f OsFile) Size() (int64, error) {
	s, e := f.Stat()
	if nil != e {
		return 0, e
	}
	return s.Size(), nil
}

func (f OsFile) AsReaderAt() io.ReaderAt { return f.File }

func (f OsFile) ToZipFileLike() (ZipFileLike, error) {
	sz, e := f.Size()
	return ZipFileLike{
		ReaderAt: f.AsReaderAt(),
		Size:     sz,
	}, e
}

type OsFilename string

func (n OsFilename) ToFile() (OsFile, error) {
	f, e := os.Open(string(n))
	return OsFile{File: f}, e
}

func (n OsFilename) ToCountInfo() (c CountInfo) {
	c.ZipFilename = string(n)
	f, e := n.ToFile()
	if nil != e {
		return c
	}
	defer f.Close()

	zfl, e := f.ToZipFileLike()
	if nil != e {
		return c
	}
	return zfl.ToCountInfo(string(n))
}

func filenames2info2json2writer(
	filenames iter.Seq[string],
	writer io.Writer,
) error {
	var enc *json.Encoder = json.NewEncoder(writer)
	for fname := range filenames {
		oname := OsFilename(fname)
		var cnt CountInfo = oname.ToCountInfo()
		e := enc.Encode(cnt)
		if nil != e {
			return e
		}
	}

	return nil
}

func reader2filenames(rdr io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		var s *bufio.Scanner = bufio.NewScanner(rdr)
		for s.Scan() {
			var filename string = s.Text()
			if !yield(filename) {
				return
			}
		}
	}
}

func stdin2names2info2json2stdout() error {
	var filenames iter.Seq[string] = reader2filenames(os.Stdin)
	var bw *bufio.Writer = bufio.NewWriter(os.Stdout)

	return errors.Join(
		filenames2info2json2writer(filenames, bw),
		bw.Flush(),
	)
}

func main() {
	e := stdin2names2info2json2stdout()
	if nil != e {
		log.Printf("%v\n", e)
	}
}
