// main.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

/*
This is just a very short example.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/keltia/ssllabs"
)

const (
	// MyVersion is for the app
	MyVersion = "0.4.0"

	// Display remote info
	InfoFmt = "SSLLabs Info\nEngine/%s Criteria/%s Max assessments/%d\nMessage: %s\n"
)

var (
	fDebug       bool
	fDetailed    bool
	fForce       bool
	fInfo        bool
	fVerbose     bool
	fShowVersion bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.BoolVar(&fForce, "F", false, "Do not use SSLLabs cache")
	flag.BoolVar(&fInfo, "I", false, "Get SSLLabs info.")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fShowVersion, "V", false, "Display version & exit.")
	flag.Parse()
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

	var cfg = ssllabs.Config{Log: level}

	if fForce {
		cfg.Force = true
	}

	// Setup client
	c, err := ssllabs.NewClient(cfg)
	if err != nil {
		log.Fatalf("error setting up client: %v", err)
	}

	if fShowVersion {
		fmt.Fprintf(os.Stderr, "%s/%s API/%s(v3)\n",
			MyName, MyVersion, ssllabs.Version())
		os.Exit(0)
	}

	if fInfo {
		info, err := c.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, InfoFmt, info.EngineVersion, info.CriteriaVersion, info.MaxAssessments, info.Messages[0])
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}

	report, err := c.GetDetailedReport(site)
	if err != nil {
		log.Fatalf("impossible to get grade for '%s'\n", site)
	}

	fmt.Fprintf(os.Stderr, "%s/%s API/%s\n\n",
		MyName, MyVersion, ssllabs.Version())

	if fDetailed {
		// Just dump the json
		fmt.Printf("%v\n", report)
	} else {
		grade, err := c.GetGrade(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s': %v\n", site, err)
		}
		d := time.Unix(report.TestTime/1000, 0).Local()
		fmt.Printf("Grade for '%s' is %s (%s)\n", site, grade, d)
	}
}
