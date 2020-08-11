#!/usr/bin/env bash

# The install script is based off of the Apache 2.0 script from Helm,
# https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3

: ${BINARY_NAME:="yangpath"}
: ${USE_SUDO:="true"}
: ${VERIFY_CHECKSUM:="true"}
: ${YANGPATH_INSTALL_DIR:="/usr/local/bin"}

# initArch discovers the architecture for this system.
initArch() {
    ARCH=$(uname -m)
    # case $ARCH in
    # armv5*) ARCH="armv5" ;;
    # armv6*) ARCH="armv6" ;;
    # armv7*) ARCH="arm" ;;
    # aarch64) ARCH="arm64" ;;
    # x86) ARCH="386" ;;
    # x86_64) ARCH="amd64" ;;
    # i686) ARCH="386" ;;
    # i386) ARCH="386" ;;
    # esac
}

# initOS discovers the operating system for this system.
initOS() {
    OS=$(echo $(uname) | tr '[:upper:]' '[:lower:]')

    case "$OS" in
    # Minimalist GNU for Windows
    mingw*) OS='windows' ;;
    esac
}

# runs the given command as root (detects if we are root already)
runAsRoot() {
    local CMD="$*"

    if [ $EUID -ne 0 -a $USE_SUDO = "true" ]; then
        CMD="sudo $CMD"
    fi

    $CMD
}

# verifySupported checks that the os/arch combination is supported for
# binary builds.
verifySupported() {
    local supported="darwin-i386\ndarwin-x86_64\nlinux-i386\nlinux-x86_64\nlinux-armv7\nlinux-aarch64\nwindows-i386\nwindows-x86_64"
    if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
        echo "No prebuilt binary for ${OS}-${ARCH}."
        echo "To build from source, go to https://github.com/hellt/yangpath"
        exit 1
    fi

    if ! type "curl" &>/dev/null && ! type "wget" &>/dev/null; then
        echo "Either curl or wget is required"
        exit 1
    fi
}

# verifyOpenssl checks if openssl is installed to perform checksum operation
verifyOpenssl() {
    if [ $VERIFY_CHECKSUM == "true" ]; then
        if ! type "openssl" &>/dev/null; then
            echo "openssl is not found. It is used to verify checksum of the downloaded archive.\nEither install openssl or provide '--skip-checksum' flag to the installer"
            exit 1
        fi
    fi
}

