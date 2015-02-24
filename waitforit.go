package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	n       = flag.Int("n", 10, "Number of times to execute command.")
	backoff = flag.String("backoff", "linear", "Backoff type. Options: simple, linear.")
	debug   = flag.Bool("d", false, "log debug messages.")
	verbose = flag.Bool("v", false, "Print subcommand output.")

	version = flag.Bool("version", false, "Print version and exit.")
)

var VERSION = "I SUCK AT VERSIONING"

func main() {
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	Debug := debugLog(*debug)

	Debug.Logf("Running Command: %s", flag.Args())

	if len(flag.Args()) < 1 {
		log.Fatal("Command is required!")
	}

	Command := flag.Args()[0]
	var Args []string
	if len(flag.Args()) > 2 {
		Args = flag.Args()[1:len(flag.Args())]
	}

	// determine backoff strategy from flags
	Strategy := Backoff(*backoff)

	var LastErr error = nil
	for i := 1; i < *n+1; i++ {
		Exec := exec.Command(Command, Args...)

		if *verbose {
			Exec.Stdout = os.Stdout
			Exec.Stderr = os.Stderr
		}

		Debug.Logf("Running %s %s", Exec.Path, Exec.Args)

		LastErr = Exec.Run()
		Debug.Logf("state: %v, %v", Exec.ProcessState, LastErr)

		if LastErr == nil {
			// This is the successful case.
			// We are done.
			break
		}

		// we don't want log collusion with this
		// but this is a useful output for most cases
		if !*debug {
			fmt.Printf(".")
		}
		Strategy(Debug, i)
	}

	if LastErr != nil {
		log.Fatal(LastErr)
	}

	Debug.Logf("exiting normally")
}
