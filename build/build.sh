#!/bin/bash
set -e
export MSYS_NO_PATHCONV=1
self_dir=$(cd $(dirname $0); pwd)
project_dir=$(dirname $self_dir)
cd $project_dir

go build -o ./debug/

