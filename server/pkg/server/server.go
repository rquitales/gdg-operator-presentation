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

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rquitales/go-presentation-server/client/crd"
	"github.com/rquitales/go-presentation-server/pkg/filepath"
	"github.com/rquitales/go-presentation-server/pkg/kubectl"
	"github.com/rquitales/go-presentation-server/pkg/socket"
)

// Serve creates a simple file server for a specified folder and serving
// address. A websocket endpoint is also created for the handling of code
// execution.
func Serve(path string, addr string) {
	pathToServe, err := filepath.IsFolder(path)
	if err != nil {
		log.Fatalf("Unable to get static file path: %s", err)
	}

	log.Printf("Serving presentation at: %s\n", addr)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	origin := &url.URL{
		Scheme: "http",
		Host:   addr,
	}

	// Handles code execution.
	mux.Handle("/socket", socket.NewHandler(origin))
	mux.HandleFunc("/crd/", handleCRD)
	mux.Handle("/", http.FileServer(http.Dir(pathToServe)))

	log.Println(server.ListenAndServe())
}

func handleCRD(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		getDetails(name, w)
	} else {
		getAllCRDNames(w)
	}
}

func getAllCRDNames(w http.ResponseWriter) {
	output, err := kubectl.GetAllCRDs()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var data struct {
		Items []crd.CRD `json:"items"`
	}
	json.Unmarshal(output, &data)

	names := make([]string, len(data.Items))
	for i, v := range data.Items {
		names[i] = v.Metadata.Name
	}

	fmt.Fprint(w, names)
}

func getDetails(name string, w http.ResponseWriter) {
	output, err := kubectl.GetCRDSpec(name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var data crd.CRD
	json.Unmarshal(output, &data)

	formatted, _ := json.MarshalIndent(parseDetails(data), "", "    ")

	fmt.Fprint(w, string(formatted))
}

type Details struct {
	Version string
	Spec    interface{}
}

// TODO(rquitales): loop through all versions, instead of hard-coding first version.
func parseDetails(data crd.CRD) Details {

	details := Details{
		Version: data.Spec.Versions[0].Name,
		Spec:    data.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties,
	}

	return details
}