# setDesiredVersion sets the desired version either to an explicit version provided by a user
# or to the latest release available on github releases
setDesiredVersion() {
    if [ "x$DESIRED_VERSION" == "x" ]; then
        # when desired version is not provided
        # get latest tag from the gh releases
        if type "curl" &>/dev/null; then
            local latest_release_url=$(curl -s https://github.com/hellt/yangpath/releases/latest | cut -d '"' -f 2)
            TAG=$(echo $latest_release_url | cut -d '"' -f 2 | awk -F "/" '{print $NF}')
            # tag with stripped `v` prefix
            TAG_WO_VER=$(echo "${TAG}" | cut -c 2-)
        elif type "wget" &>/dev/null; then
            # get latest release info and get 5th line out of the response to get the URL
            local latest_release_url=$(wget -q https://api.github.com/repos/hellt/yangpath/releases/latest -O- | sed '5q;d' | cut -d '"' -f 4)
            TAG=$(echo $latest_release_url | cut -d '"' -f 2 | awk -F "/" '{print $NF}')
            TAG_WO_VER=$(echo "${TAG}" | cut -c 2-)
        fi
    else
        TAG=$DESIRED_VERSION
        TAG_WO_VER=$(echo "${TAG}" | cut -c 2-)
    fi
}

# checkyangpathInstalledVersion checks which version of yangpath is installed and
# if it needs to be changed.
checkyangpathInstalledVersion() {
    if [[ -f "${YANGPATH_INSTALL_DIR}/${BINARY_NAME}" ]]; then
        local version=$("${YANGPATH_INSTALL_DIR}/${BINARY_NAME}" version | head -1 | awk '{print $NF}')
        if [[ "v$version" == "$TAG" ]]; then
            echo "yangpath is already at ${DESIRED_VERSION:-latest ($version)}" version
            return 0
        else
            echo "yangpath ${TAG_WO_VER} is available. Changing from version ${version}."
            return 1
        fi
    else
        return 1
    fi
}

# downloadFile downloads the latest binary package and also the checksum
# for that binary.
downloadFile() {
    YANGPATH_DIST="${BINARY_NAME}_${TAG_WO_VER}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/hellt/yangpath/releases/download/${TAG}/${YANGPATH_DIST}"
    CHECKSUM_URL="https://github.com/hellt/yangpath/releases/download/${TAG}/checksums.txt"
    YANGPATH_TMP_ROOT="$(mktemp -d)"
    YANGPATH_TMP_FILE="$YANGPATH_TMP_ROOT/$YANGPATH_DIST"
    YANGPATH_SUM_FILE="$YANGPATH_TMP_ROOT/checksums.txt"
    echo "Downloading $DOWNLOAD_URL"
    if type "curl" &>/dev/null; then
        curl -SsL "$CHECKSUM_URL" -o "$YANGPATH_SUM_FILE"
    elif type "wget" &>/dev/null; then
        wget -q -O "$YANGPATH_SUM_FILE" "$CHECKSUM_URL"
    fi
    if type "curl" &>/dev/null; then
        curl -SsL "$DOWNLOAD_URL" -o "$YANGPATH_TMP_FILE"
    elif type "wget" &>/dev/null; then
        wget -q -O "$YANGPATH_TMP_FILE" "$DOWNLOAD_URL"
    fi
}

# installFile verifies the SHA256 for the file, then unpacks and
# installs it.
installFile() {
    YANGPATH_TMP="$YANGPATH_TMP_ROOT/$BINARY_NAME"
    if [ $VERIFY_CHECKSUM == "true" ]; then
        local sum=$(openssl sha1 -sha256 ${YANGPATH_TMP_FILE} | awk '{print $2}')
        local expected_sum=$(cat ${YANGPATH_SUM_FILE} | grep -i $YANGPATH_DIST | awk '{print $1}')
        if [ "$sum" != "$expected_sum" ]; then
            echo "SHA sum of ${YANGPATH_TMP_FILE} does not match. Aborting."
            exit 1
        fi
    fi

    mkdir -p "$YANGPATH_TMP"
    tar xf "$YANGPATH_TMP_FILE" -C "$YANGPATH_TMP"
    YANGPATH_TMP_BIN="$YANGPATH_TMP/yangpath"
    echo "Preparing to install $BINARY_NAME ${TAG_WO_VER} into ${YANGPATH_INSTALL_DIR}"
    runAsRoot cp "$YANGPATH_TMP_BIN" "$YANGPATH_INSTALL_DIR/$BINARY_NAME"
    echo "$BINARY_NAME installed into $YANGPATH_INSTALL_DIR/$BINARY_NAME"
}

# fail_trap is executed if an error occurs.
fail_trap() {
    result=$?
    if [ "$result" != "0" ]; then
        if [[ -n "$INPUT_ARGUMENTS" ]]; then
            echo "Failed to install $BINARY_NAME with the arguments provided: $INPUT_ARGUMENTS"
            help
        else
            echo "Failed to install $BINARY_NAME"
        fi
        echo -e "\tFor support, go to https://github.com/hellt/yangpath"
    fi
    cleanup
    exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
    set +e
    $YANGPATH_INSTALL_DIR/$BINARY_NAME version
    if [ "$?" = "1" ]; then
        echo "$BINARY_NAME not found. Is $YANGPATH_INSTALL_DIR in your "'$PATH?'
        exit 1
    fi
    set -e
}

# help provides possible cli installation arguments
help() {
    echo "Accepted cli arguments are:"
    echo -e "\t[--help|-h ] ->> prints this help"
    echo -e "\t[--version|-v <desired_version>] . When not defined it fetches the latest release from GitHub"
    echo -e "\te.g. --version v0.1.1"
    echo -e "\t[--no-sudo]  ->> install without sudo"
    echo -e "\t[--skip-checksum]  ->> disable automatic checksum verification"
}

cleanup() {
    if [[ -d "${YANGPATH_TMP_ROOT:-}" ]]; then
        rm -rf "$YANGPATH_TMP_ROOT"
    fi
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

# Parsing input arguments (if any)
export INPUT_ARGUMENTS="${@}"
set -u
while [[ $# -gt 0 ]]; do
    case $1 in
    '--version' | -v)
        shift
        if [[ $# -ne 0 ]]; then
            export DESIRED_VERSION="v${1}"
        else
            echo -e "Please provide the desired version. e.g. --version 0.1.1"
            exit 0
        fi
        ;;
    '--no-sudo')
        USE_SUDO="false"
        ;;
    '--skip-checksum')
        VERIFY_CHECKSUM="false"
        ;;
    '--help' | -h)
        help
        exit 0
        ;;
    *)
        exit 1
        ;;
    esac
    shift
done
set +u

initArch
initOS
verifySupported
setDesiredVersion
if ! checkyangpathInstalledVersion; then
    verifyOpenssl
    downloadFile
    installFile
    testVersion
    cleanup
fi
