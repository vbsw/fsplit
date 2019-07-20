/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/osargs"
	"os"
	"strconv"
)

type cmdParser struct {
	cmdType    command
	message    string
	inputFile  string
	outputFile string
	parts      int
	bytes      int64
	lines      int64
}

type command int

const (
	none        command = 0
	info        command = 1
	split       command = 2
	concatenate command = 3
	wrong       command = 4
)

func newCmdParser() *cmdParser {
	cmd := new(cmdParser)
	cmd.cmdType = none
	return cmd
}

func (cmd *cmdParser) parseOSArgs() {
	osArgs := osargs.New()

	if len(osArgs.Args) == 0 {
		cmd.interpretZeroArguments()

	} else {
		results := parseFlaggedParameters(osArgs)
		restArgs := osArgs.Rest(results.toArray())
		restArgs, results.input = parsePathWOFlag(osArgs, restArgs, results.input)
		restArgs, results.output = parsePathWOFlag(osArgs, restArgs, results.output)

		if len(restArgs) == 0 {
			switch len(osArgs.Args) {
			case 1:
				cmd.interpretOneArgument(results)
			case 2:
				cmd.interpretTwoArguments(results)
			default:
				cmd.interpretManyArguments(results)
			}

		} else {
			cmd.cmdType = wrong
			cmd.message = "unknown argument \"" + restArgs[0].Value + "\""
		}
	}
}

func (cmd *cmdParser) interpretZeroArguments() {
	cmd.cmdType = info
	cmd.message = "Run 'fsplit --help' for usage."
}

func (cmd *cmdParser) interpretOneArgument(results *clResults) {
	if len(results.help) > 0 {
		cmd.cmdType = info
		cmd.message = "fsplit splits files into many, or combines them back to one\n\n"
		cmd.message = cmd.message + "USAGE\n"
		cmd.message = cmd.message + "  fsplit ( INFO | SPLIT-CONCATENATE )\n\n"
		cmd.message = cmd.message + "INFO\n"
		cmd.message = cmd.message + "  -h           print this help\n"
		cmd.message = cmd.message + "  -v           print version\n"
		cmd.message = cmd.message + "  --copyright  print copyright\n\n"
		cmd.message = cmd.message + "SPLIT-CONCATENATE\n"
		cmd.message = cmd.message + "  fsplit [COMMAND] INPUT-FILE [OUTPUT-FILE/-DIRECTORY]\n\n"
		cmd.message = cmd.message + "COMMAND\n"
		cmd.message = cmd.message + "  -p=N         split file into N chunks (parts)\n"
		cmd.message = cmd.message + "  -b=N[U]      split file into N bytes per chunk, U = unit (k/K, m/M or g/G)\n"
		cmd.message = cmd.message + "  -l=N         split file into N lines per chunk\n"
		cmd.message = cmd.message + "  -c           concatenate files (INPUT-FILE is only one file, the first one)"

	} else if len(results.version) > 0 {
		cmd.cmdType = info
		cmd.message = version.String()

	} else if len(results.copyright) > 0 {
		cmd.cmdType = info
		cmd.message = "copyright 2019 Vitali Baumtrok (vbsw@mailbox.org)\n"
		cmd.message = cmd.message + "distributed under the Boost Software License, version 1.0"

	} else {
		cmd.interpretInput(results)
		cmd.interpretOutput(results)

		if cmd.cmdType == none {
			cmd.cmdType = split
			cmd.parts = 2

		} else {
			cmd.setWrongArgumentUsage()
		}
	}
}

func (cmd *cmdParser) setWrongArgumentUsage() {
	cmd.cmdType = wrong
	cmd.message = "wrong argument usage"
}

func (cmd *cmdParser) interpretTwoArguments(results *clResults) {
	if results.infoAvailable() {
		cmd.setWrongArgumentUsage()

	} else if results.oneParamHasMultipleResults() {
		cmd.setWrongArgumentUsage()

		/* split */
	} else if len(results.concat) == 0 {
		cmd.interpretInput(results)
		cmd.interpretOutput(results)
		cmd.interpretParts(results)
		cmd.interpretBytes(results)
		cmd.interpretLines(results)

		if cmd.cmdType == none {
			cmd.cmdType = split
			if cmd.parts == 0 && cmd.bytes == 0 && cmd.lines == 0 {
				cmd.parts = 2
			}

		} else {
			cmd.setWrongArgumentUsage()
		}

		/* concatenate */
	} else {
		cmd.cmdType = info
		cmd.message = "concatenation not supported, yet"
	}
}

