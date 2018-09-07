// main.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

/*
This is just a very short example.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keltia/ssllabs"
)

const (
	// MyVersion is for the app
	MyVersion = "0.1.0"
)

var (
	fDebug       bool
	fDetailed    bool
	fInfo        bool
	fVerbose     bool
	fShowVersion bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.BoolVar(&fInfo, "I", false, "Get SSLLabs info.")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fShowVersion, "V", false, "Display version & exit.")
	flag.Parse()

	if fShowVersion {
		fmt.Fprintf(os.Stderr, "%s version %s API v3\n",
			MyName, ssllabs.Version())
		os.Exit(0)
	}

	if fInfo {
		fmt.Fprintf(os.Stderr, "SSLLabs server info:\n")
		return
	}

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	var level = 0

	site := flag.Arg(0)

	if fVerbose {
		level = 1
	}

	if fDebug {
		level = 2
		fVerbose = true
	}

	// Setup client
	c, err := ssllabs.NewClient(ssllabs.Config{Log: level})
	if err != nil {
		log.Fatalf("error setting up client: %v", err)
	}

	if fInfo {
		info, err := c.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		jinfo, err := json.Marshal(info)
		fmt.Fprintf(os.Stderr, "SSLLabs Info\n%#s", string(jinfo))
		os.Exit(0)
	}

	if fDetailed {

		report, err := c.GetDetailedReport(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}

		// Just dump the json
		fmt.Printf("%v\n", report)
	} else {
		fmt.Fprintf(os.Stderr, "%s Wrapper: %s API version %s\n\n",
			MyName, MyVersion, ssllabs.Version())
		grade, err := c.GetGrade(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s': %v\n", site, err)
		}
		fmt.Printf("Grade for '%s' is %s\n", site, grade)
	}
}
