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
	"github.com/vbsw/osargs"
	"github.com/vbsw/semver"
)

type command int

type commandData struct {
	cmdType command
	message string
}

var version semver.Version

const (
	info        command = 0
	split       command = 1
	concatenate command = 2
	none        command = 3
)

func main() {
	version = semver.New(0,0,1)
	cmd := parseCommandLine()

	switch cmd.cmdType {

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

func parseCommandLine() *commandData {
	osArgs := osargs.New()
	cmd := new(commandData)

	switch len(osArgs.Args) {

	case 0:
		parseZeroArguments(osArgs, cmd)

	case 1:
		parseOneArgument(osArgs, cmd)

	case 2:
		parseTwoArguments(osArgs, cmd)

	default:
		parseManyArguments(osArgs, cmd)
	}
	return cmd
}

func parseZeroArguments(osArgs *osargs.OSArgs, cmd *commandData) {
	cmd.cmdType = info
	cmd.message = "Run \"fsplit --help\" for usage."
}

func parseOneArgument(osArgs *osargs.OSArgs, cmd *commandData) {
	if osArgs.Parse("-h", "--help", "-help", "help") {
		cmd.cmdType = info
		cmd.message = "fsplit splits files into many or combines them back to one\n\n"
		cmd.message = cmd.message + "USAGE\n"
		cmd.message = cmd.message + "  fsplit [OPTIONS] [INPUT-FILE]\n\n"
		cmd.message = cmd.message + "OPTIONS\n"
		cmd.message = cmd.message + "  -i=FILE    input file\n"
		cmd.message = cmd.message + "  -o=FILE    output directory (for split) or file (for concatenate)\n"
		cmd.message = cmd.message + "  -p=N       number of chunks\n"
		cmd.message = cmd.message + "  -b=N[U]    size per chunk in bytes, U = unit (k/K, m/M or g/G)\n"
		cmd.message = cmd.message + "  -l=N       number of lines per chunk\n"
		cmd.message = cmd.message + "  -c         concatenate files (INPUT-FILE is only one file, the first one)"

	} else if osArgs.Parse("-v", "--version", "-version", "version") {
		cmd.cmdType = info
		cmd.message = version.String()

	} else if osArgs.Parse("--copyright", "-copyright", "copyright") {
		cmd.cmdType = info
		cmd.message = "copyright 2019 Vitali Baumtrok (vbsw@mailbox.org)\n"
		cmd.message = cmd.message + "distributed under the Boost Software License, version 1.0"

	} else {
		cmd.cmdType = none
		cmd.message = "unknown argument \"" + osArgs.Args[0] + "\""
	}
}

func parseTwoArguments(osArgs *osargs.OSArgs, cmd *commandData) {
	cmd.message = "too many arguments"
}

func parseManyArguments(osArgs *osargs.OSArgs, cmd *commandData) {
	cmd.message = "too many arguments"
}

func printInfo(cmd *commandData) {
	fmt.Println(cmd.message)
}

func splitFile(cmd *commandData) {
	fmt.Println("not implemented")
}

func concatenateFiles(cmd *commandData) {
	fmt.Println("not implemented")
}

func printError(cmd *commandData) {
	fmt.Println("Error:", cmd.message)
}
