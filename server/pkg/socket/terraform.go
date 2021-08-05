// Copyright 2021 Ramon Quitales
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package socket

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/tools/txtar"
)

// startKubectl saves a yaml file and runs the specified kubectl action on the yaml file,
// sending its output and end event as Messages on the provided channel.
func startTerraform(id, action, body string, dest chan<- *Message, opt *Options) *process {
	var (
		done = make(chan struct{})
		out  = make(chan *Message)
		p    = &process{out: out, done: done}
	)
	go func() {
		defer close(done)
		for m := range buffer(limiter(out, p), time.After) {
			m.Id = id
			dest <- m
		}
	}()

	err := p.startTerraform(id, action, body, opt)
	if err != nil {
		p.end(err)
		return nil
	}
	go func() {
		p.end(p.run.Wait())
	}()
	return p
}

// startTerraform saves terraform config files to the folder alongside the server binary
// location. This is done instead of using a temp folder as we likely want to persist the
// terraform state files and allow us to do a manual tf destroy in case the server unexpectedly
// crashes before the presenter can tear down the env.
func (p *process) startTerraform(id, action, body string, opt *Options) error {
	// We save the body to a yaml then kubectl apply it.

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to find binary path: %w", err)
	}

	// Don't worry about cleaning up this folder as the server admin should choose what to do
	// with it after the presentation.
	exePath := filepath.Dir(exe)
	tfPath := filepath.Join(exePath, fmt.Sprintf("terraform-%s", id))
	err = os.MkdirAll(tfPath, 0777)
	if err != nil {
		if err != os.ErrExist {
			return err
		}
	}

	if !(action == "destroy") {
		// write body to x.tf files
		a := txtar.Parse([]byte(body))
		if len(a.Comment) != 0 {
			a.Files = append(a.Files, txtar.File{Name: "main.tf", Data: a.Comment})
			a.Comment = nil
		}
		for _, f := range a.Files {
			err = ioutil.WriteFile(filepath.Join(tfPath, f.Name), f.Data, 0666)
			if err != nil {
				return err
			}
		}

		args := []string{"terraform", "init"}
		cmd := p.cmd(tfPath, args...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unable to init terraform: %w", err)
		}
	}
	// auto-approve flag required as cmds run non-interractively.
	args := []string{"terraform", action, "-auto-approve"}
	cmd := p.cmd(tfPath, args...)
	// cmd.Stdout = cmd.Stderr // send compiler output to stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	p.run = cmd
	p.kind = kubectl
	return nil
}
