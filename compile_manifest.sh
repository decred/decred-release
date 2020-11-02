#!/bin/sh

set -e
set -o pipefail

# download a file over HTTPS and dump to standard output
# arg1: the url
dl_stdout() {
	local _url=$1
	case $(uname) in
	"OpenBSD")
		ftp -VMo- ${_url}
		;;
	*)
		curl -Lf ${_url}
		;;
	esac
}

# compute SHA256 hash of standard input
# returns: the SHA256 hash
sha256_sum() {
	openssl sha256 | awk '{print $2}'
}

# download file from a url and compute its SHA256 hash
# arg1: the url
# returns: the SHA256 hash
dl_hash() {
	local _url=$1
	dl_stdout ${_url} | sha256_sum
}

[ $# -eq 1 -o $# -eq 2 ] || {
	echo "usage: $0 output [urlfile]" 2>&1
	exit 2
}
MANIFEST=$1
URLFILE=${2:-dcrinstall_manifest_urls.txt}

[ -e ${MANIFEST} ] && rm ${MANIFEST}
while read _url; do
	_hash=$(dl_hash ${_url})
	echo "${_hash}  ${_url}" >> ${MANIFEST}
done <${URLFILE}

echo "Manifest:"
cat ${MANIFEST}
echo "Manifest hash: $(sha256_sum<${MANIFEST})"
