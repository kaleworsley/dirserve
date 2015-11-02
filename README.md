# dirserve

Serve a directory with HTTP. Like `python -m SimpleHTTPServer` with extra stuff.

## Installation

```
go get -u github.com/kaleworsley/dirserve
```

## Usage

    usage: dirserve [OPTIONS] [DIRECTORY | .]

    OPTIONS

      -addr string
            address to serve on. (default "localhost:8080")
      -git-path string
            path to git binary. (default "/usr/bin/git")
