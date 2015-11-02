package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AaronO/go-git-http"
)

var (
	dirHandler     http.Handler
	gitHandler     http.Handler
	rootPath       string
	baseName       string
	psudoGitPrefix string

	realGitPrefix = "/.git"
	hasGitDir     = false
	css           = `<style>
pre {
 font-family: sans-serif;
}

pre.clone {
  font-family: monospace;
  border: 1px solid #ccc;
  background: #eee;
  float: left;
  padding: 20px;
  border-radius: 4px;
}
</style>
`
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.Host, r.URL.RequestURI())
	if hasGitDir && strings.HasPrefix(r.URL.RequestURI(), psudoGitPrefix) {
		gitPath := strings.Replace(r.URL.Path, psudoGitPrefix, realGitPrefix, 1)
		r.URL.Path = gitPath
		gitHandler.ServeHTTP(w, r)
		return
	}
	dirHandler.ServeHTTP(w, r)

	if hasGitDir {
		fmt.Fprintf(w, "<pre class=\"clone\">git clone http://%s%s</pre>", r.Host, psudoGitPrefix)
	}

	fmt.Fprintf(w, css)
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "address to serve on.")
	var gitPath = flag.String("git-path", "/usr/bin/git", "path to git binary.")

	flag.Parse()

	if len(flag.Args()) == 0 {
		rootPath, _ = os.Getwd()
	} else {
		rootPath = flag.Args()[0]
	}

	rootPath, _ = filepath.Abs(rootPath)

	stat, err := os.Stat(filepath.Join(rootPath, ".git"))
	if err == nil && stat.IsDir() {
		hasGitDir = true
	}

	baseName = filepath.Base(rootPath)

	dirHandler = http.FileServer(http.Dir(rootPath))

	if hasGitDir {
		psudoGitPrefix = fmt.Sprintf("/%s.git", baseName)
		gitHandler = &githttp.GitHttp{
			ProjectRoot:  rootPath,
			GitBinPath:   *gitPath,
			UploadPack:   true,
			ReceivePack:  true,
			EventHandler: func(e githttp.Event) {},
		}
	}
	fmt.Printf("Serving %s on %s\n", rootPath, *addr)
	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(handle)))
}
