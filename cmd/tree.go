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
	"os"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/indent"
	"github.com/openconfig/goyang/pkg/yang"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "output a yang model in a yaml like tree format",

	RunE: func(cmd *cobra.Command, args []string) error {
		ms, err := readYangDirectory(viper.GetString("yang-dir"))
		if err != nil {
			return err
		}
		errs := ms.Process()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Errorf("%v\n", err)
			}
		}
		newModules := make(map[string]*yang.Module)
		names := make([]string, 0, len(ms.Modules))
		modName := viper.GetString("tree-module")
		qPath := viper.GetString("tree-path")
		if qPath != "" {
			if !strings.HasPrefix(qPath, "/") {
				qPath = "/" + qPath
			}
			qPath = strings.TrimRight(qPath, "/")
			fmt.Printf("%s\n", qPath)
		}
		if modName != "" {
			if _, ok := ms.Modules[modName]; !ok {
				return fmt.Errorf("module %s not found", modName)
			}
			tree(os.Stdout, "", yang.ToEntry(ms.Modules[modName]), qPath)
		} else {
			for _, m := range ms.Modules {
				if newModules[m.Name] == nil {
					newModules[m.Name] = m
					names = append(names, m.Name)
				}
			}
			sort.Strings(names)
			entries := make([]*yang.Entry, 0, len(names))
			for _, n := range names {
				entries = append(entries, yang.ToEntry(newModules[n]))
			}
			for _, e := range entries {
				tree(os.Stdout, "", e, qPath)
			}
			return nil
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().StringP("module", "m", "", "module to print as a tree")
	viper.BindPFlag("tree-module", treeCmd.Flags().Lookup("module"))
	treeCmd.Flags().StringP("path", "p", "", "path to print as tree")
	viper.BindPFlag("tree-path", treeCmd.Flags().Lookup("path"))
}

func getTypeName(e *yang.Entry) string {
	if e == nil || e.Type == nil {
		return ""
	}

	return fmt.Sprintf("%s:::%s", e.Type.Name, e.Type.Base.YangType.Kind.String())
}

func write2(w io.Writer, prefix string, e *yang.Entry) {
	fmt.Printf("%s %s\n", prefix, strings.Repeat("#", 70))
	fmt.Printf("%s Name         : '%s'\n", prefix, e.Name)
	fmt.Printf("%s Kind         : '%s'\n", prefix, e.Kind.String())
	fmt.Printf("%s Type         : '%s'\n", prefix, getTypeName(e))
	fmt.Printf("%s Path         : '%s'\n", prefix, e.Path())
	fmt.Printf("%s Namespace    : '%s'\n", prefix, e.Namespace().NName())
	fmt.Printf("%s DefaultValue : '%s'\n", prefix, e.DefaultValue())
	fmt.Printf("%s ReadOnly     : '%t'\n", prefix, e.ReadOnly())
	if e.IsCase() {
		fmt.Printf("%s isCase       : '%t'\n", prefix, e.IsCase())
	}
	if e.IsChoice() {
		fmt.Printf("%s IsChoice     : '%t'\n", prefix, e.IsChoice())
	}
	if e.IsContainer() {
		fmt.Printf("%s IsContainer  : '%t'\n", prefix, e.IsContainer())
	}
	fmt.Printf("%s IsDir        : '%t'\n", prefix, e.IsDir())
	if e.IsLeaf() {
		fmt.Printf("%s IsLeaf       : '%t'\n", prefix, e.IsLeaf())
	}
	if e.IsLeafList() {
		fmt.Printf("%s IsLeafList   : '%t'\n", prefix, e.IsLeafList())
	}
	if e.IsList() {
		fmt.Printf("%s IsList       : '%t'\n", prefix, e.IsList())
	}
	if e.RPC != nil {
		fmt.Printf("%s RPC          : '%v'\n", prefix, e.RPC)
	}
	if len(e.Exts) > 0 {
		fmt.Printf("%s Extensions: \n", prefix)
		for _, ext := range e.Exts {
			if n := ext.NName(); n != "" {
				fmt.Printf("%s  %s %s;\n", prefix, ext.Kind(), n)
			} else {
				fmt.Printf("%s  %s;\n", prefix, ext.Kind())
			}
		}
	}
	switch {
	// case e.Dir == nil && e.ListAttr != nil:
	// 	fmt.Printf("%s []%s (%s)\n", prefix, e.Name, getTypeName(e))
	// 	return
	// case e.Dir == nil:
	// 	fmt.Printf("%s %s (%s)\n", prefix, e.Name, getTypeName(e))
	// 	return
	// case e.ListAttr != nil:
	// 	fmt.Printf("%s %s[%s] (%s)\n", prefix, e.Name, e.Key, getTypeName(e)) //}
	// default:
	// 	fmt.Printf("%s %s (%s)\n", prefix, e.Name, getTypeName(e)) //}
	}
	names := make([]string, 0)
	for k := range e.Dir {
		names = append(names, k)
	}
	sort.Strings(names)
	prefix += "  "
	for _, k := range names {
		write2(indent.NewWriter(w, "  "), prefix, e.Dir[k])
	}
}

