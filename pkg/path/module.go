// Copyright © 2020 Roman Dodin <dodin.roman@gmail.com>
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

package path

import (
	"github.com/openconfig/goyang/pkg/yang"
)

// GetModuleName takes a path p to the YANG file and returns the module name
func GetModuleName(p string) (string, error) {
	var name string
	ms := yang.NewModules()
	if err := ms.Read(p); err != nil {
		return "", err
	}
	for _, v := range ms.Modules {
		name = v.Name
	}
	return name, nil
}

//AddYANGDirs adds directories which have YANG files inside taking dirs as a list of directories
func AddYANGDirs(dirs []string) error {
	for _, path := range dirs {
		expanded, err := yang.PathsWithModules(path)
		if err != nil {
			return err
		}
		yang.AddPath(expanded...)
	}
	return nil
}
