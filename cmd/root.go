/*
Copyright © 2020 Roman Dodin <dodin.roman@gmail.com>

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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yangpath",
	Short: "yang path exporter",
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
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().StringSliceP("yang-dir", "y", []string{""}, "directory(ies) with YANG modules")
	viper.BindPFlag("yang-dir", rootCmd.PersistentFlags().Lookup("yang-dir"))

	// rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	// viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringP("module", "m", "", "module to export")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))
}

// snl is a string-in-list-of-strings checking func
func snl(s string, l []string) bool {
	for _, sl := range l {
		if s == sl {
			return true
		}
	}
	return false
}
