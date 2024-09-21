package gounzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func extractFileToPath(file *zip.File, path string) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	// 打开目标文件，当存在时覆盖
	fileWriter, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if os.IsNotExist(err) {
		// panic(err)
		// 创建上级目录
		if err2 := os.MkdirAll(filepath.Dir(path), file.Mode()); err2 != nil {
			return err2
		}
		fileWriter, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	}
	if err != nil {
		// panic(err)
		return err
	}
	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		return err
	}

	return nil
}

func UnZip(dst, src string, encoding encoding.Encoding, parallel int, verbose bool) (err error) {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			// panic(err)
			return err
		}
	}

	var wg = make(chan error, parallel)
	for i := 0; i < parallel; i++ {
		wg <- nil
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		// 转换文件名
		filename, err := io.ReadAll(transform.NewReader(strings.NewReader(file.Name), encoding.NewDecoder()))
		if err != nil {
			return err
		}
		path := filepath.Join(dst, string(filename))

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			if verbose {
				fmt.Printf("%s\n", path)
			}
			continue
		}

		if err := <-wg; err != nil {
			wg <- err
			break
		}
		go func(f *zip.File) {
			wg <- extractFileToPath(f, path)
			if verbose {
				fmt.Printf("%s\n", path)
			}
		}(file)
	}
	// wait
	for i := 0; i < parallel; i++ {
		if err := <-wg; err != nil {
			return err
		}
	}
	return nil
}
