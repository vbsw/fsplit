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
		nameGenerator := newFileNameGeneratorForSplit(splitter.outputFile, parts)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			outputFileName := nameGenerator.nextFileName()
			splitter.copyFile(in, outputFileName, outputSize)
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
		nameGenerator := newFileNameGeneratorForSplit(splitter.outputFile, parts)

		for splitter.bytesCopied < inputSize && splitter.err == nil {
			outputFileName := nameGenerator.nextFileName()
			splitter.copyFile(in, outputFileName, outputSize)
		}
	}
}

func (splitter *fileSplitter) splitFileByLines(lines int64) {
	inputSizes := splitter.inputSizesByLines(lines)

	if len(inputSizes) > 0 {
		var in *os.File
		in, splitter.err = os.Open(splitter.inputFile)

		if splitter.err == nil {
			defer in.Close()
			nameGenerator := newFileNameGeneratorForSplit(splitter.outputFile, int64(len(inputSizes)))

			for _, outputSize := range inputSizes {
				outputFileName := nameGenerator.nextFileName()
				splitter.copyFile(in, outputFileName, outputSize)
			}
		}
	}
}

func (splitter *fileSplitter) copyFile(in *os.File, destFileName string, copySize int64) {
	var out *os.File
	out, splitter.err = os.OpenFile(destFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)

	if splitter.err == nil {
		var written int64
		defer out.Close()

		written, splitter.err = io.CopyN(out, in, copySize)
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

func (splitter *fileSplitter) inputSizesByLines(lines int64) []int64 {
	var in *os.File
	var sizes []int64

	in, splitter.err = os.Open(splitter.inputFile)

	if splitter.err == nil {
		var linesRead int64
		var size int64
		var prevByte byte
		defer in.Close()
		sizes = make([]int64, 0, 1024)
		buffer := make([]byte, 1024*1024*8)
		bytesRead, _ := in.Read(buffer)

		for bytesRead > 0 {
			for i := 0; i < bytesRead; i++ {
				currByte := buffer[i]
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
	return sizes
}
