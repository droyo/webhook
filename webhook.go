// webhook is a github Webhook listener that runs subordinate
// scripts when a webhook is received. If webhooks for the same
// repository are received during the waiting period, they are
// coalesced into the same call.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Author    Person `json:"author"`
}

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Created     int    `json:"created_at"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
	Forks       int    `json:"forks"`
	Homepage    string `json:"homepage"`
	Issues      int    `json:"open_issues"`
	Language    string `json:"language"`
	Owner       Person `json:"owner"`
	Private     bool   `json:"private"`
	Pushed      int    `json:"pushed_at"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
}

type WebHook struct {
	Before  string     `json:"before"`
	After   string     `json:"after"`
	Ref     string     `json:"ref"`
	Commits []Commit   `json:"commits"`
	Repo    Repository `json:"repository"`
}

type Server struct {
	filter *regexp.Regexp
	prog   []string
	sync.Mutex
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
	flag.Usage = func() {
		log.Printf("Usage: %s [-a addr] [-f regex] prog ...", os.Args[0])
		os.Exit(2)
	}
}

func main() {
	var (
		addr   = flag.String("a", ":7149", "Address to listen on")
		filter = flag.String("f", ".*", "Repository URLs to listen for (regex)")
	)
	flag.Parse()
	
	if len(flag.Args()) < 1 {
		flag.Usage()
	}
	
	pattern, err := regexp.Compile(*filter)
	if err != nil {
		log.Fatal(err)
	}
	
	srv := Server{
		filter: pattern,
		prog:   flag.Args(),
	}
	
	http.Handle("/", srv)
	log.Print("Listening on ", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func (srv Server) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	var hook WebHook
	log.Printf("New connection from %s", r.RemoteAddr)
	
	if r.Method != "POST" {
		log.Printf("(%s) Rejecting non-POST request", r.RemoteAddr)
		http.Error(w, "Method not allowed", 405)
		return
	}
	input, _, err := r.FormFile("payload")
	if err != nil {
		log.Printf("(%s) Could not parse POST body: %s", r.RemoteAddr, err)
		http.Error(w, "Bad request", 400)
		return
	}
	dec := json.NewDecoder(input)
	if err := dec.Decode(&hook); err != nil {
		log.Printf("(%s) Error parsing JSON %s", r.RemoteAddr, err)
		http.Error(w, "Bad request", 400)
		return
	}
	go srv.RunHook(hook, r)
}

func (srv Server) RunHook(hook WebHook, r *http.Request) {
	url := hook.Repo.Url
	if ! srv.filter.MatchString(url) {
		log.Printf("Discarding non-matching request %s from %s", url, r.RemoteAddr)
		return
	}
	
	cmd := exec.Command(srv.prog[0], srv.prog[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env,
		fmt.Sprint("REPO_FORKS=", hook.Repo.Forks),
		fmt.Sprint("REPO_PRIVATE=", hook.Repo.Private),
		fmt.Sprint("REPO_CREATED=", hook.Repo.Created),
		fmt.Sprint("REPO_PUSHED=", hook.Repo.Pushed),
		"REPO_NAME="+hook.Repo.Name,
		"REPO_URL="+hook.Repo.Url,
		"REPO_DESCRIPTION="+hook.Repo.Description,
		"REPO_HOMEPAGE="+hook.Repo.Homepage,
		"REPO_OWNER_NAME="+hook.Repo.Owner.Name,
		"REPO_OWNER_EMAIL="+hook.Repo.Owner.Email,
		"WEBHOOK_BEFORE="+hook.Before,
		"WEBHOOK_AFTER="+hook.After,
		"WEBHOOK_REF="+hook.Ref,
	)
	if len(hook.Commits) > 0 {
		jsonPipe, err := cmd.StdinPipe()
		if err != nil {
			log.Printf("(%s) Could not create pipe: %s", r.RemoteAddr, err)
		} else {
			enc := json.NewEncoder(jsonPipe)
			go enc.Encode(hook.Commits)
		}
	}
	// We guarantee that the command will not be run at the same time.
	srv.Lock()
	defer srv.Unlock()

	log.Printf("(%s) running program %s", r.RemoteAddr, srv.prog)
	if err := cmd.Run(); err != nil {
		log.Printf("Execution of %s failed: %s", srv.prog[0], err)
	}
}
