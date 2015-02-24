package main

import "log"

type debugLog bool

func (d debugLog) Logf(msg string, args ...interface{}) {
	if d {
		log.Printf(msg, args...)
	}
}
