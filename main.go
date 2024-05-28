package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var zipfile, exdir, encoding string

func Init() {
	flag.StringVar(&exdir, "d", "", "extract files into exdir")
	flag.StringVar(&encoding, "e", "UTF-8", "specific file name encoding")

	flag.Usage = func() {
		fmt.Print(`Usage: gounzip [-d exdir] [-e encoding] zipfile
  -d exdir     extract files into exdir
  -e encoding  specific file name encoding, support UTF-8(default), Shift_JIS
  zipfile      zip file to be extracted
`)
	}
	flag.Parse()
	zipfile = flag.Arg(0)

	if zipfile == "" || len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	}
}

func extractFileToPath(file *zip.File, path string) error {
	if file.FileInfo().IsDir() {
		err := os.MkdirAll(path, file.Mode())
		if err != nil {
			return err
		}
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	// 打开目标文件，当存在时覆盖
	fileWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if os.IsNotExist(err) {
		// 创建上级目录
		if err2 := os.MkdirAll(filepath.Dir(path), file.Mode()); err2 != nil {
			return err2
		}
		fileWriter, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	}
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		return err
	}

	return nil
}

// Thanks to: https://broqiang.com/posts/archive-zip
func UnZip(dst, src string) (err error) {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		var filename string

		// 修复日文文件名乱码
		if strings.EqualFold(encoding, "Shift_JIS") {
			b, err := io.ReadAll(transform.NewReader(strings.NewReader(file.Name), japanese.ShiftJIS.NewDecoder()))
			if err != nil {
				return err
			}
			filename = string(b)
		} else {
			filename = file.Name
		}

		path := filepath.Join(dst, filename)

		err = extractFileToPath(file, path)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", path)
	}
	return nil
}

func main() {
	Init()

	err := UnZip(exdir, zipfile)
	if err != nil {
		fmt.Printf("error: %v", err.Error())
	}
}
