/*
 *          Copyright 2019, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/osargs"
)

type clResults struct {
	help      []osargs.Param
	version   []osargs.Param
	copyright []osargs.Param
	input     []osargs.Param
	output    []osargs.Param
	parts     []osargs.Param
	bytes     []osargs.Param
	lines     []osargs.Param
	concat    []osargs.Param
}

func (results *clResults) infoAvailable() bool {
	return len(results.help) > 0 || len(results.version) > 0 || len(results.copyright) > 0
}

func (results *clResults) oneParamHasMultipleResults() bool {
	for _, result := range results.toArray() {
		if len(result) > 1 {
			return true
		}
	}
	return false
}

func (results *clResults) toArray() [][]osargs.Param {
	resultsList := make([][]osargs.Param, 9)
	resultsList[0] = results.help
	resultsList[1] = results.version
	resultsList[2] = results.copyright
	resultsList[3] = results.input
	resultsList[4] = results.output
	resultsList[5] = results.parts
	resultsList[6] = results.bytes
	resultsList[7] = results.lines
	resultsList[8] = results.concat
	return resultsList
}
