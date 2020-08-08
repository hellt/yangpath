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
	"fmt"
	"strings"

	path "github.com/hellt/yangpath/pkg/path"
	"github.com/openconfig/goyang/pkg/yang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export xpath-styled paths from a given YANG module",

	RunE: func(cmd *cobra.Command, args []string) error {
		modName := viper.GetString("module")
		names, ms, errs := path.LoadAndSortModules(viper.GetStringSlice("yang-dir"), modName)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%v\n", err)
			}
		}
		if !snl(modName, names) && modName != "" {
			return fmt.Errorf("unknown module: %s", modName)
		}
		e := yang.ToEntry(ms.Modules[modName])

		paths := path.Paths(e, path.Path{}, []*path.Path{})
		if viper.GetString("path-format") == "text" {
			for _, path := range paths {
				var ps []string // path string to print as a slice

				if viper.GetString("path-with-module") == "yes" {
					ps = append(ps, path.Module)
				}
				ps = append(ps, path.XPath)
				switch viper.GetString("path-with-types") {
				case "yes":
					ps = append(ps, path.Type.Name)
				case "expanded":
					ps = append(ps, path.SType)
				}
				fmt.Println(strings.Join(ps, "    "))
			}
		}

		if viper.GetString("path-format") == "html" {
			t := viper.GetString("path-template") // path to the template file
			vars := viper.GetStringSlice("path-template-vars")
			if err := path.Template(t, paths, vars); err != nil {
				log.Fatal(err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("format", "f", "text", "paths output format. One of [text, html]")
	viper.BindPFlag("path-format", exportCmd.Flags().Lookup("format"))

	exportCmd.Flags().StringP("type", "t", "xpath", "path types, xpath or restconf")
	viper.BindPFlag("path-type", exportCmd.Flags().Lookup("type"))

	exportCmd.Flags().StringP("with-module", "", "no", "print module name")
	viper.BindPFlag("path-with-module", exportCmd.Flags().Lookup("with-module"))

	exportCmd.Flags().StringP("with-types", "", "yes", "display path type information")
	viper.BindPFlag("path-with-types", exportCmd.Flags().Lookup("with-types"))

	exportCmd.Flags().StringP("template", "", "", "path to HTML template to use instead of the default one")
	viper.BindPFlag("path-template", exportCmd.Flags().Lookup("template"))

	exportCmd.Flags().StringSliceP("template-vars", "", []string{}, "extra template variables in case a custom template is used. Key value pairs separated with ::: delimiter")
	viper.BindPFlag("path-template-vars", exportCmd.Flags().Lookup("template-vars"))
}
