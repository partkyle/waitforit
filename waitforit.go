package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	n       = flag.Int("n", 10, "Number of times to execute command.")
	backoff = flag.String("backoff", "linear", "Backoff type. Options: simple, linear.")
	debug   = flag.Bool("d", false, "log debug messages.")
	verbose = flag.Bool("v", false, "Print subcommand output.")

	addr = flag.String("addr", "", "Wait until an address is available")

	version = flag.Bool("version", false, "Print version and exit.")
)

// version value that gets replaced by the compile process
var VERSION = "I SUCK AT VERSIONING"

func main() {
	log.SetFlags(log.LstdFlags)

	flag.Parse()

	Debug := debugLog(*debug)

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	var Exec Executor

	switch {
	case *addr != "":
		Exec = &AddrExecutor{
			Debug: Debug,
			Addr:  *addr,
		}
	default:
		if len(flag.Args()) < 1 {
			log.Fatal("Command is required!")
		}

		Path := flag.Args()[0]
		var Args []string
		if len(flag.Args()) > 1 {
			Args = flag.Args()[1:len(flag.Args())]
		}

		Exec = &CommandExecutor{
			Path:  Path,
			Args:  Args,
			Debug: Debug,
		}
	}

	// determine backoff strategy from flags
	Strategy := Backoff(*backoff)

	var LastErr error = nil
	for i := 1; i < *n+1; i++ {
		LastErr = Exec.Run()
		if LastErr == nil {
			break
		}

		if !*debug {
			fmt.Printf(".")
		}

		Strategy(Debug, i)
	}

	if LastErr != nil {
		log.Fatal(LastErr)
	}

	Debug.Logf("READY")
}
