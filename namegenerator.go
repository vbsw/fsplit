/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"github.com/vbsw/checkfile"
	"strconv"
)

type fileNameGenerator struct {
	fileName     string
	suffixFormat string
	counter      int64
}

func newFileNameGeneratorForSplit(fileName string, parts int64) *fileNameGenerator {
	nameGenerator := new(fileNameGenerator)
	nameGenerator.fileName = fileName
	nameGenerator.suffixFormat = suffixFormat(len(strconv.Itoa(int(parts))))
	return nameGenerator
}

func newFileNameGeneratorForConcat(fileName string) *fileNameGenerator {
	nameGenerator := new(fileNameGenerator)
	fileNameWOSuffix, suffixLength := analizeFileName(fileName)
	nameGenerator.fileName = fileNameWOSuffix
	nameGenerator.suffixFormat = suffixFormat(suffixLength)
	return nameGenerator
}

func (nameGenerator *fileNameGenerator) nextFileName() string {
	nameGenerator.counter++
	partStr := fmt.Sprintf(nameGenerator.suffixFormat, nameGenerator.counter)
	fileNamePart := nameGenerator.fileName + "." + partStr
	return fileNamePart
}

func suffixFormat(numDigits int) string {
	numDigitsStr := strconv.Itoa(numDigits)
	format := "%0" + numDigitsStr + "d"
	return format
}

func analizeFileName(fileName string) (string, int) {
	fileNameWOSuffix := fileName
	suffixLength := 0
	dotIndex := rDotIndex(fileName)

	if dotIndex >= 0 {
		suffix := fileName[dotIndex+1:]
		part, err := strconv.Atoi(suffix)

		if err == nil && part >= 0 {
			fileNameWOSuffix = fileName[:dotIndex]
			suffixLength = len(suffix)
		}
	}
	/* guess existence of other suitable file */
	if suffixLength == 0 {
		fileNameWDot := fileName + "."
		for i := 1; i < 11; i++ {
			suffix := fmt.Sprintf(suffixFormat(i), 1)
			if checkfile.IsFile(fileNameWDot + suffix) {
				suffixLength = i
				break
			}
		}
	}
	return fileNameWOSuffix, suffixLength
}
