Exporting YANG schema paths (aka Schema node identifiers) is the prime objective of `yangpath`. Here we explain how the `export` command works and demonstrate its options.

## Required parameters
Right out of the box `yangpath` is ready to export the paths. The only parameters a user needs to provide are:

* path to a YANG file from which a user needs to export paths
* path to a directory with YANG files which are imported by a target module

#### Module to extract paths from
`yangpath` expects to receive a path to the YANG module from which it needs to extract paths. This is done with the `-m | --module` flag that takes a relative or absolute path to a file.

For example, if we cloned [OpenConfig YANG repository](https://github.com/openconfig/public) and would like to export the paths from its [interfaces module](https://github.com/openconfig/public/blob/master/release/models/interfaces/openconfig-interfaces.yang) we would specify the following path:
```bash
# current dir ~/openconfig/public
yangpath export --module release/models/interfaces/openconfig-interfaces.yang
```

#### Directory with YANG files
By using `-y | --yang-dir` flag a user specifies the path to the directory with YANG files which the target module imports. To add multiple directories use the flag several times:
```
yangpath export -y ~/dir1 -y ~/dir2 -y ~/dir3 -m ~/my_module.yang
```

It is not required (though advised for performance reasons) to specify the exact directory with the required modules, its allowed to specify the directory which nests the target directories.

Consider the following files hierarchy where the imported modules reside in the directories `dir1-3`:

```
.
├── parent-dir
    ├── dir1
    ├── dir2
    └── dir3
```
Its possible to tell `yangpath` to read the parent directory instead of specifying each of the directories separately:
```
yangpath export -y ~/parent-dir -m ~/my_module.yang
```

If the `--yang-dir` flag is not specified, it defaults to the current directory, which means that current working directory and all its subdirectories will be read. This is perfectly fine:
```bash
# ~/projects/openconfig/public
# yangpath will read every subdirectory of openconfig/public repo
# and will find all dependencies
yangpath export --module release/models/interfaces/openconfig-interfaces.yang
```
??? note "More details on directories with YANG files"
    When `yangpath` compiles the YANG module it is about to export paths from, it also needs to compile the modules that the target module imports.

    Consider the following example of a module that we would like export paths from:
    ```
    module openconfig-interfaces {

    yang-version "1";

    // namespace
    namespace "http://openconfig.net/yang/interfaces";

    prefix "oc-if";

    // import some basic types
    import ietf-interfaces { prefix ietf-if; }
    import openconfig-yang-types { prefix oc-yang; }
    import openconfig-types { prefix oc-types; }
    import openconfig-extensions { prefix oc-ext; }
    ```
    If these imported modules are not in the same directory where the target module is, a user needs to provide a path (or paths) to the directories with these imported modules.

## Default behavior
With just the above mentioned flags set, the exported paths will be printed to `stdout` with keys and types fields highlighted:

```text
❯ yangpath export --module release/models/interfaces/openconfig-interfaces.yang

[rw]  /interfaces/interface[name=*]/config/description  string
[rw]  /interfaces/interface[name=*]/config/enabled  boolean
[rw]  /interfaces/interface[name=*]/config/loopback-mode  boolean
[rw]  /interfaces/interface[name=*]/config/mtu  uint16
[rw]  /interfaces/interface[name=*]/config/name  string
[rw]  /interfaces/interface[name=*]/config/type  identityref->ietf-if:interface-type
```

Paths appear each one on a single line and consist of the following elements:

* **node state:** reflects the configuration state of a given node as per [4.2.3 of RFC 7950](https://tools.ietf.org/html/rfc7950#section-4.2.3).
    `[rw]` corresponds for the nodes for which YANG statement `config false` was *not* set  
    `[ro]` corresponds for the nodes for which YANG statement `config false` was set
* **path:** the path itself in a XPATH style with the keys preserved, starting from the root of the module
* **type:** YANG type associated with the path in the "detailed" form

## Configuration options
The `export` command is flexible, it employs some sensible defaults, allowing a user to tailor the output to their needs.

All of the configuration options are presented in the embedded help `yangpath export --help`, here we explain how these option work.

#### Node state
Node state, which is enabled by default and shows if the leaf is a configurable or not, can be turned down with the `--node-state=false`.

#### Path style
[Two path styles](about-paths.md#path-styles) are supported by `yangpath` - XPATH and RESTCONF - the selection is enabled by the `[-s | --style]` flag.

#### Type style
`yangpath` augments paths with [type information](about-paths.md#types). The `--types` flag configures the way types are displayed:

* `no`: types are not displayed
* `yes`: only type names are displayed
* `detailed` (default): both type names and enclosed values are displayed as explained [here](about-paths.md#types).

#### Color highlighting
Path keys are of [prime importance](about-paths.md#keys-in-paths) in `yangpath` export output.  
To articulate the keys in the schema path we made them highlighted with the ANSI colors. At the same time, the type information is rendered _faded_ so that each element can stay visually separated even if the path is quite long.

=== "light theme"
    ![color_light](https://gitlab.com/rdodin/pics/-/wikis/uploads/1a3cf2312b852c392501750955df5d16/image.png)
=== "dark theme"
    ![color_dark](https://gitlab.com/rdodin/pics/-/wikis/uploads/8b2ba4041c507a0606a3d88676712266/image.png)

If colors are not up to your liking, you can always turn them off by adding a flag `--no-color`.

#### Node filter
It is possible to display only state or configuration nodes by using `--only-nodes` flag that takes one of these values:

* `all` (default): both configuration and read only nodes are displayed
* `state`: only read-only nodes are displayed
* `config`: only configuration nodes are displayed

#### Module name
Although module name is likely known to a user, its possible to display the module name along each path by using `--with-module yes` flag.

#### Format
By default `yangpath` outputs the paths in text format to stdout, but it can also generate an HTML output which opens the door to some pretty cool usecases which we discuss on the [Path Browser](html-template.md) page.