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
	"encoding/json"
	"fmt"

	"sort"

	"github.com/openconfig/goyang/pkg/yang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// toJsonCmd represents the toJson command
var toJsonCmd = &cobra.Command{
	Use:   "toJson",
	Short: "yang to json",

	RunE: func(cmd *cobra.Command, args []string) error {
		ms, err := readYangDirectory(viper.GetString("yang-dir"))
		if err != nil {
			return err
		}
		newModules := make(map[string]*yang.Module)
		names := make([]string, 0, len(ms.Modules))

		for _, m := range ms.Modules {
			if newModules[m.Name] == nil {
				newModules[m.Name] = m
				names = append(names, m.Name)
			}
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Println("##### ", n)
			b, err := json.MarshalIndent(yang.ToEntry(newModules[n]), "", " ")
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			fmt.Println(string(b))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(toJsonCmd)
}
