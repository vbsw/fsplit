/*
 *      Copyright 2019 - 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package main is compiled to an executable. It splits one file into many or concatenates them back to one.
package main

import (
	"errors"
	"fmt"
	"github.com/vbsw/golib/osargs"
	"io"
	"os"
	"strconv"
)

type tParameters struct {
	help       *osargs.Result
	version    *osargs.Result
	example    *osargs.Result
	copyright  *osargs.Result
	concat     *osargs.Result
	input      *osargs.Result
	output     *osargs.Result
	parts      *osargs.Result
	bytes      *osargs.Result
	lines      *osargs.Result
	infoParams []*osargs.Result
	cmdParams  []*osargs.Result
}

type tPathGenerator struct {
	path         string
	suffixFormat string
	counter      int64
}

func main() {
	var params tParameters
	err := params.initFromOSArgs()
	if err == nil {
		if params.infoAvailable() {
			printInfo(&params)
		} else if params.concat.Available() {
			err = concatenateFiles(&params)
		} else {
			err = splitFile(&params)
		}
	}
	if err != nil {
		printError(err)
	}
}

func (params *tParameters) initFromOSArgs() error {
	args := osargs.New()
	err := params.initFromArgs(args)
	return err
}

// initFromArgs is for test purposes.
func (params *tParameters) initFromArgs(args *osargs.Arguments) error {
	var err error
	if len(args.Values) > 0 {
		delimiter := osargs.NewDelimiter(true, true, "=")
		params.help = args.Parse("-h", "--help", "-help", "help")
		params.version = args.Parse("-v", "--version", "-version", "version")
		params.example = args.Parse("-e", "--example", "-example", "example")
		params.copyright = args.Parse("--copyright", "-copyright", "copyright")
		params.concat = args.Parse("-c", "--concat", "-concat", "concat")
		params.input = args.ParsePairs(delimiter, "-i", "--input", "-input", "input")
		params.output = args.ParsePairs(delimiter, "-o", "--output", "-output", "output")
		params.parts = args.ParsePairs(delimiter, "-p", "--parts", "-parts", "parts")
		params.bytes = args.ParsePairs(delimiter, "-b", "--bytes", "-bytes", "bytes")
		params.lines = args.ParsePairs(delimiter, "-l", "--lines", "-lines", "lines")
		params.poolInfoParams()
		params.poolCmdParams()

		unparsedArgs := args.UnparsedArgs()
		unparsedArgs = params.parseInput(unparsedArgs)
		unparsedArgs = params.parseOutput(unparsedArgs)

		err = params.validateParameters(unparsedArgs)
		if err == nil && !params.infoAvailable() {
			params.ensurePartsCount()
		}
	}
	return err
}

func (params *tParameters) parseInput(unparsedArgs []string) []string {
	if !params.input.Available() {
		// just accept the first unparsed argument, if input wasn't set explicitly
		if len(unparsedArgs) > 0 {
			params.input.Values = append(params.input.Values, unparsedArgs[0])
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *tParameters) parseOutput(unparsedArgs []string) []string {
	if !params.output.Available() {
		// just accept the first unparsed argument, if output wasn't set explicitly
		if len(unparsedArgs) > 0 {
			params.output.Values = append(params.output.Values, unparsedArgs[0])
			unparsedArgs = unparsedArgs[1:]
		} else if params.input.Available() {
			params.output.Values = append(params.output.Values, params.input.Values[0])
		}
	}
	return unparsedArgs
}

func (params *tParameters) validateParameters(unparsedArgs []string) error {
	var err error
	if len(unparsedArgs) > 0 {
		unknownArg := unparsedArgs[0]
		err = errors.New("unknown argument \"" + unknownArg + "\"")
	} else {
		if params.isCompatible() {
			if !params.infoAvailable() {
				err = params.validateIODirectories()
			}
		} else {
			err = errors.New("wrong argument usage")
		}
	}
	return err
}

func (params *tParameters) poolInfoParams() {
	params.infoParams = make([]*osargs.Result, 4)
	params.infoParams[0] = params.help
	params.infoParams[1] = params.version
	params.infoParams[2] = params.example
	params.infoParams[3] = params.copyright
}

func (params *tParameters) poolCmdParams() {
	params.cmdParams = make([]*osargs.Result, 6)
	params.cmdParams[0] = params.concat
	params.cmdParams[1] = params.input
	params.cmdParams[2] = params.output
	params.cmdParams[3] = params.parts
	params.cmdParams[4] = params.bytes
	params.cmdParams[5] = params.lines
}

func (params *tParameters) infoAvailable() bool {
	return anyAvailable(params.infoParams)
}

func (params *tParameters) isCompatible() bool {
	// same parameter must not be multiple
	if isMultiple(params.infoParams) || isMultiple(params.cmdParams) {
		return false
	}
	// either info or command
	if anyAvailable(params.infoParams) && anyAvailable(params.cmdParams) {
		return false
	}
	// no mixed info parameters
	if isMixed(params.infoParams...) {
		return false
	}
	// no mixed command parameters
	if isMixed(params.concat, params.parts, params.bytes, params.lines) {
		return false
	}
	return true
}

func (params *tParameters) validateIODirectories() error {
	var err error
	if !params.input.Available() {
		err = errors.New("input file is not specified")
	} else if !params.concat.Available() {
		info, errInfo := os.Stat(params.input.Values[0])
		if errInfo == nil || !os.IsNotExist(errInfo) {
			if info != nil {
				if info.IsDir() {
					err = errors.New("input file is a directory, but must be a file")
				}
			} else {
				err = errors.New("wrong input path syntax")
			}
		} else {
			err = errors.New("input file does not exist")
		}
	}
	return err
}

func (params *tParameters) ensurePartsCount() {
	if !params.concat.Available() && !params.parts.Available() && !params.bytes.Available() && !params.lines.Available() {
		params.parts.Values = append(params.parts.Values, "2")
	}
}

func (params *tParameters) outputPathConcat(pathCalculated string) string {
	// no output path has been set explicitly
	if params.input.Values[0] == params.output.Values[0] {
		return pathCalculated
	}
	return params.output.Values[0]
}

func newPathGeneratorForConcat(path string) *tPathGenerator {
	pathGenerator := new(tPathGenerator)
	pathWOSuffix, lengthSuffix := inputPathConcat(path)
	pathGenerator.path = pathWOSuffix
	pathGenerator.suffixFormat = suffixFormat(lengthSuffix)
	return pathGenerator
}

func newPathGeneratorForSplit(path string, parts int) *tPathGenerator {
	pathGenerator := new(tPathGenerator)
	pathGenerator.path = path
	pathGenerator.suffixFormat = suffixFormat(len(strconv.Itoa(parts)))
	return pathGenerator
}

func (pathGenerator *tPathGenerator) nextPath() string {
	pathGenerator.counter++
	suffix := fmt.Sprintf(pathGenerator.suffixFormat, pathGenerator.counter)
	pathPart := pathGenerator.path + "." + suffix
	return pathPart
}

func concatenateFiles(params *tParameters) error {
	pathGenerator := newPathGeneratorForConcat(params.input.Values[0])
	inputPath := pathGenerator.nextPath()
	// check input first
	in, err := os.Open(inputPath)
	if err == nil {
		var out *os.File
		outputPath := params.outputPathConcat(pathGenerator.path)
		out, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err == nil {
			defer out.Close()
			in, err = copyFileConcat(out, in, pathGenerator)
			for err == nil {
				in, err = copyFileConcat(out, in, pathGenerator)
			}
			if os.IsNotExist(err) {
				err = nil
			}
		} else {
			in.Close()
		}
	}
	return err
}

func splitFile(params *tParameters) error {
	var err error
	if params.parts.Available() {
		err = splitFileByParts(params)
	} else if params.bytes.Available() {
		err = splitFileBySize(params)
	} else if params.lines.Available() {
		err = splitFileByLines(params)
	}
	return err
}

func splitFileByParts(params *tParameters) error {
	parts, err := interpretParts(params)
	if err == nil {
		var in *os.File
		in, err = os.Open(params.input.Values[0])
		if err == nil {
			defer in.Close()
			var copied int64
			inputSize := fileSize(params.input.Values[0])
			if parts > inputSize {
				parts = inputSize
			}
			outputSize := partSize(parts, inputSize)
			pathGenerator := newPathGeneratorForSplit(params.output.Values[0], int(parts))
			for copied < inputSize && err == nil {
				outputPath := pathGenerator.nextPath()
				copied, err = copyFileSplit(outputPath, in, outputSize, copied)
			}
		}
	}
	return err
}

func splitFileBySize(params *tParameters) error {
	outputSize, err := interpretSize(params)
	if err == nil {
		var in *os.File
		in, err = os.Open(params.input.Values[0])
		if err == nil {
			defer in.Close()
			var parts, copied int64
			inputSize := fileSize(params.input.Values[0])
			if outputSize > 0 {
				parts = inputSize / outputSize
				if parts*outputSize < inputSize {
					parts++
				}
			} else {
				parts = 2
				outputSize = partSize(parts, inputSize)
			}
			pathGenerator := newPathGeneratorForSplit(params.output.Values[0], int(parts))
			for copied < inputSize && err == nil {
				outputPath := pathGenerator.nextPath()
				copied, err = copyFileSplit(outputPath, in, outputSize, copied)
			}
		}
	}
	return err
}

func splitFileByLines(params *tParameters) error {
	lines, err := interpretLines(params)
	if err == nil {
		var inputSizes []int64
		inputSizes, err = inputSizesByLines(params.input.Values[0], lines)
		if err == nil && len(inputSizes) > 0 {
			var in *os.File
			in, err = os.Open(params.input.Values[0])
			if err == nil {
				defer in.Close()
				var copied int64
				pathGenerator := newPathGeneratorForSplit(params.output.Values[0], len(inputSizes))
				for _, outputSize := range inputSizes {
					outputPath := pathGenerator.nextPath()
					copied, err = copyFileSplit(outputPath, in, outputSize, copied)
				}
			}
		}
	}
	return err
}

func anyAvailable(results []*osargs.Result) bool {
	for _, result := range results {
		if result.Available() {
			return true
		}
	}
	return false
}

func isMultiple(paramsMult []*osargs.Result) bool {
	for _, param := range paramsMult {
		if param.Count() > 1 {
			return true
		}
	}
	return false
}

func isMixed(params ...*osargs.Result) bool {
	for i, paramA := range params {
		if paramA.Available() {
			for _, paramB := range params[i+1:] {
				if paramB.Available() {
					return true
				}
			}
			break
		}
	}
	return false
}

func suffixFormat(numDigits int) string {
	numDigitsStr := strconv.Itoa(numDigits)
	format := "%0" + numDigitsStr + "d"
	return format
}

func interpretParts(params *tParameters) (int64, error) {
	parts, err := strconv.Atoi(params.parts.Values[0])
	if err == nil {
		if parts > 1 {
			return int64(parts), nil
		}
		return 2, nil
	}
	return 0, errors.New("can't parse number of parts")
}

func interpretSize(params *tParameters) (int64, error) {
	bytes, err := parseBytes(params.bytes.Values[0])
	if err == nil {
		return bytes, nil
	}
	return bytes, errors.New("can't parse number of bytes")
}

func interpretLines(params *tParameters) (int64, error) {
	lines, err := strconv.Atoi(params.lines.Values[0])
	if err == nil {
		return int64(lines), nil
	}
	return int64(lines), errors.New("can't parse number of lines")
}

func fileSize(path string) int64 {
	fileInfo, _ := os.Stat(path)
	return fileInfo.Size()
}

func partSize(parts, inputSize int64) int64 {
	chunkSize := inputSize / parts
	// remove rounding error
	if chunkSize*parts < inputSize {
		chunkSize++
	}
	return chunkSize
}

func copyFileConcat(out, in *os.File, pathGenerator *tPathGenerator) (*os.File, error) {
	_, err := io.Copy(out, in)
	in.Close()
	if err == nil || err == io.EOF {
		inputPath := pathGenerator.nextPath()
		in, err = os.Open(inputPath)
	}
	return in, err
}

func copyFileSplit(path string, in *os.File, size, copiedTotal int64) (int64, error) {
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer out.Close()
		var copied int64
		copied, err = io.CopyN(out, in, size)
		copiedTotal += copied
		if err == io.EOF {
			err = nil
		}
	}
	return copiedTotal, err
}

func inputSizesByLines(path string, lines int64) ([]int64, error) {
	var sizes []int64
	in, err := os.Open(path)
	if err == nil {
		defer in.Close()
		var linesRead, size int64
		var bytesRead int
		var prevByte byte
		sizes = make([]int64, 0, 1024)
		buffer := make([]byte, 1024*1024*8)
		bytesRead, err = in.Read(buffer)
		if err == nil {
			for bytesRead > 0 {
				for _, currByte := range buffer[:bytesRead] {
					size++
					if currByte == '\n' || prevByte == '\r' {
						linesRead++
					}
					if linesRead == lines {
						if currByte == '\n' {
							sizes = append(sizes, size)
							size = 0
						} else {
							sizes = append(sizes, size-1)
							size = 1
						}
						linesRead = 0
					}
					prevByte = currByte
				}
				bytesRead, _ = in.Read(buffer)
			}
			if size > 0 {
				sizes = append(sizes, size)
			}
		}
	}
	return sizes, err
}

func parseBytes(bytesStr string) (int64, error) {
	if len(bytesStr) > 0 {
		var bytes64 int64
		bytesArray := []byte(bytesStr)
		lastByte := bytesArray[len(bytesArray)-1]
		if lastByte == 'k' || lastByte == 'K' || lastByte == 'm' || lastByte == 'M' || lastByte == 'g' || lastByte == 'G' {
			bytesStr = bytesStr[:len(bytesStr)-1]
		} else {
			lastByte = 0
		}
		bytes, err := strconv.Atoi(bytesStr)
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
		return bytes64, err
	}
	return 0, nil
}

func inputPathConcat(path string) (string, int) {
	pathWOSuffix := path
	lengthSuffix := 0
	indexDot := rIndex(path, '.')
	if indexDot >= 0 {
		suffix := path[indexDot+1:]
		part, err := strconv.Atoi(suffix)
		if err == nil && part >= 0 {
			pathWOSuffix = path[:indexDot]
			lengthSuffix = len(suffix)
		}
	}
	/* guess existence of other suitable file */
	if lengthSuffix == 0 {
		pathWDot := path + "."
		for i := 1; i < 11; i++ {
			suffix := fmt.Sprintf(suffixFormat(i), 1)
			if isFile(pathWDot + suffix) {
				lengthSuffix = i
				break
			}
		}
	}
	return pathWOSuffix, lengthSuffix
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return info != nil && !info.IsDir() && (err == nil || !os.IsNotExist(err))
}

