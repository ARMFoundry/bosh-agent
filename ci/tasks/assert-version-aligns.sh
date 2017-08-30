#!/usr/bin/env bash

set -e

export BASE=$(pwd)

semver=`cat ${BASE}/version-semver/number`

pushd "${BASE}/bosh-agent"
  git_branch=`git branch --list -r --contains HEAD | cut -d'/' -f2`
popd

echo "detected bosh-agent will build from branch $git_branch ..."

if [ "$git_branch" = "master" ]; then
  version_must_match='^[0-9]+\.[0-9]+\.0$'
else
  version_must_match="^${git_branch//x/[0-9]+}$"
  version_must_match="${version_must_match//./\.}"
fi

echo "will only continue if version to promote matches $version_must_match ..."

if [[ $semver =~ $version_must_match ]]; then
  echo "version $semver is appropriate for branch $git_branch -- promote will continue"
  exit 0
fi

echo "version $semver DOES NOT ALIGN with branch $git_branch -- promote step canceled!"

exit 1

