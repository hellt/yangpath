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
		wr := new(strings.Builder)
		newModules := make(map[string]*yang.Module)
		names := make([]string, 0, len(ms.Modules))
		modName := viper.GetString("tree-module")
		qPath := viper.GetString("tree-path")
		if qPath != "" {
			if !strings.HasPrefix(qPath, "/") {
				qPath = "/" + qPath
			}
			qPath = strings.TrimRight(qPath, "/")
			fmt.Fprintf(wr, "#%s\n", qPath)
		}
		if modName != "" {
			if _, ok := ms.Modules[modName]; !ok {
				return fmt.Errorf("module %s not found", modName)
			}
			tree(wr, "", yang.ToEntry(ms.Modules[modName]), qPath)
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
				tree(wr, "", e, qPath)
			}
			return nil
		}
		if viper.GetString("tree-output") != "" {
			f, err := os.Create(viper.GetString("tree-output"))
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

func init() {
	//rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().StringP("module", "m", "", "module to print as a tree")
	viper.BindPFlag("tree-module", treeCmd.Flags().Lookup("module"))

	treeCmd.Flags().StringP("path", "p", "", "path to print as tree")
	viper.BindPFlag("tree-path", treeCmd.Flags().Lookup("path"))

	treeCmd.Flags().BoolP("meta", "", false, "print some meta information for each yang object")
	viper.BindPFlag("tree-meta", treeCmd.Flags().Lookup("meta"))

	treeCmd.Flags().StringP("output", "o", "", "output file")
	viper.BindPFlag("tree-output", treeCmd.Flags().Lookup("output"))

}

func getTypeName(e *yang.Entry) string {
	if e == nil || e.Type == nil {
		return ""
	}
	return fmt.Sprintf("%s:::%s", e.Type.Name, e.Type.Base.YangType.Kind.String())
}

func tree(w io.Writer, prefix string, e *yang.Entry, path string) {
	meta := viper.GetBool("tree-meta")
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
			fmt.Fprintf(w, "%s%s:", prefix, e.Name)
			if meta {
				fmt.Fprintf(w, " # (%s)(%s)", access(e), yangType)
			}
			fmt.Fprint(w, "\n")
			keys := strings.Split(e.Key, " ")
			for i, k := range keys {
				if i == 0 {
					tree(w, prefix+"  - ", e.Dir[k], path)
					prefix += "  "
				} else {
					tree(w, prefix+"  ", e.Dir[k], path)
				}
				if meta {
					fmt.Fprintf(w, " # (%s)(%s)", access(e.Dir[k]), e.Dir[k].Kind.String())
				}
				//getTypeMeta(w, prefix, e.Type)
			}
		} else {
			fmt.Fprintf(w, "%s%s:", prefix, e.Name)
			if meta {
				fmt.Fprintf(w, " # (%s)(%s)", access(e), yangType)
			}
			// getTypeMeta(w, prefix, e.Type)
			fmt.Fprintf(w, "\n")
		}
		indent += 2
	}
	names := make([]string, 0)
	keys := strings.Split(e.Key, " ")
	for k := range e.Dir {
		if snl(k, keys) {
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

func getTypeMeta(w io.Writer, prefix string, t *yang.YangType) {
	if t == nil {
		return
	}
	fmt.Fprintf(w, "\n%s# %s\n", prefix, t.Name)
	fmt.Fprintf(w, "%s# %s", prefix, t.Root.Name)
	if t.Kind.String() != t.Root.Name {
		fmt.Fprintf(w, " (%s)", t.Kind)
	}
	fmt.Fprint(w, "\n")
	if t.Units != "" {
		fmt.Fprintf(w, "%s# units=%s\n", prefix, t.Units)
	}
	if t.Default != "" {
		fmt.Fprintf(w, "%s# default=%q\n", prefix, t.Default)
	}
	if t.FractionDigits != 0 {
		fmt.Fprintf(w, "%s# fraction-digits=%d\n", prefix, t.FractionDigits)
	}
	if len(t.Length) > 0 {
		fmt.Fprintf(w, "%s# length=%s\n", prefix, t.Length)
	}
	if t.Kind == yang.YinstanceIdentifier && !t.OptionalInstance {
		fmt.Fprintf(w, "%s# required\n", prefix)
	}
	if t.Kind == yang.Yleafref && t.Path != "" {
		fmt.Fprintf(w, "%s# path=%q\n", prefix, t.Path)
	}
	if len(t.Pattern) > 0 {
		fmt.Fprintf(w, "%s# pattern=%s\n", prefix, strings.Join(t.Pattern, "|"))
	}
	b := yang.BaseTypedefs[t.Kind.String()].YangType
	if len(t.Range) > 0 && !t.Range.Equal(b.Range) {
		fmt.Fprintf(w, "%s# range=%s\n", prefix, t.Range)
	}
	if len(t.Type) > 0 {
		fmt.Fprintf(w, "%s# union:\n", prefix)
		for _, t := range t.Type {
			getTypeMeta(w, prefix+"  ", t)
		}
		// fmt.Fprintf(w, "\n")
	}
	//fmt.Fprintf(w, "\n")
	return
}