func (cmd *cmdParser) interpretInput(results *clResults) {
	if cmd.cmdType == none {
		if len(results.input) > 0 {
			param := results.input[0]
			if fileExists(param.Value) {
				cmd.inputFile = param.Value

			} else if direcotryExists(param.Value) {
				cmd.cmdType = wrong
				cmd.message = "input is a directory, but must be a file"

			} else {
				cmd.cmdType = wrong
				cmd.message = "input file does not exist"
			}

		} else {
			cmd.cmdType = wrong
			cmd.message = "input file is not specified"
		}
	}
}

func (cmd *cmdParser) interpretOutput(results *clResults) {
	if cmd.cmdType == none {
		if len(results.output) > 0 {
			param := results.output[0]
			if fileExists(param.Value) {
				cmd.outputFile = param.Value

			} else if direcotryExists(param.Value) {
				cmd.cmdType = wrong
				cmd.message = "output is a directory, but must be a file"

			} else {
				cmd.cmdType = wrong
				cmd.message = "output file does not exist"
			}
		} else {
			cmd.outputFile = cmd.inputFile
		}
	}
}

func (cmd *cmdParser) interpretParts(results *clResults) {
	if cmd.cmdType == none {
		if len(results.parts) > 0 {
			parts, err := strconv.Atoi(results.parts[0].Value)
			if err == nil {
				cmd.parts = int(abs(int64(parts)))
			} else {
				cmd.cmdType = wrong
				cmd.message = "can't parse number of parts"
			}
		}
	}
}

func (cmd *cmdParser) interpretBytes(results *clResults) {
	if cmd.cmdType == none {
		if len(results.bytes) > 0 {
			bytes, err := parseBytes(results.bytes[0].Value)
			if err == nil {
				cmd.bytes = abs(bytes)
			} else {
				cmd.cmdType = wrong
				cmd.message = "can't parse number of parts"
			}
		}
	}
}

func (cmd *cmdParser) interpretLines(results *clResults) {
	if cmd.cmdType == none {
		if len(results.lines) > 0 {
			lines, err := strconv.Atoi(results.lines[0].Value)
			if err == nil {
				cmd.lines = abs(int64(lines))
			} else {
				cmd.cmdType = wrong
				cmd.message = "can't parse number of lines"
			}
		}
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

func (cmd *cmdParser) interpretManyArguments(results *clResults) {
	// TODO
	cmd.cmdType = info
	cmd.message = "not supported, yet"
}

func parseFlaggedParameters(osArgs *osargs.OSArgs) *clResults {
	results := new(clResults)
	ioOp := osargs.NewAsgOp("", "=")
	cmdOp := osargs.NewAsgOp(" ", "", "=")

	results.help = osArgs.Parse("-h", "--help", "-help", "help")
	results.version = osArgs.Parse("-v", "--version", "-version", "version")
	results.copyright = osArgs.Parse("--copyright", "-copyright", "copyright")
	results.input = osArgs.ParsePairs(ioOp, "-i", "--input", "-input", "input")
	results.output = osArgs.ParsePairs(ioOp, "-o", "--output", "-output", "output")
	results.parts = osArgs.ParsePairs(cmdOp, "-p", "--parts", "-parts", "parts")
	results.bytes = osArgs.ParsePairs(cmdOp, "-b", "--bytes", "-bytes", "bytes")
	results.lines = osArgs.ParsePairs(cmdOp, "-l", "--lines", "-lines", "lines")
	results.concat = osArgs.ParsePairs(cmdOp, "-c", "--concat", "-concat", "concat")

	return results
}

func parsePathWOFlag(osArgs *osargs.OSArgs, restArgs, results []osargs.Param) ([]osargs.Param, []osargs.Param) {
	if len(results) == 0 {
		for i, restArg := range restArgs {
			if stringPathLike(restArg.Value) {
				restArgs = removeResult(restArgs, i)
				results = append(results, osargs.Param{"", restArg.Value, "", restArg.Index})
				break
			}
		}
	}
	return restArgs, results
}

func removeResult(results []osargs.Param, index int) []osargs.Param {
	copy(results[index:], results[index+1:])
	return results[:len(results)-1]
}

func fileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return (err == nil || !os.IsNotExist(err)) && !fileInfo.IsDir()
}

func direcotryExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return (err == nil || !os.IsNotExist(err)) && fileInfo.IsDir()
}

func stringPathLike(path string) bool {
	bytes := []byte(path)
	for _, b := range bytes {
		if b == '.' || b == '/' || b == '\\' {
			return true
		}
	}
	return false
}

func abs(value int64) int64 {
	if value > 0 {
		return value
	}
	return -value
}
