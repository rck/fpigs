[![Build Status](https://travis-ci.org/rck/fpigs.svg?branch=master)](https://travis-ci.org/rck/fpigs)

# fpigs
Find large files fast

## Why
While it is at least a useful tool for myself, it is basically a toy project to learn golang.

# Synopsis

```
Usage of fpigs (fc85e1e-dirty):
  -c	Files from current directory only (no recursion)
  -concurrent number
    	Start this number of concurrent tree walks (values <= 0 get set to 1) (default 20)
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

# Docker
```
docker pull rck81/fpigs
docker run -v $HOME:/fpigs -w /fpigs -it --rm rck81/fpigs
```
