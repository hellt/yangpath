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
	"os"
	"sort"
	"strings"

	"github.com/google/gnxi/utils/xpath"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type path struct {
	Module       string
	Path         *gnmi.Path
	Type         string
	XPath        string
	RestconfPath string
}
type templateIntput struct {
	Paths []*path
	Vars  map[string]string
}

var defTemplate = `
<table class="table table-striped">
<thead>
  <tr>
	<th>#</th>
	<th>Module</th>
	<th>Path</th>
	<th>Leaf Type</th>
  </tr>
</thead>
<tbody>
{{range $i, $p  := .Paths}}
<tr>
	<td>{{$i}}</td>
	<td>{{$p.Module}}</td>
	<td>{{$p.XPath}}</td>
	<td>{{$p.Type}}</td>
  </tr>
{{end}}
</tbody>
</table>
`

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

		paths := walkEntry(e, "", []string{})
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

// walkEntry traverses the entry's directory Dir till the leaf node
// aggregating the xpath path
// returns a list of accumulated paths
func walkEntry(e *yang.Entry, p string, ps []string) []string {
	// fmt.Printf("walkEntry is called with p=%v and ps=%v\n", p, ps)
	// fmt.Println("current entry name:", e.Name)

	switch e.Node.(type) {
	case *yang.Module: // a module has no parent
	case *yang.Container:
		p += fmt.Sprintf("/%s", e.Name)
	case *yang.List:
		keys := strings.Split(e.Key, " ")
		var key string
		for _, k := range keys {
			key += fmt.Sprintf("[%s=*]", k)
		}
		p += fmt.Sprintf("/%s%s", e.Name, key)
	case *yang.Leaf:
		p += fmt.Sprintf("/%s", e.Name)
		// fmt.Printf("appending %v path to ps=%v\n", p, ps)
		ps = append(ps, p)
	}

	// fmt.Println("building path is", p)
	// ne is a nested entries list
	ne := make([]string, 0, len(e.Dir))

	for k := range e.Dir {
		ne = append(ne, k)
	}
	sort.Strings(ne)
	for _, k := range ne {
		ps = walkEntry(e.Dir[k], p, ps)
	}
	return ps
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

func addContainerToPath(module, prefix string, container *yang.Container, out chan *path) {
	fmt.Println("in 1")
	elementName := fmt.Sprintf("%s/%s", prefix, container.Name)
	fmt.Println(elementName)
	for _, c := range container.Container {
		fmt.Println("in 2")
		addContainerToPath(module, elementName, c, out)
	}
	for _, ls := range container.List {
		fmt.Println("in 3")
		addListToPath(module, elementName, ls, out)
	}
	for _, lf := range container.Leaf {
		sp := fmt.Sprintf("%s/%s", elementName, lf.Name)
		gnmiPath, err := xpath.ToGNMIPath(sp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "path: %s could not be changed to gnmi: %v\n", sp, err)
			continue
		}
		p := &path{
			Module: module,
			Path:   gnmiPath,
			Type:   lf.Type.Name,
		}
		out <- p
	}
}
func addListToPath(module, prefix string, ls *yang.List, out chan *path) {
	keys := strings.Split(ls.Key.Name, " ")
	keyElem := ls.Name
	for _, k := range keys {
		keyElem += fmt.Sprintf("[%s=*]", k)
	}
	elementName := fmt.Sprintf("%s/%s", prefix, keyElem)
	for _, c := range ls.Container {
		addContainerToPath(module, elementName, c, out)
	}
	for _, lls := range ls.List {
		addListToPath(module, elementName, lls, out)
	}
	for _, ch := range ls.Choice {
		for _, ca := range ch.Case {
			addCaseToPath(module, elementName, ca, out)
		}
	}
	for _, lf := range ls.Leaf {
		if lf.Name != ls.Key.Name {
			sp := fmt.Sprintf("%s/%s", elementName, lf.Name)
			gnmiPath, err := xpath.ToGNMIPath(sp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "path: %s could not be changed to gnmi: %v\n", sp, err)
				continue
			}
			p := &path{
				Module: module,
				Path:   gnmiPath,
				Type:   lf.Type.Name,
			}
			out <- p
		}
	}
}
func addCaseToPath(module, prefix string, ca *yang.Case, out chan *path) {
	for _, cont := range ca.Container {
		addContainerToPath(module, prefix, cont, out)
	}
	for _, ls := range ca.List {
		addListToPath(module, prefix, ls, out)
	}
	for _, lf := range ca.Leaf {
		sp := fmt.Sprintf("%s/%s", prefix, lf.Name)
		gnmiPath, err := xpath.ToGNMIPath(sp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "path: %s could not be changed to gnmi: %v\n", sp, err)
			continue
		}
		p := &path{
			Module: module,
			Path:   gnmiPath,
			Type:   lf.Type.Name,
		}
		out <- p
	}
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
