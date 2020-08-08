/*
 *       Copyright 2019, 2020, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"github.com/vbsw/cl"
	"os"
)

type arguments struct {
	help      []cl.Argument
	version   []cl.Argument
	copyright []cl.Argument
	concat    []cl.Argument
	input     []cl.Argument
	output    []cl.Argument
	parts     []cl.Argument
	bytes     []cl.Argument
	lines     []cl.Argument
}

func argumentsFromOSArgs() (*arguments, error) {
	var args *arguments
	var err error
	ops := []string{" ", "=", ""}
	cmdLine := cl.New(os.Args[1:])

	if len(cmdLine.Args) > 0 {
		args = new(arguments)
		args.help = cmdLine.Parse("-h", "--help", "-help", "help")
		args.version = cmdLine.Parse("-v", "--version", "-version", "version")
		args.copyright = cmdLine.Parse("--copyright", "-copyright", "copyright")
		args.concat = cmdLine.Parse("-c", "--concat", "-concat", "concat")
		args.input = cmdLine.ParsePairs(ops, "-i", "--input", "-input", "input")
		args.output = cmdLine.ParsePairs(ops, "-o", "--output", "-output", "output")
		args.parts = cmdLine.ParsePairs(ops, "-p", "--parts", "-parts", "parts")
		args.bytes = cmdLine.ParsePairs(ops, "-b", "--bytes", "-bytes", "bytes")
		args.lines = cmdLine.ParsePairs(ops, "-l", "--lines", "-lines", "lines")

		unparsedArgs := cmdLine.UnparsedArgsIndices()
		unparsedArgs = args.parseInput(cmdLine, unparsedArgs)
		unparsedArgs = args.parseOutput(cmdLine, unparsedArgs)

		if len(unparsedArgs) > 0 {
			unknownArg := cmdLine.Args[unparsedArgs[0]]
			err = errors.New("unknown argument \"" + unknownArg + "\"")
		}
	}
	return args, err
}

func (args *arguments) parseInput(cmdLine *cl.CommandLine, unparsedArgs []int) []int {
	if len(args.input) == 0 {
		// just accept the first unparsed argument, if input wasn't set explicitly
		if len(unparsedArgs) > 0 {
			index := unparsedArgs[0]
			value := cmdLine.Args[index]
			args.input = append(args.input, cl.Argument{"<none>", value, "", index})
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (args *arguments) parseOutput(cmdLine *cl.CommandLine, unparsedArgs []int) []int {
	if len(args.output) == 0 {
		// just accept the first unparsed argument, if output wasn't set explicitly
		if len(unparsedArgs) > 0 {
			index := unparsedArgs[0]
			value := cmdLine.Args[index]
			args.output = append(args.output, cl.Argument{"<none>", value, "", index})
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (args *arguments) incompatibleArguments() bool {
	opAvailable := len(args.concat) > 0 || len(args.input) > 0 || len(args.output) > 0 || len(args.parts) > 0 || len(args.bytes) > 0 || len(args.lines) > 0

	if len(args.help) > 0 && (len(args.version) > 0 || len(args.copyright) > 0 || opAvailable) {
		return true

	} else if len(args.version) > 0 && (len(args.help) > 0 || len(args.copyright) > 0 || opAvailable) {
		return true

	} else if len(args.copyright) > 0 && (len(args.help) > 0 || len(args.version) > 0 || opAvailable) {
		return true

	} else if len(args.parts) > 0 {
		return len(args.bytes) > 0 || len(args.lines) > 0 || len(args.concat) > 0

	} else if len(args.bytes) > 0 {
		return len(args.parts) > 0 || len(args.lines) > 0 || len(args.concat) > 0

	} else if len(args.lines) > 0 {
		return len(args.parts) > 0 || len(args.bytes) > 0 || len(args.concat) > 0

	} else if len(args.concat) > 0 {
		return len(args.parts) > 0 || len(args.bytes) > 0 || len(args.lines) > 0
	}
	return false
}

func (args *arguments) oneParamHasMultipleResults() bool {
	return len(args.help) > 1 || len(args.version) > 1 || len(args.copyright) > 1 || len(args.concat) > 1 || len(args.input) > 1 || len(args.output) > 1 || len(args.parts) > 1 || len(args.bytes) > 1 || len(args.lines) > 1
}

func (args *arguments) isInfo() bool {
	return len(args.help) > 0 || len(args.version) > 0 || len(args.copyright) > 0
}
