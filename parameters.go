/*
 *       Copyright 2019, 2020, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"github.com/vbsw/osargs"
	"os"
)

type parameters struct {
	help      []osargs.Parameter
	version   []osargs.Parameter
	copyright []osargs.Parameter
	concat    []osargs.Parameter
	input     []osargs.Parameter
	output    []osargs.Parameter
	parts     []osargs.Parameter
	bytes     []osargs.Parameter
	lines     []osargs.Parameter
}

func parametersFromOSArgs() (*parameters, error) {
	var params *parameters
	var err error
	args := osargs.New(" ", "=", "")

	if len(args.Str) > 0 {
		params = new(parameters)
		params.help = args.Parse("-h", "--help", "-help", "help")
		params.version = args.Parse("-v", "--version", "-version", "version")
		params.copyright = args.Parse("--copyright", "-copyright", "copyright")
		params.concat = args.Parse("-c", "--concat", "-concat", "concat")
		params.input = args.ParsePairs("-i", "--input", "-input", "input")
		params.output = args.ParsePairs("-o", "--output", "-output", "output")
		params.parts = args.ParsePairs("-p", "--parts", "-parts", "parts")
		params.bytes = args.ParsePairs("-b", "--bytes", "-bytes", "bytes")
		params.lines = args.ParsePairs("-l", "--lines", "-lines", "lines")

		unparsedArgs := args.Rest(params.help, params.version, params.copyright, params.concat, params.input, params.output, params.parts, params.bytes, params.lines)
		unparsedArgs = params.parseInput(args, unparsedArgs)
		unparsedArgs = params.parseOutput(args, unparsedArgs)

		if len(unparsedArgs) > 0 {
			unknownArg := args.Str[unparsedArgs[0]]
			err = errors.New("unknown argument \"" + unknownArg + "\"")
		}
	}
	return params, err
}

func (params *parameters) parseInput(args *osargs.Arguments, unparsedArgs []int) []int {
	if len(params.input) == 0 {
		// just accept the first unparsed argument, if input wasn't set explicitly
		if len(unparsedArgs) > 0 {
			index := unparsedArgs[0]
			value := args.Str[index]
			params.input = append(params.input, osargs.Parameter{"<none>", value, "", index})
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *parameters) parseOutput(args *osargs.Arguments, unparsedArgs []int) []int {
	if len(params.output) == 0 {
		// just accept the first unparsed argument, if output wasn't set explicitly
		if len(unparsedArgs) > 0 {
			index := unparsedArgs[0]
			value := args.Str[index]
			params.output = append(params.output, osargs.Parameter{"<none>", value, "", index})
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *parameters) incompatibleArguments() bool {
	opAvailable := len(params.concat) > 0 || len(params.input) > 0 || len(params.output) > 0 || len(params.parts) > 0 || len(params.bytes) > 0 || len(params.lines) > 0

	if len(params.help) > 0 && (len(params.version) > 0 || len(params.copyright) > 0 || opAvailable) {
		return true

	} else if len(params.version) > 0 && (len(params.help) > 0 || len(params.copyright) > 0 || opAvailable) {
		return true

	} else if len(params.copyright) > 0 && (len(params.help) > 0 || len(params.version) > 0 || opAvailable) {
		return true

	} else if len(params.parts) > 0 {
		return len(params.bytes) > 0 || len(params.lines) > 0 || len(params.concat) > 0

	} else if len(params.bytes) > 0 {
		return len(params.parts) > 0 || len(params.lines) > 0 || len(params.concat) > 0

	} else if len(params.lines) > 0 {
		return len(params.parts) > 0 || len(params.bytes) > 0 || len(params.concat) > 0

	} else if len(params.concat) > 0 {
		return len(params.parts) > 0 || len(params.bytes) > 0 || len(params.lines) > 0
	}
	return false
}

func (params *parameters) oneParamHasMultipleResults() bool {
	return len(params.help) > 1 || len(params.version) > 1 || len(params.copyright) > 1 || len(params.concat) > 1 || len(params.input) > 1 || len(params.output) > 1 || len(params.parts) > 1 || len(params.bytes) > 1 || len(params.lines) > 1
}

func (params *parameters) isInfo() bool {
	return len(params.help) > 0 || len(params.version) > 0 || len(params.copyright) > 0
}

func stringPathLike(path string) bool {
	bytes := []byte(path)
	for _, b := range bytes {
		if b == '.' || b == '/' || b == '\\' {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return (err == nil || !os.IsNotExist(err)) && !fileInfo.IsDir()
}

func direcotryExists(path string) bool {
	fileInfo, err := os.Stat(path)
	return (err == nil || !os.IsNotExist(err)) && fileInfo.IsDir()
}
