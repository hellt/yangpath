package format

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openconfig/goyang/pkg/yang"
)

// Path represents a path in the YANG tree
type Path struct {
	Module string
	// Path         *gnmi.Path
	Type         string
	XPath        string
	RestconfPath string
}
type templateIntput struct {
	Paths []*Path
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
		var key string
		for _, k := range keys {
			key += fmt.Sprintf("[%s=*]", k)
		}
		p.XPath += fmt.Sprintf("/%s%s", e.Name, key)
	case *yang.Leaf:
		p.XPath += fmt.Sprintf("/%s", e.Name)
		p.Type = e.Type.Name
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
