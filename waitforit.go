package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Version bool
	N       int
	Backoff string
	Debug   bool
	Verbose bool
	Command string
	Addr    string
	URL     string
}

// version value that gets replaced by the compile process
var VERSION = "I SUCK AT VERSIONING"

func main() {
	log.SetFlags(log.LstdFlags)

	config := &Config{}

	flag.IntVar(&config.N, "n", 10, "Number of times to execute command.")
	flag.StringVar(&config.Backoff, "backoff", "linear", "Backoff type. Options: simple, linear.")
	flag.BoolVar(&config.Debug, "d", false, "log debug messages.")
	flag.BoolVar(&config.Verbose, "v", false, "Print subcommand output.")

	flag.StringVar(&config.Command, "cmd", "", "Wait until a command has a clean exit (status == 0)")
	flag.StringVar(&config.Addr, "addr", "", "Wait until an address is available")
	flag.StringVar(&config.URL, "url", "", "Wait until a HTTP requests returns 2xx")

	flag.BoolVar(&config.Version, "version", false, "Print version and exit.")

	flag.Parse()

	Debug := debugLog(config.Debug)

	if config.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	var Exec Executor

	switch {
	case config.Addr != "":
		Exec = &AddrExecutor{
			Debug: Debug,
			Addr:  config.Addr,
		}
	case config.URL != "":
		Exec = &HTTPExecutor{
			Debug: Debug,
			URL:   config.URL,
		}
	case config.Command != "":
		Exec = &CommandExecutor{
			Path:    "sh",
			Args:    []string{"-c", config.Command},
			Debug:   Debug,
			Verbose: config.Verbose,
		}
	default:
		log.Fatal("Either 'http', 'addr' or 'cmd' is required.")
	}

	// determine backoff strategy from flags
	Strategy := Backoff(config.Backoff)

	var LastErr error = nil
	for i := 1; i < config.N+1; i++ {
		LastErr = Exec.Run()
		if LastErr == nil {
			break
		}

		if !config.Debug {
			fmt.Printf(".")
		}

		Strategy(Debug, i)
	}

	if LastErr != nil {
		log.Fatal(LastErr)
	}

	Debug.Logf("READY")
}
