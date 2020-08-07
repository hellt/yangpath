package path

import (
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
)

func TestPaths(t *testing.T) {
	type want struct {
		Module   string
		TypeName string
		XPath    string
	}
	tests := map[string]struct {
		dirs   []string
		module string
		want   []want
	}{
		"test1": {
			dirs: []string{"testdata/test1"}, module: "test1",
			want: []want{
				{
					Module:   "test1",
					TypeName: "string",
					XPath:    "/c1/l1[key1=*]/key1",
				},
				{
					Module:   "test1",
					TypeName: "string",
					XPath:    "/c1/l1[key1=*]/leaf2",
				},
			}},
		"test2": {
			dirs: []string{"testdata/test2"}, module: "test2",
			want: []want{
				{
					Module:   "test2",
					TypeName: "string",
					XPath:    "/c1/l1[key1=*]/key1",
				},
				{
					Module:   "test2",
					TypeName: "string",
					XPath:    "/c1/l1[key1=*]/leaf2",
				},
			}},
		"test3": {
			dirs: []string{"testdata/test3"}, module: "test3",
			want: []want{
				{
					Module:   "test3",
					TypeName: "string",
					XPath:    "/c1/l1[key1=*][key2=*]/key1",
				},
				{Module: "test3",
					TypeName: "age",
					XPath:    "/c1/l1[key1=*][key2=*]/key2",
				},
				{
					Module:   "test3",
					TypeName: "int64",
					XPath:    "/c1/l1[key1=*][key2=*]/leaf1",
				},
			}},
		"test4": {
			dirs: []string{"testdata/test4"}, module: "test4",
			want: []want{
				{
					Module:   "test4",
					TypeName: "identityref",
					XPath:    "/c1/leaf1"},
				{
					Module:   "test4",
					TypeName: "identityref",
					XPath:    "/c1/leaf2",
				},
			}},
		"test5": {
			dirs: []string{"testdata/test5"}, module: "test5",
			want: []want{
				{
					Module:   "test5",
					TypeName: "string",
					XPath:    "/c1/leaf1",
				},
				{
					Module:   "test5",
					TypeName: "leafref",
					XPath:    "/c1/leaf2"},
			}},
		"test6": {
			dirs: []string{"testdata/test6", "testdata/test3"}, module: "test6",
			want: []want{
				{
					Module:   "test6",
					TypeName: "enumeration",
					XPath:    "/food/chocolate",
				},
				{
					Module:   "test6",
					TypeName: "test3:age",
					XPath:    "/food/testage",
				},
				{
					Module:   "test6",
					TypeName: "empty",
					XPath:    "/food/beer",
				},
				{
					Module:   "test6",
					TypeName: "empty",
					XPath:    "/food/pretzel",
				},
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

			for i, v := range tc.want {
				if v.Module != got[i].Module {
					t.Fatalf("Module wanted %s got %s\n", v.Module, got[i].Module)
				}
				if v.XPath != got[i].XPath {
					t.Fatalf("XPATH wanted %s got %s\n", v.XPath, got[i].XPath)
				}
				if v.TypeName != got[i].Type.Name {
					t.Fatalf("Type wanted %s got %s\n", v.TypeName, got[i].Type.Name)
				}
			}

		})
	}

}
