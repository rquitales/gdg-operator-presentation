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
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// IsFolder validates that the path provided is a folder.
func IsFolder(path string) (string, error) {
	if path == "" {
		return "", errors.New("path cannot be empty")
	}

	pathToServe, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	pathStat, err := os.Stat(pathToServe)
	if err != nil {
		return "", err
	}

	if !pathStat.IsDir() {
		return "", fmt.Errorf("path must be a directory: %q", path)
	}

	return pathToServe, nil
}
