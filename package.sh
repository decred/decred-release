#!/bin/sh

# submit a package to be notarized
# returns: notarization uuid
notary_submit() {
	xcrun altool -f release/dcrinstall-${VERSION}/${EXE}.pkg \
		--notarize-app \
		--primary-bundle-id org.decred.dcrinstall.pkg \
		--asc-provider ${IDENTITY} \
		-p @keychain:${KEYCHAIN} 2>&1 \
	| perl -ne 'print if s/^RequestUUID = //'
}

# check notarization status after successful submission
# arg 1: uuid
# returns: altool output
notary_status() {
	local _uuid=$1
	xcrun altool --notarization-info ${_uuid} -p @keychain:${KEYCHAIN} 2>&1
}

# write an install script read from stdin
# arg 1: script name
installscript() {
	local _script=${SCRIPTS}/$1
	cat >${_script}
	chmod 0755 ${_script}
}

[ $(uname) = Darwin ] || {
	echo "$0 must be run from darwin" 2>&1
	exit 1
}
[ $# = 2 ] || {
	echo "usage: $0 version identity" 2>&1
	exit 2
}

VERSION=$1
IDENTITY=$2
KEYCHAIN=${KEYCHAIN:-signer}
DIST=dist/darwin
SCRIPTS=darwin/scripts
EXE=dcrinstall-darwin-amd64-${VERSION}
BUILD=release/dcrinstall-${VERSION}/${EXE}
PREFIX=${PREFIX:-/usr/local}

[ -x ${BUILD} ] || {
	echo "execute './release.sh ${VERSION}' first" 2>&1
	exit 1
}

set -ex
[ -d ${DIST} ] && rm -rf ${DIST}
[ -d ${SCRIPTS} ] && rm -rf ${SCRIPTS}
mkdir -p ${DIST}
mkdir -p ${SCRIPTS}

# prepare directory with package files
install -m 0755 ${BUILD} ${DIST}/dcrinstall
codesign -s ${IDENTITY} --options runtime ${DIST}/dcrinstall
installscript postinstall <<EOF
#!/bin/sh
echo ${PREFIX}/decred > /etc/paths.d/decred
EOF

# generate signed package for notarization
pkgbuild --identifier org.decred.dcrinstall \
	--version ${VERSION} \
	--root ${DIST} \
	--install-location ${PREFIX}/decred \
	--scripts ${SCRIPTS} \
	--sign ${IDENTITY} \
	release/dcrinstall-${VERSION}/${EXE}.pkg

# submit notarization
_uuid=$(notary_submit)

# poll notarization status until no longer in-progress
set +ex
while :; do
	sleep 60
	_date=$(date)
	_output=$(notary_status ${_uuid})
	_status=$(echo "${_output}" | perl -ne 'print if s/^\s*Status: //')
	echo "check at ${_date}: Status: ${_status}"
	case ${_status} in
	"in progress")
		continue
		;;
	"success")
		# move on to stapling
		break
		;;
	"")
		echo "warn: unknown status -- full output:\n${_output}" 2>&1
		continue
		;;
	*)
		echo "${_output}" 2>&1
		exit 1
		;;
	esac
done
set -ex

# staple package with notarization ticket
stapler staple release/dcrinstall-${VERSION}/${EXE}.pkg
