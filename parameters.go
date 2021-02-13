/*
 *      Copyright 2019 - 2021, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"github.com/vbsw/cmdl"
)

type parameters struct {
	help      *cmdl.Parameter
	version   *cmdl.Parameter
	copyright *cmdl.Parameter
	concat    *cmdl.Parameter
	input     *cmdl.Parameter
	output    *cmdl.Parameter
	parts     *cmdl.Parameter
	bytes     *cmdl.Parameter
	lines     *cmdl.Parameter
}

func parametersFromOSArgs() (*parameters, error) {
	var params *parameters
	var err error
	cl := cmdl.New()
	asgOp := cmdl.NewAsgOp(true, true, "=")

	if len(cl.Args()) > 0 {
		params = new(parameters)
		params.help = cl.NewParam().Parse("-h", "--help", "-help", "help")
		params.version = cl.NewParam().Parse("-v", "--version", "-version", "version")
		params.copyright = cl.NewParam().Parse("--copyright", "-copyright", "copyright")
		params.concat = cl.NewParam().Parse("-c", "--concat", "-concat", "concat")
		params.input = cl.NewParam().ParsePairs(asgOp, "-i", "--input", "-input", "input")
		params.output = cl.NewParam().ParsePairs(asgOp, "-o", "--output", "-output", "output")
		params.parts = cl.NewParam().ParsePairs(asgOp, "-p", "--parts", "-parts", "parts")
		params.bytes = cl.NewParam().ParsePairs(asgOp, "-b", "--bytes", "-bytes", "bytes")
		params.lines = cl.NewParam().ParsePairs(asgOp, "-l", "--lines", "-lines", "lines")

		unparsedArgs := cl.UnparsedArgs()
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
			params.input.Add("<none>", unparsedArgs[0])
			unparsedArgs = unparsedArgs[1:]
		}
	}
	return unparsedArgs
}

func (params *parameters) parseOutput(unparsedArgs []string) []string {
	if params.output.Available() {
		// just accept the first unparsed argument, if output wasn't set explicitly
		if len(unparsedArgs) > 0 {
			params.output.Add("<none>", unparsedArgs[0])
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
