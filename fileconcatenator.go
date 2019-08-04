/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"io"
	"os"
)

type fileConcatenator struct {
	inputFile   string
	outputFile  string
	bytesCopied int64
	err         error
}

func newFileConcatenator(inputFile, outputFile string) *fileConcatenator {
	concatenator := new(fileConcatenator)
	concatenator.inputFile = inputFile
	concatenator.outputFile = outputFile
	return concatenator
}

func (concatenator *fileConcatenator) concatenateFiles() {
	nameGenerator := newFileNameGeneratorForConcat(concatenator.inputFile)
	inputFile := nameGenerator.nextFileName()
	outputFile := concatenator.finalOutputFile(nameGenerator.fileName)

	if fileExists(inputFile) {
		in, inErr := os.Open(inputFile)

		if inErr == nil {
			out, outErr := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)

			if outErr == nil {
				for inErr == nil && concatenator.err == nil {
					written, err := io.Copy(out, in)
					concatenator.bytesCopied += written
					in.Close()

					if err == nil || err == io.EOF {
						inputFile = nameGenerator.nextFileName()
						in, inErr = os.Open(inputFile)

					} else {
						concatenator.err = err
					}
				}
				out.Close()

			} else {
				concatenator.err = outErr
				in.Close()
			}

		} else {
			concatenator.err = inErr
		}

	} else {
		concatenator.err = errors.New("input file not found")
	}
}

func (concatenator *fileConcatenator) finalOutputFile(finalInputFile string) string {
	if concatenator.inputFile == concatenator.outputFile {
		return finalInputFile
	}
	return concatenator.outputFile
}

func rDotIndex(str string) int {
	bytes := []byte(str)
	index := len(bytes) - 1
	for index >= 0 && bytes[index] != '.' {
		index--
	}
	return index
}
