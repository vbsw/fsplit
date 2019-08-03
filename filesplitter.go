/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"io"
	"os"
)

type fileSplitter struct {
	inputFile   string
	outputFile  string
	bytesCopied int64
	err         error
}

func newFileSplitter(inputFile, outputFile string) *fileSplitter {
	splitter := new(fileSplitter)
	splitter.inputFile = inputFile
	splitter.outputFile = outputFile
	return splitter
}

func (splitter *fileSplitter) splitFileIntoParts(parts int64) {
	var in *os.File

	in, splitter.err = os.Open(splitter.inputFile)

	if splitter.err == nil {
		defer in.Close()
		inputSize := splitter.inputFileSize()
		if parts > inputSize {
			parts = inputSize
		}
		outputSize := splitter.outputFileSize(parts, inputSize)
		nameGenerator := newFileNameGenerator(splitter.outputFile, parts)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			splitter.copyFile(in, outputSize, nameGenerator)
		}
	}
}

func (splitter *fileSplitter) splitFileBySize(outputSize int64) {
	var in *os.File

	in, splitter.err = os.Open(splitter.inputFile)

	if splitter.err == nil {
		defer in.Close()
		inputSize := splitter.inputFileSize()
		parts := splitter.outputFileParts(outputSize, inputSize)
		nameGenerator := newFileNameGenerator(splitter.outputFile, parts)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			splitter.copyFile(in, outputSize, nameGenerator)
		}
	}
}

func (splitter *fileSplitter) copyFile(in *os.File, fileSize int64, nameGenerator *fileNameGenerator) {
	var out *os.File
	fileName := nameGenerator.nextFileName()
	out, splitter.err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)

	if splitter.err == nil {
		var written int64
		defer out.Close()

		written, splitter.err = io.CopyN(out, in, fileSize)
		splitter.bytesCopied += written

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

func (splitter *fileSplitter) outputFileSize(parts int64, inputSize int64) int64 {
	chunkSize := inputSize / parts
	// remove rounding error
	if chunkSize*parts < inputSize {
		chunkSize++
	}
	return chunkSize
}

func (splitter *fileSplitter) outputFileParts(outputSize, inputSize int64) int64 {
	parts := inputSize / outputSize
	if parts*outputSize < inputSize {
		parts++
	}
	return parts
}
