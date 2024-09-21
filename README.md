# gounzip
unzip in Golang

unzip with support for filename encoding and parallel decompression.

## Usage
```text
gounzip [options...] <file>
  -d <EXDIR>    extract files into exdir
  -O <CHARSET>  specify a character encoding for filenames
  -p <NUM>      set the number of parallel jobs to run
  -v            print file names while processing
```

## Credits
Inspired by [Debian Unzip](https://tracker.debian.org/pkg/unzip).
