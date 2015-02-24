package main_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	waitforit "github.com/partkyle/waitforit"
)

func TestCommandFail(t *testing.T) {
	cmd := &waitforit.CommandExecutor{Path: "test", Args: []string{"0", "-gt", "1"}}
	err := cmd.Run()
	if err == nil {
		t.Fail()
	}
}

func TestCommandSuccess(t *testing.T) {
	cmd := &waitforit.CommandExecutor{Path: "test", Args: []string{"1", "-gt", "0"}}
	err := cmd.Run()
	if err != nil {
		t.Errorf("expected failed command: %s", err)
	}
}

func TestAddrHttpUnavailable(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	listener.Close()

	Exec := &waitforit.AddrExecutor{Addr: fmt.Sprintf("http://%s", listener.Addr().String())}

	if err := Exec.Run(); err == nil {
		t.Error("expected http test to fail when server is not running")
	}
}

func TestAddrHttp404(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	go http.Serve(listener, nil)
	defer listener.Close()

	Exec := &waitforit.AddrExecutor{Addr: fmt.Sprintf("http://%s", listener.Addr().String())}

	if err = Exec.Run(); err == nil {
		t.Error("expected http test to fail when server responds non 2xx")
	}
}

func TestAddrHttpEventually2xx(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	called := false
	var Handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if !called {
			called = true
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
	}

	go http.Serve(listener, http.HandlerFunc(Handler))
	defer listener.Close()

	Exec := &waitforit.AddrExecutor{Addr: fmt.Sprintf("http://%s", listener.Addr().String())}

	if err = Exec.Run(); err == nil {
		t.Error("expected http test to fail when server responds non 2xx")
	}

	if err = Exec.Run(); err != nil {
		t.Errorf("should have seen 2xx status: %s", err)
	}
}

func TestAddrTCPUnavailable(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	listener.Close()

	Exec := &waitforit.AddrExecutor{Addr: listener.Addr().String()}

	if err := Exec.Run(); err == nil {
		t.Fail()
	}
}

func TestAddrTCPBecomesAvailable(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	listener.Close()

	Exec := &waitforit.AddrExecutor{Addr: fmt.Sprintf("tcp://%s", listener.Addr().String())}

	if err := Exec.Run(); err == nil {
		t.Error("Expected the AddrExecutor to fail.")
	}

	listener, err = net.Listen("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	if err := Exec.Run(); err != nil {
		t.Error(err)
	}
}
