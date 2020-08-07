package path

import (
	"fmt"
	"os"
	"sort"

	"github.com/openconfig/goyang/pkg/yang"
)

// LoadAndSortModules loads and sort YANG module m
// using the scope referenced by a list of dirs
// returns module names, *yang.Modules and a list of errors encountered
func LoadAndSortModules(dirs []string, m string) ([]string, *yang.Modules, []error) {
	// for each yang directory referenced with yang-dir flag
	// perform a search for directories with YANG files inside
	for _, path := range dirs {
		expanded, err := yang.PathsWithModules(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		yang.AddPath(expanded...)
	}
	ms := yang.NewModules()

	if err := ms.Read(m); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	errs := ms.Process()

	names := make([]string, 0)
	for _, m := range ms.Modules {
		if snl(m.Name, names) {
			continue
		}
		names = append(names, m.Name)
	}
	sort.Strings(names)
	return names, ms, errs
}

// snl is a string-in-list-of-strings checking func
func snl(s string, l []string) bool {
	for _, sl := range l {
		if s == sl {
			return true
		}
	}
	return false
}
