# fsplit

[![GoDoc](https://godoc.org/github.com/vbsw/fsplit?status.svg)](https://godoc.org/github.com/vbsw/fsplit) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/fsplit)](https://goreportcard.com/report/github.com/vbsw/fsplit) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
fsplit splits one file into many or combines them back to one. fsplit is published on <https://github.com/vbsw/fsplit>.

## Copyright
Copyright 2019, Vitali Baumtrok (vbsw@mailbox.org).

fsplit is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

fsplit is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

	fsplit [OPTIONS] [INPUT-FILE]

	OPTIONS
		-i=FILE    input file
		-o=FILE    output directory (for split) or file (for concatenate)
		-p=N       number of chunks
		-b=N[U]    size per chunk in bytes, U = unit (k/K, m/M or g/G)
		-l=N       number of lines per chunk
		-c         concatenate files (INPUT-FILE is only one file, the first one)

## References
- https://golang.org/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
