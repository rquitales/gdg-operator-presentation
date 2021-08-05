// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !appengine
// +build !appengine

// Package socket implements an WebSocket-based playground backend.
// Clients connect to a websocket handler and send run/kill commands, and
// the server sends the output and exit status of the running processes.
// Multiple clients running multiple processes may be served concurrently.
// The wire format is JSON and is described by the Message type.
//
// This will not run on App Engine as WebSockets are not supported there.
package socket // import "golang.org/x/tools/playground/socket"

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	exec "golang.org/x/sys/execabs"

	"golang.org/x/net/websocket"
	"golang.org/x/tools/txtar"
)

// Environ provides an environment when a binary, such as the go tool, is
// invoked.
var Environ func() []string = os.Environ

const (
	// The maximum number of messages to send per session (avoid flooding).
	msgLimit = 1000

	// Batch messages sent in this interval and send as a single message.
	msgDelay = 10 * time.Millisecond
)

// Message is the wire format for the websocket connection to the browser.
// It is used for both sending output messages and receiving commands, as
// distinguished by the Kind field.
type Message struct {
	Id      string // client-provided unique id for the process
	Kind    string // in: "run", "kill" out: "stdout", "stderr", "end"
	Body    string
	Options *Options `json:",omitempty"`
	Path    string   //if saving a file
}

// Options specify additional message options.
type Options struct {
	Race bool // use -race flag when building code (for "run" only)
}

type runKind string

const (
	shell   runKind = "shell"
	golang  runKind = "go"
	kubectl runKind = "kubectl"
)

// NewHandler returns a websocket server which checks the origin of requests.
func NewHandler(origin *url.URL) websocket.Server {
	return websocket.Server{
		Config:    websocket.Config{Origin: origin},
		Handshake: handshake,
		Handler:   websocket.Handler(socketHandler),
	}
}

// handshake checks the origin of a request during the websocket handshake.
func handshake(c *websocket.Config, req *http.Request) error {
	// o, err := websocket.Origin(c, req)
	// if err != nil {
	// 	log.Println("bad websocket origin:", err)
	// 	return websocket.ErrBadWebSocketOrigin
	// }
	// _, port, err := net.SplitHostPort(c.Origin.Host)
	// if err != nil {
	// 	log.Println("bad websocket origin:", err)
	// 	return websocket.ErrBadWebSocketOrigin
	// }
	// ok := c.Origin.Scheme == o.Scheme && (c.Origin.Host == o.Host || c.Origin.Host == net.JoinHostPort(o.Host, port) || strings.Contains(o.Host, "localhost") || strings.Contains(o.Host, "rquitales"))
	// if !ok {
	// 	log.Println("bad websocket origin:", o)
	// 	return websocket.ErrBadWebSocketOrigin
	// }
	log.Println("accepting connection from:", req.RemoteAddr)
	return nil
}

