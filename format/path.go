package format

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
	SType  string // string representation of the Type
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
	// fmt.Printf("walkEntry is called with p=%v and ps=%v\n", p, ps)
	// fmt.Println("current entry name:", e.Name)

	switch e.Node.(type) {
	case *yang.Module: // a module has no parent
		p.Module = e.Name
	case *yang.Container:
		p.XPath += fmt.Sprintf("/%s", e.Name)
	case *yang.List:
		keys := strings.Split(e.Key, " ")
		var keyElem string
		for _, k := range keys {
			keyElem += fmt.Sprintf("[%s=*]", k)
		}
		p.XPath += fmt.Sprintf("/%s%s", e.Name, keyElem)
	case *yang.Leaf:
		p.XPath += fmt.Sprintf("/%s", e.Name)
		p.Type = e.Node.(*yang.Leaf).Type
		p.SType = e.Node.(*yang.Leaf).Type.Name
		if e.Type.IdentityBase != nil { // if the type is identityref
			p.SType += fmt.Sprintf("->%v", e.Node.(*yang.Leaf).Type.IdentityBase.Name)
		}
		if e.Type.Kind == yang.Yleafref { //handling leafref
			p.SType += fmt.Sprintf("->%v", e.Type.Path)
		}
		if e.Type.Kind == yang.Yenum { //handling enumeration types
			p.SType += fmt.Sprintf("%+q", e.Type.Enum.Names())
		}
		// fmt.Printf("appending %v path to ps=%v\n", p, ps)
		ps = append(ps, &p)
	}

	// fmt.Println("building path is", p)
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
