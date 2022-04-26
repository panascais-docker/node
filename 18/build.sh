#!/usr/bin/sh

set -e
set -u
set -o pipefail
set -x

addgroup -g 1000 node;
adduser -u 1000 -G node -s /bin/sh -D node;

apk add --no-cache libstdc++;
apk add --no-cache --virtual .build-deps curl binutils;

ARCH=;
alpineArch="$(apk --print-arch)";

case "${alpineArch##*-}" in
    amd64) ARCH=x64 ;;
    x86_64) ARCH=x64 ;;
    i386:x86-64) ARCH=x64 ;;
    aarch64) ARCH=arm64 ;;
    arm64) ARCH=arm64 ;;
    armv7l) ARCH=armv7 ;;
    *) echo "Unsupported Alpine architecture: $alpineArch"; exit 1 ;;
esac

echo "ARCH=$ARCH";

curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION.tar.xz";

tar -xf "node-v$NODE_VERSION.tar.xz";

cd "node-v$NODE_VERSION";

mkdir -p /usr/local/sbin;

ln -s /usr/bin/python3 /usr/local/sbin/python3.9;
ln -s /usr/bin/python3 /usr/local/sbin/python;

./configure;

make install -j$(getconf _NPROCESSORS_ONLN) V= ARCH=$ARCH DISTTYPE="release" VARIATION="musl";

cd ..;

strip /usr/local/bin/node;
ln -s /usr/local/bin/node /usr/local/bin/nodejs;

curl -fsSL https://get.pnpm.io/v6.js | node - add --global pnpm;
node --version;
pnpm --version;
which pnpm;
npm --version
