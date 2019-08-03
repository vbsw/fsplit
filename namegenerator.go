/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"fmt"
	"strconv"
)

type fileNameGenerator struct {
	fileName     string
	suffixFormat string
	counter      int64
}

func newFileNameGenerator(fileName string, parts int64) *fileNameGenerator {
	nameGenerator := new(fileNameGenerator)
	partsStr := strconv.Itoa(int(parts))
	digits := len(partsStr)
	digitsStr := strconv.Itoa(digits)
	nameGenerator.fileName = fileName
	nameGenerator.suffixFormat = "%0" + digitsStr + "d"
	return nameGenerator
}

func newFileNameGenerator2(originalFileName string) *fileNameGenerator {
	var nameGenerator *fileNameGenerator
	dotIndex := rDotIndex(originalFileName)
	fileName := originalFileName
	startCounter := int64(0)
	parts := int64(0)

	if dotIndex >= 0 {
		suffix := originalFileName[dotIndex+1:]
		part, err := strconv.Atoi(suffix)
		if err == nil && part >= 0 {
			fileName = originalFileName[:dotIndex]
			startCounter = int64(part - 1)
			for i := 0; i < len(suffix); i++ {
				parts = parts*10 + 9
			}
		}
	}
	if parts == 0 {
		/* TODO: analyze available files */
		parts = 9
	}
	nameGenerator = newFileNameGenerator(fileName, parts)
	nameGenerator.counter = startCounter

	return nameGenerator
}

func (nameGenerator *fileNameGenerator) nextFileName() string {
	nameGenerator.counter++
	partStr := fmt.Sprintf(nameGenerator.suffixFormat, nameGenerator.counter)
	fileNamePart := nameGenerator.fileName + "." + partStr
	return fileNamePart
}
