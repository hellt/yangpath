`yangpath` is a single static binary built for the Linux, Mac OS, Windows platforms and distributed via [Github releases](https://github.com/hellt/yangpath/releases).

### Linux/Mac OS
To download & install the latest release the following automated [installation script](https://github.com/hellt/yangpath/blob/master/install.sh) can be used:

```bash
sudo curl -sL https://github.com/hellt/yangpath/raw/master/install.sh | sudo bash
```

As a result, the latest `yangpath` version will be installed in the `/usr/local/bin` directory and the version information will be printed out.
```text
Preparing to install yangpath 0.0.1 into /usr/local/bin
yangpath installed into /usr/local/bin/yangpath
version : 0.0.1
 commit : bdaa6ab
   date : 2020-08-11T20:27:24Z
 source : https://github.com/hellt/yangpath
   docs : https://yangpath.netdevops.me
```

To upgrade run the installation script once again, it will perform the upgrade if a newer version is available.

### Windows
It is highly recommended to use [WSL](https://en.wikipedia.org/wiki/Windows_Subsystem_for_Linux) on Windows, but if its not possible, use [releases page](https://github.com/hellt/yangpath/releases) to download the windows executable file.

### Package managers
Links to the Debian and RPM packages are available in the [releases](https://github.com/hellt/yangpath/releases) section. For example, to install `yangpath v0.1.0` with `yum` issue the following:
```
yum install https://github.com/hellt/yangpath/releases/download/v0.1.0/yangpath_0.1.0-test_linux_x86_64.rpm
```

### Docker
The `yangpath` Docker image is available for each release and is tagged accordignly.  
You can pull the latest or a specific version:

```bash
# get the latest version
docker pull hellt/yangpath

# get a specific release
docker pull hellt/yangpath:0.1.0
```