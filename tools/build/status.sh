#!/bin/bash

echo HASH $(git rev-parse --verify HEAD)
echo BUILDDATE $(date '+%Y/%m/%d %H:%M:%S %Z')
