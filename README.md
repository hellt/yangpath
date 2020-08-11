<p align=center><img src=https://gitlab.com/rdodin/pics/-/wikis/uploads/f7725125d0f78721d7bb68aceca26456/yangpath_logo_paths.svg?sanitize=true/></p>

[![github release](https://img.shields.io/github/release/hellt/yangpath.svg?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://github.com/hellt/yangpath/releases/)
[![Github all releases](https://img.shields.io/github/downloads/hellt/yangpath/total.svg?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://github.com/hellt/yangpath/releases/)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-blue?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://goreportcard.com/report/github.com/hellt/yangpath)
[![Doc](https://img.shields.io/badge/Docs-yangpath.netdevops.me-blue?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://yangpath.netdevops.me)
[![build](https://img.shields.io/github/workflow/status/hellt/yangpath/Test/master?style=flat-square&labelColor=bec8d2)](https://github.com/hellt/yangpath/releases/)

---
`yangpath` is an XPATH/RESTCONF-styled schema paths exporter with superpowers.

The exported paths can be immediately used in your NETCONF/RESTCONF or gNMI applications.

Documentation available at [https://yangpath.netdevops.me](https://yangpath.netdevops.me)

**XPATH paths**
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

**RESTCONF paths**
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
    The exported paths have the [list keys present](https://yangpath.netdevops.me/about-paths#keys-in-paths).  
    Knowing the key names makes it very easy to create XPATH/RESTCONF filters targeting a particular node.  
* **Readily available for gNMI**  
    The exported paths are fully compatible with the gNMI paths, thanks to the keys being present and set to the wildcard `*` value.
* **RESTCONF-ready**  
    With a matter of a single flag value switch `yangpath` will export the paths in a [RESTCONF style](https://yangpath.netdevops.me/about-paths#path-styles). Paste them in Postman and you're good to go!
* **Type information**  
    A unique `yangpath` feature is its ability to provide [the type of a given path](https://yangpath.netdevops.me/about-paths#types). Types give additional context when you retrieve the data, but they are of utter importance for edit the configuration operations.
* **Fast**  
    Path export with `yangpath` is quite fast, working with massive models is no longer a problem!
* **User friendly**  
    As always, we strive to publish the tools which spark joy, therefore pre-built images with an effortless [installation](https://yangpath.netdevops.me/install) and a beautiful and extensive documentation comes included.

## Quick Start

#### Install
Use the following installation script to install the latest version.
```
sudo curl -sL https://github.com/hellt/yangpath/raw/master/install.sh | sudo bash
```
Alternatively, leverage the system [packages](https://yangpath.netdevops.me/install#package-managers) or [docker images](https://yangpath.netdevops.me/install#docker).

#### Export paths
To [export](https://yangpath.netdevops.me/export) the paths from a given module:
```bash
# assuming cur working dir is the root of openconfig repo
yangpath export -m release/models/interfaces/openconfig-interfaces.yang
```

#### Generate HTML path browser
To create HTML with paths out of template, leverage [templating capabilities](https://yangpath.netdevops.me/html-template) of `yangpath`.
