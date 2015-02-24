package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type Executor interface {
	Run() error
}

type AddrExecutor struct {
	Debug debugLog
	Addr  string
}

func (a *AddrExecutor) Run() error {
	url, err := url.Parse(a.Addr)
	if err != nil {
		return err
	}

	a.Debug.Logf("scheme: %s host: %s", url.Scheme, url.Host)

	switch url.Scheme {
	case "http":
		response, err := http.Get(url.String())
		if err != nil {
			a.Debug.Logf("http.Get: %s", err)
			return err
		}

		if (response.StatusCode / 200) != 1 {
			err := fmt.Errorf("expected 2xx status: %d", response.StatusCode)
			a.Debug.Logf("status: %s", err)
			return err
		}
	case "tcp":
		conn, err := net.Dial(url.Scheme, url.Host)
		if err != nil {
			a.Debug.Logf("dial: %s", err)
			return err
		}

		if err := conn.Close(); err != nil {
			a.Debug.Logf("close: %s", err)
			return err
		}
	default:
		return fmt.Errorf("invalid scheme: %s", url.Scheme)
	}

	return nil
}

type CommandExecutor struct {
	Debug debugLog
	Path  string
	Args  []string
}

func (c *CommandExecutor) Run() error {
	Exec := exec.Command(c.Path, c.Args...)

	if *verbose {
		Exec.Stdout = os.Stdout
		Exec.Stderr = os.Stderr
	}

	c.Debug.Logf("Running %s %s", Exec.Path, Exec.Args)

	return Exec.Run()
}
