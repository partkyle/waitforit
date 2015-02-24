package main

import "time"

type BackoffStrategy func(debugLog, int)

func BackoffSimple(Debug debugLog, iteration int) {
	Debug.Logf("Sleeping for %ds", 1)
	time.Sleep(1 * time.Second)
}

func BackoffLinear(Debug debugLog, iteration int) {
	Debug.Logf("Sleeping for %ds", iteration)
	time.Sleep(time.Duration(iteration) * time.Second)
}

func BackoffDefault(Debug debugLog, iteration int) {}

func Backoff(backoff string) BackoffStrategy {
	switch backoff {
	case "simple":
		return BackoffSimple

	case "linear":
		return BackoffLinear
	}

	return BackoffDefault
}
