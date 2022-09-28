#!/bin/bash
# find . -path ‘./pkg*mock*.go’ -exec rm -rf {} \;
find . -path './proto*.pb.go' -exec rm -rf {} \;
find . -path './sdk/android/bucketeer/src/main/proto/proto/*BUILD.bazel' -exec rm -rf {} \;
rm -fr sdk/android/bucketeer/build/
rm -fr sdk/android/bucketeer/src/main/proto/proto/
rm -fr vendor