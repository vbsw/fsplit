/*
 *      Copyright 2019 - 2021, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"os"
	"strconv"
)

type command struct {
	id         commandID
	message    string
	inputFile  string
	outputFile string
	parts      int64
	bytes      int64
	lines      int64
}

type commandID int

const (
	none   commandID = 0
	info   commandID = 1
	split  commandID = 2
	concat commandID = 3
	wrong  commandID = 4
)

func commandFromCommandLine() *command {
	cmd := new(command)
	params, err := parametersFromOSArgs()

	if err == nil {

		if params == nil {
			cmd.setShortInfo()

		} else if params.incompatibleParameters() {
			cmd.setWrongArgumentUsage()

		} else if params.oneParamHasMultipleResults() {
			cmd.setWrongArgumentUsage()

		} else {
			cmd.setValidCommand(params)
		}
	} else {
		cmd.id = wrong
		cmd.message = err.Error()
	}
	return cmd
}

func (cmd *command) setShortInfo() {
	cmd.id = info
	cmd.message = "Run 'fsplit --help' for usage."
}

func (cmd *command) setWrongArgumentUsage() {
	cmd.id = wrong
	cmd.message = "wrong argument usage"
}

func (cmd *command) setValidCommand(params *parameters) {
	if params.help.Available() {
		cmd.setHelp()

	} else if params.version.Available() {
		cmd.setVersion()

	} else if params.copyright.Available() {
		cmd.setCopyright()

	} else if params.concat.Available() {
		cmd.interpretInputForConcat(params)
		cmd.interpretOutput(params)
		cmd.interpretParts(params)
		cmd.interpretBytes(params)
		cmd.interpretLines(params)

		if cmd.id == none {
			cmd.id = concat
		}
	} else {
		cmd.interpretInputForSplit(params)
		cmd.interpretOutput(params)
		cmd.interpretParts(params)
		cmd.interpretBytes(params)
		cmd.interpretLines(params)
		cmd.interpretDefaultSplitParts()

		if cmd.id == none {
			cmd.id = split
		}
	}
}

func (cmd *command) setHelp() {
	cmd.id = info
	cmd.message = "fsplit splits files into many, or combines them back to one\n\n"
	cmd.message = cmd.message + "USAGE\n"
	cmd.message = cmd.message + "  fsplit ( INFO | SPLIT/CONCATENATE )\n\n"
	cmd.message = cmd.message + "INFO\n"
	cmd.message = cmd.message + "  -h, --help    print this help\n"
	cmd.message = cmd.message + "  -v, --version print version\n"
	cmd.message = cmd.message + "  --copyright   print copyright\n\n"
	cmd.message = cmd.message + "SPLIT/CONCATENATE\n"
	cmd.message = cmd.message + "  [COMMAND] INPUT-FILE [OUTPUT-FILE]\n\n"
	cmd.message = cmd.message + "COMMAND\n"
	cmd.message = cmd.message + "  -p=N          split file into N parts (chunks)\n"
	cmd.message = cmd.message + "  -b=N[U]       split file into N bytes per chunk, U = unit (k/K, m/M or g/G)\n"
	cmd.message = cmd.message + "  -l=N          split file into N lines per chunk\n"
	cmd.message = cmd.message + "  -c            concatenate files (INPUT-FILE is only one file, the first one)"
}

func (cmd *command) setVersion() {
	cmd.id = info
	cmd.message = "1.0.0"
}

func (cmd *command) setCopyright() {
	cmd.id = info
	cmd.message = "Copyright 2019, 2020, Vitali Baumtrok (vbsw@mailbox.org).\n"
	cmd.message = cmd.message + "Distributed under the Boost Software License, Version 1.0."
}

func (cmd *command) interpretInputForConcat(params *parameters) {
	if cmd.id == none {
		if params.input.Available() {
			inputFile := params.input.Values[0]

			if direcotryExists(inputFile) {
				cmd.id = wrong
				cmd.message = "input is a directory, but must be a file"

			} else {
				cmd.inputFile = inputFile
			}

		} else {
			cmd.id = wrong
			cmd.message = "input file is not specified"
		}
	}
}

func (cmd *command) interpretOutput(params *parameters) {
	if cmd.id == none {
		if params.output.Available() {
			cmd.outputFile = params.output.Values[0]
		} else {
			cmd.outputFile = cmd.inputFile
		}
	}
}

func (cmd *command) interpretParts(params *parameters) {
	if cmd.id == none {
		if params.parts.Available() {
			parts, err := strconv.Atoi(params.parts.Values[0])
			if err == nil {
				cmd.parts = abs(int64(parts))
			} else {
				cmd.id = wrong
				cmd.message = "can't parse number of parts"
			}
		}
	}
}

func (cmd *command) interpretBytes(params *parameters) {
	if cmd.id == none {
		if params.bytes.Available() {
			bytes, err := parseBytes(params.bytes.Values[0])
			if err == nil {
				cmd.bytes = abs(bytes)
			} else {
				cmd.id = wrong
				cmd.message = "can't parse number of bytes"
			}
		}
	}
}

func (cmd *command) interpretLines(params *parameters) {
	if cmd.id == none {
		if params.lines.Available() {
			lines, err := strconv.Atoi(params.lines.Values[0])
			if err == nil {
				cmd.lines = abs(int64(lines))
			} else {
				cmd.id = wrong
				cmd.message = "can't parse number of lines"
			}
		}
	}
}

func (cmd *command) interpretInputForSplit(params *parameters) {
	if cmd.id == none {
		if params.input.Available() {
			inputFile := params.input.Values[0]

			if isFile(inputFile) {
				cmd.inputFile = inputFile

			} else if isDirectory(inputFile) {
				cmd.id = wrong
				cmd.message = "input is a directory, but must be a file"

			} else {
				cmd.id = wrong
				cmd.message = "input file does not exist"
			}

		} else {
			cmd.id = wrong
			cmd.message = "input file is not specified"
		}
	}
}

func (cmd *command) interpretDefaultSplitParts() {
	if cmd.parts == 0 && cmd.bytes == 0 && cmd.lines == 0 {
		cmd.parts = 2
	}
}

func parseBytes(bytesStr string) (int64, error) {
	var bytes64 int64
	var err error

	if len(bytesStr) > 0 {
		var bytes int
		bytesArray := []byte(bytesStr)
		lastByte := bytesArray[len(bytesArray)-1]

		if lastByte == 'k' || lastByte == 'K' || lastByte == 'm' || lastByte == 'M' || lastByte == 'g' || lastByte == 'G' {
			bytesStr = bytesStr[:len(bytesStr)-1]

		} else {
			lastByte = 0
		}
		bytes, err = strconv.Atoi(bytesStr)

		if err == nil {
			switch lastByte {
			case 'k':
				bytes64 = int64(bytes) * 1024
			case 'K':
				bytes64 = int64(bytes) * 1000
			case 'm':
				bytes64 = int64(bytes) * 1024 * 1024
			case 'M':
				bytes64 = int64(bytes) * 1000 * 1000
			case 'g':
				bytes64 = int64(bytes) * 1024 * 1024 * 1024
			case 'G':
				bytes64 = int64(bytes) * 1000 * 1000 * 1000
			default:
				bytes64 = int64(bytes)
			}
		}
	}
	return bytes64, err
}

func direcotryExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return (err == nil || !os.IsNotExist(err)) && fileInfo.IsDir()
}

func abs(value int64) int64 {
	if value > 0 {
		return value
	}
	return -value
}
