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

func Handle(w http.ResponseWriter, r *http.Request) {
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
	var port = flag.Int("port", 8080, "port to serve on.")
	var enableGit = flag.Bool("git", true, "serve git repo (if present).")
	var gitPath = flag.String("git-path", "/usr/bin/git", "path to git binary.")

	flag.Parse()

	if len(flag.Args()) == 0 {
		rootPath, _ = os.Getwd()
	} else {
		rootPath = flag.Args()[0]
	}

	rootPath, _ := filepath.Abs(rootPath)

	if *enableGit {
		stat, err := os.Stat(filepath.Join(rootPath, ".git"))
		if err == nil && stat.IsDir() {
			hasGitDir = true
		}
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
	fmt.Printf("Serving %s on :%d\n", rootPath, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), http.HandlerFunc(Handle)))
}
