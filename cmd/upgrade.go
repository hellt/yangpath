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

package cmd

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var downloadURL = "https://github.com/hellt/yangpath/raw/master/install.sh"

// upgradeCmd represents the version command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade yangpath",

	Run: func(cmd *cobra.Command, args []string) {
		f, err := ioutil.TempFile("", "yangpath")
		defer os.Remove(f.Name())
		if err != nil {
			log.Fatalf("Failed to create temp file %s\n", err)
		}
		downloadFile(downloadURL, f)

		c := exec.Command("bash", f.Name())
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	},
}

// downloadFile will download a file from a URL and write its content to a file
func downloadFile(url string, file *os.File) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	versionCmd.AddCommand(upgradeCmd)
}
