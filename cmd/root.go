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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ygdlv",
	Short: "yang delve",

	RunE: func(cmd *cobra.Command, args []string) error {
		names, ms, err := loadAndSortModules(viper.GetString("yang-dir"))
		if err != nil {
			return err
		}
		wr := new(strings.Builder)
		modName := viper.GetString("module")
		if !snl(modName, names) && modName != "" {
			return fmt.Errorf("unknown module: %s", modName)
		}
		qPath := viper.GetString("path")
		if qPath != "" {
			if !strings.HasPrefix(qPath, "/") {
				qPath = "/" + qPath
			}
			qPath = strings.TrimRight(qPath, "/")
			fmt.Fprintf(wr, "#%s\n", qPath)
		}
		if modName != "" {
			toYaml(wr, "", yang.ToEntry(ms.Modules[modName]), qPath)
		} else {
			for _, n := range names {
				toYaml(wr, "", yang.ToEntry(ms.Modules[n]), qPath)
			}
		}
		outFile := viper.GetString("output")
		if outFile != "" {
			f, err := os.Create(outFile)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.WriteString(wr.String())
			if err != nil {
				return err
			}
			return nil
		}
		fmt.Println(wr.String())
		return nil
	},
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
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yangdelve.yaml)")

	rootCmd.PersistentFlags().StringP("yang-dir", "y", "", "yang directory")
	viper.BindPFlag("yang-dir", rootCmd.PersistentFlags().Lookup("yang-dir"))

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringP("format", "f", "yaml", "output format")
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))

	rootCmd.PersistentFlags().StringP("module", "m", "", "module to export")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	rootCmd.PersistentFlags().StringP("path", "p", "", "path to object to export")
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))

	rootCmd.PersistentFlags().StringP("output", "o", "", "output file name")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	rootCmd.PersistentFlags().BoolP("meta", "", false, "print entries metadata")
	viper.BindPFlag("meta", rootCmd.PersistentFlags().Lookup("meta"))
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

	//	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
func readYangDirectory(dir string) (*yang.Modules, error) {
	yFiles := make([]string, 0)
	fChan := make(chan string, 0)
	defer close(fChan)
	go func() {
		for {
			select {
			case fName := <-fChan:
				yFiles = append(yFiles, fName)
			}
		}
	}()
	err := getYangfiles(dir, fChan)
	if err != nil {
		return nil, err
	}
	ms := yang.NewModules()
	for _, f := range yFiles {
		err = ms.Read(f)
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}
func loadAndSortModules(dir string) ([]string, *yang.Modules, error) {
	ms, err := readYangDirectory(dir)
	if err != nil {
		return nil, nil, err
	}
	errs := ms.Process()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("%v\n", err)
		}
	}
	names := make([]string, 0, len(ms.Modules))
	for _, m := range ms.Modules {
		if snl(m.Name, names) {
			continue
		}
		names = append(names, m.Name)
	}
	sort.Strings(names)
	return names, ms, nil
}

func toYaml(w io.Writer, prefix string, e *yang.Entry, path string) {
	if e == nil {
		return
	}
	indent := 0
	meta := viper.GetBool("meta")
	keys := strings.Split(e.Key, " ")
	if path == "" || (path != "" && strings.HasPrefix(e.Path(), path+"/")) || (path == e.Path()) {
		fmt.Fprintf(w, "%s%s:", prefix, e.Name)
		if meta {
			fmt.Fprintf(w, " # (%s)", access(e))
			// getTypeMeta(w, prefix, e.Type)
		}
		fmt.Fprint(w, "\n")
		switch {
		case e.IsList():
			for i, k := range keys {
				if i == 0 {
					fmt.Fprintf(w, "%s  - %s:", prefix, k)
					indent += 2
				} else {
					fmt.Fprintf(w, "%s%s:", prefix, k)
				}
				if meta {
					fmt.Fprintf(w, " # (%s)", access(e))
				}
				fmt.Fprint(w, "\n")
			}
		}
		indent += 2
	}
	names := make([]string, 0)
	for n := range e.Dir {
		if !snl(n, keys) {
			names = append(names, n)
		}
	}
	prefix += strings.Repeat(" ", indent)
	sort.Strings(names)
	for _, k := range names {
		toYaml(w, prefix, e.Dir[k], path)
	}
}

func getYangfiles(dir string, fChan chan string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".yang" {
			fChan <- filepath.Join(dir, f.Name())
		}
		if f.IsDir() {
			err = getYangfiles(filepath.Join(dir, f.Name()), fChan)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
