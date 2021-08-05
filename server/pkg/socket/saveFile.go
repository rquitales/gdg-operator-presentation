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
	"io/ioutil"
	"time"
)

// startKubectl saves a yaml file and runs the specified kubectl action on the yaml file,
// sending its output and end event as Messages on the provided channel.
func startSaveFile(id, path, body string, dest chan<- *Message, opt *Options) *process {
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

	err := p.startSaveFile(path, body, opt)
	if err != nil {
		p.end(err)
		return nil
	}
	go func() {
		p.end(nil)
	}()
	return p
}

// startKubectl saves a yaml and kubectl apply/create/deletes it, sending its output to p.out,
// and stores the running *exec.Cmd in the run field.
func (p *process) startSaveFile(path string, body string, opt *Options) error {

	err := ioutil.WriteFile(path, []byte(body), 0666)
	if err != nil {
		return err
	}
	p.run = nil
	return nil
}
