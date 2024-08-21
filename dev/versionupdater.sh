#!/usr/bin/env bash

version=$(cat VERSION.txt)
IFS='.' read -r -a version_parts <<< "$version"
major=${version_parts[0]}
minor=${version_parts[1]}
patch=${version_parts[2]}

new_patch=$((patch + 1))
new_version="$major.$minor.$new_patch"

echo "$new_version" > VERSION.txt
