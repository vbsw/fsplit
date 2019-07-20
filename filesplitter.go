/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type fileSplitter struct {
	inputFile   string
	outputFile  string
	bytesCopied int64
	partCounter int
	err         error
}

func newFileSplitter(inputFile, outputFile string) *fileSplitter {
	splitter := new(fileSplitter)
	splitter.inputFile = inputFile
	splitter.outputFile = outputFile
	return splitter
}

func (splitter *fileSplitter) splitFileIntoParts(parts int) {
	var in *os.File
	splitter.bytesCopied = 0
	splitter.partCounter = 1

	in, splitter.err = os.Open(splitter.inputFile)

	if splitter.err == nil {
		defer in.Close()
		inputSize := splitter.inputFileSize()
		outputSize := splitter.outputFileSize(parts, inputSize)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			splitter.copyFile(in, parts, outputSize)
		}
	}
}

func (splitter *fileSplitter) splitFileBySize(outputSize int64) {
	var in *os.File
	splitter.bytesCopied = 0
	splitter.partCounter = 1

	in, splitter.err = os.Open(splitter.inputFile)

	if splitter.err == nil {
		defer in.Close()
		inputSize := splitter.inputFileSize()
		parts := splitter.outputParts(outputSize, inputSize)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			splitter.copyFile(in, parts, outputSize)
		}
	}
}

func (splitter *fileSplitter) copyFile(in *os.File, parts int, fileSize int64) {
	var out *os.File
	fileName := outputFileName(splitter.inputFile, parts, splitter.partCounter)
	out, splitter.err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)

	if splitter.err == nil {
		var written int64
		defer out.Close()

		written, splitter.err = io.CopyN(out, in, fileSize)
		splitter.bytesCopied += written
		splitter.partCounter++

		if splitter.err == io.EOF {
			splitter.err = nil
		}
	}
}

func (splitter *fileSplitter) inputFileSize() int64 {
	fileInfo, _ := os.Stat(splitter.inputFile)
	inputSize := fileInfo.Size()
	return inputSize
}

func (splitter *fileSplitter) outputFileSize(parts int, inputSize int64) int64 {
	chunkSize := inputSize / int64(parts)
	// remove rounding error
	if chunkSize*int64(parts) < inputSize {
		chunkSize++
	}
	return chunkSize
}

func (splitter *fileSplitter) outputParts(outputSize, inputSize int64) int {
	parts := inputSize / outputSize
	if parts*outputSize < inputSize {
		parts++
	}
	return int(parts)
}

func outputFileName(fileNameOriginal string, parts, partNumber int) string {
	partsStr := strconv.Itoa(parts)
	digits := len(partsStr)
	digitsStr := strconv.Itoa(digits)
	format := "%0" + digitsStr + "d"
	partStr := fmt.Sprintf(format, partNumber)
	fileNamePart := fileNameOriginal + "." + partStr
	return fileNamePart
}
