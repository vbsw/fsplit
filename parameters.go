/*
 *      Copyright 2019 - 2021, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"github.com/vbsw/osargs"
)

type parameters struct {
	help      *osargs.Result
	version   *osargs.Result
	copyright *osargs.Result
	concat    *osargs.Result
	input     *osargs.Result
	output    *osargs.Result
	parts     *osargs.Result
	bytes     *osargs.Result
	lines     *osargs.Result
}

func parametersFromOSArgs() (*parameters, error) {
	var params *parameters
	var err error
	args := osargs.New()
	delimiter := osargs.NewDelimiter(true, true, "=")

	if len(args.Values) > 0 {
		params = new(parameters)
		params.help = args.Parse("-h", "--help", "-help", "help")
		params.version = args.Parse("-v", "--version", "-version", "version")
		params.copyright = args.Parse("--copyright", "-copyright", "copyright")
		params.concat = args.Parse("-c", "--concat", "-concat", "concat")
		params.input = args.ParsePairs(delimiter, "-i", "--input", "-input", "input")
		params.output = args.ParsePairs(delimiter, "-o", "--output", "-output", "output")
		params.parts = args.ParsePairs(delimiter, "-p", "--parts", "-parts", "parts")
		params.bytes = args.ParsePairs(delimiter, "-b", "--bytes", "-bytes", "bytes")
		params.lines = args.ParsePairs(delimiter, "-l", "--lines", "-lines", "lines")

		unparsedArgs := args.UnparsedArgs()
		unparsedArgs = params.parseInput(unparsedArgs)
		unparsedArgs = params.parseOutput(unparsedArgs)

		if len(unparsedArgs) > 0 {
			unknownArg := unparsedArgs[0]
			err = errors.New("unknown argument \"" + unknownArg + "\"")
		}
	}
	return params, err
}

func (params *parameters) parseInput(unparsedArgs []string) []string {
	if !params.input.Available() {
		// just accept the first unparsed argument, if input wasn't set explicitly
		if len(unparsedArgs) > 0 {
			params.input.Values = append(params.input.Values, unparsedArgs[0])
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *parameters) parseOutput(unparsedArgs []string) []string {
	if params.output.Available() {
		// just accept the first unparsed argument, if output wasn't set explicitly
		if len(unparsedArgs) > 0 {
			params.output.Values = append(params.output.Values, unparsedArgs[0])
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *parameters) incompatibleParameters() bool {
	opAvailable := params.concat.Available() || params.input.Available() || params.output.Available() || params.parts.Available() || params.bytes.Available() || params.lines.Available()

	if params.help.Available() && (params.version.Available() || params.copyright.Available() || opAvailable) {
		return true

	} else if params.version.Available() && (params.help.Available() || params.copyright.Available() || opAvailable) {
		return true

	} else if params.copyright.Available() && (params.help.Available() || params.version.Available() || opAvailable) {
		return true

	} else if params.parts.Available() {
		return params.bytes.Available() || params.lines.Available() || params.concat.Available()

	} else if params.bytes.Available() {
		return params.parts.Available() || params.lines.Available() || params.concat.Available()

	} else if params.lines.Available() {
		return params.parts.Available() || params.bytes.Available() || params.concat.Available()

	} else if params.concat.Available() {
		return params.parts.Available() || params.bytes.Available() || params.lines.Available()
	}
	return false
}

func (params *parameters) oneParamHasMultipleResults() bool {
	return params.help.Count() > 1 || params.version.Count() > 1 || params.copyright.Count() > 1 || params.concat.Count() > 1 || params.input.Count() > 1 || params.output.Count() > 1 || params.parts.Count() > 1 || params.bytes.Count() > 1 || params.lines.Count() > 1
}

func (params *parameters) isInfo() bool {
	return params.help.Available() || params.version.Available() || params.copyright.Available()
}
