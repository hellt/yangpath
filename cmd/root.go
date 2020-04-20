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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yangdelve",
	Short: "A brief description of your application",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yangdelve.yaml)")

	rootCmd.PersistentFlags().StringP("yang-dir", "", "", "yang directory")
	rootCmd.PersistentFlags().BoolP("debug", "", false, "debug")
	viper.BindPFlag("yang-dir", rootCmd.PersistentFlags().Lookup("yang-dir"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".yangdelve" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".yangdelve")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func readYangDirectory(dir string) (*yang.Modules, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	yfiles := make([]os.FileInfo, 0)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yang" {
			yfiles = append(yfiles, f)
		}
	}
	ms := yang.NewModules()
	for _, f := range yfiles {
		ms.Read(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}
