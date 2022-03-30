#!/usr/bin/env bash
set -e
mkdir -p rubygem/build/linux-amd64
mkdir -p rubygem/build/darwin-all
mkdir -p rubygem/build/freebsd-amd64
cp dist/ejson2env_linux_amd64/ejson2env rubygem/build/linux-amd64/ejson2env
cp dist/ejson2env_darwin_all/ejson2env rubygem/build/darwin-all/ejson2env
cp dist/ejson2env_freebsd_amd64/ejson2env rubygem/build/freebsd-amd64/ejson2env
cp LICENSE.txt rubygem/LICENSE.txt
mkdir -p rubygem/man
bundle install
bundle exec ronn -r --pipe man/ejson2env.1.ronn | gzip > rubygem/man/ejson2env.1.gz
mkdir -p rubygem/lib/ejson2env
echo -e "module EJSON2ENV\n  VERSION = \"${VERSION}\"\nend" >rubygem/lib/ejson2env/version.rb
cd rubygem
gem build ejson2env.gemspec

# get release_id for use in actions/upload-release-asset step
release_id=$(gh api repos/${GITHUB_REPOSITORY}/releases/tags/${GITHUB_REF_NAME} | jq -r .id)
echo ::set-output name=release_id::$release_id
