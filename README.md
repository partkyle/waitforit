# waitforit

A simple go program to run a command and wait for it to successfully execute.

It is meant for applicaitons that require services to be available before they can start.

## Usage

### Command

Waits until a command returns a 0 exit status.

```
waitforit -cmd "curl www.google.com"
```

### Addr

Waits until a network is reachable.

```
waitforit -addr "localhost:8080"
```

### HTTP

Waits until a url responds with a 2xx request

```
waitforit -http "http://localhost:8080/health"
```
