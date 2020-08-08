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

package path

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/spf13/viper"
)

// Path represents a path in the YANG tree
type Path struct {
	Module string
	Type   *yang.Type
	XPath  string
	SType  string        // string representation of the Type
	Config yang.TriState // type of the node (config or read-only)
}

// templateInput holds HTML template variables
// Paths is a list of Path data
// Vars is a user-defined map of k/v pairs used in the template
type templateIntput struct {
	Paths []*Path
	Vars  map[string]string
}

// default template used in Template
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
	<td>{{$p.Type.Name}}</td>
  </tr>
{{end}}
</tbody>
</table>
`

// Paths recursively traverses the entry's e directory Dir till the leaf node
// populating Path structure along the way
// returns a list of pointers to the individual Path
func Paths(e *yang.Entry, p Path, ps []*Path) []*Path {
	switch e.Node.(type) {
	case *yang.Module: // a module has no parent
		p.Module = e.Name
	case *yang.Container:
		p.XPath += fmt.Sprintf("/%s", e.Name)
		if e.Config != yang.TSUnset {
			p.Config = e.Config
		}
	case *yang.List:
		if e.Config != yang.TSUnset {
			p.Config = e.Config
		}
		var keyElem string
		if e.Key != "" { // for key-less lists skip the keyElem creation
			keys := strings.Split(e.Key, " ")
			for _, k := range keys {
				keyElem += fmt.Sprintf("[%s=*]", k)
			}
		}
		p.XPath += fmt.Sprintf("/%s%s", e.Name, keyElem)
	case *yang.LeafList:
		if e.Config != yang.TSUnset {
			p.Config = e.Config
		}
	case *yang.Leaf:
		if e.Config != yang.TSUnset {
			p.Config = e.Config
		}
		p.XPath += fmt.Sprintf("/%s", e.Name)
		p.Type = e.Node.(*yang.Leaf).Type
		p.SType = e.Node.(*yang.Leaf).Type.Name

		// if the immediate type is identityref
		if e.Node.(*yang.Leaf).Type.IdentityBase != nil {
			p.SType += fmt.Sprintf("->%v", e.Node.(*yang.Leaf).Type.IdentityBase.Name)
		}

		//handling leafref
		if e.Type.Kind == yang.Yleafref {
			p.SType += fmt.Sprintf("->%v", e.Type.Path)
		}

		//handling enumeration types
		if e.Type.Kind == yang.Yenum {
			p.SType += fmt.Sprintf("%+q", e.Type.Enum.Names())
		}

		//handling union types
		if e.Type.Kind == yang.Yunion {
			var u []string // list of union types
			for _, ut := range e.Node.(*yang.Leaf).Type.Type {
				switch {
				case ut.IdentityBase != nil:
					u = append(u, fmt.Sprintf("identityref->%v", ut.IdentityBase.Name))
				case ut.YangType.Kind == yang.Yenum:
					u = append(u, fmt.Sprintf("enumeration%+q", ut.YangType.Enum.Names()))
				default:
					u = append(u, ut.Name)
				}

			}
			p.SType += fmt.Sprintf("{%v}", strings.Join(u, " "))
		}
		ps = append(ps, &p)
	}

	// ne is a nested entries list
	ne := make([]string, 0, len(e.Dir))

	for k := range e.Dir {
		ne = append(ne, k)
	}
	sort.Strings(ne)
	for _, k := range ne {
		ps = Paths(e.Dir[k], p, ps)
	}
	return ps
}

// Template take template t, paths ps and template variables vars
// and renders template to stdout
func Template(t string, ps []*Path, vars []string) error {
	// template body as string
	var outTemplate string
	switch {
	case t != "":
		data, err := ioutil.ReadFile(viper.GetString("path-template"))
		if err != nil {
			return err
		}
		outTemplate = string(data)
	default:
		outTemplate = defTemplate
	}

	tmpl, err := template.New("output-template").Parse(outTemplate)
	if err != nil {
		return err
	}
	input := &templateIntput{
		Paths: ps,
		Vars:  make(map[string]string),
	}
	for _, v := range vars {
		vk := strings.Split(v, ":::")
		if len(vk) < 2 {
			log.Printf("ignoring variable %s", v)
			continue
		}
		input.Vars[vk[0]] = strings.Join(vk[1:], ":::")
	}
	err = tmpl.Execute(os.Stdout, input)
	if err != nil {
		return err
	}
	return nil
}
