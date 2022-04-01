#!/bin/bash

VERSION="$1"

if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo "Invalid version number"
	exit 1
fi

header_match="	return CreateWithIdentity\(apiKey, \"modernmt-go\", "
header_ver="	return CreateWithIdentity\(apiKey, \"modernmt-go\", \"${VERSION}\"\)"
sed -i -E "/$header_match/s/.*/$header_ver/" modernmt.go
