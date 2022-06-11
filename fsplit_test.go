/*
 *      Copyright 2021, 2022 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/golib/osargs"
	"testing"
)

func TestParseOSArgsA(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)
	err := params.initFromArgs(args)
	if err != nil {
		t.Error(err.Error())
	}

	args.Values = []string{"--help", "--version"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"--help", "cp"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"asdf", "qwer"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("missing command not recognized")
	}
}

func TestParseOSArgsB(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{"-p=5"}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)

	err := params.initFromArgs(args)
	if err == nil {
		t.Error("unspecified input directory not recognized")
	}

	args.Values = []string{"-c"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("unspecified input directory not recognized")
	}

	args.Values = []string{"-p=5", "a file that hopefully does not exist"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("not existing input directory not recognized")
	}

	args.Values = []string{"-c", "a file that hopefully does not exist"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		err = concatenateFiles(params)
		if err == nil {
			t.Error("not existing input directory not recognized")
		}
	} else {
		t.Error("check on file existence too early")
	}
}

func TestParseOSArgsC(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{"-p=5", "--help"}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)

	err := params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"-c", "-c"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"-c", "-p=3", "./a.txt"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"-h"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err != nil {
		t.Error("valid parameter not recognized: " + err.Error())
	}
}
