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
	"log"
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
)

func TestPaths(t *testing.T) {
	type want struct {
		Module   string
		TypeName string
		XPath    string
		SType    string
	}
	tests := map[string]struct {
		dirs   []string
		module string
		want   []want
	}{
		"test1": {
			dirs: []string{"testdata/test1"}, module: "testdata/test1/test1.yang",
			want: []want{
				{
					Module:   "test1",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/l1[key1=*]/key1",
				},
				{
					Module:   "test1",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/l1[key1=*]/leaf2",
				},
			}},
		"test2": {
			dirs: []string{"testdata/test2"}, module: "testdata/test2/test2.yang",
			want: []want{
				{
					Module:   "test2",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/l1[key1=*]/key1",
				},
				{
					Module:   "test2",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/l1[key1=*]/leaf2",
				},
			}},
		"test3": {
			dirs: []string{"testdata/test3"}, module: "testdata/test3/test3.yang",
			want: []want{
				{
					Module:   "test3",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/l1[key1=*][key2=*]/key1",
				},
				{Module: "test3",
					TypeName: "age",
					SType:    "age",
					XPath:    "/c1/l1[key1=*][key2=*]/key2",
				},
				{
					Module:   "test3",
					TypeName: "int64",
					SType:    "int64",
					XPath:    "/c1/l1[key1=*][key2=*]/leaf1",
				},
			}},
		"test4": {
			dirs: []string{"testdata/test4"}, module: "testdata/test4/test4.yang",
			want: []want{
				{
					Module:   "test4",
					TypeName: "identityref",
					SType:    "identityref->test4:IDENTITY2",
					XPath:    "/c1/leaf1"},
				{
					Module:   "test4",
					TypeName: "identityref",
					SType:    "identityref->IDENTITY1",
					XPath:    "/c1/leaf2",
				},
			}},
		"test5": {
			dirs: []string{"testdata/test5"}, module: "testdata/test5/test5.yang",
			want: []want{
				{
					Module:   "test5",
					TypeName: "string",
					SType:    "string",
					XPath:    "/c1/leaf1",
				},
				{
					Module:   "test5",
					TypeName: "leafref",
					SType:    "leafref->../leaf1",
					XPath:    "/c1/leaf2"},
			}},
		"test6": {
			dirs: []string{"testdata/test6", "testdata/test3"}, module: "testdata/test6/test6.yang",
			want: []want{
				{
					Module:   "test6",
					TypeName: "enumeration",
					SType:    "enumeration[\"dark\" \"milk\"]",
					XPath:    "/food/chocolate",
				},
				{
					Module:   "test6",
					TypeName: "test3:age",
					SType:    "test3:age",
					XPath:    "/food/testage",
				},
				{
					Module:   "test6",
					TypeName: "empty",
					SType:    "empty",
					XPath:    "/food/beer",
				},
				{
					Module:   "test6",
					TypeName: "empty",
					SType:    "empty",
					XPath:    "/food/pretzel",
				},
			}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			modName, err := GetModuleName(tc.module)
			if err != nil {
				log.Fatal(err)
			}
			if err = AddYANGDirs(tc.dirs); err != nil {
				log.Fatal(err)
			}

			e, errs := yang.GetModule(modName)
			for _, err := range errs {
				log.Fatalf("%v\n", err)
			}
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
				if v.SType != got[i].SType {
					t.Fatalf("Type wanted %s got %s\n", v.TypeName, got[i].Type.Name)
				}
			}

		})
	}

}
