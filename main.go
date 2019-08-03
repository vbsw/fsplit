/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package fsplit is compiled to an executable. It splits one file into many or concatenates them back to one.
package main

import (
	"fmt"
	"github.com/vbsw/semver"
)

var version semver.Version

func main() {
	version = semver.New(0, 1, 0)
	cmd := newCmdParser()

	cmd.parseOSArgs()

	switch cmd.cmdType {
	case none:
		cmd.message = "unknown state"
		printError(cmd)
	case info:
		printInfo(cmd)
	case split:
		splitFile(cmd)
	case concatenate:
		concatenateFiles(cmd)
	default:
		printError(cmd)
	}
}

func printInfo(cmd *cmdParser) {
	fmt.Println(cmd.message)
}

func splitFile(cmd *cmdParser) {
	splitter := newFileSplitter(cmd.inputFile, cmd.outputFile)

	if cmd.parts > 0 {
		splitter.splitFileIntoParts(cmd.parts)

	} else if cmd.bytes > 0 {
		splitter.splitFileBySize(cmd.bytes)

	} else if cmd.lines > 0 {
		cmd.message = "split by lines is not supported, yet"
		printInfo(cmd)
	}
	if splitter.err != nil {
		cmd.message = splitter.err.Error()
		printError(cmd)
	}
}

func concatenateFiles(cmd *cmdParser) {
	concatenator := newFileConcatenator(cmd.inputFile, cmd.outputFile)
	concatenator.concatenateFiles()

	if concatenator.err != nil {
		cmd.message = concatenator.err.Error()
		printError(cmd)
	}
}

func printError(cmd *cmdParser) {
	fmt.Println("error:", cmd.message)
}
