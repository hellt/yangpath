/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/hellt/yangform/format"
	"github.com/openconfig/goyang/pkg/yang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "generate xpath or restconf style paths from yang files",

	RunE: func(cmd *cobra.Command, args []string) error {
		modName := viper.GetString("module")
		names, ms, errs := format.LoadAndSortModules(viper.GetStringSlice("yang-dir"), modName)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%v\n", err)
			}
		}
		if !snl(modName, names) && modName != "" {
			return fmt.Errorf("unknown module: %s", modName)
		}
		e := yang.ToEntry(ms.Modules[modName])

		paths := format.Paths(e, format.Path{}, []*format.Path{})
		if viper.GetString("path-format") == "text" {
			for _, path := range paths {
				var ps string // path string to print

				if viper.GetString("path-with-module") == "yes" {
					ps += fmt.Sprintf("%s    ", path.Module)
				}
				ps += fmt.Sprintf("%s", path.XPath)
				if viper.GetString("path-with-types") == "yes" {
					ps += fmt.Sprintf("    %s", path.Type.Name)
				}
				fmt.Println(ps)
			}
		}

		if viper.GetString("path-format") == "html" {
			t := viper.GetString("path-template") // path to the template file
			vars := viper.GetStringSlice("path-template-vars")
			format.Template(t, paths, vars)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	pathCmd.Flags().StringP("format", "f", "text", "paths output format. One of [text, html]")
	viper.BindPFlag("path-format", pathCmd.Flags().Lookup("format"))

	pathCmd.Flags().StringP("type", "t", "xpath", "path types, xpath or restconf")
	viper.BindPFlag("path-type", pathCmd.Flags().Lookup("type"))

	pathCmd.Flags().StringP("with-module", "", "no", "print module name")
	viper.BindPFlag("path-with-module", pathCmd.Flags().Lookup("with-module"))

	pathCmd.Flags().StringP("with-types", "", "yes", "display path type information")
	viper.BindPFlag("path-with-types", pathCmd.Flags().Lookup("with-types"))

	pathCmd.Flags().StringP("template", "", "", "path to HTML template to use instead of the default one")
	viper.BindPFlag("path-template", pathCmd.Flags().Lookup("template"))

	pathCmd.Flags().StringSliceP("template-vars", "", []string{}, "extra template variables in case a custom template is used. Key value pairs separated with ::: delimiter")
	viper.BindPFlag("path-template-vars", pathCmd.Flags().Lookup("template-vars"))
}