func rIndex(str string, b byte) int {
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == b {
			return i
		}
	}
	return -1
}

func printInfo(params *tParameters) {
	if params.help == nil {
		printShortInfo()
	} else if params.help.Available() {
		printHelp()
	} else if params.version.Available() {
		printVersion()
	} else if params.example.Available() {
		printExample()
	} else if params.copyright.Available() {
		printCopyright()
	} else {
		printShortInfo()
	}
}

func printShortInfo() {
	fmt.Println("Run 'fsplit --help' for usage.")
}

func printHelp() {
	message := "\nUSAGE\n"
	message += "  fsplit ( INFO | SPLIT/CONCATENATE )\n\n"
	message += "INFO\n"
	message += "  -h, --help    print this help\n"
	message += "  -v, --version print version\n"
	message += "  --copyright   print copyright\n\n"
	message += "SPLIT/CONCATENATE\n"
	message += "  [COMMAND] INPUT-FILE [OUTPUT-FILE]\n\n"
	message += "COMMAND\n"
	message += "  -p=N          split file into N parts (chunks)\n"
	message += "  -b=N[U]       split file into N bytes per chunk, U = unit (k/K, m/M or g/G)\n"
	message += "  -l=N          split file into N lines per chunk\n"
	message += "  -c            concatenate files (INPUT-FILE is only one file, the first one)"
	fmt.Println(message)
}

func printVersion() {
	fmt.Println("1.0.3")
}

func printExample() {
	message := "\nEXAMPLES\n"
	message += "   ... not available"
	fmt.Println(message)
}

func printCopyright() {
	message := "Copyright 2019 - 2022, Vitali Baumtrok (vbsw@mailbox.org).\n"
	message += "Distributed under the Boost Software License, Version 1.0."
	fmt.Println(message)
}

func printError(err error) {
	fmt.Println("error:", err.Error())
}
