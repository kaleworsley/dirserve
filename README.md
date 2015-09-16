# dirserve

Serve a directory with HTTP. Like `python -m SimpleHTTPServer` with extra stuff.

## Installation

```
go get -u github.com/kaleworsley/dirserve
```

## Usage

    Usage of dirserve:
      -addr string
            address to serve on. (default "localhost:8080")
      -git
            serve git repo (if present). (default true)
      -git-path string
            path to git binary. (default "/usr/bin/git")
