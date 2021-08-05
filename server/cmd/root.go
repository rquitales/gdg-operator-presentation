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

package cmd

import (
	"fmt"
	"os"

	"github.com/rquitales/go-presentation-server/cmd/server"
	"github.com/spf13/cobra"
)

var (
	folder, addr string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "present",
	Short: "Start a simple fileserver with a websocket for code execution.",
	Run:   server.Serve(&folder, &addr),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&folder, "folder", "", "path to folder containing static assets")
	rootCmd.Flags().StringVar(&addr, "address", "localhost:8080", "the address to serve on")
	rootCmd.MarkFlagRequired("folder")
}
