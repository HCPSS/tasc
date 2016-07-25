package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gosuri/uilive"
)

const (
	// BANNER lets people know what this is.
	BANNER = `  _______        _____  _____
 |__   __|/\    / ____|/ ____|
    | |  /  \  | (___ | |
    | | / /\ \  \___ \| |
    | |/ ____ \ ____) | |____
    |_/_/    \_\_____/ \_____|

A tool of assembling source code.

Version %s
Copyright (C) 2016 Howard County Public Schools
Distributed under the terms of the MIT license
Written by Brendan Anderson

`
	// VERSION of the application.
	VERSION = "v0.2.1"
)

var (
	manifest         Manifest
	destinationDir   string
	extraParams      map[string]string
	manifestFilename string

	version bool
)

func init() {
	var extraParamsJSON string

	flag.StringVar(&manifestFilename, "manifest", "tasc-manifest.yml",
		"Name of the manifest file.")
	flag.StringVar(&destinationDir, "destination", "./",
		"Where to build the project")
	flag.StringVar(&extraParamsJSON, "params", "{}",
		"A JSON encoded string with extra parameters.")

	flag.BoolVar(&version, "version", false, "Print the version.")
	flag.BoolVar(&version, "v", false, "Print the version.")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("tasc %s\n", VERSION)
		os.Exit(0)
	}

	// extraParamsJson is a JSON encoded string, so we need to decode it.
	extraParams := make(map[string]string)
	json.Unmarshal([]byte(extraParamsJSON), &extraParams)

	// Add some constants to params
	extraParams["manifest_dir"] = filepath.Dir(manifestFilename)
	extraParams["destination_dir"] = destinationDir

	// Load the manifest
	err := manifest.Load(manifestFilename, extraParams)
	if err != nil {
		panic(err)
	}
}

func main() {
	tasc := Tasc{manifest: manifest, destination: destinationDir}

	c := make(chan string)
	go tasc.Assemble(c)

	writer := uilive.New()
	writer.Start()

	for s := range c {
		fmt.Fprintf(writer, s)
		writer.Flush()
	}

	writer.Flush()
	writer.Stop()

	results := tasc.Patch()

	// Report on the success/failure of patches.
	if len(results) > 0 {
		numSuccess := len(results.GetSuccess())
		if numSuccess > 0 {
			fmt.Printf("%d patches successfully applied.\n", numSuccess)
		}

		numFailed := len(results.GetFailed())
		if len(results.GetFailed()) > 0 {
			fmt.Printf(
				"%d patches failed to apply. Errors are listed below:\n",
				numFailed,
			)
			for _, r := range results.GetFailed() {
				fmt.Printf("%s: %s\n", r.Patcher.GetSource(), r.Error.Error())
			}
		}
	}
}
