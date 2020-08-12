It doesn't matter what color you banners are, `yangpath` is vendor-agnostic and YANG centric. Simply put, it digests YANG files from multiple vendors with no hiccups.

Here we demonstrate how `yangpath` can be used with models from different vendors and standard organizations.

## OpenConfig
Likely the most popular vendor-agnostic provider of YANG modules is OpenConfig.

Source: [openconfig/public](https://github.com/openconfig/public)

```bash
# assuming cur working dir is the root of the repo
yangpath export -m release/models/interfaces/openconfig-interfaces.yang
```

## IETF
The foundational IETF models

Source: [YangModels/yang](https://github.com/YangModels/yang)

```bash
# assuming cur working dir is the root of the repo
yangpath export -m  standard/ietf/RFC/ietf-interfaces@2018-02-20.yang
```

## Nokia
By the way, the paths extracted with `yangpath` are published at [hellt/nokia-yangtree](https://github.com/hellt/nokia-yangtree).

Source: [nokia/7x50_YangModels](https://github.com/nokia/7x50_YangModels)

```bash
# assuming cur working dir is the root of the repo
yangpath export -y YANG -m YANG/nokia-combined/nokia-conf-combined.yang
```

## Arista
Arista uses a subset of OpenConfig modules and does not provide IETF modules inside their repo. So make sure you have IETF models somewhere where you can reference it.

Source: [aristanetworks/yang](https://github.com/aristanetworks/yang)

```bash
# assuming cur working dir is the root of the repo
# notice the second import where we specify path to the IETF models from OC repo
yangpath export -y EOS-4.23.2F/openconfig/public/release/models \
                -y ~/projects/openconfig/public/third_party/ietf/ \
                -m EOS-4.23.2F/openconfig/public/release/models/interfaces/openconfig-interfaces.yang
```

## Cisco

Source: [YangModels/yang](https://github.com/YangModels/yang)

```bash
# assuming cur working dir is the root of the repo
yangpath export -y standard/ietf/ -m vendor/cisco/xr/711/Cisco-IOS-XR-mpls-ldp-cfg.yang
```

## Juniper

Source: [Juniper/yang](https://github.com/Juniper/yang)

Unfortunately, the underlying library that `yangpath` uses, have troubles reading the YANG files directories that Juniper has in their repo.

