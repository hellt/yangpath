Knowing the schema paths (aka [schema node identifiers](https://tools.ietf.org/html/rfc7950#section-6.5)) are instrumental for any activity involving model-driven management interfaces. Whatever the management task is, a user ends up being in need to either create, read, update or delete data.

The YANG modelled data follows a tree-like hierarchy where each node can be uniquely identified with a schema path. Without having a path to a piece of YANG modelled data its impossible to manipulate it, therefore its mandatory to have one.

Take a look at this examples where different management interfaces use paths to get different YANG modelled data:

=== "gNMI"
    gNMI[^1] can use XPATH-like paths to access the YANG modelled data.
    ```
    gnmic -a 10.1.0.11:57400 -u admin -p admin \
          get --path /state/system/platform
    ```
=== "NETCONF"
    NETCONF uses XPATH filtering or Subtree filtering to access the data. Here is an example of how an XPATH filter can be used within a NETCONF RPC:
    ```
    <rpc message-id="101"
         xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
    <get-config>
        <source>
        <running/>
        </source>
        <!-- get the user named fred -->
        <filter xmlns:t="http://example.com/schema/1.2/config"
                type="xpath"
                select="/top/users/user[name='fred']"/>
        </get-config>
    </rpc>
    ```
=== "RESTCONF"
    RESTCONF embeds the path information in its URI or YANG-PATCH target field. In any case, the path information must be present:
    ```
    GET https://restconf-server:8545/restconf/data/network-device-mgr:network-devices/network-device=192.168.1.11/root/nokia-conf:configure/policy-options/prefix-list=my-prefix-list
    ```

## How to get schema paths?
The schema paths (or YANG paths as we call it here) is not something you can find in a YANG module itself:

```cpp
// tiny YANG module
module test3 {

  yang-version "1";
  namespace "https://hellt/yangpath/test3";

  prefix "test3";

  typedef age {
    type uint16 {
      range 1..100;
    }
  }

  container c1 {
    list l1 {
      key "key1 key2";
      leaf key1 {
        type string;
      }
      leaf key2 {
        type age;
      }
      leaf leaf1 {
        type int64;
      }
    }
  }
}
```
As you see, the path information is not present in the module. But for the given compact and simple module its quite easy to derive the paths by just browsing the contents with a naked eye. We could also you a tree representation of the module to have a better view:

```
pyang -f tree pkg/path/testdata/test3/test3.yang
module: test3
  +--rw c1
     +--rw l1* [key1 key2]
        +--rw key1     string
        +--rw key2     age
        +--rw leaf1?   int64
```

From here, its not that hard to come up with the schema path for these three leafs:
```
/c1/l1/key1
/c1/l1/key2
/c1/l1/leaf1
```

Unfortunately, this approach is not practical when working with the real-life YANG modules which are of hundred lines of code with multiple cross-references and encapsulations. 

`yangpath` mission is to help with this task at hand. It [exports](export.md) the paths from the given YANG module in [XPATH or RESTCONF style](about-paths.md#path-styles).

```
❯ yangpath export -m pkg/path/testdata/test3/test3.yang
[rw]  /c1/l1[key1=*][key2=*]/key1  string
[rw]  /c1/l1[key1=*][key2=*]/key2  age
[rw]  /c1/l1[key1=*][key2=*]/leaf1  int64
```

## Keys in paths
If you noticed, the paths that `yangpath` provided differ from the ones we extracted ourselves by just looking at the module tree representation. The extra piece here is the list keys that are present in the paths.

The reason for the keys to be present is to make paths more universally applicable. When the key information is missing, you loose the granularity of the path. By looking at the path in the `/c1/l1/leaf1` form there is no way to tell if the list `l1` has keys, and if it does, how many and what are their names?

For that reason, `yangpath` adds keys to the lists, making the paths complete.

If the key information is indeed not needed, a user can easily delete the key elements from the path `/c1/l1[key1=*][key2=*]/key1 -> /c1/l1/key1`.

## Path styles
By default `yangpath` exports paths in XPATH style, but it is also possible to display the paths in a RESTCONF style, all it takes is a single flag switch:

```
❯ yangpath export -s restconf -m pkg/path/testdata/test3/test3.yang
[rw]  /c1/l1=key1,key2/key1  string
[rw]  /c1/l1=key1,key2/key2  age
[rw]  /c1/l1=key1,key2/leaf1  int64
```

As with the XPATHs, we try to give you an idea about the keyed lists, by adding key names towards the list elements of the paths. To make this RESTCONF paths to work you need to substitute the key names in the path to the actual values of these keys. In case it is desired to get the list nodes for all keys, just remove the keys from the path: `/c1/l1/key1`. Easy!

## Types
Another addition is the leaf `type` information. Knowing the type of the YANG node that the path is pointing to is very important. It allows a user to know which values are applicable to that particular YANG node.

`yangpath` does some extra job by _expanding_ the type information for the leafs.

If the leaf has the basic YANG type (such as `string`, `int32`, etc) it is displayed as is:
<table>
<th>YANG type</th>
<th>Path type</th>
<tr>
<td>
```
leaf key1 {
    type string;
}
```
</td>
<td>
`string`
</td>
</tr>
</table>

#### Enumeration
If the leaf is of `enumeration` type, the values of enumeration will be displayed:
<table>
<th>YANG type</th>
<th>Path type</th>
<tr>
<td>
```
leaf admin-status {
  type enumeration {
    enum UP {
      description
        "Ready to pass packets.";
    }
    enum DOWN {
      description
        "Not ready to pass packets.";
    }
```
</td>
<td>
`enumeration["DOWN" "UP"]`
</td>
</tr>
</table>

#### Leafref
If the leaf is of `leafref` type, the path the leafref has is displayed:
<table>
<th>YANG type</th>
<th>Path type</th>
<tr>
<td>
```
leaf index {
    type leafref {
      path "../config/index";
    }
    description
      "The index number of the subinterface -- used to address
      the logical interface";
}
```
</td>
<td>
`leafref->../config/index`
</td>
</tr>
</table>

#### Identityref
If the leaf is of `identityref` type, the referenced identity is displayed:
<table>
<th>YANG type</th>
<th>Path type</th>
<tr>
<td>
```
leaf type {
  type identityref {
    base ietf-if:interface-type;
  }
```
</td>
<td>
`identityref->ietf-if:interface-type`
</td>
</tr>
</table>

#### Union
If the leaf is of `union` type, the embedded types are displayed[^2]:
<table>
<th>YANG type</th>
<th>Path type</th>
<tr>
<td>
```
type union {
  type oc-inet:ip-address;
  type string;
}
```
</td>
<td>
`union{oc-inet:ip-address string}`
</td>
</tr>
</table>

[^1]: examples uses [gNMIc CLI client](https://gnmic.kmrd.dev)
[^2]: not implemented for leaflists of union type