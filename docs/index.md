<nbsp/>
<p style="text-align:center;">![headline](images/yangpath_logo_paths_m.svg)</p>

[![github release](https://img.shields.io/github/release/hellt/yangpath.svg?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://github.com/hellt/yangpath/releases/)
[![Github all releases](https://img.shields.io/github/downloads/hellt/yangpath/total.svg?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://github.com/hellt/yangpath/releases/)

---

`yangpath` is an XPATH/RESTCONF-styled schema paths exporter with superpowers.

The exported paths can be immediately used in your NETCONF/RESTCONF or gNMI applications.

=== "XPATH paths"
    ```text
    $ yangpath export --yang-dir ~/openconfig/public/ \
                      --module ~/openconfig/public/release/models/interfaces/openconfig-interfaces.yang
    [rw]  /interfaces/interface[name=*]/config/description  string
    [rw]  /interfaces/interface[name=*]/config/enabled  boolean
    [rw]  /interfaces/interface[name=*]/config/loopback-mode  boolean
    [rw]  /interfaces/interface[name=*]/config/mtu  uint16
    [rw]  /interfaces/interface[name=*]/config/name  string
    [rw]  /interfaces/interface[name=*]/config/type  identityref->ietf-if:interface-type
    <SNIPPED>
    ```
=== "RESTCONF paths"
    ```text
    $ yangpath export --yang-dir ~/openconfig/public/ \
                      --module ~/openconfig/public/release/models/interfaces/openconfig-interfaces.yang
                      --style restconf
    [rw]  /interfaces/interface=name/config/description  string
    [rw]  /interfaces/interface=name/config/enabled  boolean
    [rw]  /interfaces/interface=name/config/loopback-mode  boolean
    [rw]  /interfaces/interface=name/config/mtu  uint16
    [rw]  /interfaces/interface=name/config/name  string
    [rw]  /interfaces/interface=name/config/type  identityref->ietf-if:interface-type
    ```

## Features
* **Preserved list keys**  
    The exported paths have the [list keys present](about-paths.md#keys-in-paths).  
    Knowing the key names makes it very easy to create XPATH/RESTCONF filters targeting a particular node.  
* **Readily available for gNMI**  
    The exported paths are fully compatible with the gNMI paths, thanks to the keys being present and set to the wildcard `*` value.
* **RESTCONF-ready**  
    With a matter of a single flag value switch `yangpath` will export the paths in a [RESTCONF style](about-paths.md#path-styles). Paste them in Postman and you're good to go!
* **Type information**  
    A unique `yangpath` feature is its ability to provide [the type of a given path](about-paths.md#types). Types give additional context when you retrieve the data, but they are of utter importance for edit the configuration operations.
* **Fast**  
    Path export with `yangpath` is quite fast, working with massive models is no longer a problem!
* **User friendly**  
    As always, we strive to publish the tools which spark joy, therefore pre-built images with an effortless [installation](#install) and a beautiful and extensive documentation comes included.

## Quick Start

#### Install
Use the following installation script to install the latest version.
```
sudo curl -sL https://github.com/hellt/yangpath/raw/master/install.sh | sudo bash
```
Alternatively, leverage the system [packages](install.md#package-managers) or [docker images](install.md#docker).

#### Export paths
To [export](export.md) the paths from a given module[^1]:
```bash
# assuming cur working dir is the root of openconfig repo
yangpath export -m release/models/interfaces/openconfig-interfaces.yang
```

#### Generate HTML path browser
To create HTML with paths out of template, leverage [templating capabilities](html-template.md) of `yangpath`.

[^1]: in this example paths from the Openconfig interfaces module are exported