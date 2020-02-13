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
)

func main() {
	cmd := commandFromCommandLine()

	switch cmd.id {
	case none:
		cmd.message = "unknown state"
		printError(cmd)
	case info:
		printInfo(cmd)
	case split:
		splitFile(cmd)
	case concat:
		concatenateFiles(cmd)
	default:
		printError(cmd)
	}
}

func printInfo(cmd *command) {
	fmt.Println(cmd.message)
}

func splitFile(cmd *command) {
	splitter := newFileSplitter(cmd.inputFile, cmd.outputFile)

	if cmd.parts > 0 {
		splitter.splitFileIntoParts(cmd.parts)

	} else if cmd.bytes > 0 {
		splitter.splitFileBySize(cmd.bytes)

	} else if cmd.lines > 0 {
		splitter.splitFileByLines(cmd.lines)
	}
	if splitter.err != nil {
		cmd.message = splitter.err.Error()
		printError(cmd)
	}
}

func concatenateFiles(cmd *command) {
	concatenator := newFileConcatenator(cmd.inputFile, cmd.outputFile)
	concatenator.concatenateFiles()

	if concatenator.err != nil {
		cmd.message = concatenator.err.Error()
		printError(cmd)
	}
}

func printError(cmd *command) {
	fmt.Println("error:", cmd.message)
}
