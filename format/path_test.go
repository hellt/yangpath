package format

import (
	"log"
	"reflect"
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
)

func TestPaths(t *testing.T) {
	tests := map[string]struct {
		dirs   []string
		module string
		want   []*Path
	}{
		"test1": {
			dirs: []string{"testdata/test1"}, module: "test1",
			want: []*Path{
				{Module: "test1", Type: "string", XPath: "/c1/l1[key1=*]/key1"},
				{Module: "test1", Type: "string", XPath: "/c1/l1[key1=*]/leaf2"},
			}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, ms, errs := LoadAndSortModules(tc.dirs, tc.module)
			if len(errs) > 0 {
				for _, err := range errs {
					t.Errorf("%v\n", err)
				}
			}
			e := yang.ToEntry(ms.Modules[tc.module])

			got := Paths(e, Path{}, []*Path{})
			log.Printf("%#v\n", got)
			if !reflect.DeepEqual(tc.want, got) {
				// t.Fatalf("expected: %+v, got: %+v", *tc.want[0], *got)
				for i, v := range tc.want {
					if v.Module != got[i].Module {
						t.Logf("Module wanted %s got %s\n", v.Module, got[i].Module)
					}
					if v.XPath != got[i].XPath {
						t.Logf("XPATH wanted %s got %s\n", v.XPath, got[i].XPath)
					}
					if v.Type != got[i].Type {
						t.Logf("Type wanted %s got %s\n", v.Type, got[i].Type)
					}
				}
				t.Fatalf("Objects not equal!")
			}
		})
	}

}