// socketHandler handles the websocket connection for a given present session.
// It handles transcoding Messages to and from JSON format, and starting
// and killing processes.
func socketHandler(c *websocket.Conn) {
	in, out := make(chan *Message), make(chan *Message)
	errc := make(chan error, 1)

	// Decode messages from client and send to the in channel.
	go func() {
		dec := json.NewDecoder(c)
		for {
			var m Message
			if err := dec.Decode(&m); err != nil {
				errc <- err
				return
			}
			in <- &m
		}
	}()

	// Receive messages from the out channel and encode to the client.
	go func() {
		enc := json.NewEncoder(c)
		for m := range out {
			if err := enc.Encode(m); err != nil {
				errc <- err
				return
			}
		}
	}()
	defer close(out)

	// Start and kill processes and handle errors.
	proc := make(map[string]*process)
	for {
		select {
		case m := <-in:
			switch m.Kind {
			case "run":
				log.Println("running code snippet from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				var wd string
				if p, ok := proc[m.Id]; ok {
					wd = p.wd
				}
				proc[m.Id] = startProcess(m.Id, m.Body, out, m.Options, wd)
			case "saveFile":
				proc[m.Id].Kill()
				proc[m.Id] = startSaveFile(m.Id, m.Path, m.Body, out, m.Options)
			case "kubectlApply":
				log.Println("running kubectl apply from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				proc[m.Id] = startKubectl(m.Id, "apply", m.Body, out, m.Options)
			case "kubectlCreate":
				log.Println("running kubectl create from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				proc[m.Id] = startKubectl(m.Id, "create", m.Body, out, m.Options)
			case "kubectlDelete":
				log.Println("running kubectl delete from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				proc[m.Id] = startKubectl(m.Id, "delete", m.Body, out, m.Options)
			case "terraformApply":
				log.Println("running terraform apply from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				proc[m.Id] = startTerraform(m.Id, "apply", m.Body, out, m.Options)
			case "terraformDestroy":
				log.Println("running terraform destroy from:", c.Request().RemoteAddr)
				proc[m.Id].Kill()
				proc[m.Id] = startTerraform(m.Id, "destroy", m.Body, out, m.Options)
			case "kill":
				proc[m.Id].Kill()
			}
		case err := <-errc:
			if err != io.EOF {
				// A encode or decode has failed; bail.
				log.Println(err)
			}
			// Shut down any running processes.
			for _, p := range proc {
				p.Kill()
			}
			return
		}
	}
}

// process represents a running process.
type process struct {
	out  chan<- *Message
	done chan struct{} // closed when wait completes
	run  *exec.Cmd
	path string
	kind runKind
	wd   string
}

// startProcess builds and runs the given program, sending its output
// and end event as Messages on the provided channel.
func startProcess(id, body string, dest chan<- *Message, opt *Options, wd string) *process {
	var (
		done = make(chan struct{})
		out  = make(chan *Message)
		p    = &process{out: out, done: done, wd: wd}
	)
	go func() {
		defer close(done)
		for m := range buffer(limiter(out, p), time.After) {
			m.Id = id
			dest <- m
		}
	}()
	var err error
	if path, args := shebang(body); path != "" {
		err = p.startProcess(path, args, body)
	} else {
		err = p.start(body, opt)
	}
	if err != nil {
		p.end(err)
		return nil
	}
	go func() {
		p.end(p.run.Wait())
	}()
	return p
}

// end sends an "end" message to the client, containing the process id and the
// given error value. It also removes the binary, if present.
func (p *process) end(err error) {
	if p.path != "" {
		defer os.RemoveAll(p.path)
	}
	m := &Message{Kind: "end"}
	if err != nil {
		m.Body = err.Error()
	}
	p.out <- m
	close(p.out)
}

// A killer provides a mechanism to terminate a process.
// The Kill method returns only once the process has exited.
type killer interface {
	Kill()
}

// limiter returns a channel that wraps the given channel.
// It receives Messages from the given channel and sends them to the returned
// channel until it passes msgLimit messages, at which point it will kill the
// process and pass only the "end" message.
// When the given channel is closed, or when the "end" message is received,
// it closes the returned channel.
func limiter(in <-chan *Message, p killer) <-chan *Message {
	out := make(chan *Message)
	go func() {
		defer close(out)
		n := 0
		for m := range in {
			switch {
			case n < msgLimit || m.Kind == "end":
				out <- m
				if m.Kind == "end" {
					return
				}
			case n == msgLimit:
				// Kill in a goroutine as Kill will not return
				// until the process' output has been
				// processed, and we're doing that in this loop.
				go p.Kill()
			default:
				continue // don't increment
			}
			n++
		}
	}()
	return out
}

// buffer returns a channel that wraps the given channel. It receives messages
// from the given channel and sends them to the returned channel.
// Message bodies are gathered over the period msgDelay and coalesced into a
// single Message before they are passed on. Messages of the same kind are
// coalesced; when a message of a different kind is received, any buffered
// messages are flushed. When the given channel is closed, buffer flushes the
// remaining buffered messages and closes the returned channel.
// The timeAfter func should be time.After. It exists for testing.
func buffer(in <-chan *Message, timeAfter func(time.Duration) <-chan time.Time) <-chan *Message {
	out := make(chan *Message)
	go func() {
		defer close(out)
		var (
			tc    <-chan time.Time
			buf   []byte
			kind  string
			flush = func() {
				if len(buf) == 0 {
					return
				}
				out <- &Message{Kind: kind, Body: safeString(buf)}
				buf = buf[:0] // recycle buffer
				kind = ""
			}
		)
		for {
			select {
			case m, ok := <-in:
				if !ok {
					flush()
					return
				}
				if m.Kind == "end" {
					flush()
					out <- m
					return
				}
				if kind != m.Kind {
					flush()
					kind = m.Kind
					if tc == nil {
						tc = timeAfter(msgDelay)
					}
				}
				buf = append(buf, m.Body...)
			case <-tc:
				flush()
				tc = nil
			}
		}
	}()
	return out
}

// Kill stops the process if it is running and waits for it to exit.
func (p *process) Kill() {
	if p == nil || p.run == nil {
		return
	}

	if p.kind == shell {
		// Explicitly kill process group ID if running shell commands.
		syscall.Kill(-p.run.Process.Pid, syscall.SIGKILL)
	} else {
		p.run.Process.Kill()
	}

	<-p.done // block until process exits
}

// shebang looks for a shebang ('#!') at the beginning of the passed string.
// If found, it returns the path and args after the shebang.
// args includes the command as args[0].
func shebang(body string) (path string, args []string) {
	body = strings.TrimSpace(body)
	if !strings.HasPrefix(body, "#!") {
		return "", nil
	}
	if i := strings.Index(body, "\n"); i >= 0 {
		body = body[:i]
	}
	fs := strings.Fields(body[2:])
	return fs[0], fs
}

// startProcess starts a given program given its path and passing the given body
// to the command standard input.
func (p *process) startProcess(path string, args []string, body string) error {
	cmdString := strings.Split(body, "\n")
	cmdSlice := strings.Split(cmdString[1], " ")
	log.Println(cmdSlice)

	if len(cmdSlice) > 0 {
		switch cmdSlice[0] {
		case "cd":
			if len(cmdSlice) < 2 {
				return errors.New("provide path!")
			}

			homedir, _ := os.UserHomeDir()
			loc := strings.Replace(cmdSlice[1], "~", homedir, -1)
			if loc[0] != byte('/') {
				currDir, _ := os.Getwd()
				if p.wd != "" {
					p.wd = filepath.Join(p.wd, loc)
				} else {
					p.wd = filepath.Join(currDir, loc)
				}
				body = strings.Replace(body, loc, p.wd, -1)
			} else {
				p.wd = loc
			}

		}
	}

	if p.wd != "" {
		err := os.Chdir(p.wd)
		if err != nil {
			return err
		}

	}

	cmd := &exec.Cmd{
		Path:   path,
		Args:   args,
		Stdin:  strings.NewReader(body),
		Stdout: &messageWriter{kind: "stdout", out: p.out},
		Stderr: &messageWriter{kind: "stderr", out: p.out},
	}

	log.Println("body: ", body)

	// Assign a process group ID that all child processes will belong to.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	p.run = cmd
	p.kind = shell
	return nil
}

// start builds and starts the given program, sending its output to p.out,
// and stores the running *exec.Cmd in the run field.
func (p *process) start(body string, opt *Options) error {
	// We "go build" and then exec the binary so that the
	// resultant *exec.Cmd is a handle to the user's program
	// (rather than the go tool process).
	// This makes Kill work.

	path, err := ioutil.TempDir("", "present-")
	if err != nil {
		return err
	}
	p.path = path // to be removed by p.end

	out := "prog"
	if runtime.GOOS == "windows" {
		out = "prog.exe"
	}
	bin := filepath.Join(path, out)

	// write body to x.go files
	a := txtar.Parse([]byte(body))
	if len(a.Comment) != 0 {
		a.Files = append(a.Files, txtar.File{Name: "prog.go", Data: a.Comment})
		a.Comment = nil
	}
	hasModfile := false
	for _, f := range a.Files {
		err = ioutil.WriteFile(filepath.Join(path, f.Name), f.Data, 0666)
		if err != nil {
			return err
		}
		if f.Name == "go.mod" {
			hasModfile = true
		}
	}

	// build x.go, creating x
	args := []string{"go", "build", "-tags", "OMIT"}
	if opt != nil && opt.Race {
		p.out <- &Message{
			Kind: "stderr",
			Body: "Running with race detector.\n",
		}
		args = append(args, "-race")
	}
	args = append(args, "-o", bin)
	cmd := p.cmd(path, args...)
	if !hasModfile {
		cmd.Env = append(cmd.Env, "GO111MODULE=off")
	} else {
		cmd := p.cmd(path, "go", "mod", "tidy")
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = p.cmd("", bin)
	if opt != nil && opt.Race {
		cmd.Env = append(cmd.Env, "GOMAXPROCS=2")
	}
	if err := cmd.Start(); err != nil {
		// If we failed to exec, that might be because they built
		// a non-main package instead of an executable.
		// Check and report that.
		if name, err := packageName(body); err == nil && name != "main" {
			return errors.New(`executable programs must use "package main"`)
		}
		return err
	}
	p.run = cmd
	p.kind = golang
	return nil
}

// cmd builds an *exec.Cmd that writes its standard output and error to the
// process' output channel.
func (p *process) cmd(dir string, args ...string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = Environ()
	cmd.Stdout = &messageWriter{kind: "stdout", out: p.out}
	cmd.Stderr = &messageWriter{kind: "stderr", out: p.out}
	return cmd
}

func packageName(body string) (string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "prog.go",
		strings.NewReader(body), parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return f.Name.String(), nil
}

// messageWriter is an io.Writer that converts all writes to Message sends on
// the out channel with the specified id and kind.
type messageWriter struct {
	kind string
	out  chan<- *Message
}

func (w *messageWriter) Write(b []byte) (n int, err error) {
	w.out <- &Message{Kind: w.kind, Body: safeString(b)}
	return len(b), nil
}

// safeString returns b as a valid UTF-8 string.
func safeString(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var buf bytes.Buffer
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		b = b[size:]
		buf.WriteRune(r)
	}
	return buf.String()
}
