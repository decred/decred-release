#!/bin/sh

set -e
set -o pipefail

usage() {
	echo "usage: $0 output manifest1 [manifest2...]" 2>&1
	exit 2
}

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

find_url() {
	NAME="$1" perl -lane 'print if /${ENV{NAME}}$/' <${URLFILE}
}

[ $# -gt 1 ] || usage
OUTPUT=$1
URLFILE=dcrinstall_manifest_urls.txt
shift
FILES="$@"

# sort filenames
export FILES
FILES=$(perl -e 'print join("\n", sort(split(/\s+/, $ENV{FILES})))')

[ -e ${OUTPUT} ] && rm ${OUTPUT}
for _file in ${FILES}; do
	_what=$(basename ${_file})
	_url=$(find_url ${_what})
	[ -z "${_url}" ] && {
		echo "${_what} not found in url file" 2>&1
		exit 1
	}
	_hash=$(sha256_sum <${_file})
	[ ${_hash} = $(dl_hash ${_url}) ] || {
		echo "download of ${_what} does not match local file" 2>&1
		exit 1
	}
	echo "${_hash}  ${_url}" >> ${OUTPUT}
done

echo "Manifest:"
cat ${OUTPUT}
echo "Manifest hash: $(sha256_sum<${OUTPUT})"
