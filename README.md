[![Build Status](https://travis-ci.org/rck/fpigs.svg?branch=master)](https://travis-ci.org/rck/fpigs)
[![Docker Automated build](https://img.shields.io/docker/automated/rck/fpigs.svg)](https://hub.docker.com/r/rck81/fpigs/)

# fpigs
Find large files fast

## Why
While it is at least a useful tool for myself, it is basically a toy project to learn golang.

# Synopsis

```
Usage: fpigs [OPTION]... [STARTDIR]...
  -concurrent number
    	Start this number of concurrent tree walks (values <= 0 get set to 1) (default 20)
  -d depth
    	Recursion depth (negative values mean unlimited) (default -1)
  -i regex
    	Ignore files/directories matching regex (can be used multiple times)
  -n number
    	Print this number of largest files (default 10)
  -o	Print file names only (e.g., for xargs scripting)
  -u unit
    	Print sizes in specified unit (B, KiB, MiB, GiB, TiB, KB, MB, GB, TB) (default GiB)
  -version
    	Print version and exit
```

# Example output

```
$ ./fpigs -u KB -i ".git$"
Press <return> to stop processing

2467.8 KB: fpigs
[##############################################################################]

35.1 KB: LICENSE
[#]

5.2 KB: main.go
[]

1.4 KB: flags.go
[]

1.0 KB: flags_test.go
[]

0.7 KB: README.md
[]

0.7 KB: Makefile
[]

7 largest files: 2512.0 KB / 2512.0 KB total
```

# Regex filtering
Filters are applied when:
- A directory is traversed (This allows cutting of ".git$")
- A file is inserted in the current list of the N largest files.

# Releases
Pre-built binaries are provided [here](https://github.com/rck/fpigs/releases/latest). Please note that these
binaries are automatically built by [Travis-CI](https://travis-ci.org). Your decision if you trust them.

# Docker
```
docker pull rck81/fpigs
docker run -v $HOME:/fpigs -w /fpigs -it --rm rck81/fpigs
```
