package main

import (
	"fmt"
	"net"
	"net/http"
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
	conn, err := net.Dial("tcp", a.Addr)
	if err != nil {
		a.Debug.Logf("dial: %s", err)
		return err
	}

	if err := conn.Close(); err != nil {
		a.Debug.Logf("close: %s", err)
		return err
	}

	return nil
}

type HTTPExecutor struct {
	Debug debugLog
	URL   string
}

func (h *HTTPExecutor) Run() error {
	response, err := http.Get(h.URL)
	if err != nil {
		h.Debug.Logf("http.Get: %s", err)
		return err
	}

	if (response.StatusCode / 200) != 1 {
		err := fmt.Errorf("expected 2xx status: %d", response.StatusCode)
		h.Debug.Logf("status: %s", err)
		return err
	}

	return nil
}

type CommandExecutor struct {
	Debug   debugLog
	Path    string
	Args    []string
	Verbose bool
}

func (c *CommandExecutor) Run() error {
	Exec := exec.Command(c.Path, c.Args...)

	if c.Verbose {
		Exec.Stdout = os.Stdout
		Exec.Stderr = os.Stderr
	}

	c.Debug.Logf("Running %s %s", Exec.Path, Exec.Args)

	return Exec.Run()
}