func tree(w io.Writer, prefix string, e *yang.Entry, path string) {
	if e.RPC != nil {
		return
	}
	indent := 0
	if path == "" || (path != "" && strings.HasPrefix(e.Path(), path)) {
		yangType := e.Kind.String()
		switch {
		case e.IsContainer():
			yangType = "Container"
		case e.IsLeafList():
			yangType = "LeafList"
		case e.IsList():
			yangType = "List"
		case e.IsChoice():
			yangType = "Choice"
		case e.IsCase():
			yangType = "Case"
		case e.IsLeaf():
			yangType = "Leaf"
		}

		if e.Key != "" {
			fmt.Fprintf(w, "%s%s: # (%s)(%s)\n", prefix, e.Name, access(e), yangType)
			for i, k := range strings.Split(e.Key, " ") {
				if i == 0 {
					fmt.Fprintf(w, "%s  - %s:  # isKey (%s)(%s) type=%s", prefix, k, access(e.Dir[k]), e.Dir[k].Kind.String(), getTypeName(e.Dir[k]))
					prefix += "  "
				} else {
					fmt.Fprintf(w, "%s  %s:  # isKey (%s)(%s) type=%s", prefix, k, access(e.Dir[k]), e.Dir[k].Kind.String(), getTypeName(e.Dir[k]))
				}
				if e.IsLeaf() {
					if e.DefaultValue() != "" {
						fmt.Fprintf(w, "  default=%s", e.DefaultValue())
					}
				}
				if e.Type != nil {
					fmt.Fprintf(w, "  yangType=%s", getTypeName(e))
				}
				fmt.Fprintf(w, "\n")
			}
		} else {
			fmt.Fprintf(w, "%s%s: # (%s)(%s)", prefix, e.Name, access(e), yangType)
			if e.IsLeaf() {
				if e.DefaultValue() != "" {
					fmt.Fprintf(w, "  default=%s", e.DefaultValue())
				}
			}
			if e.Type != nil {
				fmt.Fprintf(w, "  yangType=%s", getTypeName(e))
			}
			fmt.Fprintf(w, "\n")
		}
		indent += 2
	}
	names := make([]string, 0)
	for k := range e.Dir {
		if snl(k, strings.Split(e.Key, " ")) {
			continue
		}
		names = append(names, k)
	}
	prefix += strings.Repeat(" ", indent)
	sort.Strings(names)
	for _, k := range names {
		tree(w, prefix, e.Dir[k], path)
	}
}

func snl(s string, l []string) bool {
	for _, sl := range l {
		if s == sl {
			return true
		}
	}
	return false
}

func access(e *yang.Entry) string {
	access := "rw"
	switch {
	case e.RPC != nil:
		access = "rpc"
	case e.ReadOnly():
		access = "ro"
	}
	return access
}
