#!/bin/bash

set -euo pipefail

VERSIONS="7.1.4 8.0.10 9.0.0-M7"
DISTS="webprofile plus plume"

mkdir tmp

for VERSION in $VERSIONS; do
	for DIST in $DISTS; do
    		curl -L --silent -o tmp/apache-tomee-${VERSION}-${DIST}.tar.gz "https://downloads.apache.org/tomee/tomee-${VERSION}/apache-tomee-${VERSION}-${DIST}.tar.gz"
		echo "$VERSION / $DIST"
		sha256sum tmp/apache-tomee-${VERSION}-${DIST}.tar.gz
	done
done

rm -fr tmp
