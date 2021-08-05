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

package filepath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const dirName = "presentation-server-tests"

func initTestDirectory(numFiles int) (string, error) {
	td, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		return "", err
	}

	for i := 0; i < numFiles; i++ {
		filename := filepath.Join(td, fmt.Sprintf("file_%d.txt", i))

		_, err := os.Create(filename)
		if err != nil {
			return td, fmt.Errorf("unable to create temporary file: %s", filename)
		}
	}

	return td, nil
}

func TestIsFolder(t *testing.T) {
	tempDir, err := initTestDirectory(2)
	if err != nil {
		t.Fatalf("unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			"Happy path - is folder",
			tempDir,
			tempDir,
			false,
		},
		{
			"Not a folder",
			filepath.Join(tempDir, "file_1.txt"),
			"",
			true,
		},
		{
			"Fake folder",
			filepath.Join(tempDir, "/fake-dir/"),
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsFolder(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}
