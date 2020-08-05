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
	"strings"

	"github.com/hellt/yangform/format"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "generate xpath or restconf style paths from yang files",

	RunE: func(cmd *cobra.Command, args []string) error {
		names, ms, err := loadAndSortModules(viper.GetStringSlice("yang-dir"))
		if err != nil {
			return err
		}
		modName := viper.GetString("module")
		if !snl(modName, names) && modName != "" {
			return fmt.Errorf("unknown module: %s", modName)
		}
		e := yang.ToEntry(ms.Modules[modName])

		paths := format.Paths(e, "", []string{})
		for _, path := range paths {
			fmt.Println(path)
		}
		// paths := make([]*path, 0)
		// pc := make(chan *path, 0)
		// go func() {
		// 	if modName != "" {
		// 		// fmt.Printf("%+v\n", ms.Modules[modName].Container)
		// 		spew.Dump(yang.ToEntry(ms.Modules[modName]))
		// 		for _, c := range ms.Modules[modName].Container {
		// 			addContainerToPath(modName, "", c, pc)
		// 		}
		// 	} else {
		// 		for _, mn := range names {
		// 			for _, c := range ms.Modules[mn].Container {
		// 				addContainerToPath(mn, "", c, pc)
		// 			}
		// 		}
		// 	}
		// 	close(pc)
		// }()
		// for p := range pc {
		// 	p.XPath = gnmiPathToXPath(p.Path)
		// 	//p.RestconfPath = gnmiPathToRestconfPath(p.Path)
		// 	if viper.GetString("format") == "text" {
		// 		fmt.Printf("%s | %s | %s\n", p.Module, p.XPath, p.Type)
		// 	}
		// 	//fmt.Printf("%s | %s | %s\n", p.Module, p.RestconfPath, p.Type)
		// 	paths = append(paths, p)
		// }
		// if viper.GetString("format") == "html" {
		// 	outTemplate := defTemplate
		// 	if viper.GetString("path-template") != "" {
		// 		data, err := ioutil.ReadFile(viper.GetString("path-template"))
		// 		if err != nil {
		// 			return err
		// 		}
		// 		outTemplate = string(data)
		// 	}

		// 	tmpl, err := template.New("output-template").Parse(outTemplate)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	input := &templateIntput{
		// 		Paths: paths,
		// 		Vars:  make(map[string]string),
		// 	}
		// 	for _, v := range viper.GetStringSlice("path-template-vars") {
		// 		vk := strings.Split(v, ":::")
		// 		if len(vk) < 2 {
		// 			log.Printf("ignoring variable %s", v)
		// 			continue
		// 		}
		// 		input.Vars[vk[0]] = strings.Join(vk[1:], ":::")
		// 	}
		// 	err = tmpl.Execute(os.Stdout, input)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	pathCmd.Flags().StringP("type", "t", "xpath", "path types, xpath or restconf")
	viper.BindPFlag("path-type", pathCmd.Flags().Lookup("type"))
	pathCmd.Flags().StringP("template", "", "", "path to golang html template to use instead of the default one")
	viper.BindPFlag("path-template", pathCmd.Flags().Lookup("template"))
	pathCmd.Flags().StringSliceP("template-vars", "", []string{}, "extra template variables in case a custom template is used for html output")
	viper.BindPFlag("path-template-vars", pathCmd.Flags().Lookup("template-vars"))
}

func gnmiPathToXPath(p *gnmi.Path) string {
	if p == nil {
		return ""
	}
	pathElems := make([]string, 0, len(p.GetElem()))
	for _, pe := range p.GetElem() {
		elem := ""
		if pe.GetName() != "" {
			elem += pe.GetName()
		}
		if pe.GetKey() != nil {
			for k, v := range pe.GetKey() {
				elem += fmt.Sprintf("[%s=%s]", k, v)
			}
		}
		pathElems = append(pathElems, elem)
	}
	return "/" + strings.Join(pathElems, "/")
}
func gnmiPathToRestconfPath(p *gnmi.Path) string {
	if p == nil {
		return ""
	}
	pathElems := make([]string, 0, len(p.GetElem()))
	for _, pe := range p.GetElem() {
		elem := ""
		if pe.GetName() != "" {
			elem += pe.GetName()
		}
		if pe.GetKey() != nil {
			for k, v := range pe.GetKey() {
				elem += fmt.Sprintf("%s=%s", k, v)
			}
		}
		pathElems = append(pathElems, elem)
	}
	return strings.Join(pathElems, "/")
}
