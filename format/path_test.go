package format

import (
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
		"test2": {
			dirs: []string{"testdata/test2"}, module: "test2",
			want: []*Path{
				{Module: "test2", Type: "string", XPath: "/c1/l1[key1=*]/key1"},
				{Module: "test2", Type: "string", XPath: "/c1/l1[key1=*]/leaf2"},
			}},
		"test3": {
			dirs: []string{"testdata/test3"}, module: "test3",
			want: []*Path{
				{Module: "test3", Type: "string", XPath: "/c1/l1[key1=*][key2=*]/key1"},
				{Module: "test3", Type: "age", XPath: "/c1/l1[key1=*][key2=*]/key2"},
				{Module: "test3", Type: "int64", XPath: "/c1/l1[key1=*][key2=*]/leaf1"},
			}},
		"test4": {
			dirs: []string{"testdata/test4"}, module: "test4",
			want: []*Path{
				{Module: "test4", Type: "identityref -> test4:IDENTITY2", XPath: "/c1/leaf1"},
				{Module: "test4", Type: "identityref -> IDENTITY1", XPath: "/c1/leaf2"},
			}},
		"test5": {
			dirs: []string{"testdata/test5"}, module: "test5",
			want: []*Path{
				{Module: "test5", Type: "string", XPath: "/c1/leaf1"},
				{Module: "test5", Type: "leafref -> ../leaf1", XPath: "/c1/leaf2"},
			}},
		"test6": {
			dirs: []string{"testdata/test6"}, module: "test6",
			want: []*Path{
				{Module: "test6", Type: "enumeration: [dark milk]", XPath: "/food/chocolate"},
				{Module: "test6", Type: "empty", XPath: "/food/beer"},
				{Module: "test6", Type: "empty", XPath: "/food/pretzel"},
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
			// spew.Dump(e)
			got := Paths(e, Path{}, []*Path{})
			if !reflect.DeepEqual(tc.want, got) {
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
