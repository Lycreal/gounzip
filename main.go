package main

import (
	"flag"
	"fmt"
	gounzip "gounzip/gozunip"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

var zipfile, exdir, charset string
var parallel int
var verbose bool
var enc encoding.Encoding

const usage = `Usage: gounzip [options...] <file>
  -d <EXDIR>    extract files into exdir
  -O <CHARSET>  specify a character encoding for DOS, Windows and OS/2 archives
  -p <NUM>      set the number of parallel jobs to run
  -v            print file names while processing

unzip with support for filename encoding and parallel decompression.
source code: https://github.com/Lycreal/gounzip
`

func Init() {
	flag.StringVar(&exdir, "d", "", "extract files into exdir")
	flag.StringVar(&charset, "O", "UTF-8", "specific file name encoding")
	flag.IntVar(&parallel, "p", 1, "set the number of parallel jobs to run")
	flag.BoolVar(&verbose, "v", false, "print file names while processing")

	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()
	zipfile = flag.Arg(0)

	if zipfile == "" || len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	enc, err = ianaindex.MIB.Encoding(charset)
	if err != nil || enc == nil {
		fmt.Printf("error: Invalid charset.\n")
		os.Exit(1)
	}
}

func main() {
	Init()
	//numcpu
	err := gounzip.UnZip(exdir, zipfile, enc, parallel, verbose)
	if err != nil {
		fmt.Printf("error: %v", err.Error())
	}
}
