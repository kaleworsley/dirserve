# dirserve

Serve a directory with HTTP. Like `python -m SimpleHTTPServer` with extra stuff.

## Installation

```
go get -u github.com/kaleworsley/dirserve
```

## Usage

    Usage of dirserve:
      -git=true: serve git repo (if present).
      -git-path="/usr/bin/git": path to git binary.
      -bind-address="": address to bind to, defaults to all.
      -port=8080: port to serve on.
