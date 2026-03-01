#!/bin/bash

#SETUP
build() {
    cd ..
    go build
    cd e2e
    mv ../buckmate buckmate-executable
}

setup() {
    log_deep "STARTING"
    mkdir buckmate-test-directory
    cd buckmate-test-directory
}

cleanup() {
    cd ..
    rm -rf buckmate-test-directory
}

full_cleanup() {
    rm -rf buckmate-test-directory
    rm -rf buckmate-executable
}

compare_and_clean() {
    if diff -r "$1" "$2" --exclude=".gitignore"; then
        log_deep "SUCCESS"
        cleanup
    else
        log_deep "FAILURE"
        cleanup
        exit 1
    fi
}

log() {
    echo "${FUNCNAME[1]}: $1"
}

log_deep() {
    echo "${FUNCNAME[2]}: $1"
}

#TESTS
local_to_local() {
    setup
    cp -R ../tests/local-to-local/* ./
    ../buckmate-executable --path buckmate apply
    compare_and_clean "../result" "result"
}

s3_to_s3() {
    setup
    cp -R ../tests/s3-to-s3/* ./
    ../buckmate-executable --path preload apply
    ../buckmate-executable --path buckmate apply
    ../buckmate-executable --path result apply
    compare_and_clean "../result" "result/data"
}

local_to_s3() {
    setup
    cp -R ../tests/local-to-s3/* ./
    ../buckmate-executable --path buckmate apply
    ../buckmate-executable --path result apply
    compare_and_clean "../result" "result/data"
}

s3_to_local() {
    setup
    cp -R ../tests/s3-to-local/* ./
    ../buckmate-executable --path preload apply
    ../buckmate-executable --path buckmate apply
    compare_and_clean "../result" "result"
}

dry_local() {
    setup
    cp -R ../tests/local-to-local/* ./
    result=$(../buckmate-executable --path buckmate apply --dry 2>&1)
    tmp_path=$(echo "$result" | awk '{for(i=1;i<=NF;i++) if($i ~ /^\/tmp\//) print $i}')
    compare_and_clean "../result" "$tmp_path"
}

dry_remote() {
    setup
    cp -R ../tests/local-to-s3/* ./
    result=$(../buckmate-executable --path buckmate apply --dry 2>&1)
    tmp_path=$(echo "$result" | awk '{for(i=1;i<=NF;i++) if($i ~ /^\/tmp\//) print $i}')
    compare_and_clean "../result" "$tmp_path"
}

#
cache_control_metadata() {
    setup
    cp -R ../tests/local-to-s3-cache-metadata/* ./
    ../buckmate-executable --path buckmate apply
    aws s3api head-object --bucket BUCKET_2 --key index.html >>tmp1
    aws s3api head-object --bucket BUCKET_2 --key common-file.json >>tmp2
    tmp1CacheControl=$(cat tmp1 | jq '.CacheControl')
    tmp1MetadataKey=$(cat tmp1 | jq '.Metadata."some-metadata-key"')
    tmp2CacheControl=$(cat tmp2 | jq '.CacheControl')
    tmp2MetadataKey=$(cat tmp2 | jq '.Metadata."some-metadata-key"')
    if [ "$tmp1CacheControl" != "\"no-cache\"" ]; then
        log "Wrong cache control on index.html"
        exit 1
    fi
    if [ "$tmp1MetadataKey" != "\"some-metadata-value\"" ]; then
        log "Wrong metadata on index.html"
        exit 1
    fi
    if [ "$tmp2CacheControl" != "\"no-store\"" ]; then
        log "Wrong cache control on common-file.json"
        exit 1
    fi
    if [ "$tmp2MetadataKey" != "null" ]; then
        log "Wrong metadata on common-file.json"
        exit 1
    fi
    cleanup
    log "SUCCESS"
}

keep_previous_remote() {
    setup
    cp -R ../tests/local-to-s3-keep-previous/* ./
    ../buckmate-executable --path preload apply
    ../buckmate-executable --path buckmate apply
    mkdir result
    aws s3 cp s3://BUCKET_2 result --recursive
    compare_and_clean "../result-keep-previous" "result"
}

keep_previous_local() {
    setup
    cp -R ../tests/local-to-local-keep-previous/* ./
    ../buckmate-executable --path buckmate apply
    compare_and_clean "../result-keep-previous" "result"
}

local_to_local_dev() {
    setup
    cp -R ../tests/local-to-local/* ./
    ../buckmate-executable --path buckmate --env dev apply
    compare_and_clean "../result-dev" "result"
}

local_to_s3_dev() {
    setup
    cp -R ../tests/local-to-s3/* ./
    ../buckmate-executable --path buckmate --env dev apply
    ../buckmate-executable --path result apply
    compare_and_clean "../result-dev" "result/data"
}

s3_to_s3_dev() {
    setup
    cp -R ../tests/s3-to-s3/* ./
    ../buckmate-executable --path preload apply
    ../buckmate-executable --path buckmate --env dev apply
    ../buckmate-executable --path result apply
    compare_and_clean "../result-dev" "result/data"
}

full_cleanup

build

local_to_local
s3_to_s3
local_to_s3
s3_to_local
dry_local
dry_remote
cache_control_metadata
keep_previous_remote
keep_previous_local
local_to_local_dev
local_to_s3
s3_to_s3_dev

full_cleanup
